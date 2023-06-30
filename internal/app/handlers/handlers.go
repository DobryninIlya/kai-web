package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errBadID        = errors.New("missing or incorrect id")
	errUserNotFound = errors.New("user not found")
	errBadPayload   = errors.New("payload is incorrect")
	errCantCreate   = errors.New("cant create this")
)

func errorHandler(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, r, code, map[string]string{"error": err.Error()})
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		byteData, ok := data.([]byte)
		if !ok {
			json.NewEncoder(w).Encode(data)
			return
		}
		w.Write(byteData)
	}
}
