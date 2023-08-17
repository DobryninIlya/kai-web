const lessons = document.querySelectorAll('.schedule_item');
const popup = document.getElementById('popup');

console.log(lessons)
lessons.forEach((lesson) => {
    lesson.addEventListener('click', () => {
        const info = lesson.getAttribute('id');
        popup.innerHTML = info;
        popup.style.display = 'block';
    });
});

popup.addEventListener('click', () => {
    popup.style.display = 'none';
});