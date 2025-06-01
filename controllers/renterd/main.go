package renterd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	constants "renterd-remote/constant"
	"renterd-remote/middlewares"
	models "renterd-remote/models"
	"renterd-remote/responseUtils"
	utils "renterd-remote/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// Transfert function to renterd
func ReverseProxy(c *gin.Context) {
	director := func(req *http.Request) {
		//r := c.Request

		req.URL.Scheme = "http"
		req.URL.Host = os.Getenv("RENTERD_ADDRESS")
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+os.Getenv("RENTERD_KEY"))))
		// Golang camelcases headers
		//req.Header["my-header"] = []string{r.Header.Get("my-header")}
		//delete(req.Header, "My-Header")

	}

	proxy := &httputil.ReverseProxy{Director: director}
	rec := httptest.NewRecorder()

	// Ajout d'une gestion des erreurs pour capturer les erreurs de proxy
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, err.Error(), constants.ReverseProxyError)
		//c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": constants.ReverseProxyError})
		//Transfert response to encrypt middelware
		middlewares.EncryptResponse(rec, c)
	}

	proxy.ServeHTTP(rec, c.Request)
	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)
}

type BackupStruct struct {
	Database string `json:"database"`
	Path     string `json:"path"`
}

// Send SQlLite database to mobile app
func SaveSqliteDb(c *gin.Context) {
	rec := httptest.NewRecorder()

	// Get the SQLite database file path from the request body
	path := utils.GetDefaultSqliteBackupPath()

	// Body of the request
	payloadData := BackupStruct{
		Database: "main",
		Path:     path,
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON:", err)
		return
	}

	url := "http://" + os.Getenv("RENTERD_ADDRESS") + "/api/bus/system/sqlite3/backup"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatal(err)
		return
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+os.Getenv("RENTERD_KEY"))))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return
	}
	defer res.Body.Close()

	// Check if the file exists
	// Vérifiez si le fichier existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		println("File does not exist:", err.Error())
		// Handle the error as needed
		responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
		return
	}

	// Utilisez c.FileAttachment pour envoyer le fichier
	// Récupérer le fichier depuis la requête
	_, err = os.Stat(path)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
		return
	}

	src, err := os.Open(path)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
		return
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
		return
	}

	// Définissez les en-têtes pour indiquer que c'est un téléchargement de fichier
	rec.Header().Add("Content-Disposition", "attachment; filename=\"renterd.sqlite3.bak\"")
	rec.Header().Add("Content-Type", "application/octet-stream")
	rec.Header().Add("Cache-Control", "no-cache")
	rec.Header().Add("Pragma", "no-cache")
	rec.Header().Add("Expires", "0")

	rec.Body.Write(content)
	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)
}

// Restore SQlLite database from mobile app
func RestoreSqliteDb(c *gin.Context) {
	rec := httptest.NewRecorder()
	// Get the SQLite database file path from the request body
	path := utils.GetDefaultSqliteBackupPath()

	// Récupérer le fichier depuis la requête
	buffer := make([]byte, c.Request.ContentLength)
	_, err := c.Request.Body.Read(buffer)
	if err != nil && err != io.EOF {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, err.Error(), constants.BadRequest)
		return
	}
	file := bytes.NewReader(buffer)

	// Create the destination file
	dst, err := os.Create(path)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
		return
	}
	defer dst.Close()

	// Copy the contents of the downloaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
		return
	}

	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(http.StatusOK)
	rec.Body.Write([]byte(`{"message":` + constants.SqlliteRestoreSuccessMessage + `}`))

	// Transfert response to encrypt middleware
	middlewares.EncryptResponse(rec, c)
}

func GetShareLink(c *gin.Context) {
	rec := httptest.NewRecorder()

	bodyAsBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, err.Error(), constants.BadRequest)
		return
	}
	//bodyParams := make(map[string]interface{})
	var bodyParams models.BucketObject
	json.Unmarshal(bodyAsBytes, &bodyParams)

	// Check if the path is empty
	if bodyParams.Key == "" || bodyParams.Bucket == "" {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, constants.BadRequest, constants.BadRequest)
		return
	}

	url := "http://" + os.Getenv("RENTERD_ADDRESS") + "/api/bus/object/" + bodyParams.Key + "?bucket=" + bodyParams.Bucket
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+os.Getenv("RENTERD_KEY"))))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return
	}

	respBodyBytes, err := io.ReadAll(res.Body)
	if err != nil || string(respBodyBytes) == "object not found\n" {
		responseUtils.ErrorResponse(rec, c, http.StatusNotFound, constants.NotObjectFoundError, constants.NotObjectFoundError)
		return
	}

	respBody := make(map[string]interface{})
	json.Unmarshal(respBodyBytes, &respBody)
	defer res.Body.Close()

	fileKey := strings.Split(respBody["key"].(string), "/")
	dataOfLink := "{\"bucket\": \"" + bodyParams.Bucket + "\", \"key\": \"" + bodyParams.Key + "\", \"filename\": \"" + fileKey[len(fileKey)-1] + "\"}"
	encrypt, err := utils.GetAESEncrypted([]byte(dataOfLink))
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, constants.BadRequest, constants.BadRequest)
		return
	}

	link := "\"\\renterd\\sharefile\\" + string(base64.URLEncoding.EncodeToString([]byte(encrypt))) + "\""
	fmt.Println("Link : ", link)
	responseUtils.SuccessJsonResponse(rec, c, http.StatusOK, map[string]any{"Link": link}, constants.ShareLinkSuccessMessage)

	// Définissez les en-têtes pour indiquer que c'est un téléchargement de fichier
	/*rec.Header().Add("Content-Type", "application/json")
	//fmt.Println("Encrypt : ", encrypt)
	link := "{\"Link\": \"\\renterd\\sharefile\\" + string(base64.URLEncoding.EncodeToString([]byte(encrypt))) + "\"}"
	fmt.Println("Link : ", link)

	responseUtils.SuccessJsonResponse(rec, c, http.StatusOK, map[string]any{"Link": link}, constants.ShareLinkSuccessMessage)
	rec.Body.Write([]byte(link))
	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)*/
}

func GetShareFile(c *gin.Context) {
	rec := httptest.NewRecorder()
	key := c.Param("key")
	if key == "" {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, constants.BadRequest, constants.BadRequest)
		return
	}

	decodedKey, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, err.Error(), constants.BadRequest)
		return
	}

	decryptParams, err := utils.GetAESDecrypted(string(decodedKey))
	if err != nil {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, err.Error(), constants.BadRequest)
		return
	}

	//bodyParams := make(map[string]interface{})
	var bodyParams models.BucketObject
	json.Unmarshal(decryptParams, &bodyParams)

	// Check if the path is empty
	if bodyParams.Key == "" || bodyParams.Bucket == "" {
		responseUtils.ErrorResponse(rec, c, http.StatusBadRequest, constants.BadRequest, constants.BadRequest)
		return
	}

	url := "http://" + os.Getenv("RENTERD_ADDRESS") + "/api/worker/object/" + bodyParams.Key + "?bucket=" + bodyParams.Bucket
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+os.Getenv("RENTERD_KEY"))))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return
	}

	respBodyBytes, err := io.ReadAll(res.Body)
	if err != nil || res.StatusCode == http.StatusNotFound {
		responseUtils.ErrorResponse(rec, c, http.StatusNotFound, err.Error(), constants.NotObjectFoundError)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+bodyParams.FileName)
	c.Data(http.StatusOK, "application/octet-stream", respBodyBytes)
}
