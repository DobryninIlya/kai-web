// auth_error_block = document.querySelector(".auth_error_block");
// let urlParams = new URLSearchParams(window.location.search);
// var scoreExam = "?"
// for (const [key, value] of urlParams.entries()) {
//     scoreExam = scoreExam + key + "=" + value + "&"
// }
// scoreExam= scoreExam.slice(0, -1)
var score = ""
menu_score.addEventListener("click", scoreListShow)


function scoreListShow() {
    if (score == "") {
        score_block.insertAdjacentHTML('beforeend', loaderHTML);
        fetch(`/web/scoretable?${urlParams}`)
            .then(response => {
                if (response.status == 404) {
                    auth_error_block.style.display = "flex"
                    return
                }
                response.text().then(html => {
                    response.ok
                    score = html
                    score_block.innerHTML = score
                    let examItems = document.querySelectorAll('.score_elem');


                    examItems.forEach((item) => {
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

                    })
                })

            })
    }
};


