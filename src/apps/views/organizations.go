package views

import (
	"net/http"

	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func organizationsGroup(router *gin.Engine) {
	g := router.Group("organizations")
	g.Use(LoginRequired())

	g.GET("", func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		orgs, err := models.GetUserOrganizations(u.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"organizations": orgs})
	})

	g.GET("/:id", func(c *gin.Context) {
		org, err := models.GetOrganization(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		}
		c.JSON(http.StatusOK, org)
	})
}
