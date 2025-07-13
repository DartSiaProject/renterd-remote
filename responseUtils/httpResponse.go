package responseUtils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"renterd-remote/middlewares/encryptMiddleware"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(rec *httptest.ResponseRecorder, c *gin.Context, statusCode int, errorMessage string, message string) {
	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(statusCode)
	rec.Body.Write([]byte(`{"message":"` + message + `", "error": "` + errorMessage + `"}`))

	//Transfert response to encrypt middelware
	encryptMiddleware.EncryptResponse(rec, c)
}

func SuccessJsonResponse(rec *httptest.ResponseRecorder, c *gin.Context, statusCode int, data map[string]any, message string) {
	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(statusCode)

	response := map[string]any{
		"message": message,
		"data":    data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		rec.WriteHeader(http.StatusInternalServerError)
		rec.Body.Write([]byte(`{"message":"` + message + `", "error": "Failed to marshal data"}`))
		fmt.Println("Error marshaling response data:", err)
		return
	}

	rec.Body.Write(responseJSON)

	//Transfert response to encrypt middelware
	encryptMiddleware.EncryptResponse(rec, c)
}
