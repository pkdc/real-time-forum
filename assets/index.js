// import RegisterForm from "./reg.js";
import loginForm from "./login.js";
import logoutBtn from "./logout.js";

const formDiv = document.querySelector(".form-div");

formDiv.append(logoutBtn);
formDiv.append(loginForm);

document.addEventListener("DOMContentLoaded", function() {
    formDiv.classList.add("show");
});
