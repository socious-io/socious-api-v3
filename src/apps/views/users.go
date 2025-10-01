package views

import (
	"context"
	"net/http"
	"socious/src/apps/models"
	"socious/src/apps/utils"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")

	g.GET("", LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		c.JSON(http.StatusOK, user)
	})

	g.GET("/by-username/:username", func(c *gin.Context) {
		username := c.Param("username")

		user, err := models.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u := new(models.PublicUser)
		utils.Copy(user, u)

		c.JSON(http.StatusOK, u)
	})

	g.PUT("/wallets", LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		form := new(WalletForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		wallet := &models.Wallet{
			Network: form.Network,
			Address: form.Address,
			Testnet: form.Testnet,
			UserID:  user.ID,
		}

		if err := wallet.Upsert(context.Background()); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, wallet)
	})
}
