package renterd

import (
	"renterd-remote/controllers/renterd"

	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {
	//Declare group
	renterd_Routes := route.Group("/")

	//Add Middleware
	//renterd_Routes.Use(middlewares.JwtAuthMiddleware())
	//renterd_Routes.Use(middlewares.DecryptRequest())

	//Add routes to router Group
	renterd_Routes.GET("/api/bus/accounts", renterd.ReverseProxy)
}
