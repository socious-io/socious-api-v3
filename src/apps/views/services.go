package views

import (
	"context"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"
	"socious/src/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func servicesGroup(router *gin.Engine) {
	g := router.Group("services")
	g.Use(auth.LoginRequired())

	g.GET("", paginate(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")

		services, total, err := models.GetServices(u.(*models.User).ID, page.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": services,
			"total":   total,
		})
	})
	g.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")

		s, err := models.GetService(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})
	g.POST("", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")

		form := new(ServiceForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s := new(models.Project)
		utils.Copy(form, s)
		s.IdentityID = u.(*models.User).ID
		if err := s.CreateService(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	g.PATCH("/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		form := new(ServiceForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s := new(models.Project)
		utils.Copy(form, s)
		s.ID = uuid.MustParse(id)
		if err := s.UpdateService(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	g.DELETE("/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		s, err := models.GetService(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := s.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
