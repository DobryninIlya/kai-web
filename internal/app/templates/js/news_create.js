// получаем элементы формы и область предпросмотра
var form = document.getElementById('form');
var titleInput = document.getElementById('title');
var descriptionTextarea = document.getElementById('description');
var bodyTextarea = document.getElementById('body');
var tagInput = document.getElementById('tag');
var imageInput = document.getElementById('image');
var previewDiv = document.getElementById('preview');

// обработчик изменения поля url
imageInput.addEventListener('input', function() {
    // получаем html код из полей
    var url_preview = '<img src="' + imageInput.value + '">';
    var title = '<div class="header_text">' + titleInput.value + '</div> <hr>';
    var body = '<div class="news_body">' + bodyTextarea.value + '</div>';
    // отображаем его в области справа
    previewDiv.innerHTML = url_preview + title + body;
});

// обработчик изменения поля заголовка
titleInput.addEventListener('input', function() {
    // получаем html код из полей
    var title = '<div class="header_text">' + titleInput.value + '</div> <hr>';
    var body = '<div class="news_body">' + bodyTextarea.value + '</div>';
    // отображаем его в области справа
    previewDiv.innerHTML = title + body;
});

// обработчик изменения поля тела текста
bodyTextarea.addEventListener('input', function() {
    // получаем html код из полей
    var title = '<div class="header_text">' + titleInput.value + '</div><hr>';
    // var description = '<p>' + descriptionTextarea.value + '</p>';
    var body = '<div class="news_body">' + bodyTextarea.value + '</div>';
    // отображаем его в области справа
    previewDiv.innerHTML = title + body;
});

// обработчик отправки формы
form.addEventListener('submit', function(event) {
    event.preventDefault(); // отменяем стандартное поведение формы
    // создаем объект с данными формы
    var data = {
        header: titleInput.value,
        description: descriptionTextarea.value,
        body: bodyTextarea.value,
        tag: tagInput.value,
        preview_url: imageInput.value
    };
    // отправляем данные на сервер через fetch
    fetch('/api/news/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
        .then(function(response) {
            if (response.ok) {
                // если сервер вернул успешный ответ, выводим сообщение об успехе
                alert('Новость создана!');
                response.json().then(function(data) {
                    var newsId = data.result.id;
                    window.location.href = '/api/news/' + newsId;
                });
                previewDiv.innerHTML = ""
                titleInput.value = ""
                descriptionTextarea.value = ""
                tagInput.value = ""
                imageInput.value = ""
            } else {
                // если произошла ошибка, выводим сообщение об ошибке
                alert('Ошибка при создании новости');
            }
        })
        .catch(function(error) {
            // если произошла ошибка при выполнении запроса, выводим сообщение об ошибке
            alert('Ошибка при создании новости: ' + error.message);
        });
});