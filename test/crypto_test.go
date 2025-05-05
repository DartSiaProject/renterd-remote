package test

import (
	"renterd-remote/utils"
	"testing"

	"github.com/stretchr/testify/assert"

	//"fmt"
	"os"
	model "renterd-remote/test/models"
)

func mock() (model.CryptoTest, error) {
	testValues := model.CryptoTest{
		Key:         "418D125DC1475C7817f31185981e4ebe",
		Iv:          "300372804b842889",
		TextDecrypt: "Test avec plusieurs caracteres spéciaux !@#$%^&*()_+",
		/*TextDecrypt:`{
			"Email":    "falseEmail@false.cm",
			"Password": "FalsePassword"
		}`*/
		TextEncrypt: "G1BB9AY/jE1lf1B3DBbBE0FBv5KvPe4+5dJaYoYXfD6igWYzCdEcLKd68+Wr9ZH4U9BvoWZXDMCdJYnaQX53Mw==",
	}

	os.Setenv("USER_KEY", testValues.Key)
	os.Setenv("USER_IV", testValues.Iv)
	return testValues, nil
}

func TestCryptoAlgo(t *testing.T) {
	testValues, err := mock()
	assert.NoError(t, err, "Error setting up test values")

	encrypted, err := utils.GetAESEncrypted([]byte(testValues.TextDecrypt))
	assert.NoError(t, err, "Error encrypting text")
	assert.Equal(t, testValues.TextEncrypt, encrypted, "Encrypted text does not match expected")

	decrypted, err := utils.GetAESDecrypted(encrypted)
	assert.NoError(t, err, "Error decrypting text")
	assert.Equal(t, testValues.TextDecrypt, string(decrypted), "Decrypted text does not match original")

	//fmt.Println("This is an original:", testValues.textDecrypt)
	//fmt.Println("This is an encrypted:", encrypted)
}
