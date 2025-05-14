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

	//Save sqlite db
	renterd_Routes.GET("/renterd/savedb", renterd.SaveSqliteDb)

	//Restore sqlite db
	renterd_Routes.POST("/renterd/restoredb", renterd.RestoreSqliteDb)

	//Generate File Share Link
	renterd_Routes.POST("/renterd/sharelink", renterd.GetShareLink)

	//Get File Share Content
	renterd_Routes.GET("/renterd/sharefile/:key", renterd.GetShareFile)
}
