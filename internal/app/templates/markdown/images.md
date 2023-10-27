
# Загрузка фотографий

## Необходима авторизация по токену

### 1. POST api/images/tasks

Загружает картинку к заданию и возвращает ссылку на нее
```json
{
  "result": {
    "url": "https://schedule-bot.kai.ru/images/tasks/4201/a0d36c4c-15a6-4423-ad6c-a439d2a73fbb.png"
  }
}
```

HTTP:
```http request
POST /api/image/tasks?token={{TOKEN}} HTTP/1.1
Host: localhost:8283
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Length: 211

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="image"; filename="/C:/....../file.png"
Content-Type: image/png

(data)
------WebKitFormBoundary7MA4YWxkTrZu0gW--
```