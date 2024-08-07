package users

import (
	"socious/src/apps/auth"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/profile", auth.LoginRequired(), getProfile)
	}
}
