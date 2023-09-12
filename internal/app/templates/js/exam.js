let urlParams = new URLSearchParams(window.location.search);
var paramsExam = "?"
for (const [key, value] of urlParams.entries()) {
    paramsExam = paramsExam + key + "=" + value + "&"
}
paramsExam= paramsExam.slice(0, -1)
var exams = ""
open_exam.addEventListener("click", examListShow)
menu_exam.addEventListener("click",  examListShow)
function examListShow () {
    if (exams == "") {
        schedule_exam_block.insertAdjacentHTML('beforeend', loaderHTML);
        fetch(`/web/exam${paramsExam}`)
            .then(response => response.text())
            .then(html => {
                exams = html
                schedule_exam_block.innerHTML = exams
            })
    }
};
