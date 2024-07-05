package auth

import (
	"renterd-remote/controllers/auth"

	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {
	//Declare group
	user := route.Group("/auth")

	//Add routes to router Group
	user.GET("/login", auth.Login)
	//user.POST("/logout", auth.GetAlbums)
}
