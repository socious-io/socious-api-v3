package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	ssoGroup(r)
	authGroup(r)
	projectsGroup(r)
	contractsGroup(r)
	usersGroup(r)
	mediaGroup(r)
}
