// import chatForm from "./chat.js";
import RegisterForm from "./reg.js";
import loginForm from "./login.js";
import Post from "./post.js";
import logoutBtn from "./logout.js";
import { chatForm } from "./chat.js";

const chatBox = document.querySelector(".col-1")
const loginArea = document.querySelector("#userPopUpPOne")
const loginInputs = loginArea.firstElementChild
const regArea = document.querySelector("#userPopUpPTwo")
const reginInputs = regArea.firstElementChild
const body = document.querySelector("body")
const logoutLi = document.querySelector("#logoutBtn")
const postPage = document.querySelector(".postPage")
logoutLi.appendChild(logoutBtn)
loginInputs.append(loginForm)
reginInputs.append(RegisterForm)
loginArea.append(loginInputs)
regArea.append(reginInputs)
postPage.append(Post.PostForm)
postPage.append(Post.DisplayPost)
chatBox.append(chatForm)
// const chatBox = document.createElement("div")
// chatBox.append(chatForm);
// body.append(chatBox);
// body.append(PostForm)
// body.append(LogoutButton)
// const chatBox = document.createElement("div")
// chatBox.append(chatForm);
// body.append(chatBox);
