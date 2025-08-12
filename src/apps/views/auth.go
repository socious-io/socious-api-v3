package views

import (
	"context"
	"log"
	"net/http"
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/socious-io/goaccount"
)

func authGroup(router *gin.Engine) {
	g := router.Group("auth")

	g.POST("", func(c *gin.Context) {
		form := new(AuthForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionPolicies := []goaccount.PolicyType{}
		if form.OrgOnboarding {
			sessionPolicies = append(sessionPolicies, goaccount.PolicyTypeEnforceOrgCreation)
		}

		session, authURL, err := goaccount.StartSession(
			form.RedirectURL,
			form.AuthMode,
			sessionPolicies,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"session":  session,
			"auth_url": authURL,
		})
	})

	g.POST("/session", func(c *gin.Context) {
		form := new(SessionForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := goaccount.GetSessionToken(form.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var (
			connect *models.OauthConnect
			user    = new(models.User)
			ctx     = c.MustGet("ctx").(context.Context)
		)

		goaccountUser, err := token.GetUserProfile()
		user = models.GetTransformedUser(ctx, *goaccountUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := user.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := user.AttachMedia(ctx, *goaccountUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if connect, err = models.GetOauthConnectByMUI(user.ID.String(), models.OauthConnectedProvidersSociousID); err != nil {
			connect = &models.OauthConnect{
				Provider:       models.OauthConnectedProvidersSociousID,
				AccessToken:    token.AccessToken,
				RefreshToken:   &token.RefreshToken,
				MatrixUniqueID: user.ID.String(),
				IdentityId:     user.ID,
			}
		}

		orgs, err := token.GetMyOrganizations()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, o := range orgs {
			org := models.GetTransformedOrganization(ctx, o)
			if err := org.Upsert(ctx, user.ID); err != nil {
				log.Println(err.Error(), o)
			}
			if err := org.AttachMedia(ctx, o, user.ID); err != nil {
				log.Println(err.Error(), o)
			}
		}

		if err := connect.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jwt, err := goaccount.GenerateFullTokens(user.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, jwt)
	})

	g.POST("/refresh", func(c *gin.Context) {
		form := new(RefreshForm)
		if err := c.Bind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims, err := goaccount.VerifyToken(form.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jwt, err := goaccount.GenerateFullTokens(claims.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, jwt)
	})
}
