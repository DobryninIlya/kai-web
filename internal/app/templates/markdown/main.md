
# Документация к API КапиПары
---
<center>
#### [Аворизация по токену](doc/autorization)
#### [Новости](news/?count=10)
#### [Условия публикации новостей](doc/news_content)
#### [Загрузка фотографий](doc/images)
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

### 3. POST **api/registration**
Регистрирует нового API клиента и возвращает его токен.
Если пользователь зарегистрирован - возвращает токен.

Принимаемый payload:
```json
{
	"login": "NameSM", // Логин пользователя
	"password": "password123" // Пароль пользователя (от личного кабинета)
	
}
```

В случае неверного пароля возвращается:

```json
{
	"error": "wrong password"
}
```

При успешном ответе возвращается токен:
```json
{
	"result": {
		"token":"301437403bdf09fa95d14f568375178b59c3aff51dad0a97241a4bc1ed4cce33"
	}
}
```
#### Длина токена: 64

Ошибки:

[Аворизация по токену](doc/autorization)

### 4. GET **/api/token/whoiam** _(с авторизацией)_
Возвращает известную информацию о владельце токена
```json
{
	"result": {
		"uid": "Tjf8dC9mZ6V8cLzGltKSCrhrLdq1       ",
		"token": "ce20ff5c18c65a5fdsfsdfsdfsdfdsfaqr28b26f9ed263e8eb00fe7e28e49d9705",
		"create_date": {
			"Time": "2023-11-16T00:00:00Z",
			"Status": 2,
			"InfinityModifier": 0
		},
		"name": "Иванов Иван Иванович",
		"groupname": 4215,
		"login": "Login"
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
# DEPRECATED | Секция не поддерживается. Используйте api/auth/attestation

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

# **Секция /api/auth**

### 15. GET **/api/auth/personal** 
Получает ФИО по токену

```json
{
  "result": {
    "user": {
      "FirstName": "Ivan",
      "LastName": "Ivanov",
      "MiddleName": "Ivanovich"
    }
  }
}

```


### 16. GET **/api/auth/group**
Получает номер группы по токену

```json
{
  "result": {
	"group": "4215"
  }
}
```


### 17. GET **/api/auth/attestation** 
Получает баллы аттестации по токену

```json
{
  "result": {
    "attestation": [
      {
        "Number": 1,
        "Name": "Технологическая (проектно-технологическая) практика",
        "Assessments": [
          {
            "YourScore": 20,
            "MaxScore": 20
          },
			... 10 assessments
```

### 19. POST **/api/auth/profile_photo**
Загружает фотографию профиля по токену

Необходимо передать multipart/form-data с полем file, содержащим файл изображения

Пример ответа:
```json
{
	"result": {
		"status": "ok"
	}
}
```


### 18. GET **/api/auth/profile_photo**
Получает url на фотографию профиля по токену

Полученный путь необходимо прибавлять к https://kai.ru

Пример ответа:
```json
{
	"result": {
		"photo_url": "/image/user_male_portrait?img_id=12532825&img_id_token=KfVwudcwBczXPoUNBk8nC5fSWWM%3D&t=1701531843682"
	}
}
```
