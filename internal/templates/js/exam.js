urlParams = new URLSearchParams(window.location.search);
var paramsExam = "?"
for (const [key, value] of urlParams.entries()) {
    paramsExam = paramsExam + key + "=" + value + "&"
}
paramsExam= paramsExam.slice(0, -1)
var exams = ""
menu_exam.addEventListener("click",  function () {
    if (exams == "") {
        fetch(`/web/exam${paramsExam}`)
            .then(response => response.text())
            .then(html => {
                exams = html
                scheduleExamBlock.innerHTML = exams
            })
    }
})
