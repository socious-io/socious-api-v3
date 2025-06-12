package views

import (
	"net/http"
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
)

func identitiesGroup(router *gin.Engine) {
	g := router.Group("identities")
	g.Use(LoginRequired())

	g.GET("", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		identity := c.MustGet("identity").(*models.Identity)

		identities, err := models.GetAllIdentities(user.ID, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, identities)
	})

	// g.GET("/:id", func(c *gin.Context) {
	// 	id := uuid.MustParse(c.Param("id"))
	// 	identity := c.MustGet("identity").(*models.Identity)

	// 	if identity.ID != id {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, identity)
	// })
}
