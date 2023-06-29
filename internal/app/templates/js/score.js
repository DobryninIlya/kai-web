import { showLoader, hideLoader } from './loader.js';
const auth_error_block = document.querySelector(".auth_error_block");
let urlParams = new URLSearchParams(window.location.search);
var scoreExam = "?"
for (const [key, value] of urlParams.entries()) {
    scoreExam = scoreExam + key + "=" + value + "&"
}
scoreExam= scoreExam.slice(0, -1)
var score = ""
menu_score.addEventListener("click",  function () {
    if (score == "") {
        showLoader();
        fetch(`/web/scoretable${scoreExam}`)
            .then(response => {
                if (response.status == 404) {
                    auth_error_block.style.display = "flex"
                    return
                }
                response.text() .then(html => {
                    response.ok
                    score = html
                    score_block.innerHTML = score
                })
            })
        hideLoader();
    }
})


