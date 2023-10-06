package sqlstore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"main/internal/app/model"
	_ "main/internal/app/store"
	"net/http"
)

// ApiRepository реализует работу API с хранилищем базы данных
type ApiRepository struct {
	store *Store
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
func (r ApiRepository) RegistrationToken(client *model.ApiClient) (string, error) {
	newToken := string(r.generateToken())
	err := r.store.db.QueryRow("INSERT INTO public.api_tokens(device_id, device_tag, token) VALUES ($1, $2, $3) RETURNING id",
		client.DeviceId,
		client.DeviceTag,
		newToken,
	).Scan(&client.Id)
	return newToken, err
}

func (r ApiRepository) CheckToken(tokenStr string) (int, error, int) {
	var id int
	err := r.store.db.QueryRow("SELECT id FROM public.api_tokens WHERE token=$1",
		tokenStr,
	).Scan(&id)
	if err != nil || id == 0 {
		return 0, errors.New("bad token"), http.StatusForbidden
	}
	return id, nil, 200
}

// GetTokenInfo получает информацию о владельце токена
func (r ApiRepository) GetTokenInfo(tokenStr string) (model.ApiClient, error) {
	var res model.ApiClient
	err := r.store.db.QueryRow("SELECT id, device_id, device_tag, create_date FROM public.api_tokens WHERE token=$1",
		tokenStr,
	).Scan(
		&res.Id,
		&res.DeviceId,
		&res.DeviceTag,
		&res.CreateDate,
	)
	if err != nil || res.Id == 0 {
		return model.ApiClient{}, errors.New("bad token")
	}
	return res, nil
}
