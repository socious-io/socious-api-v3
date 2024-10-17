package apps

import (
	"context"
	"fmt"
	"socious/src/apps/views"
	"socious/src/config"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve() {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		c.Set("ctx", ctx)
		c.Next()
	})

	views.Init(router)

	router.Run(fmt.Sprintf("127.0.0.1:%d", config.Config.Port))
}
