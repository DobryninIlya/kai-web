package image_handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func NewPostTaskPhotoHandler(log *logrus.Logger, filePath string, store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const funcPath = "handlers.api.getMe.NewWhoIAmHandler"

		url := r.URL.Query()
		token := url.Get("token")
		client, err, _ := store.API().CheckToken(token)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}
		fileExt := filepath.Ext(header.Filename)
		fileName := uuid.New().String() + fileExt
		path := filepath.Join(filePath, "groups", "tasks", strconv.Itoa(client.Groupname))
		os.Mkdir(path, 0755)
		path = filepath.Join(path, fileName)
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			log.Logf(
				logrus.ErrorLevel,
				"%s : Ошибка записи файла: %v",
				funcPath,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusInternalServerError, err)
			return
		}

		urlPath := fmt.Sprintf("http://%s/image/groups/tasks/%s/%s", client.Groupname, fileName)
		h.RespondAPI(w, r, http.StatusOK, struct {
			URL string `json:"url"`
		}{
			URL: urlPath,
		})
	}
}
