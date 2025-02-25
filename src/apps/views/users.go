package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/lib"
	"socious/src/apps/models"
	"socious/src/config"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("/", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		//Fetching Socious ID token
		oauthConnect, err := models.GetOauthConnectByIdentityId(user.ID, models.OauthConnectedProvidersSociousId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get User's information from Socious ID
		response, err := lib.HTTPRequest(lib.HTTPRequestOptions{
			Endpoint: fmt.Sprintf("%s/users", config.Config.SSO.Host),
			Method:   lib.HTTPRequestMethodGet,
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", oauthConnect.AccessToken),
			},
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user = new(models.User)
		if err := json.Unmarshal(response, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

}
