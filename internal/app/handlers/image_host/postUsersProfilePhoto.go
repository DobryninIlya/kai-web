package image_handler

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const imgSize = 48

func NewPostUserProfilePhotoHandler(log *logrus.Logger, filePath string, store sqlstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const funcPath = "handlers.api.getMe.NewPostUserProfilePhotoHandler"

		url := r.URL.Query()
		token := url.Get("token")
		client, err, _ := store.API().CheckToken(token)
		client.UID = strings.TrimSpace(client.UID)
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
		//fileName := uuid.New().String() + fileExt
		fileName := client.UID + fileExt
		path := filepath.Join(filePath, "users")
		os.Mkdir(path, 0755)
		err = makeThumbnail(filepath.Join(path, "thumb_"+fileName), file, header)
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
		path = filepath.Join(path, "user_"+fileName)
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

		urlPath := fmt.Sprintf("https://%s/image/users/user_%s", r.Host, fileName)
		urlPathThumb := fmt.Sprintf("https://%s/image/users/thumb_%s", r.Host, fileName)
		h.RespondAPI(w, r, http.StatusOK, struct {
			URL   string `json:"url"`
			Thumb string `json:"thumbnail"`
		}{
			URL:   urlPath,
			Thumb: urlPathThumb,
		})
	}
}

func convertMultipartFileToFile(file multipart.File, fileHeader *multipart.FileHeader) (*os.File, error) {
	defer file.Close()

	header := fileHeader.Header
	fileName := header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(fileName)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(params["filename"])

	tempFile, err := os.CreateTemp("", "tempfile")
	if err != nil {
		return nil, err
	}
	// Копирование содержимого multipart.File во временный файл
	file.Seek(0, io.SeekStart)
	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	// Закрытие временного файла перед переименованием
	tempFile.Close()

	// Переименование временного файла с правильным расширением
	newTempFile := tempFile.Name() + ext
	err = os.Rename(tempFile.Name(), newTempFile)
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, err
	}

	// Открытие переименованного файла
	finalFile, err := os.Open(newTempFile)
	if err != nil {
		os.Remove(newTempFile)
		return nil, err
	}

	return finalFile, nil
}

func makeThumbnail(path string, multipartFile multipart.File, header *multipart.FileHeader) error {
	file, err := convertMultipartFileToFile(multipartFile, header)
	if err != nil {
		log.Println(err)
		return err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		return err
	}
	if err != nil {
		return err
	}
	thumb := imaging.Thumbnail(img, imgSize, imgSize, imaging.Lanczos)
	err = imaging.Save(thumb, path)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
