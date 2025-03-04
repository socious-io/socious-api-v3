package views

import (
	"context"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"

	"github.com/gin-gonic/gin"
	"github.com/socious-io/goaccount"
)

func authGroup(router *gin.Engine) {
	g := router.Group("auth")

	g.POST("/register", func(c *gin.Context) {
		form := new(auth.RegisterForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := new(models.User)
		utils.Copy(form, u)
		if form.Password != nil {
			password, _ := auth.HashPassword(*form.Password)
			u.Password = &password
		}

		ctx, _ := c.Get("ctx")
		if err := u.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokens, err := auth.GenerateFullTokens(u.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
	})

	g.GET("/login", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		_, entrypoint, err := goaccount.StartSession(redirect_url)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, entrypoint)
	})

	g.POST("/token", func(c *gin.Context) {

		code, status := c.Query("code"), c.Query("status")
		ctx, _ := c.Get("ctx")

		if status != "success" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authentication status is not success"})
			return
		}

		//Get the token from Socious ID
		sessionToken, err := goaccount.GetSessionToken(code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get User's information from Socious ID
		sessionUser := new(models.User)
		err = sessionToken.GetUserProfile(&sessionUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := models.GetUserByEmail(sessionUser.Email)
		if err != nil {
			//Try to create user if doesn't exist
			user := new(models.User)
			utils.Copy(sessionUser, user)
			err = user.Create(ctx.(context.Context))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		//Create token for front end to communicate with this platform
		tokens, err := auth.GenerateFullTokens(user.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		oauthConnect, err := models.GetOauthConnectByEmail(user.Email, models.OauthConnectedProvidersSociousId)

		if err != nil && oauthConnect == nil {
			oauthConnect = &models.OauthConnect{
				AccessToken:    sessionToken.AccessToken,
				RefreshToken:   &sessionToken.RefreshToken,
				MatrixUniqueId: tokens["access_token"].(string),
				Provider:       models.OauthConnectedProvidersSociousId,
				IdentityId:     user.ID,
			}
			err = oauthConnect.Create(ctx.(context.Context))
		} else if oauthConnect != nil {
			oauthConnect.MatrixUniqueId = tokens["access_token"].(string)
			oauthConnect.AccessToken = sessionToken.AccessToken
			oauthConnect.RefreshToken = &sessionToken.RefreshToken
			err = oauthConnect.Update(ctx.(context.Context))
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
		return

	})

}
