
# Документация к API КапиПары
---
##### [Аворизация по токену](doc/autorization)
---
## *Основной* адрес вызова API:

> schedule-bot.kai.ru/api


Успешный ответ:

>	{
		"result": {
		полученный ответ (здесь также json)
		}
	}


Ответ с ошибкой: 

>	{
	"error":"missing or incorrect id"
	}

Вместе с ошибкой также передается соответствующий статус **код** ответа.

---

# **Секция /schedule**

### 1.  GET **/api/schedule/_{groupid}_** 
Получает сохраненное расписание группы полностью без каких либо искажений.
Формат ответа структуры Schedule, которая потом маршалится в json:
```go
	type Schedule struct {
		Day3 []Lesson `json:"3,omitempty"`  
		Day2 []Lesson `json:"2,omitempty"`  
		Day1 []Lesson `json:"1,omitempty"`  
		Day6 []Lesson `json:"6,omitempty"`  
		Day5 []Lesson `json:"5,omitempty"`  
		Day4 []Lesson `json:"4,omitempty"`  
	}

	type Lesson struct {  
		PrepodNameEnc string `json:"prepodNameEnc"`  
		DayDate string `json:"dayDate"`  
		AudNum string `json:"audNum"`  
		DisciplName string `json:"disciplName"`  
		BuildNum string `json:"buildNum"`  
		OrgUnitName string `json:"orgUnitName"`  
		DayTime string `json:"dayTime"`  
		DayNum string `json:"dayNum"`  
		Potok string `json:"potok"`  
		PrepodName string `json:"prepodName"`  
		DisciplNum string `json:"disciplNum"`  
		OrgUnitId string `json:"orgUnitId"`  
		PrepodLogin string `json:"prepodLogin"`  
		DisciplType string `json:"disciplType"`  
		DisciplNameEnc string `json:"disciplNameEnc"`  
	}
```
### 2. GET **/api/schedule/_{groupid}_/by_margin** 
Получает расписание на конкретный день с соответствующим отступом.

Отступ задается с помощью параметра *?margin={int}*
В случае отсутствия параметра margin в запросе, отступ приравнивается к 0.

В ответ приходит список из Lesson

---

# **Секция /get_token**

### 3. GET **/api/get_token** 
Регистрирует нового API клиента и возвращает его токен.
Принимаемый payload:
```json
{
    "device_id": "adbcdef123", // unique constraint len=16
    "device_tag": "SM-1234" // len=16
}
```
При успешном ответе возвращается:
```json
{
	"result": {
		"token":"301437403bdf09fa95d14f568375178b59c3aff51dad0a97241a4bc1ed4cce33"
	}
}
```
#### Длина токена: 64

Ошибки:

 **unique constraint failed for one of the field** - *поле device_id уже имеется в базе данных*

 **the length of one of parameters is too much** - *поле device_id или device_tag превышает допустимую длину*

[Аворизация по токену](doc/autorization)

# **Секция /feedback**

### 3. POST **/api/feedback/** 
Обрабатывает фидбек от пользователя и отправляет его в соответствующий телеграм-чат

Принимаемый payload:

```json
{
    "version": "alpha-0.1", // Версия приложения
    "text": "the best app" // len=16
}
```