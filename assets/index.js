// import RegisterForm from "./reg.js";
import loginForm from "./login.js";

const formDiv = document.querySelector(".form-div");

formDiv.append(loginForm);

document.addEventListener("DOMContentLoaded", function() {
    formDiv.classList.add("show");
});