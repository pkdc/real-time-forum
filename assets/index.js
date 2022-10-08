import RegisterForm from "./reg.js";
import loginForm from "./login.js";

const formDiv = document.querySelector(".form-div");
const loginBut = document.querySelector(".login-btn");
const regBut = document.querySelector(".reg-btn");

// formDiv.append(loginForm);
regBut.append(RegisterForm)
loginBut.append(loginForm)

document.addEventListener("DOMContentLoaded", function() {
    formDiv.classList.add("show");
});
