package tg_api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
)

const tgTemplate = "https://api.telegram.org/bot%s/%s%s"
const tgSendMethod = "sendMessage?"

type APItg struct {
	tgToken    string
	tgTemplate string
}

func NewAPItg() *APItg {
	return &APItg{
		tgToken:    os.Getenv("TG_TOKEN"),
		tgTemplate: tgTemplate,
	}
}

func (s APItg) SendMessageTG(log *logrus.Logger, uId int64, message string, buttons string, threadId int) bool {
	if uId == 0 {
		log.Printf("Попытка отправить сообщение некорректному айди")
		return false
	}
	params := fmt.Sprintf("chat_id=%v&text=%v&reply_markup=%s&parse_mode=markdown&message_thread_id=%v",
		uId,
		url.QueryEscape(message),
		url.QueryEscape(buttons),
		threadId,
	)
	url := fmt.Sprintf(s.tgTemplate, s.tgToken, tgSendMethod, params)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка API. Отправка сообщений: %v", err)
		return false
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("TG: При этом возникла ошибка: %v", err)
		return false
	}
	defer resp.Body.Close()
	return true
}
