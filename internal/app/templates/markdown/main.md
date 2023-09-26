
# Документация к API КапиПары
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

В ответ приходит список из [Lesson](#L39) 