package driver

import (
	"log"
	"net/http"

	"github.com/Sn0wo2/go-qq-bot/platform/qqbot/handler"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Request)
	log.Println("Listening on :9000")
	return http.ListenAndServe(":9000", mux)
}
