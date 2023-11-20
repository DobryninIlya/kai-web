var url = new URL(window.location.href);

// Удаляем параметр loading
url.searchParams.delete('loading');

let tgWebAppStartParamURL = url.searchParams.get("tgWebAppStartParam");
tgWebAppStartParamURL = tgWebAppStartParamURL.replace(/---/g, '/');

// Проверяем, есть ли параметр "tgWebAppStartParam"
if (tgWebAppStartParamURL) {
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
// Загружаем страницу из нового URL
fetch(url.toString())
    .then(response => response.text())
    .then(content => {
        document.body.innerHTML = content;
    });

var loadingText = document.getElementById("loading_text_p");
var dots = 0;
setInterval(function() {
    dots = (dots + 1) % 4;
    loadingText.textContent = "Загрузка" + ".".repeat(dots);
}, 500);