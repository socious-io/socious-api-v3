package views

import (
	"context"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"

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
		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")

		form := new(ContractForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contract := new(models.Contract)
		utils.Copy(form, contract)
		contract.ProviderID = u.(*models.User).ID

		if err := contract.Create(ctx.(context.Context)); err != nil {
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
}
