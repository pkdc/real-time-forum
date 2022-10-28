import RegisterForm from "./reg.js";
import loginForm from "./login.js";
import PostForm from "./post.js";


const formDiv = document.querySelector(".form-div");
const loginBut = document.querySelector(".login-btn");
const regBut = document.querySelector(".reg-btn");
const formPage = document.querySelector(".form-page")
const openBut = document.createElement("button")
const openButDiv = document.createElement("div")
const closeBut = document.createElement("button")
const closeButDiv = document.createElement("div")
const loginOrReg= document.querySelector(".login-or-reg")
const body= document.querySelector("body")
openButDiv.setAttribute("id", "openButDiv" )
openButDiv.append(openBut)
closeButDiv.append(closeBut)
regBut.append(RegisterForm)
loginBut.append(loginForm)
openBut.setAttribute("id", "openPageButton")
openBut.textContent= "Login / Register"
closeBut.textContent= String.fromCodePoint(0x274C)
body.append(openButDiv)
openBut.addEventListener("click", function(){
    console.log("openpge")
    formPage.style.height = "100%";
});
closeBut.addEventListener("click", function(){
    console.log("openpge")
    formPage.style.height = "0%";
});


loginOrReg.append(closeButDiv)
document.addEventListener("DOMContentLoaded", function() {
    formDiv.classList.add("show");
});
body.append(PostForm)
