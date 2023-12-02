package authorization

import (
	"errors"
	"fmt"
	"log"
	"main/internal/app/model"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	SessionLiveTime = time.Minute * 90
	AuthEndpoint    = "https://kai.ru/main?p_p_id=58&p_p_lifecycle=1&p_p_state=normal&p_p_mode=view&_58_struts_action=%2Flogin%2Flogin"
	AboutPage       = "https://kai.ru/group/guest/common/about-me"
	AttestationPage = "https://kai.ru/group/guest/student/attestacia"
	MyGroupPage     = "https://kai.ru/group/guest/student/moa-gruppa"
)

var (
	ErrWrongPassword       = errors.New("wrong password")
	ErrAttestationNotFound = errors.New("attestation table not found")
)

type AuthorizationInterface interface {
	GetCookiesByPassword(login, password string) ([]*http.Cookie, error)
	GetAboutInfo(uid string, client model.ApiRegistration) (model.SiteUserInfo, error)
	GetGroupNum(uid string, apiClient model.ApiRegistration) (int, error)
	GetAttestationList(uid string, client model.ApiRegistration) ([]model.Discipline, error)
	SetCookies(key string, cookies []*http.Cookie)
	GetCookies(key string) (cookie, bool)
	GetAuthorizedCookies(uid string, client model.ApiRegistration) (cookie, error)
	SetAttestations(key string, list []model.Discipline)
	GetAttestations(key string) ([]model.Discipline, bool)
	GetProfilePhotoURL(uid string, client model.ApiRegistration) (string, error)
}

type Authorization struct {
	Cookies      SafeMap
	Attestations AttestationCache
}

type SafeMap struct {
	mu sync.Mutex
	m  map[string]cookie
}

type cookie struct {
	Cookie     []*http.Cookie
	LastUpdate time.Time
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]cookie),
	}
}

func (r *Authorization) SetCookies(key string, cookies []*http.Cookie) {
	r.Cookies.mu.Lock()
	defer r.Cookies.mu.Unlock()
	cookies = append(cookies, &http.Cookie{
		Name:  "COOKIE_SUPPORT",
		Value: "true",
		Raw:   "COOKIE_SUPPORT=true;",
	})
	r.Cookies.m[key] = cookie{Cookie: cookies, LastUpdate: time.Now()}
}

func (r *Authorization) GetCookies(key string) (cookie, bool) {
	r.Cookies.mu.Lock()
	defer r.Cookies.mu.Unlock()
	value, ok := r.Cookies.m[key]
	if !ok || value.LastUpdate.Add(SessionLiveTime).Before(time.Now()) {
		return cookie{}, false
	}
	return value, ok
}

func (r *Authorization) GetAuthorizedCookies(uid string, client model.ApiRegistration) (cookie, error) {
	cookies, ok := r.GetCookies(uid)
	if !ok {
		passwordDecrypted, err := Decrypt(string(client.EncryptedPassword))
		cookiesList, err := r.GetCookiesByPassword(client.Login, passwordDecrypted)
		if err != nil {
			return cookie{}, err
		}
		r.SetCookies(uid, cookiesList)

		return r.GetAuthorizedCookies(uid, client)
	}
	return cookie{
		Cookie: cookies.Cookie,
	}, nil
}

func getInitialCookies() (string, error) {
	// Создание GET-запроса
	req, err := http.NewRequest("GET", "https://kai.ru/c", nil)
	if err != nil {
		return "", err
	}

	// Отправка GET-запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Получение Cookie из ответа
	cookies := resp.Cookies()
	cookieHeader := ""
	for _, cookiePart := range cookies {
		cookieHeader += cookiePart.Name + "=" + cookiePart.Value + "; "
	}

	return cookieHeader, nil
}

func getAuthRequest(login, password, cookies string) (*http.Request, error) {
	if cookies == "" {
		cookies, _ = getInitialCookies()
	}
	method := "POST"
	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)
	payload := strings.NewReader(fmt.Sprintf("_58_formDate=1699709626094&_58_saveLastPath=false&_58_redirect=&_58_doActionAfterLogin=false&_58_login=%v&_58_password=%v", login, password))
	req, err := http.NewRequest(method, AuthEndpoint, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	return req, nil
}

// GetCookiesByPassword возвращает куки, полученные после авторизации по логину и паролю (либо ошибку)
func (r *Authorization) GetCookiesByPassword(login, password string) ([]*http.Cookie, error) {
	cookies, err := getInitialCookies()
	if err != nil {
		return nil, err
	}
	req, err := getAuthRequest(login, password, cookies)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	for resp.Request.Response != nil {
		resp = resp.Request.Response
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cookiesList := resp.Cookies()
	//cookieHeader := GetCookiesHeader(cookiesList)
	var authCookie bool
	if len(cookiesList) >= 7 {
		authCookie = true
	} else {
		log.Println("Куки файлы авторизации не найдены, количество: ", len(cookiesList))
		log.Println(login, password)
		log.Println(cookies)
	}
	if authCookie {
		//user, err := r.GetAboutInfo(cookieHeader)
		//r.GetGroupNum(cookieHeader)
		//if err != nil || user.FirstName == "" {
		//	return nil, err
		//}

		return cookiesList, nil
	}
	return nil, ErrWrongPassword
}

func GetCookiesHeader(cookies []*http.Cookie) string {
	cookieHeader := ""
	for _, cookie := range cookies {
		cookieHeader += cookie.Name + "=" + cookie.Value + "; "
	}
	return cookieHeader
}

func (r *Authorization) SetAttestations(key string, list []model.Discipline) {
	r.Attestations.mu.Lock()
	defer r.Attestations.mu.Unlock()
	r.Attestations.m[key] = list
}

func (r *Authorization) GetAttestations(key string) ([]model.Discipline, bool) {
	r.Attestations.mu.Lock()
	defer r.Attestations.mu.Unlock()
	value, ok := r.Attestations.m[key]
	return value, ok
}

func NewAuthorization() *Authorization {
	return &Authorization{
		Cookies:      *NewSafeMap(),
		Attestations: *NewAttestationCache(),
	}
}
