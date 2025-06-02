package driver

import (
	"log"
	"net/http"
)

func RunWebHook(handler func(w http.ResponseWriter, r *http.Request)) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	log.Println("Listening on :9000")
	if err := http.ListenAndServe(":9000", mux); err != nil {
		panic(err)
	}
}
