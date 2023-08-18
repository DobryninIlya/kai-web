const select = document.getElementById("role_select");
const result = document.getElementById("result");
const inputGroup = document.getElementById("input_group");
const submitButton = document.querySelector(".submit_button");
const errorText = document.querySelector(".error_form");
const registrationButton = document.querySelector(".registration_button");
const authButton = document.querySelector(".auth_button");

const currentUrl = window.location.href;
const hash = new URL(currentUrl).hash; // или currentUrl.substring(currentUrl.indexOf("#"))
console.log(hash);




select.addEventListener("change", function() {
    const selectedOption = select.options[select.selectedIndex];
    var placeholder = "Введите группу"
    if (select.selectedIndex == 2) {
        placeholder = "Введите логин bb"
    }
    const newDiv = "<input type=\"number\" class=\"input_group\" id=\"input_group\" placeholder=\"" + placeholder + "\">"
    result.innerHTML = newDiv
});

submitButton.addEventListener("click",  (event) => {
   if (select.selectedIndex == 2 || select.selectedIndex == 1) {
       errorText.textContent = "Временно недоступна регистрация преподавателей и родителей"
       return
   }
    makeRegistration()

});

function makeRegistration() {
    var params = new URLSearchParams(window.location.search);
    const data = {
        role : select.selectedIndex,
        identificator : inputGroup.value,
        vk_id: parseInt(params.get('vk_user_id'))
    }

    fetch("/web/registration"+window.location.search, {
        method: "POST",
        headers: {
            "Content-Type" : "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(data => {
            try {
                if (data.hasOwnProperty('error')) {
                    errorText.textContent = data['error']
                    return
                }
            } catch (e) {
                console.log(e)
                errorText.textContent = "Произошла неизвестная ошибка."
                return
            }
            location.reload()

        })
        .catch(error => console.error(error));

}

window.onload = function () {
    const inputBox = document.getElementById("main");
    inputBox.style.display = "none"

}

registrationButton.addEventListener("click", () => {
    const inputBox = document.getElementById("main");
    inputBox.style.display = "block"
    const infoText = document.getElementById("info_text");
    infoText.style.display = "none"
})



authButton.addEventListener("click",  function(event) {
    event.preventDefault(); // отменяем стандартное действие кнопки
    fetch("/web/verification" + window.location.search)
        .then(response => response.text())
        .then(html => {
            document.body.innerHTML = html
            const script = document.createElement('script');
            script.src = '/static/js/verification.js';
            script.onload = () => {
                getFac() // вызываем функцию getFac() после загрузки скрипта
            };
            document.body.appendChild(script);
            getFac();
        })
        .catch(error => console.log(error));


})

window.onload = function () {
    if (hash=="#reset_group") {

        const inputBox = document.getElementById("main");
        inputBox.style.display = "block"
        const infoText = document.getElementById("info_text");
        infoText.style.display = "none"
    }
};
