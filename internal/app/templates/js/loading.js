var url = new URL(window.location.href);

// Удаляем параметр loading
url.searchParams.delete('loading');

let tgWebAppStartParamURL = url.searchParams.get("tgWebAppStartParam");
// Проверяем, есть ли параметр "tgWebAppStartParam"

if (tgWebAppStartParamURL) {
    tgWebAppStartParamURL = tgWebAppStartParamURL.replace(/---/g, '/');
    // Разбиваем параметр на отдельные значения
    const paramPairs = tgWebAppStartParamURL.split("___");

    // Преобразуем каждое значение в отдельный URL-параметр
    paramPairs.forEach(pair => {
        const [key, value] = pair.split("=");
        let decodedValue = decodeURIComponent(value);
        if (key == "loading") {
            decodedValue = false;
        }
        url.searchParams.append(key, decodedValue);
        console.log(key + " = " + decodedValue);
    });
    url.searchParams.delete('loading');
    url.searchParams.delete('tgWebAppStartParam');
}

function loadPage() {
    fetch(url.toString())
        .then(response => {
            if (response.ok) {
                return response.text();
            } else {
                throw new Error("Ошибка загрузки страницы");
            }
        })
        .then(content => {
            document.body.innerHTML = content;
        })
        .catch(error => {
            console.error(error);
            var loadingText = document.getElementById("loading_text_p");
            loadingText.textContent = "В этот раз потребуется чуть больше времени...";
            setTimeout(loadPage, 5000);
        // TODO добавить флаг выхода из этого цикла рекурсии
        });
}

var loadingText = document.getElementById("loading_text_p");
var dots = 0;
setInterval(function() {
    dots = (dots + 1) % 4;
    loadingText.textContent = "Загрузка" + ".".repeat(dots);
}, 500);

loadPage();