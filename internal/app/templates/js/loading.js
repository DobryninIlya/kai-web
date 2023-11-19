var url = new URL(window.location.href);

// Удаляем параметр loading
url.searchParams.delete('loading');

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