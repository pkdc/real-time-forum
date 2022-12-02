import userListSocket from "./userList.js";
import { chatSocket } from "./chat.js";
// console.log(userListSocket);
let regSocket = null; 
const userList = document.querySelector(".user-list");
const splitScreen = document.querySelector(".container")
const navbar = document.querySelector(".navbar")
const displayMsgDiv = document.createElement("div");
const displayMsg = document.createElement("h2");

let RnameInput = null;
let RLastnameInput = null;
let RNicknameInput = null;
let RAgeInput = null;
let REmailInput = null;
let RpwInput = null;
let GenderOpt1 = null;
let GenderOpt2 = null;
let GenderOpt3 = null;
let GenderOpt4 = null;


document.addEventListener("DOMContentLoaded", function() {
    regSocket = new WebSocket("ws://localhost:8080/regWs/");
    console.log("JS attempt to connect to reg");
    regSocket.onopen = () => console.log("connected-reg");
    regSocket.onclose = () => console.log("Bye-reg");
    regSocket.onerror = (err) => console.log("Error!-reg",err);
    regSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        console.log({resp});

        if (resp.label === "Greet") {
            console.log(resp.content);
            navbar.children[0].style.display = "block"
            navbar.children[1].style.display = "block"
            navbar.children[2].style.display = "none"
        } else if (resp.label === "reg") {
            console.log("uid: ",resp.cookie.uid, "sid: ", resp.cookie.sid, "age: ", resp.cookie.max_age);
            document.cookie = `session=${resp.cookie.sid}; max-age=${resp.cookie.max_age}`;
            navbar.children[0].style.display = "none"
            navbar.children[1].style.display = "none"
            navbar.children[2].style.display = "block"
            const signPage = document.querySelector("#userPopUpPOne")
            signPage.style.display= "none"
            splitScreen.style.display= "flex"
            console.log("msg: ", resp.content);

            if (resp.pass) {
                // hide the login and reg btn, show the logout btn
                navbar.children[0].style.display = "none"
                navbar.children[1].style.display = "none"
                navbar.children[2].style.display = "block"

                // clear input fields
                RnameInput.value = "";
                RLastnameInput.value = "";
                RNicknameInput.value = "";
                RAgeInput.value = "";
                REmailInput.value = "";
                RpwInput.value = "";
                 GenderOpt1.checked = 0;
                 GenderOpt2.checked = 0;
                 GenderOpt3.checked = 0;
                 GenderOpt4.checked = 0;

                // close the popup
                const regPopup = document.querySelector("#userPopUpPTwo");
                regPopup.style.display = "none";
                
                // update user list after a user reg
                let uListPayload = {};
                uListPayload["label"] = "login-reg-update";
                uListPayload["cookie_value"] = resp.cookie.sid;
                console.log("reg UL sending: ", uListPayload);
                userListSocket.send(JSON.stringify(uListPayload));

                // store chatWs after a new-user reg
                let chatPayload = {};
                chatPayload["label"] = "new-user";
                chatPayload["sender_id"] = resp.cookie.uid;
                console.log("reg chat: ", chatPayload);
                chatSocket.send(JSON.stringify(chatPayload));
            } else {
                displayMsgDiv.classList.add("display-msg");
                displayMsg.id = "reg-msg";
                displayMsg.textContent = `${resp.content}`;
                displayMsgDiv.append(displayMsg);
            }
        }
    }
});

const regHandler = function(e) {
    e.preventDefault();
    const formFields = new FormData(e.target);
    const payloadObj = Object.fromEntries(formFields.entries());
    payloadObj["label"] = "reg";
    console.log({payloadObj});
    regSocket.send(JSON.stringify(payloadObj));

    displayMsg.textContent = "";
};

// reg form//
const RegisterForm = document.createElement("form");
RegisterForm.className = "formPage"
RegisterForm.addEventListener("submit", regHandler);

