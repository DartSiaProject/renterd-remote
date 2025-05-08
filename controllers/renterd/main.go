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
	utils "renterd-remote/utils"

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
		//println("Error in reverse proxy:", err.Error())

		//log.Printf("[ERROR] Proxy connection failed: %s", err.Error())

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": constants.ReverseProxyError})
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
		rec.Header().Set("Content-Type", "application/json")
		rec.WriteHeader(http.StatusInternalServerError)
		rec.Body.Write([]byte(`{"message":` + constants.InternalServerError + `, "error": "` + err.Error() + `"}`))

		//Transfert response to encrypt middelware
		middlewares.EncryptResponse(rec, c)
		return
	}

	// Définissez les en-têtes pour indiquer que c'est un téléchargement de fichier
	rec.Header().Add("Content-Disposition", "attachment; filename=\"renterd.sqlite3.bak\"")
	rec.Header().Add("Content-Type", "application/octet-stream")
	rec.Header().Add("Cache-Control", "no-cache")
	rec.Header().Add("Pragma", "no-cache")
	rec.Header().Add("Expires", "0")

	// Utilisez c.FileAttachment pour envoyer le fichier
	// Récupérer le fichier depuis la requête
	_, err = os.Stat(path)
	if err != nil {
		rec.Header().Set("Content-Type", "application/json")
		rec.WriteHeader(http.StatusInternalServerError)
		rec.Body.Write([]byte(`{"message":` + constants.InternalServerError + `, "error": "` + err.Error() + `"}`))
		return
	}

	src, err := os.Open(path)
	if err != nil {
		rec.Header().Set("Content-Type", "application/json")
		rec.WriteHeader(http.StatusInternalServerError)
		rec.Body.Write([]byte(`{"message":` + constants.InternalServerError + `, "error": "` + err.Error() + `"}`))
		return
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		rec.Header().Set("Content-Type", "application/json")
		rec.WriteHeader(http.StatusInternalServerError)
		rec.Body.Write([]byte(`{"message":` + constants.InternalServerError + `, "error": "` + err.Error() + `"}`))
		return
	}

	rec.Body.Write(content)
	//Transfert response to encrypt middelware
	middlewares.EncryptResponse(rec, c)
}

// Restore SQlLite database from mobile app
func RestoreSqliteDb(c *gin.Context) {
	// Get the SQLite database file path from the request body
	path := utils.GetDefaultSqliteBackupPath()

	// Récupérer le fichier depuis la requête
	buffer := make([]byte, c.Request.ContentLength)
	_, err := c.Request.Body.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"message": constants.BadRequest, "error": err.Error()})
		return
	}
	file := bytes.NewReader(buffer)

	// Create the destination file
	dst, err := os.Create(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constants.InternalServerError, "error": err.Error()})
		return
	}
	defer dst.Close()

	// Copy the contents of the downloaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constants.InternalServerError, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.SqlliteRestoreSuccessMessage})
}
