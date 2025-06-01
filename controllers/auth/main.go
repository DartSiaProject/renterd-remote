package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	constants "renterd-remote/constant"
	"renterd-remote/models"
	"renterd-remote/responseUtils"
	"renterd-remote/services/auth"
	"renterd-remote/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var user models.User
	rec := httptest.NewRecorder()

	// Call BindJSON to bind the received JSON to
	// User data
	if err := c.BindJSON(&user); err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, err.Error(), constants.BadRequest)
		//c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": constants.BadRequest})
		return
	}

	if user.Email == os.Getenv("USER_EMAIL") && utils.CreateSecretKey(user.Email, user.Password) == os.Getenv("USER_KEY") {
		tokenString, err := auth.CreateToken(user.Email)
		if err != nil {
			responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
			//c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": constants.InternalServerError})
			return
		}

		responseUtils.SuccessJsonResponse(rec, c, http.StatusOK, map[string]any{"AccessToken": tokenString}, constants.SuccessMessage)
		/*c.IndentedJSON(http.StatusOK, gin.H{
			"message": constants.SuccessMessage,
			"data":    map[string]any{"AccessToken": tokenString},
		})*/
		return
	} else {
		responseUtils.ErrorResponse(rec, c, http.StatusUnauthorized, constants.InvalidCredentials, constants.Unauthorized)
		return
		//c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": constants.InvalidCredentials})
	}
}
