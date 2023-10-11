package mailer

import (
	"github.com/sirupsen/logrus"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tg_api"
	"main/internal/app/vk_api"
	"sync"
	"time"
)

const goroutines = 10
const vkAPIQueryLimit = 20 // 20 per second

type Mailing struct {
	tgClients []int64
	vkClients []int64
	vk        *vk_api.APIvk
	tg        *tg_api.APItg
	storage   *sqlstore.MailingRepository
	log       *logrus.Logger
	vkLocker  sync.Locker
	tgLocker  sync.Locker
}

func NewMailing(storage sqlstore.StoreInterface, log *logrus.Logger) *Mailing {
	return &Mailing{
		tgClients: make([]int64, 0),
		vkClients: make([]int64, 0),
		storage:   storage.Mail(),
		vkLocker:  &sync.Mutex{},
		tgLocker:  &sync.Mutex{},
		log:       log,
		vk:        vk_api.NewAPIvk(),
		tg:        tg_api.NewAPItg(),
	}
}

func (m *Mailing) GetVkClients() ([]int64, error) {
	var err error
	m.vkClients, err = m.storage.GetVKRecipients()
	if err != nil {
		m.log.Errorf("Ошибка получения списка клиентов для рассылки: %v", err)
		return nil, err
	}
	return m.vkClients, nil
}

func (m *Mailing) SendMailVK(message string, buttons string) {
	if len(m.vkClients) == 0 {
		m.log.Printf("Попытка отправить сообщение пустому списку оповещения")
		return
	}
	offset := 100
	recipients := make(chan []int64, 10)
	defer close(recipients)
	m.log.Logf(
		logrus.InfoLevel,
		"Начинаем отправку сообщений всем пользователям из списка. Всего %v пользователей",
		len(m.vkClients),
	)
	for i := 0; i < len(m.vkClients); i += offset { // Заполняем канал срезами по 100 пользователей для отправки
		currentOffsetUsers := make([]int64, 10)
		leftOffset := i
		if leftOffset+offset > len(m.vkClients) {
			currentOffsetUsers = m.vkClients[leftOffset:]
		} else {
			currentOffsetUsers = m.vkClients[leftOffset : leftOffset+offset]
		}
		recipients <- currentOffsetUsers
	}

	for i := 0; i < goroutines; i++ { // Создаем горутины и отправляем сообщения в них методом отправки по 100 пользователей
		go func() {
			for {
				select {
				case recipient := <-recipients:
					m.vkLocker.Lock()
					m.vk.SendMessageVKids(m.log, recipient, message, buttons)
					time.Sleep(time.Second / vkAPIQueryLimit)
					m.vkLocker.Unlock()
				}
			}
		}()
	}
}
