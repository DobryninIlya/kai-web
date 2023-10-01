package sqlstore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"main/internal/app/model"
	_ "main/internal/app/store"
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
