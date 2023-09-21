document.addEventListener('click', function(event) {

    if (event.target.classList.contains('remove-btn')) {
        var contentRemover = event.target.closest('.content-remover');
        var confirmButtons = contentRemover.querySelector('.confirm-buttons');
        event.target.classList.add('hidden');
        confirmButtons.classList.remove('hidden');
    } else if (event.target.classList.contains('cancel-btn')) {
        console.log("cancel btn clicked")
        var contentRemover = event.target.closest('.content-remover');
        var removeBtn = contentRemover.querySelector('.remove-btn');
        var confirmButtons = contentRemover.querySelector('.confirm-buttons');
        removeBtn.classList.remove('hidden');
        confirmButtons.classList.add('hidden');
    } else if (event.target.classList.contains('confirm-btn')) {
        var contentRemover = event.target.closest('.content-remover');
        var confirmButtons = contentRemover.querySelector('.confirm-buttons');
        var removeBtn = contentRemover.querySelector('.remove-btn');
        // var deleteLessonInfo = contentRemover.querySelector('.delete_lesson_info');
        var popupBody = event.target.closest('.popup_body');
        var return_delete_lesson_info = popupBody.querySelector('.return_delete_lesson_info');
        confirmButtons.classList.add('hidden');
        removeBtn.classList.add('hidden');
        const contentId = contentRemover.getAttribute('content-id');
        const uniqstring = contentRemover.getAttribute('uniqstring');
        return_delete_lesson_info.innerHTML = "Обрабатываю запрос...";
        deleteLesson(contentId, uniqstring)
            .then(message => {
                return_delete_lesson_info.innerHTML = message;
                contentRemover.innerHTML = "";
            });
    } else if (event.target.classList.contains('return-btn')) {
        var popupBody = event.target.closest('.popup_body');
        var return_delete_lesson_info = popupBody.querySelector('.return_delete_lesson_info');
        var contentRemover = popupBody.querySelector('.content-remover');
        var returnBtn = popupBody.querySelector('.return-btn');
        const contentId = contentRemover.getAttribute('content-id');
        const uniqstring = contentRemover.getAttribute('uniqstring');
        returnBtn.classList.add('hidden');
        return_delete_lesson_info.innerHTML = "Обрабатываю запрос...";
        returnLesson(contentId, uniqstring)
            .then(message => {
                return_delete_lesson_info.innerHTML = message;
            });
    }
});


function deleteLesson(lessonId, uniqstring) {
    return fetch('/web/delete_lesson' + window.location.search + '&lesson_id=' + lessonId + '&uniqstring=' + uniqstring, {
        method: 'POST'
    })
    .then(response => {
        if (!response.ok) {
            if (response.status == 403) {
                throw new Error('Для продолжения авторизуйтесь по студенческому билету');
            }
            throw new Error('Error deleting lesson');
        }
        console.log('Lesson deleted successfully');

        return 'Занятие помечено как удаленное для вашей группы. Изменения отобразятся при следующем запуске'
    })
        .catch(error => {
            console.error(error);
            return error.message;
        });
}

function returnLesson(lessonId, uniqstring) {
    return fetch('/web/return_lesson' + window.location.search + '&lesson_id=' + lessonId + '&uniqstring=' + uniqstring, {
        method: 'POST'
    })
        .then(response => {
            if (!response.ok) {
                if (response.status == 403) {
                    throw new Error('Для продолжения авторизуйтесь по студенческому билету');
                }
                throw new Error('Не вернуть занятие. Вероятно, занятие уже возвращено.');
            }
            console.log('Lesson returning successfully');

            return 'Занятие возвращено в расписание для вашей группы. Изменения отобразятся при следующем запуске'
        })
        .catch(error => {
            console.error(error);
            return error.message;
        });
}