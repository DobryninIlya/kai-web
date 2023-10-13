package api_handler

import (
	"io"
	"net/http"
)

func NewEchoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.makeRegistration.NewWhoIAmHandler"
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		w.Write(body)
	}
}
