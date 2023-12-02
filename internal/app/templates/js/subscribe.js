document.getElementById("next_payment").addEventListener("click", function() {
    var successBody = document.querySelector(".advantages_text");
    var subscribePay = document.querySelector(".subscribe_pay");

    // Скрываем элемент .success_body
    successBody.style.display = "none";

    // Отображаем элемент .subscribe_pay
    subscribePay.style.display = "block";
});

document.getElementById("make_payment").addEventListener("click", function() {
    // Получаем значения параметров client_id и tariff из URL
    var urlParams = new URLSearchParams(window.location.search);
    var clientId = urlParams.get('client_id');
    var requestUrl = ""
    var tariff = document.querySelector('input[name="tariff"]:checked').value;
    if (clientId) {
        requestUrl = '/payments/request?level=' + tariff + '&client_id=' + clientId;

    } else {
        var tgWebAppStartParamURL = urlParams.get('tgWebAppStartParam');
        tgWebAppStartParamURL = tgWebAppStartParamURL.replace(/=/, '%3D');
        requestUrl = '/payments/request?tgWebAppStartParam=' + tgWebAppStartParamURL + '___level%3D' + tariff;

    }
    alert(requestUrl)

    // Формируем URL для запроса
    // var requestUrl = '/payments/request?level=' + tariff + '&client_id=' + clientId;

    // Отправляем запрос на сервер
    fetch(requestUrl)
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            // Обрабатываем результат
            var resultUrl = data.result;
            // Перенаправляем пользователя по ссылке из result
            window.location.href = resultUrl;
        })
        .catch(function(error) {
            console.log('Ошибка:', error);
        });
});