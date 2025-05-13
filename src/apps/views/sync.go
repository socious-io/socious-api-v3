package views

import (
	"context"
	"log"
	"net/http"
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
)

func syncGroup(router *gin.Engine) {
	g := router.Group("sync")
	g.Use(AccountCenterRequired())

	g.PUT("", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)

		form := new(SyncForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.GetTransformedUser(ctx, form.User)
		if err := user.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, o := range form.Organizations {
			organization := models.GetTransformedOrganization(ctx, o)
			if err := organization.Upsert(ctx, user.ID); err != nil {
				log.Println(err.Error(), o)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

}
