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
	"main/internal/app/authorization"
	"main/internal/app/firebase"
	"main/internal/app/formatter"
	"main/internal/app/model"
	"main/internal/app/openai"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrBadNews                = errors.New("новость не подходит для публикации")
	ErrBadPhoto               = errors.New("новость не подходит для публикации, не содержит фото")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserNotFoundInFirebase = errors.New("user not found in firebase")
	ErrBadMobileUserInfo      = errors.New("incorrect mobile user info")
	ErrAlreadyRegistered      = errors.New("user already registered")
)

type API interface {
}

// ApiRepository реализует работу API с хранилищем базы данных
type ApiRepository struct {
	store            *Store
	ConfirmationCode string
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ApiRepositoryInterface
type ApiRepositoryInterface interface {
	RegistrationToken(ctx context.Context, client *model.ApiClient, firebase firebase.FirebaseAPIInterface) (string, error)
	CheckToken(tokenStr string) (model.ApiRegistration, error, int)
	GetClient(tokenStr string) (model.ApiRegistration, error)
	CheckSecret(secret string) (bool, error, int)
	GetTokenInfo(tokenStr string) (model.ApiClient, error)
	GetNewsById(id int) (model.News, error)
	MakeNews(news model.News) (int, error)
	GetNewsPreviews(count, offset int) ([]model.News, error)
	AddAuthor(groupId int) bool
	ParseNews(update model.VKUpdate, log *logrus.Logger, openai *openai.ChatGPT) error
	SaveMobileUserInfo(user model.ApiClient) error
	SetConfirmationCode(code string)
	GetConfirmationCode() string
	RegistrationUserByPassword(ctx context.Context, client *model.ApiRegistration, firebase firebase.FirebaseAPIInterface, auth authorization.AuthorizationInterface, login, password string) (string, error)
	SaveTelegramAuth(user model.TelegramUser) error
}

const (
	tokenLength   = 32
	maxTagLength  = 30
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
func (r ApiRepository) RegistrationToken(ctx context.Context, client *model.ApiClient, firebase firebase.FirebaseAPIInterface) (string, error) {
	ctx, _ = context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	fbUser, err := firebase.GetFirebaseUser(ctx, client.UID)
	// TODO: Добавить сохранение данных пользователя в базу данных из возвращаемой выше функции
	if err != nil || len(fbUser.UID) == 0 {
		return "", ErrUserNotFoundInFirebase
	}
	if len(fbUser.UID) == 0 {
		return "", ErrUserNotFound
	}
	newToken := string(r.generateToken())
	err = r.store.db.QueryRow("INSERT INTO public.api_clients(uid, device_tag, token) VALUES ($1, $2, $3) RETURNING uid",
		client.UID,
		client.DeviceTag,
		newToken,
	).Scan(&client.UID)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "ограничение уникальности") {
			err = r.store.db.QueryRow("SELECT token FROM public.api_clients WHERE uid=$1",
				client.UID,
			).Scan(&newToken)
		} else {
			return "", err
		}
	}
	if err = r.SaveMobileUserInfo(*client); err != nil {
		return "", err
	}

	return newToken, err
}

