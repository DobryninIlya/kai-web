document.addEventListener('DOMContentLoaded', function() {
    var openFormBtn = document.querySelector('.open_form_btn');
    var authForm = document.querySelector('.auth_form');

    openFormBtn.addEventListener('click', function() {
        openFormBtn.style.display = 'none';
        authForm.style.display = 'block';
    });
});

var phrases = [
    'Капибара проверяет ваш пароль.',
    'Капибара проверяет ваш пароль..',
    'Капибара проверяет ваш пароль...',
    'Решаем, можно ли вас пускать.',
    'Решаем, можно ли вас пускать..',
    'Решаем, можно ли вас пускать...',
    'Капибара пьет латте.',
    'Капибара пьет латте..',
    'Капибара пьет латте...',
    'Смотрим в базу данных.',
    'Смотрим в базу данных..',
    'Смотрим в базу данных...',
    'Иногда это занимает немного времени.',
    'Иногда это занимает немного времени..',
    'Иногда это занимает немного времени...',
    'Иногда это может продлиться даже минуту.',
    'Иногда это может продлиться даже минуту..',
    'Иногда это может продлиться даже минуту...',
];

document.getElementById('main_auth_form').addEventListener('submit', function(e) {
    e.preventDefault(); // Отменяем стандартное поведение формы
    var formData = new FormData(this);
    var currentIndex = 0; // Индекс текущей фразы
    var intervalId; // Идентификатор интервала

    // Функция для обновления текста в элементе "status"
    function updateStatusText(text) {
        document.getElementById('status').innerHTML = text;
    }

    // Функция для выполнения запроса
    function performRequest() {
        document.getElementById('main_auth_form_btn').disabled = true;
        fetch('/portal/authorization/telegram' + formData.get('params'), {
            method: 'POST',
            body: formData,
            redirect: 'follow' // Указываем, что fetch должен следовать за редиректами
        })
            .then(response => {
                if (!response.ok && response.status != 302) {
                    throw response;
                }
                if (response.redirected || response.status == 302) {
                    return response.json()// Выполняем переадресацию на URL из ответа
                } else {
                    return response
                }
            })
            .then(data => {
                if (data.result) {
                    document.getElementById('result').innerHTML = "Ждите, перенаправляю..."
                    // TODO: убрать эту штуку и сделать нормально:
                    document.getElementById('result').style = "background-color: #4CAF50; padding: 10px;"
                    window.location.href = data.result['redirect_url'];
                    return;
                }
                if (data.error) {
                    document.getElementById('result').innerHTML = data.error;
                } else {
                    document.getElementById('result').innerHTML = data;
                }
            })
            .catch(error => {
                console.error('Error:', error);
                if (error.status == 302) {
                    window.location.href = error.url;
                    return;
                }
                error.json().then(errorMessage => {
                    intervalId = clearInterval(intervalId); // Останавливаем интервал
                    updateStatusText('');
                    document.getElementById('main_auth_form_btn').disabled = false;
                    document.getElementById('result').innerHTML = errorMessage.error;
                    document.getElementById('result').style = "padding: 10px;";
                });
            });
    }

    // Функция для отображения фраз с паузой
    function displayPhrasesWithDelay() {
        updateStatusText(phrases[currentIndex]);
        currentIndex++;
        if (currentIndex === phrases.length) {
            currentIndex = 0; // Возвращаемся к началу массива, если достигли его конца
        }
    }

    intervalId = setInterval(displayPhrasesWithDelay, 750); // Запускаем отображение фраз с паузой
    performRequest(); // Запускаем выполнение запроса
});