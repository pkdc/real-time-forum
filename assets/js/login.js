import userListSocket from "./userList.js";
import { chatSocket } from "./chat.js";
let loginSocket = null;
let nameInput = null;
let pwInput = null;
const navbar = document.querySelector(".navbar")
const displayMsgDiv = document.createElement("div");
const displayMsg = document.createElement("h2");
// const logout = document.querySelector("#logout")
// console.log(userListSocket);
document.addEventListener("DOMContentLoaded", function() {
    loginSocket = new WebSocket("ws://localhost:8080/loginWs/");
    console.log("JS attempt to connect to login");
    loginSocket.onopen = () => console.log("login connected");
    loginSocket.onclose = () => console.log("Bye login");
    loginSocket.onerror = (err) => console.log("login ws Error!");
    loginSocket.onmessage = (msg) => {
        // display msg
        const resp = JSON.parse(msg.data);

        if (resp.label === "greet") {
            console.log(resp.content);
            navbar.children[0].style.display = "block"
            navbar.children[1].style.display = "block"
            navbar.children[2].style.display = "none"
        } else if (resp.label === "login") {
            console.log("uid: ",resp.cookie.uid, "sid: ", resp.cookie.sid, "age: ", resp.cookie.max_age);
            document.cookie = `session=${resp.cookie.sid}; max-age=${resp.cookie.max_age}`;

            // update user list after a user login

            if (resp.pass) {
                const splitScreen = document.querySelector(".container")
                const signPage = document.querySelector("#userPopUpPOne")
                signPage.style.display= "none"
                splitScreen.style.display= "flex"

                // hide the login and reg btn, show the logout btn
                navbar.children[0].style.display = "none"
                navbar.children[1].style.display = "none"
                navbar.children[2].style.display = "block"

                // clear input fields
                nameInput.value = "";
                pwInput.value = "";

                // close the popup
                const loginPopup = document.querySelector("#userPopUpPOne");
                loginPopup.style.display = "none";

                // update user list after a user login
                let uListPayload = {};
                uListPayload["label"] = "login-reg-update";
                uListPayload["cookie_value"] = resp.cookie.sid;
                console.log("login UL sending: ", uListPayload);
                userListSocket.send(JSON.stringify(uListPayload));

                // user is online and avalible to chat
                let chatPayloadObj = {};
                chatPayloadObj["label"] = "user-online";
                console.log(`login chat uid: ${resp.cookie.uid}`);
                chatPayloadObj["sender_id"] = (resp.cookie.uid).toString();
                console.log("login chat: ", chatPayloadObj);
                chatSocket.send(JSON.stringify(chatPayloadObj));
            } else {
                // error msg
                displayMsgDiv.classList.add("display-msg");
                displayMsg.id = "login-msg";
                displayMsg.textContent = `${resp.content}`;
                displayMsgDiv.append(displayMsg);
            }
        }
    }
});

// logout.addEventListener("click", logouthandler)
const loginHandler = function (e) {
    e.preventDefault();
    const formFields = new FormData(e.target);
    const payloadObj = Object.fromEntries(formFields.entries());
    payloadObj["label"] = "login";
    console.log({ payloadObj });
    loginSocket.send(JSON.stringify(payloadObj));

    displayMsg.textContent = "";
};


const loginForm = document.createElement("form");
loginForm.className = "formPage"
loginForm.addEventListener("submit", loginHandler);

// login form
// name label
const nameLabelDiv = document.createElement('div');
const nameLabel = document.createElement('label');
nameLabel.textContent = "Please Enter Your Nickname or Email:";
nameLabel.setAttribute("for", "name");
nameLabelDiv.append(nameLabel);
// name input
const nameInputDiv = document.createElement('div');
nameInput = document.createElement('input');
nameInput.setAttribute("type", "text");
nameInput.setAttribute("name", "name");
nameInput.setAttribute("id", "name");
nameInput.setAttribute("placeholder", "eg: deathstar123 or abc@def.com")
nameInputDiv.append(nameInput);

// pw label
const pwLabelDiv = document.createElement('div');
const pwLabel = document.createElement('label');
pwLabel.textContent = "Please Enter Your Password:";
pwLabel.setAttribute("for", "pw");
pwLabelDiv.append(pwLabel);
// password input
const pwInputDiv = document.createElement('div');
pwInput = document.createElement('input');
pwInput.setAttribute("type", "password");
pwInput.setAttribute("name", "pw");
pwInput.setAttribute("id", "pw");
pwInputDiv.append(pwInput);

const loginSubmitDiv = document.createElement('div');
const loginSubmit = document.createElement("button");
loginSubmit.textContent = "Login";
loginSubmit.setAttribute("type", "submit");
loginSubmitDiv.append(loginSubmit);

loginForm.append(displayMsgDiv, nameLabelDiv, nameInputDiv, pwLabelDiv, pwInputDiv, loginSubmitDiv);
export default loginForm;