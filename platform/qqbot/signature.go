package handler

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Sn0wo2/go-qq-bot/config"
)

func Signature(w http.ResponseWriter, r *http.Request) error {
	sigHex := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")

	if sigHex == "" || timestamp == "" {
		log.Println("missing signature headers")
		return fmt.Errorf("missing signature headers")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body for verify failed: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	var msg bytes.Buffer
	msg.WriteString(timestamp)
	msg.Write(body)

	seed := config.Instance.BotSecret
	for len(seed) < ed25519.SeedSize {
		seed += seed
	}
	seed = seed[:ed25519.SeedSize]
	reader := strings.NewReader(seed)

	publicKey, _, err := ed25519.GenerateKey(reader)
	if err != nil {
		return fmt.Errorf("generate public key failed: %w", err)
	}

	sig, err := hex.DecodeString(sigHex)
	if err != nil || len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		return fmt.Errorf("invalid signature")
	}

	if !ed25519.Verify(publicKey, msg.Bytes(), sig) {
		return fmt.Errorf("signature verify failed")
	}

	return nil
}
