package views

import (
	"context"
	"fmt"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"

	"github.com/socious-io/gopay"
	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func contractsGroup(router *gin.Engine) {
	g := router.Group("contracts")
	g.Use(auth.LoginRequired())

	g.GET("", paginate(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")

		contracts, total, err := models.GetContracts(u.(*models.User).ID, page.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": contracts,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")

		s, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	g.POST("", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		u := c.MustGet("user").(*models.User)

		form := new(ContractForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contract := new(models.Contract)
		utils.Copy(form, contract)
		contract.ProviderID = u.ID

		if err := contract.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, contract)
	})

	g.PATCH("/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")
		id := c.Param("id")

		form := new(ContractForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if contract.ProviderID != u.(*models.User).ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}

		utils.Copy(form, contract)

		if err := contract.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, contract)
	})

	g.POST("/:id/deposit", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		id := c.Param("id")
		ctx, _ := c.Get("ctx")

		form := new(ContractDepositForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Fetching Contract
		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Fetching Client
		client, err := models.GetUser(contract.ClientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("client fetch error : %v", err)})
			return
		}

		//Start a payment session
		payment, err := gopay.New(gopay.PaymentParams{
			Tag:         contract.Name,
			Description: *contract.Description,
			Ref:         contract.ID.String(),
			Type:        gopay.PaymentType(*contract.PaymentType),
			Currency:    gopay.Currency(*contract.Currency),
			TotalAmount: contract.TotalAmount,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var source_account, destination_account *string
		if *contract.PaymentType == models.PaymentModeTypeFiat {

			//Set Source account
			customerCard, err := models.GetCard(*form.CardID, contract.ProviderID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find corresponding Stripe customer"})
			}
			source_account = customerCard.Customer

			//Set Destination account
			oauthConnect, err := models.GetOauthConnectByIdentityId(contract.ClientID, models.OauthConnectedProvidersStripeJp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find corresponding Stripe account"})
			}
			destination_account = &oauthConnect.MatrixUniqueId

			payment.SetToFiatMode(string(oauthConnect.Provider))
		} else {
			source_account = client.WalletAddress
			if source_account == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing wallet address on client"})
			}
			payment.SetToCryptoMode(*contract.CryptoCurrency, float64(contract.CurrencyRate))
		}

		//Add Payment Identities
		if _, err := payment.AddIdentity(gopay.IdentityParams{
			ID:       identity.ID,
			RoleName: "assigner",
			Account:  *source_account,
			Amount:   0,
		}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Only fiat payment needs second payment identity
		if *contract.PaymentType == models.PaymentModeTypeFiat {
			if _, err := payment.AddIdentity(gopay.IdentityParams{
				ID:       identity.ID,
				RoleName: "assignee",
				Account:  *destination_account,
				Amount:   float64(contract.TotalAmount),
			}); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		//Enroll the payment
		if *contract.PaymentType == models.PaymentModeTypeFiat {
			err = payment.Deposit()
		} else {
			err = payment.ConfirmDeposit(*form.TxID)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Updating contract
		contract.PaymentID = &payment.ID
		contract.Status = models.ContractStatusSinged
		err = contract.Update(ctx.(context.Context))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, contract)
	})
}
