document.getElementById("go-back-btn").addEventListener("click", function() {
    window.history.back();
});

// Получаем все кнопки "open_desc_button"
const buttons = document.querySelectorAll('.open_desc_button');

// Обходим все кнопки и добавляем обработчик события "click"
buttons.forEach((button) => {
    button.addEventListener('click', () => {
        // Получаем элемент "attestation_block_assessment"
        const assessment = button.parentNode.parentNode.querySelector('.attestation_block_assessment');
        // Получаем изображение "arrow-bottom-btn"
        const arrow = button.querySelector('.arrow-bottom-btn');

        // Если элемент "attestation_block_assessment" скрыт, то отображаем его и меняем изображение стрелки
        if (assessment.style.display === 'none') {
            assessment.style.display = 'block';
            arrow.src = '/static/img/icon-arrow-top.svg';
        } else {
            // Иначе скрываем элемент "attestation_block_assessment" и меняем изображение стрелки
            assessment.style.display = 'none';
            arrow.src = '/static/img/icon-arrow-bottom.svg';
        }
    });
});

window.addEventListener('beforeunload', function (e) {
    document.getElementById('loading').style.display = 'block';
});

// Скрыть элемент загрузки, когда страница полностью загрузится
window.addEventListener('load', function (e) {
    document.getElementById('loading').style.display = 'none';
});