package renterd

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"renterd-remote/middlewares"

	"github.com/gin-gonic/gin"
)

// Transfert function to renterd
func ReverseProxy(c *gin.Context) {
	director := func(req *http.Request) {
		r := c.Request

		req.URL.Scheme = "http"
		req.URL.Host = os.Getenv("RENTERD_ADDRESS")
		req.Header["my-header"] = []string{r.Header.Get("my-header")}
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+os.Getenv("RENTERD_KEY"))))
		// Golang camelcases headers
		delete(req.Header, "My-Header")
	}

	proxy := &httputil.ReverseProxy{Director: director}
	rec := httptest.NewRecorder()
	proxy.ServeHTTP(rec, c.Request)

	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)
}
