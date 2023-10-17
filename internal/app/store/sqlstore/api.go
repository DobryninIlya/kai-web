package sqlstore

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"main/internal/app/firebase"
	"main/internal/app/formatter"
	"main/internal/app/model"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrBadNews  = errors.New("новость не подходит для публикации")
	ErrBadPhoto = errors.New("новость не подходит для публикации, не содержит фото")
)

// ApiRepository реализует работу API с хранилищем базы данных
type ApiRepository struct {
	store            *Store
	ConfirmationCode string
}

const (
	tokenLength   = 32
	maxTagLength  = 25
	minTextLength = 100
)

type token string

// generateToken генерирует уникальный криптоустойчивый токен
func (r ApiRepository) generateToken() token {
	// генерируем случайную строку заданной длины
	tokenBytes := make([]byte, tokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		log.Printf("Ошибка генерации токена: %s", err)
		return ""
	}
	tokenValue := base64.URLEncoding.EncodeToString(tokenBytes)

	// хэшируем строку токена
	hashBytes := sha256.Sum256([]byte(tokenValue))
	hashValue := fmt.Sprintf("%x", hashBytes)

	resultToken := token(hashValue)

	return resultToken

}

// RegistrationToken регистрирует нового апи-клиента и возвращает уникальный токен
func (r ApiRepository) RegistrationToken(ctx context.Context, client *model.ApiClient, firebase *firebase.FirebaseAPI) (string, error) {
	ctx, _ = context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	fbUser, err := firebase.GetFirebaseUser(ctx, client.UID)
	// TODO: Добавить сохранение данных пользователя в базу данных из возвращаемой выше функции
	if err != nil || len(fbUser.UID) == 0 {
		return "", err
	}
	if len(fbUser.UID) == 0 {
		return "", errors.New("bad uid")
	}
	newToken := string(r.generateToken())
	err = r.store.db.QueryRow("INSERT INTO public.api_clients(uid, device_tag, token) VALUES ($1, $2, $3) RETURNING uid",
		client.UID,
		client.DeviceTag,
		newToken,
	).Scan(&client.UID)
	return newToken, err
}

func (r ApiRepository) CheckToken(tokenStr string) (model.ApiClient, error, int) {
	var client model.ApiClient
	err := r.store.db.QueryRow("SELECT uid, device_tag, create_date FROM public.api_clients WHERE token=$1",
		tokenStr,
	).Scan(&client.UID, &client.DeviceTag, &client.CreateDate)
	if err != nil || len(client.UID) == 0 {
		return model.ApiClient{}, errors.New("bad token"), http.StatusForbidden
	}
	return client, nil, 200
}

func (r ApiRepository) CheckSecret(secret string) (bool, error, int) {
	osSecret := os.Getenv("WEB_KAI_SECRET")
	if len(osSecret) == 0 || secret != osSecret {
		return false, errors.New("bad secret"), http.StatusForbidden
	}
	return true, nil, 200
}

// GetTokenInfo получает информацию о владельце токена
func (r ApiRepository) GetTokenInfo(tokenStr string) (model.ApiClient, error) {
	var res model.ApiClient
	err := r.store.db.QueryRow("SELECT uid, device_tag, create_date FROM public.api_clients WHERE token=$1",
		tokenStr,
	).Scan(
		&res.UID,
		&res.DeviceTag,
		&res.CreateDate,
	)
	if err != nil || len(res.UID) == 0 {
		return model.ApiClient{}, errors.New("bad token")
	}
	return res, nil
}

func (r ApiRepository) GetNewsById(id int) (model.News, error) {
	var news model.News
	news.Id = id
	err := r.store.db.QueryRow("SELECT n.header, n.description, n.body, n.date, n.preview_url, a.name FROM public.news AS n LEFT JOIN public.news_authors AS a ON n.author = a.id WHERE n.id=$1",
		id,
	).Scan(
		&news.Header,
		&news.Description,
		&news.Body,
		&news.Date,
		&news.PreviewURL,
		&news.AuthorName,
	)
	return news, err
}

