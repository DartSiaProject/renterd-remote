package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// GetAESDecrypted decrypts given text in AES 256 CBC
func GetAESDecrypted(encrypted string) ([]byte, error) {
	key := os.Getenv("USER_KEY")
	iv := os.Getenv("USER_IV")

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext []byte) (string, error) {
	key := os.Getenv("USER_KEY")
	iv := os.Getenv("USER_IV")

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

// Hash data using SHA-256
func HashData(plaintext string) (string, error) {

	hash := sha256.New()
	hash.Write([]byte(plaintext))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Create key for AES256 algorithm
func CreateSecretKey(email, password string) string {
	data, _ := HashData(email + password)

	return strings.ToUpper(data[0:16]) + data[16:32]
}

// Create IV for AES256 algorithm
func CreateIV(email, password string) string {
	data, _ := HashData((email + password))

	return strings.ToUpper(data[32:40]) + data[40:48]
}

func Test() {
	plainText := "{Accept:[*/*], Accept-Encoding:[gzip, deflate, br], Authorization:[Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjAzOTU2MDksInVzZXJuYW1lIjoidGVzdEB0ZXN0LmNvbSJ9.npJHWNHxEYNa3XV12a4TweFlQHlp1mcUO-_DmUz4SPg], Cache-Control:[no-cache], Connection:[keep-alive], Content-Length:[46], Content-Type:[application/json], Postman-Token:[5b2982a6-0fb8-4fb3-893e-3800fe1ac662], User-Agent:[PostmanRuntime/7.39.0]}"
	fmt.Println("This is an original:", plainText)

	encrypted, err := GetAESEncrypted([]byte(plainText))

	if err != nil {
		fmt.Println("Error during encryption", err)
	}

	fmt.Println("This is an encrypted:", encrypted)

	decrypted, err := GetAESDecrypted(encrypted)

	if err != nil {
		fmt.Println("Error during decryption", err)
	}
	fmt.Println("This is a decrypted:", string(decrypted))

	hash, _ := HashData("test@test.com@Bcabc12")
	fmt.Println("This is a Hash: ", hash)
}
