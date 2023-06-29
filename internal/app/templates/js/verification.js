const selectName = document.querySelector('#name_select');
const selectFac = document.querySelector('#fac_select');
const selectGroup = document.querySelector('#group_select');
const selectKurs = document.querySelector('#kurs_select');
const selectZach = document.querySelector('#input_group');
const confirmBtn = document.querySelector('.submit_button');
const errorForm = document.querySelector('.error_form');

// selectFac.selectedIndex = -1;
// selectKurs.selectedIndex = -1;



window.addEventListener('load', () => {
    getFac();
});

function getFac() {
    selectKurs.selectedIndex = -1;
    fetch('/web/attestation/get_fac')
        .then(response => response.json())
        .then(data => {
            const options = data.result;
            for (const key in options) {
                const option = document.createElement('option');
                option.value = key;
                option.text = options[key];
                selectFac.appendChild(option);
            }
            selectFac.selectedIndex =-1;
        });
}
function getGroups(p_fac, p_kurs) {
    selectGroup.innerHTML = ''; // очищаем все значения option
    fetch(`/web/attestation/get_groups?p_fac=${p_fac}&p_kurs=${p_kurs}`)
        .then(response => response.json())
        .then(data => {
            const options = data.result;
            for (const key in options) {
                const option = document.createElement('option');
                option.text = key;
                option.value = options[key];
                selectGroup.appendChild(option);
            }
            selectGroup.selectedIndex = -1;
        });
}

function getPerson(p_fac, p_kurs, p_group) {
    selectName.innerHTML = ''; // очищаем все значения option
    fetch(`/web/attestation/get_person?p_fac=${p_fac}&p_kurs=${p_kurs}&p_group=${p_group}`, { method: 'GET' })
        .then(response => response.json())
        .then(data => {
            const options = data.result;
            for (const key in options) {
                const option = document.createElement('option');
                option.value = key;
                option.text = options[key];
                selectName.appendChild(option);
            }
            selectName.selectedIndex = -1;
        });
}

selectFac.addEventListener('change', () => {
    selectName.innerHTML = ""
    selectGroup.innerHTML = ""
    selectKurs.selectedIndex = -1;
    selectGroup.selectedIndex = -1;
    // selectFac.innerHTML = ""
    // getFac()
});

selectKurs.addEventListener('change', () => {
    selectName.innerHTML = ""
    const selectedFac = selectFac.options[selectFac.selectedIndex];
    const valueFac = selectedFac.value;
    const selectedKurs = selectKurs.options[selectKurs.selectedIndex];
    const valueKurs = selectedKurs.value;
    getGroups(valueFac, valueKurs)

});

selectGroup.addEventListener('change', () => {
    const selectedFac = selectFac.options[selectFac.selectedIndex];
    const valueFac = selectedFac.value;
    const selectedKurs = selectKurs.options[selectKurs.selectedIndex];
    const valueKurs = selectedKurs.value;
    const selectedGroup = selectGroup.options[selectGroup.selectedIndex];
    const valueGroup = selectedGroup.value;
        getPerson(valueFac, valueKurs, valueGroup)
});


async function sendVerificationData()  {
    const data = {
        faculty: parseInt(selectFac.options[selectFac.selectedIndex].value),
        course: parseInt(selectKurs.options[selectKurs.selectedIndex].value),
        group: parseInt(selectGroup.options[selectGroup.selectedIndex].value),
        student: parseInt(selectName.options[selectName.selectedIndex].value),
        id: parseInt(selectZach.value),
        groupname: parseInt(selectGroup.options[selectGroup.selectedIndex].text)
    };

    // fetch("/web/verification/done" + window.location.search, {
    //     method: "POST",
    //     headers: {
    //         "Content-Type": "application/json"
    //     },
    //     body: JSON.stringify(data)
    // })
    //     .then(response => {
    //
    //         if (response.ok) {
    //             console.log("Verification data sent successfully!");
    //             return true
    //         } else if (response.status == 404) {
    //             errorForm.innerHTML = "Номер зачетной книжки неверный!"
    //         }
    //         else {
    //             console.error("Error sending verification data!");
    //         }
    //         return false
    //     })
    //     .catch(error => {
    //         console.error("Error sending verification data:", error);
    //     });

    try {
        const response = await fetch("/web/verification/done" + window.location.search, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });
        if (response.ok) {
            console.log("Verification data sent successfully!");
            return true
        } else if (response.status == 404) {
            errorForm.innerHTML = "Номер зачетной книжки неверный!"
        }
        else {
            console.error("Error sending verification data!");
        }
       return false
    } catch (error) {
        console.error(error);
        return false; // возвращаем false в случае ошибки
    }

    return false
}

confirmBtn.addEventListener("click", async function() {
    const result = await sendVerificationData();
    console.log(result)
    if (result) {
        location.reload()
    }

})