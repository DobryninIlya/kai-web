
# Документация к API КапиПары
---
<center>
#### [Аворизация по токену](doc/autorization)
#### [Новости](news/?count=10)
#### [Условия публикации новостей](doc/news_content)
</center>

---
## *Основной* адрес вызова API:

> ### _schedule-bot.kai.ru/api_


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

# **Секция /schedule** _(с авторизацией)_

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


# **Секция /token** 

### 3. GET **/api/token** 
Регистрирует нового API клиента и возвращает его токен.
Принимаемый payload:
```json
{
    "uid": "adbcdef123", // unique constraint max-len=35
    "device_tag": "SM-1234" // len=16
}
```
uid - идентификатор пользователя от Firebase

device_tag - тэг устройства от производителя

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

 **unique constraint failed for one of the field** - *поле uid уже имеется в базе данных*

 **the length of one of parameters is too much** - *поле uid или device_tag превышает допустимую длину*

 **user not found** - *ошибка получения данных от Firebase: данные о пользователе не найдены*

[Аворизация по токену](doc/autorization)

### 4. GET **/api/token/whoiam** _(с авторизацией)_
Возвращает известную информацию о владельце токена
```json
{
	"result":
	{
		"uid":"dsadasdasdasdasdasds",
		"device_tag":"SM-1234         ",
		"create_date":
		{
			"Time":"2023-10-06T00:00:00Z",
			"Status":2,
			"InfinityModifier":0
		}
	}
}
```

---

# **Секция api/feedback** _(с авторизацией)_

### 5. POST **/api/feedback/** 
Обрабатывает фидбек от пользователя и отправляет его в соответствующий телеграм-чат

Принимаемый payload:

```json
{
    "version": "alpha-0.1", // Версия приложения
    "text": "the best app" // len=16
}
```

---

# **Секция api/attestation** _(с авторизацией)_

### 6. GET **/api/attestation/** 
Получает баллы аттестации

Принимаемые URL параметры:

| **Параметр** | **Описание**  |
|--------------|---------------|
| p_fac        | Факультет     |
| p_kurs       | Курс          |
| p_group      | Группа        |
| p_stud       | Студент       |
| p_zach       | Номер зачетки |

В ответ возвращается список BRS:

```json
{
  "result": {
    "scores": [
      {
        "index": 1,
        "name": "Теория вероятностей и математическая статистика (экз.)",
        "scoreCurrent1": 0,
        "scoreMax1": 0,
        "scoreCurrent2": 0,
        "scoreMax2": 0,
        "scoreCurrent3": 0,
        "scoreMax3": 0,
        "previouslyScore": 0,
        "additionalScore": 0,
        "debt": 0,
        "final": 0,
        "result": ""
      },
........
```

### 7. GET **/api/attestation/faculties** 

Получает список факультетов

В ответ возвращается список факультетов: 
```json
{
  "result": {
    "faculties": {
      "1": "ИАНТЭ",
      "2": "ФМФ",
      "28": "ВШПИТ и ИИэП",
      "3": "ИАЭП",
      "4": "ИКТЗИ",
      "5": "ИРЭТ"
    }
  }
}
```

### 8. GET **/api/attestation/groups** 

Получает список групп

Обязательные URL параметры: p_fac, p_kurs

В ответ возвращается список групп:
```json
{
  "result": {
    "groups": {
      "4201": "9062",
.....
```

### 9. GET **/api/attestation/persons** 

Получает список ФИО студентов

Обязательные URL параметры: p_fac, p_kurs, p_group

Возвращает json с ассоциативным массивом, где ключ - номер студента, значение - ФИО студента:
```json
{
  "result": {
    "persons": {
      "161400": "Москаль Илон Маск",
.....
```

---

### 10. GET **/api/week**
Получает текущие настройки четности. Для корректности высчитывания расписания, необходимо прибавлять полученный результат к остатку от деления номера текущей недели на 2.

Формат ответа json:

```json
{
  "result": {
    "week_parity": 0
  }
}
```

# **Секция /api/groups** _(с авторизацией)_
 
### 11. GET **/api/groups/{groupName}** 
Получает id группы по ее номеру

Формат ответа json:

```json
{
  "result": {
	  "group_id": 24691
  }
}
```

# **Секция /api/news**

### 12. GET **/api/news/{newsId}** 
Отрисовывает и возвращает html страницу с новостью.
Ориентирован на мобильный просмотр.
[Пример](news/1)

### 13. POST **/api/news**
Создает новость. Пока без описания.

[Новостной редактор](news/create)

### 14. GET **/api/news/previews**
Получате новостные превью в порядке убывания от самого последнего

Возможные параметры: 

**count** - число получаемых постов (обязательный)

**offset** - отступ от начала. Если не указан, то 0

Возвращает json со списком объектов News (только без длинного тэга body)

```json
{
  "result": {
    "news": [
      {
        "header": "3 | TEST HEADER, SOME INTRESTIG                                                 ",
        "description": "3 | very very very INTRESTING and AWESOME description                                                                                                                                                                                                     ",
        "date": null
      },
	.....
```