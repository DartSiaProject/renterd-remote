package context

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// mock gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// mock getrequest
func MockJsonGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)

	// set path params
	c.Params = params

	// set query params
	c.Request.URL.RawQuery = u.Encode()
}

func EnvSetter() (closer func()) {
	originalEnvs := map[string]string{}
	envs := map[string]string{
		"SERVER_ADDRESS": "localhost",
		"SERVER_PORT":    "8000",
		"JWT_SECRET":     "1hwGE8Y6nHbPVRA9",

		"GIN_MODE": "release",

		"USER_EMAIL": "test@test.com",
		"USER_KEY":   "418D125DC1475C7817f31185981e4ebe",
		"USER_IV":    "300372804b842889",

		"RENTERD_ADDRESS": "localhost:19980",
		"RENTERD_KEY":     "Test@12345",
	}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}
