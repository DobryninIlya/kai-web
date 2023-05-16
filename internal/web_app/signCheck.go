package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
)

func getSortedMap(params url.Values) []string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	// Сортируем ключи по алфавиту
	sort.Strings(keys)
	return keys
}

func IsCorrectSign(url url.Values) bool {
	sortedUrl := getSortedMap(url)
	sign := url["sign"][0]
	resultUrl := ""
	for _, val := range sortedUrl {
		if !strings.HasPrefix(val, "vk_") {
			continue
		}
		resultUrl += val + "=" + url[val][0] + "&"
	}
	resultUrl = resultUrl[:len(resultUrl)-1]
	key := os.Getenv("SECRET_KEY")
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(resultUrl))
	//signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	signature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return reflect.DeepEqual(signature[:len(signature)-1], sign)
}
