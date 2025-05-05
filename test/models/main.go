package test

type SimpleResponse struct {
	Message string `json:"message"`
}

type EncryptResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type Crypto struct {
	Key         string `json:"Key"`
	TextDecrypt string `json:"TextDecrypt"`
	TextEncrypt string `json:"TextEncrypt"`
}

type CryptoTest struct {
	Key         string `json:"key"`
	Iv          string `json:"iv"`
	TextDecrypt string `json:"textDecrypt"`
	TextEncrypt string `json:"textEncrypt"`
}
