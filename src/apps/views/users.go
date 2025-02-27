package views

import (
	"context"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"

	"github.com/gin-gonic/gin"
	sociousid "github.com/socious-io/go-socious-id"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("/", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		ctx, _ := c.Get("ctx")

		//Fetching Socious ID token
		oauthConnect, err := models.GetOauthConnectByIdentityId(user.ID, models.OauthConnectedProvidersSociousId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get User's information from Socious ID
		userSociousID := new(models.User)
		err = sociousid.GetUserProfile(oauthConnect.AccessToken, &userSociousID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Updating user on local
		utils.Copy(userSociousID, user)
		err = user.UpdateProfile(ctx.(context.Context))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)

	})

	g.PUT("/hook", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")

		form := new(models.User)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Fetching Socious ID token
		oauthConnect, err := models.GetOauthConnectByEmail(form.Email, models.OauthConnectedProvidersSociousId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get User's information from Socious ID
		userSociousID := new(models.User)
		err = sociousid.GetUserProfile(oauthConnect.AccessToken, &userSociousID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Updating user on local
		user := new(models.User)
		utils.Copy(userSociousID, user)
		err = user.UpdateProfile(ctx.(context.Context))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

}
