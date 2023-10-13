package vk_api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

const vkTemplate = "https://api.vk.com/method/%s?v=5.131&%s"
const vkSendMethod = "messages.send"

type APIvk struct {
	vkToken    string
	vkTemplate string
}

func NewAPIvk() *APIvk {
	return &APIvk{
		vkToken:    os.Getenv("VK_TOKEN"),
		vkTemplate: vkTemplate,
	}
}


// SendMessageVKids отправляет сообщение всем пользователям из списка
func (r APIvk) SendMessageVKids(log *logrus.Logger, uId []int64, message string, buttons string) bool {
	if len(uId) == 0 {
		log.Printf("Попытка отправить сообщение пустому списку оповещения")
		return false
	}
	ids := ""
	for _, v := range uId {
		ids += fmt.Sprintf("%v,", v)
	}
	ids = ids[:len(ids)-1]
	randomInt := rand.Int31()
	params := fmt.Sprintf("random_id=%v&peer_ids=%s&access_token=%s&disable_mentions=0&message=%s&keyboard=%s",
		randomInt,
		ids,
		r.vkToken,
		url.QueryEscape(message),
		buttons,
	)
	url := fmt.Sprintf(r.vkTemplate, vkSendMethod, params)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка API. Отправка сообщений: %v", err)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("VK: При этом возникла ошибка: %v", err)
	}
	defer resp.Body.Close()
	return true
}
