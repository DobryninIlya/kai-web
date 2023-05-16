const select = document.getElementById("role_select");
const result = document.getElementById("result");
const inputGroup = document.getElementById("input_group");
const submitButton = document.querySelector(".submit_button");
const errorText = document.querySelector(".error_form");

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
        vk_id: params.get('vk_user_id')
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

