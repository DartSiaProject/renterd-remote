package models

type httpRequestHeader struct {
	Content_Type              []string `json:"content-type" `
	Authorization             []string `json:"authorization" `
	Accept                    []string `json:"accept" `
	Content_Length            []string `json:"Ccontent-length" `
	User_Agent                []string `json:"user-agent" `
	Date                      []string `json:"date" `
	Referrer_Policy           []string `json:"referrer_policy" `
	Strict_Transport_Security []string `json:"strict-transport-security" `
	X_Content_Type_Options    []string `json:"x_content_type_Options" `
	X_Forwarded_Proto         []string `json:"x-forwarded-proto" `
	X_Xss_Protection          []string `json:"x-xss-protection" `
	Transfer_Encoding         []string `json:"transfer-encoding" `
	Last_Modified             []string `json:"last-modified" `
	ETag                      []string `json:"eTag" `
}
