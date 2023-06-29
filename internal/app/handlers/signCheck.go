package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"sort"
	"strings"
)

//func getSortedMap(params url.Values) []string {
//	keys := make([]string, 0, len(params))
//	for k := range params {
//		keys = append(keys, k)
//	}
//
//	// Сортируем ключи по алфавиту
//	sort.Strings(keys)
//	return keys
//}

//func IsCorrectSign(url url.Values) bool {
//	sortedUrl := getSortedMap(url)
//	sign := url["sign"][0]
//	resultUrl := ""
//	for _, val := range sortedUrl {
//		if !strings.HasPrefix(val, "vk_") {
//			continue
//		}
//		resultUrl += val + "=" + url[val][0] + "&"
//	}
//	resultUrl = resultUrl[:len(resultUrl)-1]
//	key := os.Getenv("SECRET_KEY")
//	keys := make([]string, 0, len(url))
//	for k := range params {
//		keys = append(keys, k)
//	}
//	VerifyLaunchParams(resultUrl, "xENraXo5RzenYDOxAaLZ")
//	key = "xENraXo5RzenYDOxAaLZ"
//	hash := hmac.New(sha256.New, []byte(key))
//	hash.Write([]byte(resultUrl))
//	//signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
//	signature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
//	return reflect.DeepEqual(signature[:len(signature)-1], sign)
//}

type queryParameter struct {
	Key   string
	Value string
}

const SignInvalid = "подпись невалидна"

func VerifyLaunchParams(querySearch string, secretKey string) error {
	var searchIndex = strings.Index(querySearch, "?")

	// Необходимо удалить всё, что находится до search части в случае, если
	// эта часть существует.
	if searchIndex >= 0 {
		querySearch = querySearch[searchIndex+1:]
	}

	var (
		// Отфильтрованные параметры запуска. Мы используем именно
		// слайс по той причине, что позже нам будет необходимым этот слайс
		// отсортировать по возрастанию ключа параметра.
		query []queryParameter
		// Подпись, которая была сгенерирована сервером ВКонтакте и основана на
		// параметрах из query.
		sign string
	)

	// Разделяем параметры запуска на вхождения, разделенные знаком "&".
	for _, part := range strings.Split(querySearch, "&") {
		var keyAndValue = strings.Split(part, "=")
		var key = keyAndValue[0]
		var value string

		if len(keyAndValue) > 1 {
			value = keyAndValue[1]
		}

		// Мы обрабатываем только те ключи, которые начинаются с префикса "vk_".
		// Все остальные ключи в создании подписи не участвуют.
		if strings.HasPrefix(key, "vk_") {
			query = append(query, queryParameter{key, value})
		} else if key == "sign" {
			// Если ключ равен "sign", то в значении записана подпись параметров
			// запуска.
			sign = value
		}
	}

	// В случае, если подпись параметров не удалось найти, либо параметров с
	// префиксом "vk_" передано не было, мы считаем параметры запуска невалидными.
	if sign == "" || len(query) == 0 {
		return errors.New(SignInvalid)
	}

	// Сортируем параметры запуска по порядку их возрастания.
	sort.SliceStable(query, func(a int, b int) bool {
		return query[a].Key < query[b].Key
	})

	// Далее снова превращаем параметры запуска в единую строку.
	var queryString = ""

	for idx, param := range query {
		if idx > 0 {
			queryString += "&"
		}
		queryString += param.Key + "=" + url.PathEscape(param.Value)
	}

	// Далее нам необходимо вычислить хэш SHA-256.
	var hashCreator = hmac.New(sha256.New, []byte(secretKey))
	hashCreator.Write([]byte(queryString))

	var hash = base64.URLEncoding.EncodeToString(hashCreator.Sum(nil))

	// Далее по правилам создания параметров запуска ВКонтакте, необходимо
	// произвести ряд замен символов.
	hash = strings.ReplaceAll(hash, "+", "-")
	hash = strings.ReplaceAll(hash, "/", "_")
	hash = strings.ReplaceAll(hash, "=", "")

	if sign != hash {
		return errors.New(SignInvalid)
	}
	return nil
}
