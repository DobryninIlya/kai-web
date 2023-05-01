const sliderContainer = document.querySelector('.slider-container');
const prevBtn = document.querySelector('.prev-btn');
const nextBtn = document.querySelector('.next-btn');
const scheduleBlock = document.getElementById('lesson_list');

let currentPage = 1;

function loadPage(page) {
    fetch(`/web/get_lesson/${page}`)
        .then(response => response.text())
        .then(html => {
            scheduleBlock.innerHTML = html.trim();
        });
}

loadPage(currentPage);

prevBtn.addEventListener('click', () => {
    if (currentPage > 1) {
        currentPage--;
        sliderContainer.scrollTo({
            left: sliderContainer.scrollLeft - sliderContainer.offsetWidth,
            behavior: 'smooth'
        });
    }
});

nextBtn.addEventListener('click', () => {
    currentPage++;
    loadPage(currentPage);
    sliderContainer.scrollTo({
        left: sliderContainer.scrollLeft + sliderContainer.offsetWidth,
        behavior: 'smooth'
    });
});
