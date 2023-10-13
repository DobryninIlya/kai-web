package sqlstore

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
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
	if err != nil || len(client.Token) == 0 || client.Token != tokenStr {
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
