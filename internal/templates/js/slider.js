// Получаем контейнер и кнопки для перемещения
document.addEventListener("DOMContentLoaded", function() {
const slider = document.querySelector('.schedule_block');
const prevBtn = document.querySelector('.prev-btn');
const nextBtn = document.querySelector('.next-btn');
const scheduleBlock = document.getElementById('lesson_list');

// Устанавливаем начальное значение смещения
let offset = 0;

// Обработчик клика на кнопку "Вперед"
nextBtn.addEventListener('click', () => {
    // Делаем GET запрос к серверу, передавая текущее смещение
    // fetch(`/get_data?offset=${offset}`)
    fetch(`/web/get_lesson/${offset}`)
        .then(response => response.text())
        .then(html => {
            // Создаем новый контейнер и добавляем его в слайдер
            // const newContainer = document.createElement('div');
            scheduleBlock.innerHTML = ""
            scheduleBlock.innerHTML = html.trim();

            // Увеличиваем смещение на ширину контейнера
            // offset += scheduleBlock.offsetWidth;
            offset += 1
        })
});

// Обработчик клика на кнопку "Назад"
prevBtn.addEventListener('click', () => {
    // Если смещение равно 0, то ничего не делаем
    if (offset === 0) return;

    // Удаляем последний контейнер из слайдера
    const lastContainer = slider.lastChild;
    slider.removeChild(lastContainer);

    // Уменьшаем смещение на ширину контейнера
    offset -= lastContainer.offsetWidth;
});
});
