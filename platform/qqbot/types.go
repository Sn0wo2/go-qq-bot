package types

type QQBot struct {
}

type APIAccessPayload struct {
	AppId        string `json:"appId"`
	ClientSecret string `json:"clientSecret"`
}

type HTTPCallBackACK struct {
	OP uint `json:"op"`
}

type HTTPCallBackValidation struct {
	HTTPCallBackACK
	PlainToken string `json:"plain_token"`
	Signature  string `json:"signature"`
}
