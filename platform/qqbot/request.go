package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

func Request(w http.ResponseWriter, r *http.Request) {
	if err := Signature(w, r); err != nil {
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
		Validation(w, r, payload)
	default:
		Dispatch(w, r, payload)
	}
}
