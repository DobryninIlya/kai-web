const sliderContainer = document.querySelector('.schedule_block');
const scheduleExamBlock = document.querySelector('.schedule_exam_block');
const scoreBlock = document.querySelector('.score_block');
const teacherBlock = document.querySelector('.teacher_block');
const sliderWrapper = document.querySelector('.schedule-wrapper');
const prevButton = document.querySelector('.arrow-left');
const nextButton = document.querySelector('.arrow-right');
const find = document.querySelector('.schedule_block');
const authErrorBlock = document.querySelector(".auth_error_block");
const enterButton = document.querySelector(".enter_button");

const popup = document.getElementById('popup');
const popupOverlay = document.createElement('div');
popupOverlay.classList.add('popup-overlay');

const loaderHTML = "<p class=\"loader_text\">Подождите.<br>Идет загрузка данных....</p>"
let slideIndex = 0;
let isFetching = false;
let newElements = false;
const slideWidth = sliderContainer.offsetWidth;
var elem;

const params = new URLSearchParams(window.location.search);
const user_id = params.get('vk_user_id');

var controlParamsString = "?"
for (const [key, value] of params.entries()) {
    controlParamsString = controlParamsString + key + "=" + value + "&"
}
controlParamsString= controlParamsString.slice(0, -1)

var today = new Date();



// tomorrow.setDate(tomorrow.getDate() + 1);

const daysOfWeek = ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота', 'Воскресенье', 'Понедельник'];
const monthsOfYear = ['Января', 'Февраля', 'Марта', 'Апреля', 'Мая', 'Июня', 'Июля', 'Августа', 'Сентября', 'Октября', 'Ноября', 'Декабря'];


function nextButtonFunction() {
    if (slideWidth*sliderWrapper.childElementCount >= sliderContainer.offsetWidth && sliderWrapper.childElementCount-1 < 10) {
        newElements = loadSlides();
    }
    if (newElements || slideIndex+1 < 10) {
        slideIndex++;
        sliderWrapper.style.transform = `translate(${-slideIndex * (slideWidth)}px)`;
        newElements = false

        resizeWrapper()
    }
};

function prevButtonFunction() {
    if (slideIndex > 0) {
        slideIndex--;
        sliderWrapper.style.transform = `translate(${-slideIndex * (slideWidth)}px)`;
    }
    resizeWrapper()
};


// Функция для загрузки новых слайдов
 function loadSlides(margin = null) {
    if (isFetching && sliderWrapper.childElementCount > 1) {
        return false;
    }
    isFetching = true
     var pageMargin
     if (margin == null) {
         pageMargin = sliderWrapper.childElementCount
     } else {
         pageMargin = margin
     }

     fetch(`/web/get_lesson/${user_id}${controlParamsString}&margin=${pageMargin}`)
         .then(response => response.text())
         .then(html => {
             if (slideIndex >=10) {
                 isFetching = false
                 return true
             }
             var index = sliderWrapper.childElementCount;
             var currDate = new Date()
             var dayWeekDelta = currDate.getDate() + index
             if (dayWeekDelta > 6) {
                 dayWeekDelta = dayWeekDelta - 7
             }
             currDate.setDate(today.getDate()+ index);
             var dayOfWeek = daysOfWeek[currDate.getDay()];
             var monthOfYear = monthsOfYear[currDate.getMonth()];
             var formattedDate = `${dayOfWeek}, ${currDate.getDate()} ${monthOfYear}`;


             html = " <div class=\"schedule\" id=\"" + index + "\">\n" +
                 "    <div class=\"schedule_header\">\n" +
                 "<button class=\"arrow-left\ arrow \" onclick=\"prevButtonFunction()\"></button>" +
                 "      <p class=\"lesson_date\">" + formattedDate + "</p>\n" +
                 "<button class=\"arrow-right arrow \" onclick=\"nextButtonFunction()\"></button>" +
                 "    </div>\n" +
                 "    <div class=\"lesson_list\" id=\"lesson_list\">" + html + "</div>\n" +
                 "  </div>"
             sliderWrapper.insertAdjacentHTML('beforeend', html);

             elem = find.querySelector('[id="' + index + '"]');
             let elemAfter = elem.querySelector('.lesson_list')
             // find.style.height = elemAfter.offsetHeight + "px" // TODO продолжить тут. - не меняется высота слайдера нормально при свайпах
             var currentSlideHeight = elemAfter.offsetHeight
             if (elemAfter.querySelector('.lesson_none')) {
                 currentSlideHeight = 80
             }
             elem.style.height = `${currentSlideHeight + 55}px`;

             isFetching = false; // устанавливаем флаг в false, чтобы указать, что запрос завершился

             // Если запрос первый, делаем еще один
             if (margin == 0) {
                 loadSlides(1);
             }

             // Всплывающее окно с доп инфой о расписании
             let scheduleItems = document.querySelectorAll('.schedule_item');



             scheduleItems.forEach((item) => {
                 item.addEventListener('click', () => {
                     const hiddenPayload = item.querySelector('.hidden_payload').innerHTML
                     popup.innerHTML = hiddenPayload;
                     popup.classList.add('popup-visible');
                     popup.classList.add('active');
                     popup.classList.remove('unactive');

                     document.body.appendChild(popupOverlay);
                     setTimeout(() => {
                         popup.style.opacity = '1';
                     }, 10);
                 });
                 // });


             popup.addEventListener('click', () => {
                 // popup.style.display = 'none';
                 popup.style.opacity = '0';
                 popup.classList.remove('active');
                 popup.classList.add('unactive');
                  // удаляем затемнение из DOM
                 setTimeout(() => {
                     popup.classList.remove('popup-visible');
                     document.body.removeChild(popupOverlay);
                 }, 300);
             });
             });
         })
         .catch(error => {
             console.error(error);
             isFetching = false; // устанавливаем флаг в false, чтобы указать, что запрос завершился
         });

    return true;
}

function resizeWrapper() {
    elem = find.querySelector('[id="' + slideIndex + '"]');
    let elemAfter = elem.querySelector('.lesson_list')
    // find.style.height = elemAfter.offsetHeight + "px"
    var currentSlideHeight = elemAfter.offsetHeight
    if (elemAfter.querySelector('.lesson_none')) {
        currentSlideHeight = 80
    }
    sliderWrapper.style.height = `${currentSlideHeight + 55}px`;
}





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
    if (slideWidth*sliderWrapper.childElementCount >= sliderContainer.offsetWidth && sliderWrapper.childElementCount-1 < 10 && startX - currentX>50) {
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

        sliderWrapper.style.transform = `translate(${-slideIndex * (slideWidth )}px)`;
        startX = null;
        currentX = null;
    }
    resizeWrapper()
});

enterButton.addEventListener('click', function (event) {
    fetch("/web/verification" + window.location.search)
        .then(response => response.text())
        .then(html => {
            scoreBlock.innerHTML = html
            const script = document.createElement('script');
            script.src = '/static/js/verification.js';
            script.onload = () => {
                getFac() // вызываем функцию getFac() после загрузки скрипта
                authErrorBlock.style.display = "none"
            };
            scoreBlock.appendChild(script);
            getFac();
        })
        .catch(error => console.log(error));
});


// function showLoader() {
//     loaderContainer.style.display = 'block';
// }
//
// function hideLoader() {
//     loaderContainer.style.display = 'none';
// }



window.onload = function () {
    loadSlides(0)

    // loadSlides(1)

};







// Загружаем первые слайды
