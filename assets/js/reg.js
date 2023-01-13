import userListSocket from "./userList.js";
import { chatSocket } from "./chat.js";
import { createProfile, updateChat } from "./login.js";
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
document.addEventListener("DOMContentLoaded", function () {
    regSocket = new WebSocket("ws://localhost:8080/regWs/");
    console.log("JS attempt to connect to reg");
    regSocket.onopen = () => console.log("connected-reg");
    regSocket.onclose = () => console.log("Bye-reg");
    regSocket.onerror = (err) => console.log("Error!-reg", err);
    regSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);

        if (resp.label === "Greet") {
            navbar.children[0].style.display = "block"
            navbar.children[1].style.display = "block"
            navbar.children[2].style.display = "none"
        } else if (resp.label === "reg") {
            if (resp.pass) {
                //create profile
                let user = JSON.parse(resp.content)

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
                const logPopup = document.querySelector("#userPopUpPOne")
                logPopup.style.display = "block"

                const signInName = document.querySelector("#signInName")
                signInName.value = user.nickname

                // update user list after a user reg
                // let uListPayload = {};
                // uListPayload["label"] = "login-reg-update";
                // // uListPayload["cookie_value"] = resp.cookie.sid;
                // console.log("reg UL sending: ", uListPayload);
                // userListSocket.send(JSON.stringify(uListPayload));

                // user is online and avalible to chat
                // let chatPayloadObj = {};
                // chatPayloadObj["label"] = "user-online";
                // console.log(`reg chat uid: ${resp.cookie.uid}`);
                // chatPayloadObj["sender_id"] = (resp.cookie.uid);
                // console.log("reg chat: ", chatPayloadObj);
                // chatSocket.send(JSON.stringify(chatPayloadObj));
            } else {
                displayMsgDiv.classList.add("display-msg");
                displayMsg.id = "reg-msg";
                displayMsg.textContent = `${resp.content}`;
                displayMsgDiv.append(displayMsg);
            }
        }
    }
});

const regHandler = function (e) {
    e.preventDefault();
    const formFields = new FormData(e.target);
    const payloadObj = Object.fromEntries(formFields.entries());
    payloadObj["label"] = "reg";
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
RnameInput.setAttribute("required", "");
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
RLastnameInput.setAttribute("required", "");
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
RNicknameInput.setAttribute("required", "");
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
RAgeInput.setAttribute("required", "");
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
REmailInput.setAttribute("required", "");
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
RpwInput.setAttribute("required", "");
RpwInputDiv.append(RpwInput);

//gender
const GenderTitle = document.createElement('div');
GenderTitle.className = "GenderTitle"
GenderTitle.textContent = "Please Select Your Gender"
// const RnameLabel = document.createElement('label');
// RnameLabel.textContent = "Select you Gender:";
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
RgenderDiv.append(GenderOpt1, GenderOpt2, GenderOpt3, GenderOpt4)

const pictureTitle = document.createElement("div");
pictureTitle.className = "picturetitle"
pictureTitle.textContent = "Please Select A Profil Picture"
const selectforPp = document.createElement("select");
selectforPp.setAttribute("id", "selectPP");
selectforPp.setAttribute("name", "pp_option");
const PreviewPPDiv = document.createElement("div")
const PreviewPP = document.createElement("img")
PreviewPP.className = "userProfilImageReg"
PreviewPP.src = "./assets/images/0.png"
PreviewPPDiv.append(PreviewPP)
selectforPp.onchange = function () {
    let sourc = selectforPp.options[selectforPp.selectedIndex].value
    let nmb = sourc.replace("option ", "")
    PreviewPP.src = "./assets/images/" + nmb + ".png"
}

for (let i = 0; i < 14; i++) {
    let opt = document.createElement("option");
    opt.setAttribute("name", "pp_option");
    opt.setAttribute("value", i);
    opt.textContent = "option " + i;
    selectforPp.appendChild(opt);

}


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
    GenderTitle,
    RgenderDiv,
    pictureTitle,
    selectforPp,
    PreviewPPDiv,
    regSubmitDiv);
export default RegisterForm;