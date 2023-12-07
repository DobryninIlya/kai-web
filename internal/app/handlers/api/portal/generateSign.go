package api_handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

type queryParameter struct {
	Key   string
	Value string
}

const SignInvalid = "подпись невалидна"

func GetSignForURLParams(params url.Values, secretKey string) string {
	var (
		query []queryParameter
	)

	for key, value := range params {
		if key == "sign" || key == "secret" || key == "loading" || key[0] == '_' {
			continue
		}
		query = append(query, queryParameter{key, value[0]})
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

	hash = strings.ReplaceAll(hash, "+", "-")
	hash = strings.ReplaceAll(hash, "/", "_")
	hash = strings.ReplaceAll(hash, "=", "")

	return hash
}

func GetSignForStringParams(params string, secretKey string) string {
	var hashCreator = hmac.New(sha256.New, []byte(secretKey))
	hashCreator.Write([]byte(params))

	var hash = base64.URLEncoding.EncodeToString(hashCreator.Sum(nil))

	hash = strings.ReplaceAll(hash, "+", "-")
	hash = strings.ReplaceAll(hash, "/", "_")
	hash = strings.ReplaceAll(hash, "=", "")

	return hash
}
