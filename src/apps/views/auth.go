package views

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/lib"
	"socious/src/apps/models"
	"socious/src/apps/utils"
	"socious/src/config"

	"github.com/gin-gonic/gin"
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

		//Create auth session
		response, err := lib.HTTPRequest(lib.HTTPRequestOptions{
			Endpoint: fmt.Sprintf("%s/auth/session", config.Config.SSO.Host),
			Method:   lib.HTTPRequestMethodPost,
			Body: map[string]any{
				"client_id":     config.Config.SSO.ID,
				"client_secret": config.Config.SSO.Secret,
				"redirect_url":  redirect_url, //NOTE: if needs redirection within backend fmt.Sprintf("%s/auth/login/callback", config.Config.Host),
			},
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authSession := new(auth.AuthSessionResponse)
		if err := json.Unmarshal(response, &authSession); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf(
				"%s/auth/session/%s?auth_mode=login",
				config.Config.SSO.Host,
				authSession.AuthSession.ID,
			),
		)
	})

	g.POST("/token", func(c *gin.Context) {

		code, status := c.Query("code"), c.Query("status")
		ctx, _ := c.Get("ctx")

		if status != "success" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authentication status is not success"})
			return
		}

		//Get the token from Socious ID
		response, err := lib.HTTPRequest(lib.HTTPRequestOptions{
			Endpoint: fmt.Sprintf("%s/auth/session/token", config.Config.SSO.Host),
			Method:   lib.HTTPRequestMethodPost,
			Body: map[string]any{
				"client_id":     config.Config.SSO.ID,
				"client_secret": config.Config.SSO.Secret,
				"code":          code,
			},
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sessionToken := new(auth.SessionTokenResponse)
		if err := json.Unmarshal(response, &sessionToken); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get User's information from Socious ID
		response, err = lib.HTTPRequest(lib.HTTPRequestOptions{
			Endpoint: fmt.Sprintf("%s/users", config.Config.SSO.Host),
			Method:   lib.HTTPRequestMethodGet,
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", sessionToken.AccessToken),
			},
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := new(models.User)
		if err := json.Unmarshal(response, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Create token for front end to communicate with this platform
		tokens, err := auth.GenerateFullTokens(user.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		oauthConnect, err := models.GetOauthConnectByIdentityId(user.ID, models.OauthConnectedProvidersSociousId)

		if err != nil && oauthConnect == nil {
			oauthConnect = &models.OauthConnect{
				AccessToken:    sessionToken.AccessToken,
				RefreshToken:   &sessionToken.RefreshToken,
				MatrixUniqueId: tokens["access_token"].(string),
				Provider:       models.OauthConnectedProvidersSociousId,
				IdentityId:     user.ID,
			}
			err = oauthConnect.Create(ctx.(context.Context))
		} else if err != nil && oauthConnect != nil {
			oauthConnect.MatrixUniqueId = tokens["access_token"].(string)
			oauthConnect.AccessToken = sessionToken.AccessToken
			oauthConnect.RefreshToken = &sessionToken.RefreshToken
			err = oauthConnect.Update(ctx.(context.Context))
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
		return

	})

}
