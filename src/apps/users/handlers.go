package users

import (
	"github.com/gin-gonic/gin"
)

func getProfile(c *gin.Context) {
	id, _ := c.Get("user_id")
	u, _ := Get(id.(string))
	c.JSON(200, u)
}
