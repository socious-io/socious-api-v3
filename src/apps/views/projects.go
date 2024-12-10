package views

import (
	"context"
	"net/http"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

func projectsGroup(router *gin.Engine) {
	g := router.Group("projects")
	g.Use(auth.LoginRequired())

	g.GET("/services", paginate(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")
		pagination := page.(database.Paginate)
		pagination.Filters = []database.Filter{
			{
				Key:   "kind",
				Value: "SERVICE",
			},
		}

		services, total, err := models.GetProjects(u.(*models.User).ID, pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": services,
			"total":   total,
		})
	})

	g.GET("/services/:id", func(c *gin.Context) {
		id := c.Param("id")

		s, err := models.GetProject(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	g.POST("/services", func(c *gin.Context) {
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
		s.Kind = models.ProjectKindService
		s.CommitmentHoursLower, s.CommitmentHoursHigher = &form.TotalHours, &form.TotalHours
		s.PaymentRangeLower, s.PaymentRangeHigher = &form.Price, &form.Price
		if err := s.Create(ctx.(context.Context), form.WorkSamples); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	g.PATCH("/services/:id", func(c *gin.Context) {
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
		s.CommitmentHoursLower, s.CommitmentHoursHigher = &form.TotalHours, &form.TotalHours
		s.PaymentRangeLower, s.PaymentRangeHigher = &form.Price, &form.Price
		if err := s.Update(ctx.(context.Context), form.WorkSamples); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	g.DELETE("/services/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		s, err := models.GetProject(uuid.MustParse(id))
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