// name label
const RnameLabelDiv = document.createElement('div');
const RnameLabel = document.createElement('label');
RnameLabel.textContent = "Please Enter Your First Name :";
RnameLabel.setAttribute("for", "name");
RnameLabelDiv.append(RnameLabel);
// name input
const RnameInputDiv = document.createElement('div');
RnameInput = document.createElement('input');
RnameInput.setAttribute("type", "text");
RnameInput.setAttribute("name", "name");
RnameInput.setAttribute("id", "name");
RnameInput.setAttribute("placeholder", "eg: Nick");
RnameInputDiv.append(RnameInput);
//last name label
const RLastnameLabelDiv = document.createElement('div');
const RLastnameLabel = document.createElement('label');
RLastnameLabel.textContent = "Please Enter Your Last Name :";
RLastnameLabel.setAttribute("for", "lastname");
RLastnameLabelDiv.append(RLastnameLabel);
// last name input
const RLastnameInputDiv = document.createElement('div');
RLastnameInput = document.createElement('input');
RLastnameInput.setAttribute("type", "text");
RLastnameInput.setAttribute("name", "lastname");
RLastnameInput.setAttribute("id", "lastname");
RLastnameInput.setAttribute("placeholder", "eg: Smith");
RLastnameInputDiv.append(RLastnameInput);
// Nickname label
const RNicknameLabelDiv = document.createElement('div');
const RNicknameLabel = document.createElement('label');
RNicknameLabel.textContent = "Please Enter Your Nickname :";
RNicknameLabel.setAttribute("for", "nickname");
RNicknameLabelDiv.append(RNicknameLabel);
// nickname input
const RNicknameInputDiv = document.createElement('div');
RNicknameInput = document.createElement('input');
RNicknameInput.setAttribute("type", "text");
RNicknameInput.setAttribute("name", "nickname");
RNicknameInput.setAttribute("id", "nickname");
RNicknameInput.setAttribute("placeholder", "eg:deathstar123 ");
RNicknameInputDiv.append(RNicknameInput);
//  Age label
const RAgeLabelDiv = document.createElement('div');
const RAgeLabel = document.createElement('label');
RAgeLabel.textContent = "Please Enter Your Date of Birth :";
RAgeLabel.setAttribute("for", "age");
RAgeLabelDiv.append(RAgeLabel);
// age input
const RAgeInputDiv = document.createElement('div');
RAgeInput = document.createElement('input');
RAgeInput.setAttribute("type", "date");
RAgeInput.setAttribute("name", "age");
RAgeInput.setAttribute("id", "age");
RAgeInputDiv.append(RAgeInput);
//  Gender label

//  E-mail label
const REmailLabelDiv = document.createElement('div');
const REmailLabel = document.createElement('label');
REmailLabel.textContent = "Please Enter Your e-mail :";
REmailLabel.setAttribute("for", "email");
REmailLabelDiv.append(REmailLabel);
// email input
const REmailInputDiv = document.createElement('div');
REmailInput = document.createElement('input');
REmailInput.setAttribute("type", "email");
REmailInput.setAttribute("name", "email");
REmailInput.setAttribute("id", "email");
REmailInput.setAttribute("placeholder", "eg: deathstar@123.com");
REmailInputDiv.append(REmailInput);
// pw label
const RpwLabelDiv = document.createElement('div');
const RpwLabel = document.createElement('label');
RpwLabel.textContent = "Please Enter Your Password:";
RpwLabel.setAttribute("for", "pw");
RpwLabelDiv.append(RpwLabel);
// password input
const RpwInputDiv = document.createElement('div');
RpwInput = document.createElement('input');
RpwInput.setAttribute("type", "password");
RpwInput.setAttribute("name", "pw");
RpwInput.setAttribute("id", "pw");
RpwInputDiv.append(RpwInput);

