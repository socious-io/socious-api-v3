package apps

import (
	"fmt"
	"socious/src/apps/users"
	"socious/src/config"

	"github.com/gin-gonic/gin"
)

func Serve() {
	router := gin.Default()
	users.Register(router)
	router.Run(fmt.Sprintf("127.0.0.1:%d", config.Config.Port))
}
