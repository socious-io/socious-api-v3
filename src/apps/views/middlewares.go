package views

import (
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/config"
	"strconv"
	"strings"

	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
)

func paginate() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		if page < 1 {
			page = 1
		}
		if limit > 100 || limit < 1 {
			limit = 10
		}
		filters := make([]database.Filter, 0)
		for key, values := range c.Request.URL.Query() {
			if strings.Contains(key, "filter.") && len(values) > 0 {
				filters = append(filters, database.Filter{
					Key:   strings.Replace(key, "filter.", "", -1),
					Value: values[0],
				})
			}
		}

		c.Set("paginate", database.Paginate{
			Limit:   limit,
			Offet:   (page - 1) * limit,
			Filters: filters,
		})
		c.Set("limit", limit)
		c.Set("page", page)
		c.Next()

	}
}

func sociousIdSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		//Fetching Socious ID token
		oauthConnect, err := models.GetOauthConnectByEmail(user.Email, models.OauthConnectedProvidersSociousId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("socious_id_session", oauthConnect.SociousIdSession())
	}
}

func AccountCenterRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("x-account-center-id")
		secret := c.Request.Header.Get("x-account-center-secret")
		hash, _ := auth.HashPassword(secret)

		if id != config.Config.GoAccounts.ID || auth.CheckPasswordHash(secret, hash) != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account center required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
