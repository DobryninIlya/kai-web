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
	vkChan    chan []int64
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
		vkChan:    make(chan []int64, 10),
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

func (m *Mailing) loadRecipients() {
	offset := 100
	for i := 0; i < len(m.vkClients); i += offset { // Заполняем канал срезами по 100 пользователей для отправки
		currentOffsetUsers := make([]int64, 100)
		leftOffset := i
		if leftOffset+offset > len(m.vkClients) {
			currentOffsetUsers = m.vkClients[leftOffset:]
		} else {
			currentOffsetUsers = m.vkClients[leftOffset : leftOffset+offset]
		}
		m.vkChan <- currentOffsetUsers
	}
	close(m.vkChan)
}

func (m *Mailing) SendMailVK(message string, buttons string) int {
	if len(m.vkClients) == 0 {
		m.log.Printf("Попытка отправить сообщение пустому списку оповещения")
		return 0
	}
	go m.loadRecipients()

	m.log.Logf(
		logrus.InfoLevel,
		"Подготовка получателей начата. Всего: %v пользователей",
		len(m.vkClients),
	)
	for i := 0; i < goroutines; i++ { // Создаем горутины и отправляем сообщения в них методом отправки по 100 пользователей
		go func() {
			for recipient := range m.vkChan {
				m.vkLocker.Lock()
				start := time.Now()
				var j int
				var result bool
				for result = m.vk.SendMessageVKids(m.log, recipient, message, buttons); !result; {
					j++
					time.Sleep(time.Second * 10)
					m.log.Logf(
						logrus.WarnLevel,
						"Результат отправки сообщения не получен. Ожидаем: %v сек. Попытка: %v",
						10,
						j,
					)
					result = m.vk.SendMessageVKids(m.log, recipient, message, buttons)
					if j >= 5 {
						m.log.Logf(
							logrus.WarnLevel,
							"Результат отправки сообщения не получен. Попытка: %v. Провал",
							j,
						)
						break
					}
				}
				if result {
					m.log.Logf(
						logrus.InfoLevel,
						"Сообщение отправлено в %v",
						recipient,
					)
				}
				sleepTime := time.Second/vkAPIQueryLimit - time.Since(start)
				if sleepTime < 0 {
					sleepTime = -sleepTime
				}
				time.Sleep(sleepTime)
				m.vkLocker.Unlock()
			}
		}()
	}
	return len(m.vkClients)
}
