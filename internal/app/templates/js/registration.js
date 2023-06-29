const select = document.getElementById("role_select");
const result = document.getElementById("result");
const inputGroup = document.getElementById("input_group");
const submitButton = document.querySelector(".submit_button");
const errorText = document.querySelector(".error_form");
const registrationButton = document.querySelector(".registration_button");
const authButton = document.querySelector(".auth_button");

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
   if (select.selectedIndex == 2) {
       errorText.textContent = "Временно недоступна регистрация преподавателей"
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

    fetch("/web/registration", {
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

    // fetch('/web/verification' + window.location.search)
    //     .then(response => response.text())
    //     .then(html => {
    //         const parser = new DOMParser();
    //         const doc = parser.parseFromString(html, 'text/html');
    //         const script = doc.querySelector('script');
    //         const code = script.innerHTML;
    //         const head = document.querySelector('head');
    //         const body = document.querySelector('body');
    //         head.innerHTML = doc.head.innerHTML;
    //         body.innerHTML = doc.body.innerHTML;
    //         const newScript = document.createElement('script');
    //         newScript.innerHTML = code;
    //         body.appendChild(newScript);
    //     });

})
