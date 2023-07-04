// let urlParams = new URLSearchParams(window.location.search);
var teacher = "?"
for (const [key, value] of urlParams.entries()) {
    teacher = teacher + key + "=" + value + "&"
}
teacher= teacher.slice(0, -1)
var teacherResult = ""
menu_teachers.addEventListener("click",  function () {
    if (teacherResult == "") {
        teacher_block.insertAdjacentHTML('beforeend', loaderHTML);
        fetch(`/web/teacher${teacher}`)
            .then(response => {
                if (!response.ok) {
                    console.log("Ошибка получения списка преподавателей")
                    return
                }
                response.text() .then(html => {
                    response.ok
                    teacher_block.innerHTML = html
                })
            })

    }
})
