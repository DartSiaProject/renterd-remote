package auth

import (
	"net/http"
	"os"
	constants "renterd-remote/constant"
	"renterd-remote/models"
	"renterd-remote/services/auth"
	"renterd-remote/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var user models.User

	// Call BindJSON to bind the received JSON to
	// User data
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": constants.BadRequest})
		return
	}

	if user.Email == os.Getenv("USER_EMAIL") && utils.CreateSecretKey(user.Email, user.Password) == os.Getenv("USER_KEY") {
		tokenString, err := auth.CreateToken(user.Email)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": constants.InternalServerError})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": constants.SuccessMessage,
			"data":    map[string]any{"AccessToken": tokenString},
		})
		return
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": constants.InvalidCredentials})
	}
}
