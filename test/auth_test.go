package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	auth "renterd-remote/controllers/auth"
	authService "renterd-remote/services/auth"
	testContext "renterd-remote/test/context"
	"testing"

	responseModels "renterd-remote/test/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Define a custom struct to match the expected JSON response
var response responseModels.SimpleResponse

func TestLoginHandler(t *testing.T) {

	//Setup environment variables
	closer := testContext.EnvSetter()
	t.Cleanup(closer)

	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "falseEmail@false.cm",
		Password: "FalsePassword",
	}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(out))
	assert.NoError(t, err, "creating request should not fail")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		auth.Login(c)
	})

	// Set the request method and URL parameters
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "handler should return 200 OK")

	// Check the content type
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"),
		"content type should be application/json")

	// Check the response body
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "unmarshaling response should not fail")
	assert.Equal(t, "Invalid credentials", response.Message, "message should match")

	// Check the token creation
	token, err := authService.CreateToken(body.Email)

	assert.NoError(t, err, "creating token should not fail")
	assert.NotEmpty(t, token, "token should not be empty")

	//Check the token generate
	decodeTokenError := authService.VerifyToken(token)
	assert.NoError(t, decodeTokenError, "token should be valid")
}