func (r ApiRepository) SaveMobileUserInfo(user model.ApiClient) error {
	err := r.store.db.QueryRow("INSERT INTO public.mobile_users(uid, name, faculty, idcard, groupname) VALUES ($1, $2, $3, $4, $5) RETURNING uid",
		user.UID,
		user.Name,
		user.Faculty,
		user.IDCard,
		user.Groupname,
	).Scan(&user.UID)
	if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "ограничение уникальности") {
		err = r.store.db.QueryRow("UPDATE public.mobile_users SET name=$1, faculty=$2, idcard=$3, groupname=$4 WHERE uid=$5 RETURNING uid",
			user.Name,
			user.Faculty,
			user.IDCard,
			user.Groupname,
			user.UID,
		).Scan(&user.UID)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (r ApiRepository) CheckToken(tokenStr string) (model.ApiRegistration, error, int) {
	var client model.ApiRegistration
	var name, login sql.NullString
	var groupname sql.NullInt32
	err := r.store.db.QueryRow(
		"SELECT c.uid, c.device_tag, c.create_date, mu.name, mu.groupname, pw.login, pw.encrypted_password FROM public.api_clients AS c JOIN public.mobile_users AS mu ON c.uid = mu.uid JOIN public.mobile_user_password AS pw ON c.uid = pw.uid WHERE token=$1",
		tokenStr,
	).Scan(&client.UID, &client.DeviceTag, &client.CreateDate, &name, &groupname, &login, &client.EncryptedPassword)
	client.Login = login.String
	client.Name = name.String
	client.Groupname = int(groupname.Int32)
	client.Token = tokenStr
	if err != nil || len(client.UID) == 0 {
		return model.ApiRegistration{}, errors.New("bad token"), http.StatusForbidden
	}
	return client, nil, 200
}

func (r ApiRepository) GetClient(tokenStr string) (model.ApiRegistration, error) {
	var client model.ApiRegistration
	var login sql.NullString
	err := r.store.db.QueryRow(
		"SELECT login, encrypted_password FROM public.mobile_user_password WHERE uid=$1",
		tokenStr,
	).Scan(&login, &client.EncryptedPassword)
	client.Login = login.String
	if err != nil || len(client.Login) == 0 {
		return model.ApiRegistration{}, errors.New("bad token")
	}
	return client, nil
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
	err := r.store.db.QueryRow("SELECT n.header, n.description, n.body, n.date, n.preview_url, a.name, n.ai_correct FROM public.news AS n LEFT JOIN public.news_authors AS a ON n.author = a.id WHERE n.id=$1",
		id,
	).Scan(
		&news.Header,
		&news.Description,
		&news.Body,
		&news.Date,
		&news.PreviewURL,
		&news.AuthorName,
		&news.AICorrect,
	)
	return news, err
}

func (r ApiRepository) MakeNews(news model.News) (int, error) {
	var id int
	err := r.store.db.QueryRow("INSERT INTO public.news (header, description, body, preview_url, tag, author, ai_correct) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		news.Header,
		news.Description,
		news.Body,
		news.PreviewURL,
		news.Tag,
		news.Author,
		news.AICorrect,
	).Scan(
		&id,
	)
	return id, err
}

func (r ApiRepository) GetNewsPreviews(count, offset int) ([]model.News, error) {
	rows, err := r.store.db.Query("SELECT n.id, n.header, n.description, n.date, n.preview_url, n.tag, a.name FROM public.news AS n JOIN public.news_authors AS a ON n.author = a.id ORDER BY id DESC LIMIT $1 OFFSET $2",
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
		log.Printf("Не удалось создать автора: %v", err)
		return false
	}
	return true
}

func (r ApiRepository) ParseNews(update model.VKUpdate, log *logrus.Logger, openai *openai.ChatGPT) error {
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
	newsParams, err := openai.GenerateAnswer(body, 5)

	if err != nil {
		log.Logf(
			logrus.WarnLevel,
			"%v : Ошибка генерации новости: %v",
			path,
			err,
		)
	}
	if err == nil && newsParams.Header != "" && newsParams.Description != "" {
		news.AICorrect = true
	}
	var header string
	if newsParams.Header == "" {
		if header, err = formatter.GetHeader(body); err != nil {
			return err
		} else {
			news.Header = header
		}
	} else {
		news.Header = newsParams.Header
	}
	images := formatter.GetImages(update.Object.Attachments)
	news.Body = strings.ReplaceAll(body, "\n", "<br>") + images
	news.Tag = formatter.GetTagsInText(body)
	if len(news.Tag) > maxTagLength {
		news.Tag = ""
	}

	if newsParams.Description == "" {
		if description, err := formatter.GetDescription(body, header); err != nil {
			return err
		} else {
			news.Description = description
		}
	} else {
		news.Description = newsParams.Description
	}

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

func (r ApiRepository) SetConfirmationCode(code string) {
	r.ConfirmationCode = code
}

func (r ApiRepository) GetConfirmationCode() string {
	return r.ConfirmationCode
}

func (r ApiRepository) RegistrationUserByPassword(ctx context.Context, client *model.ApiRegistration, firebase firebase.FirebaseAPIInterface,
	auth authorization.AuthorizationInterface, login, password string) (string, error) {

	ctxFirebase, _ := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	fbUser, err := firebase.GetFirebaseUser(ctxFirebase, client.UID)
	// TODO: Добавить сохранение данных пользователя в базу данных из возвращаемой выше функции
	if err != nil || len(fbUser.UID) == 0 {
		return "", ErrUserNotFoundInFirebase
	}
	if len(fbUser.UID) == 0 {
		return "", ErrUserNotFound
	}
	cookies, err := auth.GetCookiesByPassword(login, password)
	if err != nil {
		return "", err
	}
	auth.SetCookies(fbUser.UID, cookies)
	aboutInfo, err := auth.GetAboutInfo(fbUser.UID, *client)
	if err != nil {
		return "", err
	}
	client.Name = aboutInfo.FirstName + " " + aboutInfo.LastName + " " + aboutInfo.MiddleName
	authorization.Encrypt(&client.EncryptedPassword, password)
	group, err := auth.GetGroupNum(fbUser.UID, *client)
	if err != nil {
		return "", err
	}
	client.Groupname = group
	tokenStr, err := r.saveMobileUser(ctx, client)
	if err != nil {
		return "", err
	}
	auth.SetCookies(fbUser.UID, cookies)
	return tokenStr, nil
}

func (r ApiRepository) saveMobileUser(ctx context.Context, client *model.ApiRegistration) (string, error) {
	tx, err := r.store.db.Begin()
	if err != nil {
		panic(err)
	}

	// Запросы на вставку данных в таблицы
	insertMobileUserPassword := `INSERT INTO public.mobile_user_password (uid, login, encrypted_password) VALUES ($1, $2, $3)`
	insertMobileUser := `INSERT INTO public.mobile_users (uid, name, groupname) VALUES ($1, $2, $3)`
	insertAPIClient := `INSERT INTO public.api_clients (uid, device_tag, token) VALUES ($1, $2, $3)`

	// Данные для вставки
	uid := client.UID
	login := client.Login
	encryptedPassword := client.EncryptedPassword
	name := client.Name
	groupname := client.Groupname
	deviceTag := client.DeviceTag
	newToken := r.generateToken()

	// Вставка данных в таблицу mobile_user_password
	_, err = tx.ExecContext(ctx, insertMobileUserPassword, uid, login, encryptedPassword)
	if err != nil {
		// Ошибка при вставке данных, отмена транзакции
		tx.Rollback()
		return "", err
	}

	// Вставка данных в таблицу mobile_users
	_, err = tx.ExecContext(ctx, insertMobileUser, uid, name, groupname)
	if err != nil {
		// Ошибка при вставке данных, отмена транзакции
		tx.Rollback()
		return "", err
	}

	// Вставка данных в таблицу api_clients
	_, err = tx.ExecContext(ctx, insertAPIClient, uid, deviceTag, newToken)
	if err != nil {
		// Ошибка при вставке данных, отмена транзакции
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return string(newToken), nil
}

func (r ApiRepository) SaveTelegramAuth(user model.TelegramUser) error {
	_, err := r.store.db.Query("INSERT INTO public.mobile_user_password (uid, login, encrypted_password) VALUES ($1, $2, $3)",
		user.UID,
		user.Login,
		user.EncryptedPassword,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "ограничение уникальности") {
			return ErrAlreadyRegistered
		}
		return err
	}
	return nil
}