func (r ApiRepository) MakeNews(news model.News) (int, error) {
	var id int
	err := r.store.db.QueryRow("INSERT INTO public.news (header, description, body, preview_url, tag, author) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		news.Header,
		news.Description,
		news.Body,
		news.PreviewURL,
		news.Tag,
		news.Author,
	).Scan(
		&id,
	)
	return id, err
}

func (r ApiRepository) GetNewsPreviews(count, offset int) ([]model.News, error) {
	rows, err := r.store.db.Query("SELECT n.id, n.header, n.description, n.date, n.preview_url, n.tag, a.name FROM public.news AS n LEFT JOIN public.news_authors AS a ON n.author = a.id ORDER BY id DESC LIMIT $1 OFFSET $2",
		count,
		offset,
	)
	if err != nil {
		return nil, err
	}
	result := make([]model.News, 0)
	for rows.Next() {
		var news model.News
		var previewUrl, tag sql.NullString
		err := rows.Scan(&news.Id, &news.Header, &news.Description, &news.Date, &previewUrl, &tag, &news.AuthorName)
		if previewUrl.Valid {
			news.PreviewURL = previewUrl.String
		} else {
			news.PreviewURL = ""
		}
		if tag.Valid {
			news.Tag = tag.String
		} else {
			news.Tag = ""
		}
		if err != nil {
			return nil, err
		}
		result = append(result, news)
	}
	return result, err
}

func (r ApiRepository) AddAuthor(groupId int) bool {
	err := r.store.db.QueryRow("INSERT INTO public.news_authors (id, name) VALUES ($1, 'Неизвестный')",
		groupId,
	)
	if err != nil {
		return false
	}
	return true
}

func (r ApiRepository) ParseNews(update model.VKUpdate, log *logrus.Logger) error {
	const path = "internal.app.store.sqlstore.api.ParseNews"
	var news model.News
	// Не создаем новость, если некорректный айди, текст меньше 100 символов, ни одной картинки, тип публикации не "post"
	// пост помечен как рекламный
	if update.Object.FromId == 0 || len(update.Object.Text) < minTextLength || len(update.Object.Attachments) == 0 ||
		update.Object.PostType != "post" || update.Object.MarkedAsAds != 0 {
		return ErrBadNews
	}
	body := update.Object.Text
	err := formatter.ValidateText(body)
	if err != nil {
		log.Logf(
			logrus.WarnLevel,
			"%v : Ошибка валидации новости: %v",
			path,
			err,
		)
		return ErrBadNews
	}
	header, err := formatter.GetHeader(body)
	if err != nil {
		return err
	}

	news.Header = header
	images := formatter.GetImages(update.Object.Attachments)
	news.Body = strings.ReplaceAll(body, "\n", "<br><br>") + images
	news.Tag = formatter.GetTagsInText(body)
	if len(news.Tag) > maxTagLength {
		news.Tag = ""
	}
	description, err := formatter.GetDescription(body, header)
	news.Description = description
	if err != nil {
		return err
	}
	if len(update.Object.Attachments) == 0 { // Если нет вложений
		return ErrBadPhoto
	}
	if update.Object.Attachments[0].Type != "photo" { // Если первое вложение не фото
		return ErrBadPhoto
	}
	// Получаем ссылку на картинку, первую в подборке с самым большим (последним) размером
	news.PreviewURL = update.Object.Attachments[0].Photo.Sizes[len(update.Object.Attachments[0].Photo.Sizes)-1].Url
	news.Author = update.Object.FromId
	id, err := r.MakeNews(news)
	if err != nil {
		return err
	}
	log.Logf(
		logrus.WarnLevel,
		"%v : Новость успешно сохранена в базе. ID : %v",
		path,
		id,
	)
	return nil
}
