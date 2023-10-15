package sqlstore

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"main/internal/app/firebase"
	"main/internal/app/model"
	"net/http"
	"os"
	"time"
)

// ApiRepository реализует работу API с хранилищем базы данных
type ApiRepository struct {
	store            *Store
	ConfirmationCode string
}

const tokenLength = 32

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
	err := r.store.db.QueryRow("SELECT header, description, body, date, preview_url FROM public.news WHERE id=$1",
		id,
	).Scan(
		&news.Header,
		&news.Description,
		&news.Body,
		&news.Date,
		&news.PreviewURL,
	)
	return news, err
}

func (r ApiRepository) MakeNews(news model.News) (int, error) {
	var id int
	err := r.store.db.QueryRow("INSERT INTO public.news (header, description, body, preview_url, tag) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		news.Header,
		news.Description,
		news.Body,
		news.PreviewURL,
		news.Tag,
	).Scan(
		&id,
	)
	return id, err
}

func (r ApiRepository) GetNewsPreviews(count, offset int) ([]model.News, error) {
	rows, err := r.store.db.Query("SELECT header, description, date, preview_url, tag FROM public.news ORDER BY id DESC LIMIT $1 OFFSET $2",
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
		err := rows.Scan(&news.Header, &news.Description, &news.Date, &previewUrl, &tag)
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
