package views

import (
	"net/http"
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(LoginRequired())

	g.GET("", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		c.JSON(http.StatusOK, user)
	})
}
