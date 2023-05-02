const sliderContainer = document.querySelector('.schedule_block');
const sliderWrapper = document.querySelector('.schedule-wrapper');
const prevButton = document.querySelector('.prev-button');
const nextButton = document.querySelector('.next-button');
const find = document.querySelector('.schedule_block');

let slideIndex = 0;
const slideWidth = sliderContainer.offsetWidth;
var elem;

const params = new URLSearchParams(window.location.search);
const user_id = params.get('vk_user_id');

// Функция для загрузки новых слайдов
async function loadSlides() {
    const response = await fetch(`/web/get_lesson/${user_id}?margin=${sliderWrapper.childElementCount}`);
    var html = await response.text();
    if (slideIndex >=10) {
        return
    }
    var index = sliderWrapper.childElementCount

    html = " <div class=\"schedule\" id=\"" + index + "\">\n" +
        "    <div class=\"schedule_header\">\n" +
        "      <p class=\"lesson_date\">%v</p>\n" +
        "    </div>\n" +
        "    <div class=\"lesson_list\" id=\"lesson_list\">" + html + "</div>\n" +
        "  </div>"
    sliderWrapper.insertAdjacentHTML('beforeend', html);

    elem = find.querySelector('[id="' + index + '"]');
    elemAfter = elem.querySelector('.lesson_list')
    // find.style.height = elemAfter.offsetHeight + "px" // TODO продолжить тут. - не меняется высота слайдера нормально при свайпах
    var currentSlideHeight = elemAfter.offsetHeight
    elem.style.height = `${currentSlideHeight + 55}px`;
}

// Кнопка "Вперед"
nextButton.addEventListener('click', () => {
    slideIndex++;
    sliderWrapper.style.transform = `translate(${-slideIndex * slideWidth}px)`;

    // elem = find.querySelector('[id="' + slideIndex + '"]');
    // elemAfter = elem.querySelector('.lesson_list')
    // find.style.height = elemAfter.offsetHeight + "px" // TODO продолжить тут. - не меняется высота слайдера нормально при свайпах
    // var currentSlideHeight = elemAfter.offsetHeight
    // find.style.height = `${currentSlideHeight + 55}px`;

    // Если достигнут конец слайдера, загружаем новые слайды
    if (sliderWrapper.offsetWidth - sliderContainer.offsetWidth < sliderContainer.scrollLeft) {
        loadSlides();
    }
    if (slideWidth*sliderWrapper.childElementCount >= sliderContainer.offsetWidth) {
        loadSlides();
    }
});

// Кнопка "Назад"
prevButton.addEventListener('click', () => {
    if (slideIndex > 0) {
        slideIndex--;
        sliderWrapper.style.transform = `translate(${-slideIndex * slideWidth}px)`;


    }
});

// Обработчик свайпа
let startX = null;
let currentX = null;
sliderWrapper.addEventListener('touchstart', (event) => {
    startX = event.touches[0].pageX;
});
sliderWrapper.addEventListener('touchmove', (event) => {
    event.preventDefault();
    currentX = event.touches[0].pageX;
    const delta = startX - currentX;
    sliderWrapper.style.transform = `translate(${-slideIndex * slideWidth - delta}px)`;
});
sliderWrapper.addEventListener('touchend', () => {
    if (slideWidth*sliderWrapper.childElementCount >= sliderContainer.offsetWidth && sliderWrapper.childElementCount-1 < 10) {
        loadSlides();
    }
    if (currentX !== null) {
        const delta = startX - currentX;
        const threshold = slideWidth / 3;
        if (delta > threshold && slideIndex < sliderWrapper.childElementCount - 1 ) {
            slideIndex++;
        } else if (delta < -threshold && slideIndex > 0) {
            slideIndex--;
        }

        sliderWrapper.style.transform = `translate(${-slideIndex * slideWidth}px)`;
        startX = null;
        currentX = null;
    }
});

loadSlides()


// Загружаем первые слайды
