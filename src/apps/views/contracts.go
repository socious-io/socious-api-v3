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
		identity := c.MustGet("identity").(*models.Identity)
		page, _ := c.Get("paginate")

		contracts, total, err := models.GetContracts(identity.ID, page.(database.Paginate))
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
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(ContractForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contract := new(models.Contract)
		utils.Copy(form, contract)
		contract.ProviderID = identity.ID
		contract.ClientID = form.ClientID
		if err := contract.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, contract)
	})

	g.PATCH("/:id", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx, _ := c.Get("ctx")

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

		if contract.ProviderID != identity.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}

		utils.Copy(form, contract)

		if err := contract.Update(ctx.(context.Context), []uuid.UUID{}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, contract)
	})

	g.POST("/:id/sign", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx, _ := c.Get("ctx")

		id := c.Param("id")

		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if contract.ClientID != identity.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Just client can sign the contract"})
			return
		}
		contract.Status = models.ContractStatusSinged

		if err := contract.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, contract)
	})

	g.POST("/:id/cancel", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx, _ := c.Get("ctx")

		id := c.Param("id")

		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if contract.ProviderID == identity.ID {
			contract.Status = models.ContractStatusProviderCanceled
		} else if contract.ClientID == identity.ID {
			contract.Status = models.ContractStatusClientCanceled
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provider or Client don't match identity"})
			return
		}

		if err := contract.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, contract)
	})

	g.POST("/:id/apply", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx, _ := c.Get("ctx")

		id := c.Param("id")

		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if contract.ClientID != identity.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Just client can set the contract to applied"})
			return
		}
		contract.Status = models.ContractStatusApplied

		if err := contract.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, contract)
	})

	g.POST("/:id/complete", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx, _ := c.Get("ctx")

		id := c.Param("id")

		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if contract.ProviderID != identity.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Just provider can complete the contract"})
			return
		}
		contract.Status = models.ContractStatusCompleted

		if err := contract.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, contract)
	})

	g.DELETE("/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		contract := models.Contract{
			ID: uuid.MustParse(id),
		}
		err := contract.Delete(ctx.(context.Context))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.POST("/:id/deposit", func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		ctx, _ := c.Get("ctx")

		id := c.Param("id")

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
		provider, err := models.GetIdentity(contract.ProviderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("provider fetch error : %v", err)})
			return
		}

		//Determine Currency
		var currency gopay.Currency
		if contract.Currency == nil && *contract.PaymentType == models.PaymentModeTypeFiat {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Currency is nil in Fiat payment : %v", err)})
			return
		} else if contract.Currency == nil {
			//Default payment is set not to prevent the runtime from crashing while its empty
			currency = gopay.JPY
		} else {
			currency = gopay.Currency(*contract.Currency)
		}

		//Start a payment session
		payment, err := gopay.New(gopay.PaymentParams{
			Tag:         contract.Name,
			Description: *contract.Description,
			Ref:         contract.ID.String(),
			Type:        gopay.PaymentType(*contract.PaymentType),
			Currency:    currency,
			TotalAmount: contract.TotalAmount,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var sourceAccount, destinationAccount *string
		if *contract.PaymentType == models.PaymentModeTypeFiat {
			//Set Source account
			card, err := models.GetCard(*form.CardID, contract.ProviderID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find corresponding Stripe customer"})
				return
			}
			sourceAccount = card.Customer

			//Set Destination account
			oauthConnect, err := models.GetOauthConnectByIdentityId(contract.ClientID, models.OauthConnectedProvidersStripeJp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find corresponding Stripe account"})
				return
			}
			destinationAccount = &oauthConnect.MatrixUniqueId

			payment.SetToFiatMode(string(oauthConnect.Provider))
		} else {
			walletAddress, ok := provider.MetaMap["wallet_address"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing wallet address on provider"})
				return
			}
			sourceAccount = &walletAddress
			payment.SetToCryptoMode(*contract.CryptoCurrency, float64(contract.CurrencyRate))
		}

		//Add Payment Identities
		if _, err := payment.AddIdentity(gopay.IdentityParams{
			ID:       identity.ID,
			RoleName: "assigner",
			Account:  *sourceAccount,
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
				Account:  *destinationAccount,
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
			err = payment.ConfirmDeposit(*form.TxID, form.Meta)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Updating contract
		contract.PaymentID = &payment.ID
		err = contract.Update(ctx.(context.Context), []uuid.UUID{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, contract)
	})

	g.PATCH("/:id/requirements", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		form := new(ContractRequirementsForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		contract, err := models.GetContract(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contract.RequirementDescription = &form.RequirementDescription
		if err := contract.Update(ctx.(context.Context), form.RequirementFiles); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusAccepted, contract)

	})

}
