package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	ssoGroup(r)
	servicesGroup(r)
}
