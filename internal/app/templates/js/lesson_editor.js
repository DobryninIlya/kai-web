document.addEventListener('click', function(event) {
    // Проверяем, что кликнули по элементу с классом .remove-btn

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
        contentRemover.classList.add('hidden');
    }
});

