package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	encryptBody, err := utils.GetAESEncrypted(res.Body.Bytes())
	if err != nil {
		fmt.Println(constants.BodyRequestEncryptionError, " : z ", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
		return err
	}

	c.Header("Header", encryptHeader)
	c.IndentedJSON(res.Result().StatusCode, gin.H{"data": encryptBody})
	return nil
}

// Decrypt request's data
func DecryptRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := c.Request.URL.Query()["params"]
		//fmt.Println("params :", c.Request.URL.Query())
		if len(params) > 0 {
			//Decrypt params in format "params1=value&params2=value
			decryptParams, err := utils.GetAESDecrypted(params[0])
			if err != nil {
				fmt.Println(constants.HeaderRequestDecryptionError, " : ", err)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
				return
			}

			//Set request params
			c.Request.URL, _ = url.Parse(c.Request.URL.Path + "?" + string(decryptParams))
		}

		//Get Header request using the field Header in request header
		//header := c.Request.Header.Values("Header")
		header := c.Request.URL.Query()["header"]
		if len(header) > 0 {
			//Descrypt Header request
			decryptHeader, err := utils.GetAESDecrypted(header[0])
			if err != nil {
				fmt.Println(constants.HeaderRequestDecryptionError, " : ", err)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
				return
			}

			//Set Header
			c.Request.Header.Del("Content-Type")
			c.Request.Header.Add("Content-Type", utils.StringToJSON(string(decryptHeader)).ContentType)
		}

		//fmt.Println("Test : ", string(decryptHeader))
		//Get body
		body, _ := io.ReadAll(c.Request.Body)

		// Map to hold the JSON data using fields Data
		var bodyData map[string]string

		if len(body) > 0 {
			//Descrypt body
			err1 := json.Unmarshal(body, &bodyData)
			if err1 != nil {
				fmt.Println(constants.BodyRequestDecryptionError, " : ", err1)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
				return
			}

			decryptBody, err := utils.GetAESDecrypted(bodyData["data"])
			if err != nil {
				fmt.Println(constants.BodyRequestDecryptionError, " : ", err)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": constants.Unauthorized, "message": constants.Unauthorized})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(decryptBody))

			c.Request.Header.Del("Content-Length")
			c.Request.Header.Add("Content-Length", fmt.Sprintf("%d", len(decryptBody)))
			c.Request.ContentLength = int64(len(decryptBody))

			/*err1 := auth.VerifyToken(auth.ExtractToken(c))
			if err1 != nil {
				c.String(http.StatusUnauthorized, constants.Unauthorized)
				c.Abort()
			}*/
		}
		c.Next()
	}
}
