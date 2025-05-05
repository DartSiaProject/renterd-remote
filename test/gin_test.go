package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	testContext "renterd-remote/test/context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// code
func GetId(c *gin.Context) {

	//fmt.Println(c.Query("foo")) //will print "bar" while running test
	//fmt.Println(c.Param("id"))  // will print "1" while running test

	id, _ := strconv.Atoi(c.Param("id"))
	c.JSON(http.StatusOK, id)
}

// test
func TestGetId(t *testing.T) {
	w := httptest.NewRecorder()

	ctx := testContext.GetTestGinContext(w)

	//configure path params
	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	// configure query params
	u := url.Values{}
	u.Add("foo", "bar")

	testContext.MockJsonGet(ctx, params, u)

	GetId(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

	got, _ := strconv.Atoi(w.Body.String())

	assert.Equal(t, 1, got)
}
