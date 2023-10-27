package image_handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	h "main/internal/app/handlers/web_app"
	"net/http"
	"os"
	"path/filepath"
)

func NewPostPhotoHandler(log *logrus.Logger, filePath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const funcPath = "handlers.api.getMe.NewWhoIAmHandler"
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := header.Filename
		path := filepath.Join(filePath, filename)
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("http://%s/image/%s", r.Host, filename)
		fmt.Fprint(w, url)
		h.RespondAPI(w, r, http.StatusOK, struct {
			URL string `json:"url"`
		}{
			URL: url,
		})
	}
}
