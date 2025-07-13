package encryptMiddleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	constants "renterd-remote/constant"
	"renterd-remote/utils"

	"github.com/gin-gonic/gin"
)

// Encrypt request's data
func EncryptResponse(res *httptest.ResponseRecorder, c *gin.Context) error {
	encryptHeader, err := utils.GetAESEncrypted([]byte(utils.HttpHeaderMapToString(res.Result().Header)))
	if err != nil {
		fmt.Println(constants.HeaderRequestEncryptionError, " : ", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
		return err
	}

	responseBody := res.Body.Bytes()
	if res.Result().Header.Get("Content-Type") == "application/json" {
		if len(responseBody) > 0 && responseBody[len(responseBody)-1] == '\n' {
			// Remove the newline character
			responseBody = responseBody[:len(responseBody)-1]
		}
	}

	encryptBody, err := utils.GetAESEncrypted(responseBody)
	if err != nil {
		fmt.Println(constants.BodyRequestEncryptionError, " : z ", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
		return err
	}

	c.Header("Header", encryptHeader)
	c.IndentedJSON(res.Result().StatusCode, gin.H{"data": encryptBody})
	return nil
}
