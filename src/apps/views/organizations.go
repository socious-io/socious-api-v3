package views

import (
	"socious/src/apps/auth"

	"github.com/gin-gonic/gin"
)

func organizationsGroup(router *gin.Engine) {
	g := router.Group("organizations")
	g.Use(auth.LoginRequired())

	// g.GET("", paginate(), func(c *gin.Context) {
	// 	user := c.MustGet("user").(*models.User)
	// 	page, _ := c.Get("paginate")

	// })

	// g.GET("/:id", func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)

	// })

	// g.GET("/my", paginate(), func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)

	// })

	// g.POST("", auth.LoginRequired(), sociousIdSession(), paginate(), func(c *gin.Context) {
	// 	ctx := c.MustGet("ctx").(context.Context)
	// 	user := c.MustGet("user").(*models.User)
	// 	sociousIdSession := c.MustGet("socious_id_session").(goaccount.SessionToken)

	// 	form := new(OrganizationUpdateForm)
	// 	if err := c.ShouldBindJSON(form); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	//Updating user on local
	// 	userId := user.ID
	// 	utils.Copy(form, organization)

	// 	err := user.UpdateProfile(ctx, &sociousIdSession)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, user)
	// })

	// g.PUT("/:id", paginate(), func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)
	// 	page, _ := c.Get("paginate")

	// })

	// g.DELETE("/:id", paginate(), func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)
	// 	page, _ := c.Get("paginate")

	// })

	// g.POST("/:id/member/:user_id", paginate(), func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)
	// 	page, _ := c.Get("paginate")

	// })

	// g.DELETE("/:id/member/:user_id", paginate(), func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)
	// 	page, _ := c.Get("paginate")

	// })
}
