package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"renterd-remote/controllers/renterd"
	"renterd-remote/middlewares/decryptMiddleware"
	"renterd-remote/middlewares/encryptMiddleware"
	testContext "renterd-remote/test/context"

	models "renterd-remote/models"
	responseModels "renterd-remote/test/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var responseEncrypt responseModels.EncryptResponse

func TestEncryptTrafficHandler(t *testing.T) {

	//Setup environment variables
	closer := testContext.EnvSetter()
	t.Cleanup(closer)

	body := struct {
		Name   string            `json:"name"`
		Policy map[string]string `json:"policy"`
	}{
		Name: "o5x38u7186vs-xr7aymiru",
		Policy: map[string]string{
			"publicReadAccess": "true",
		},
	}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/bus/buckets", bytes.NewBuffer(out))
	assert.NoError(t, err, "creating request should not fail")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r

		encryptMiddleware.EncryptResponse(rr, c)
	})

	// Set the request method and URL parameters
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "handler should return 401 unauthorized")

	// Check the content type
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"),
		"content type should be application/json")

	// Check the response body
	err = json.Unmarshal(rr.Body.Bytes(), &responseEncrypt)
	assert.NoError(t, err, "unmarshaling response should not fail")
	assert.Equal(t, "", responseEncrypt.Data, "message should match")
}

func TestDecryptTrafficHandler(t *testing.T) {
	//Setup environment variables
	closer := testContext.EnvSetter()
	t.Cleanup(closer)

	body := struct {
		Data string `json:"data"`
	}{
		Data: "z2cZP727pGiwyPCzFDj927lvxLlasKP907a9H+IldUlY3HcA4sbTv4d9rT7P4gpSm1vyG5ZZAKd5vTNwWv4eR54N8WZWxab3Cw36nzuV9fM=",
	}

	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	// Set the request header
	engine.Use(decryptMiddleware.DecryptRequest())

	// Create a request
	req, err := http.NewRequestWithContext(ctx, "POST", "/", bytes.NewBuffer(out))
	assert.NoError(t, err, "creating request should not fail")

	// Set the request header
	req.Header.Set("Header", "FfUPsIplgP+5L9SOgGOQ2w==")

	// Set the request method and URL parameters
	engine.ServeHTTP(w, req)

	bodyBytes, _ := io.ReadAll(req.Body)

	// Check the content type
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"),
		"content type should be text/plain")

	// Check the response body
	var responseData models.User
	err = json.Unmarshal(bodyBytes, &responseData)
	assert.NoError(t, err, "unmarshaling response should not fail")

	assert.Equal(t, `{
		"Email":    "falseEmail@false.cm",
		"Password": "FalsePassword"
	}`, string(bodyBytes), "message should match")
}

func TestReverseProxyHandler(t *testing.T) {

	//Setup environment variables
	closer := testContext.EnvSetter()
	t.Cleanup(closer)

	body := struct {
		Name   string            `json:"name"`
		Policy map[string]string `json:"policy"`
	}{
		Name: "o5x38u7186vs-xr7aymiru",
		Policy: map[string]string{
			"publicReadAccess": "true",
		},
	}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/bus/buckets", bytes.NewBuffer(out))
	assert.NoError(t, err, "creating request should not fail")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r

		renterd.ReverseProxy(c)
	})

	// Set the request method and URL parameters
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadGateway, rr.Code, "handler should return 502 Bad Gateway")

}
