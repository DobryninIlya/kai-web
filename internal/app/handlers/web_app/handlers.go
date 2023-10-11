package web_app

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrBadID            = errors.New("missing or incorrect id")
	ErrUserNotFound     = errors.New("user not found")
	ErrRecordNotFound   = errors.New("record not found")
	ErrBadPayload       = errors.New("payload is incorrect")
	ErrCantCreated      = errors.New("cant create this")
	ErrInternal         = errors.New("internal server error")
	ErrUniqueConstraint = errors.New("unique constraint failed for one of the field")
	ErrLongData         = errors.New("the length of one of parameters is too much")
)

func RespondAPI(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	Respond(w, r, code, map[string]interface{}{"result": data})
}

func ErrorHandlerAPI(w http.ResponseWriter, r *http.Request, code int, err error) {
	Respond(w, r, code, map[string]string{"error": err.Error()})
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, err error) {
	Respond(w, r, code, map[string]string{"error": err.Error()})
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		byteData, ok := data.([]byte)
		if !ok {
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			enc.Encode(data)
			return
		}
		w.Write(byteData)
	}
}