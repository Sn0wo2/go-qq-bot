package handler

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Sn0wo2/go-qq-bot/config"
	"github.com/Sn0wo2/go-qq-bot/platform/qqbot/types"
	"github.com/tidwall/gjson"
)

func Validation(w http.ResponseWriter, r *http.Request, payload gjson.Result) {
	d := payload.Get("d")
	plainToken := d.Get("plain_token").Str
	eventTs := d.Get("event_ts").Str

	if plainToken == "" || eventTs == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	seed := config.Instance.BotSecret
	for len(seed) < ed25519.SeedSize {
		seed += seed
	}
	seed = seed[:ed25519.SeedSize]
	reader := strings.NewReader(seed)

	_, privateKey, err := ed25519.GenerateKey(reader)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var msg bytes.Buffer
	msg.WriteString(eventTs)
	msg.WriteString(plainToken)

	validation, err := json.Marshal(types.HTTPCallBackValidation{HTTPCallBackACK: types.HTTPCallBackACK{OP: 12}, PlainToken: plainToken, Signature: hex.EncodeToString(ed25519.Sign(privateKey, msg.Bytes()))})
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(validation)
}
