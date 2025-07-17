package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	authGroup(r)
	projectsGroup(r)
	contractsGroup(r)
	usersGroup(r)
	organizationsGroup(r)
	syncGroup(r)
	identitiesGroup(r)
}
