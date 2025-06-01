package responseUtils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"renterd-remote/middlewares"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(rec *httptest.ResponseRecorder, c *gin.Context, statusCode int, errorMessage string, message string) {
	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(statusCode)
	rec.Body.Write([]byte(`{"message":` + message + `, "error": "` + errorMessage + `"}`))

	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)
}

func SuccessJsonResponse(rec *httptest.ResponseRecorder, c *gin.Context, statusCode int, data map[string]any, message string) {
	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(statusCode)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		rec.WriteHeader(http.StatusInternalServerError)
		rec.Body.Write([]byte(`{"message":` + message + `, "error": "Failed to marshal data"}`))
	} else {
		rec.Body.Write([]byte(`{"message":` + message + `, "data": ` + string(dataJSON) + `}`))
	}

	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)
}
