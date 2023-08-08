const menu_exam = document.getElementById("menu_exam");
const menu_score = document.getElementById("menu_score");
const menu_schedule = document.getElementById("menu_schedule");
const menu_teachers = document.getElementById("menu_teachers");
const menu_services = document.getElementById("menu_services");

const mail_to_developer = document.getElementById("mail_to_developer");
const mail_to_improvement = document.getElementById("mail_to_improvement");

const schedule_block = document.querySelector(".schedule_block");
const main_block = document.querySelector(".main_block");
const schedule_exam_block = document.querySelector(".schedule_exam_block");
const score_block = document.querySelector(".score_block");
const teacher_block = document.querySelector(".teacher_block");
const services_block = document.querySelector(".services_block");
const auth_error_block = document.querySelector(".auth_error_block");

const action = document.getElementById("action_id");
const address = "https://kai.ru/raspisanie?p_p_id=pubStudentSchedule_WAR_publicStudentSchedule10&p_p_lifecycle=2&p_p_resource_id=examSchedule&groupId="
function hideAll() {
    schedule_block.style.display = "none";
    main_block.style.opacity = "0"
    schedule_exam_block.style.display = "none";
    score_block.style.display = "none"
    teacher_block.style.display = "none"
    services_block.style.display = "none"
    auth_error_block.style.display = "none"
}


menu_exam.addEventListener("click",  function () {
    hideAll();
    main_block.style.opacity = "1"
    schedule_exam_block.style.display = "block";
})

menu_schedule.addEventListener("click",  function () {
    hideAll();
    main_block.style.opacity = "1"
    schedule_block.style.display = "block"

})

menu_score.addEventListener("click",  function () {
    hideAll();
    main_block.style.opacity = "1"
    score_block.style.display = "block"

})

menu_teachers.addEventListener("click",  function () {
    hideAll();
    main_block.style.opacity = "1"
    teacher_block.style.display = "block"

})

menu_services.addEventListener("click",  function () {
    hideAll();
    main_block.style.opacity = "1"
    services_block.style.display = "block"

})

mail_to_improvement.addEventListener("click", function (event) {
    event.preventDefault();
    window.location.href = "https://m.vk.com/topic-182372147_43301544";;
})
mail_to_developer.addEventListener("click", function (event) {
    event.preventDefault();
    window.location.href = "https://t.me/dobryninilya";
})
