package api_handler

import (
	"errors"
	"github.com/sirupsen/logrus"
	"main/internal/app/authorization"
	h "main/internal/app/handlers/web_app"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore"
	"net/http"
	u "net/url"
	"strconv"
	"strings"
)

func NewAuthTelegramHandler(store sqlstore.StoreInterface, log *logrus.Logger, secretKey string, auth authorization.AuthorizationInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.api.portal.NewAuthTelegramHandler"
		var user model.TelegramUser
		url := r.URL.Query()
		tgID := url.Get("tg_id")
		redirectURL := url.Get("redirect_url")
		login := r.FormValue("login")
		password := r.FormValue("password")
		user.UID = "tg" + tgID
		tgIDInt, err := strconv.Atoi(tgID)
		if err != nil {
			log.Logf(
				logrus.WarnLevel,
				"%s : Ошибка преобразования: %v",
				path,
				err.Error(),
			)
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, h.ErrInternal)
			return
		}
		user.TelegramID = int64(tgIDInt)
		user.Login = login
		user.Password = password
		authorization.Encrypt(&user.EncryptedPassword, password)
		_, err = auth.GetCookiesByPassword(user.Login, user.Password)
		if err != nil {
			log.Logf(
				logrus.WarnLevel,
				"%s : Ошибка авторизации: %v",
				path,
				err.Error(),
			)
			if err == authorization.ErrWrongPassword {
				h.ErrorHandlerAPI(w, r, http.StatusBadRequest, errors.New("Неверный пароль"))
				return
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		err = store.API().SaveTelegramAuth(user)
		if err != nil {
			log.Logf(
				logrus.WarnLevel,
				"%s : Ошибка сохранения: %v",
				path,
				err.Error(),
			)
			if redirectURL != "" {
				params := strings.Split(redirectURL, "?")
				var queryParams u.Values
				if len(params) > 1 {
					// Если в URL есть параметры
					queryParams, _ = u.ParseQuery(params[1])
					delete(queryParams, "loading")
				} else {
					// Если в URL нет параметров
					queryParams = make(u.Values)
				}

				queryParams.Add("tg_id", tgID)
				urlParams := queryParams.Encode()

				sign := GetSignForStringParams(urlParams, secretKey)

				redirectURL += "?" + urlParams + "&sign=" + sign

				if loadingParam := queryParams.Get("loading"); loadingParam != "" {
					redirectURL += "&loading=" + loadingParam
				}
				redirectURL = params[0] + "?" + urlParams + "&sign=" + sign
				result := struct {
					RedirectURL string `json:"redirect_url"`
				}{
					RedirectURL: redirectURL,
				}
				h.RespondAPI(w, r, http.StatusFound, result)
				return
			}
			h.ErrorHandlerAPI(w, r, http.StatusBadRequest, err)
			return
		}
		if redirectURL != "" {
			params := strings.Split(redirectURL, "?")
			var queryParams u.Values
			if len(params) > 1 {
				// Если в URL есть параметры
				queryParams, _ = u.ParseQuery(params[1])
				delete(queryParams, "loading")
			} else {
				// Если в URL нет параметров
				queryParams = make(u.Values)
			}

			queryParams.Add("tg_id", tgID)
			urlParams := queryParams.Encode()

			sign := GetSignForStringParams(urlParams, secretKey)

			redirectURL += "?" + urlParams + "&sign=" + sign

			if loadingParam := queryParams.Get("loading"); loadingParam != "" {
				redirectURL += "&loading=" + loadingParam
			}
			redirectURL = params[0] + "?" + urlParams + "&sign=" + sign
			result := struct {
				RedirectURL string `json:"redirect_url"`
			}{
				RedirectURL: redirectURL,
			}
			h.RespondAPI(w, r, http.StatusFound, result)
			return
		}
	}
}
