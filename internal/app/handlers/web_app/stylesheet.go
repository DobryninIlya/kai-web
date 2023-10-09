package web_app

import (
	"main/internal/app/tools"
	"net/http"
)

func NewStyleSheetHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := tools.GetMainStylesheet()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
