package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Sn0wo2/go-qq-bot/driver"
	"github.com/Sn0wo2/go-qq-bot/handler"

	"github.com/tidwall/gjson"
)

func main() {
	driver.RunWebHook(handleRequest)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if err := handler.Signature(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	fmt.Println(string(body))

	payload := gjson.ParseBytes(body)

	switch payload.Get("op").Int() {
	case 13:
		handler.Validation(w, r, payload)
	default:
		handler.Dispatch(w, r, payload)
	}
}