//gender
const RgenderDiv = document.createElement('select');
RgenderDiv.setAttribute("name", "gender_option")
GenderOpt1 = document.createElement("option");
GenderOpt2 = document.createElement("option");
GenderOpt3 = document.createElement("option");
GenderOpt4 = document.createElement("option");
GenderOpt1.setAttribute("name", "gender_option");
GenderOpt2.setAttribute("name", "gender_option");
GenderOpt3.setAttribute("name", "gender_option");
GenderOpt4.setAttribute("name", "gender_option");
GenderOpt1.setAttribute("value", "Prefer not");
GenderOpt2.setAttribute("value", "Female");
GenderOpt3.setAttribute("value", "Male");
GenderOpt4.setAttribute("value", "Other");
GenderOpt1.textContent = "Prefer not";
GenderOpt2.textContent = "Female";
GenderOpt3.textContent = "Male";
GenderOpt4.textContent = "Other";
RgenderDiv.setAttribute("id", "genderOption");
RgenderDiv.append(GenderOpt1,GenderOpt2,GenderOpt3,GenderOpt4)

// -----------------------
// const RgenderDiv = document.createElement('div');
// const RgenderOptionDiv = document.createElement('div');
// const RgenderLabel = document.createElement("label");
// RgenderLabel.textContent= "Please choose your gender";
// RgenderLabel.setAttribute("for","gender");
// RgenderDiv.append(RgenderLabel);
// RgenderInputOpt1= document.createElement("input");
// RgenderInputOpt2= document.createElement("input");
// RgenderInputOpt3= document.createElement("input");
// RgenderInputOpt4= document.createElement("input");
// const RgenderLabelOpt1= document.createElement("label");
// const RgenderLabelOpt2= document.createElement("label");
// const RgenderLabelOpt3= document.createElement("label");
// const RgenderLabelOpt4= document.createElement("label");
// RgenderInputOpt1.setAttribute("type","radio");
// RgenderInputOpt2.setAttribute("type","radio");
// RgenderInputOpt3.setAttribute("type","radio");
// RgenderInputOpt4.setAttribute("type","radio");
// RgenderInputOpt1.setAttribute("name","gender_option");
// RgenderInputOpt2.setAttribute("name","gender_option");
// RgenderInputOpt3.setAttribute("name","gender_option");
// RgenderInputOpt4.setAttribute("name","gender_option");
// RgenderInputOpt1.setAttribute("id","male");
// RgenderInputOpt2.setAttribute("id","female");
// RgenderInputOpt3.setAttribute("id","other");
// RgenderInputOpt4.setAttribute("id","prefernot");
// RgenderInputOpt1.setAttribute("value","male");
// RgenderInputOpt2.setAttribute("value","female");
// RgenderInputOpt3.setAttribute("value","other");
// RgenderInputOpt4.setAttribute("value","prefernot");
// RgenderLabelOpt1.setAttribute("for","male");
// RgenderLabelOpt2.setAttribute("for","female");
// RgenderLabelOpt3.setAttribute("for","other");
// RgenderLabelOpt4.setAttribute("for","prefernot");
// RgenderLabelOpt1.textContent= "Male";
// RgenderLabelOpt2.textContent= "Female";
// RgenderLabelOpt3.textContent= "Other";
// RgenderLabelOpt4.textContent= "Prefer not to say";
// RgenderOptionDiv.append(
//     RgenderInputOpt1,RgenderLabelOpt1,
//     RgenderInputOpt2,RgenderLabelOpt2,
//     RgenderInputOpt3,RgenderLabelOpt3,
//     RgenderInputOpt4,RgenderLabelOpt4);

// RgenderOptionDiv.setAttribute("id", "gender");
//----------------------
const regSubmitDiv = document.createElement('div');
const regSubmit = document.createElement("button");
regSubmit.textContent = "Register";
regSubmit.setAttribute("type", "submit");
regSubmitDiv.append(regSubmit);
//append
RegisterForm.append(
    displayMsgDiv,
    RnameLabelDiv,
    RnameInputDiv,
    RLastnameLabelDiv,
    RLastnameInputDiv,
    RNicknameLabelDiv,
    RNicknameInputDiv,
    RAgeLabelDiv,
    RAgeInputDiv,
    REmailLabelDiv,
    REmailInputDiv,
    RpwLabelDiv,
    RpwInputDiv,
    RgenderDiv,
    regSubmitDiv);
export default RegisterForm;