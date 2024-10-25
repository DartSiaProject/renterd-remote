package main

import (
	"os"
	"renterd-remote/config"
	"renterd-remote/controllers/renterd"
	"renterd-remote/middlewares"
	"renterd-remote/routes/auth"
	renterdRoutes "renterd-remote/routes/renterd"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	//utils.Test()
	if os.Getenv("USER_EMAIL") == "" || os.Getenv("USER_KEY") == "" {
		config.InitApp()
		LaunchWebServer()
	} else {
		LaunchWebServer()
	}
}

func LaunchWebServer() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	//Create router
	router := gin.Default()

	//Add Routes
	auth.Routes(router)
	renterdRoutes.Routes(router)

	router.Use(middlewares.DecryptRequest())
	//Redirect all route to renterd
	router.NoRoute(renterd.ReverseProxy)

	server_address := os.Getenv("SERVER_ADDRESS")
	server_port := os.Getenv("SERVER_PORT")
	if server_address != "" && server_port != "" {
		router.RunTLS(server_address+":"+server_port, "./config/ssl/server.pem", "./config/ssl/server.key")
	} else {
		router.RunTLS("localhost:8080", "./config/ssl/server.pem", "./config/ssl/server.key")
	}
}
