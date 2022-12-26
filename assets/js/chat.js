export const chatSocket = new WebSocket("ws://localhost:8080/chatWs/");
export var targetUserId = null
export const chatForm = document.createElement("form");
const msgArea = document.querySelector(".msgArea")
chatForm.id = "chat-form";
chatForm.addEventListener("submit", function (e) {
    e.preventDefault();
    // add msg
    // send msg to ws
});
const chatInputDiv = document.createElement("div");
chatInputDiv.id = "chat-input-div";
const chatInput = document.createElement("input");
chatInputDiv.append(chatInput);

let typingLeftDotsArr, typingRightDotsArr;
let leftDot1Pos = 0;
let rightDot1Pos = 0;
let leftDot2Pos = 0;
let rightDot2Pos = 0;
let leftDot3Pos = 0;
let rightDot3Pos = 0;
// let dotSize = 50;

// const sendBtn = document.createElement("button");
// sendBtn.textContent = "Send";
// sendBtn.id = "send-btn";
// chatForm.append(chatInputDiv, sendBtn);


document.addEventListener("DOMContentLoaded", function (e) {

    // chatSocket = new WebSocket("ws://localhost:8080/chatWs/");
    console.log("JS attempt to connect to chat");
    chatSocket.onopen = () => console.log("chat connected");
    chatSocket.onclose = () => console.log("Bye chat");
    chatSocket.onerror = (err) => console.log("chat ws Error!");
    chatSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        if (resp.label === "created_room") {
            // console.log(`chat room created between ${resp.sender_id} and ${resp.receiver_id}`);
        } else if (resp.label === "msgIncoming") {
            console.log("recievedChatMsg", resp)
            let msgrow = document.createElement("div")
            let msgtext = document.createElement("p")
            msgrow.className = "msg-row"
            msgtext.className = "msg-text"
            msgtext.textContent = resp.content
            msgrow.append(msgtext)
            msgArea.append(msgrow)
            targetUserId= resp.contactID
        //     var userlist = document.querySelector(".user-list")
        //    var targetUser = document.querySelector(`#li${resp.contactID}`)
        //     console.log("targetuser:", targetUser)
        //     userlist.insertBefore(targetUser, userlist.firstChild)
        } else if (resp.label === "sender-typing") {
            // display typing-in-progress
            // const chatBox = document.querySelector(".chatbox");
            const typingDiv = document.querySelector(".typing-div");
            const typingText = document.querySelector(".typing-text");
            console.log(`${resp.sender} with id ${resp.userID} is typing... to ${resp.contactID}`);
            typingText.textContent = `${resp.sender} is typing`;
            // typingDiv.style.opacity = 1;
            // setTimeout(() => typingDiv.style.opacity = 0, 5000);
            typingDiv.classList.add("show");
            // typingDotsArr = [...document.querySelectorAll(".typing-dots")];
            typingLeftDotsArr = [...document.querySelectorAll(".typing-left-dots")];
            typingRightDotsArr = [...document.querySelectorAll(".typing-right-dots")];

            const animationID = requestAnimationFrame(animateDots);
            setTimeout(() => {
                typingDiv.classList.remove("show");
                cancelAnimationFrame(animationID);
            }, 5000);
        }
    }
})


export const genTypingDiv = function() {
    const typingDiv = document.createElement("div");
    const typingText = document.createElement("p");

    const typingLeftDotsDiv = document.createElement("div");
    typingLeftDotsDiv.classList.add("typing-left-dots-div");
    const typingLeftDot1 = document.createElement("p");
    const typingLeftDot2 = document.createElement("p");
    const typingLeftDot3 = document.createElement("p");
    typingLeftDot1.classList.add("typing-dots");
    typingLeftDot2.classList.add("typing-left-dots");
    typingLeftDot3.classList.add("typing-dots");
    typingLeftDot1.classList.add("typing-left-dots");
    typingLeftDot2.classList.add("typing-dots");
    typingLeftDot3.classList.add("typing-left-dots");
    typingLeftDotsDiv.append(typingLeftDot1, typingLeftDot2, typingLeftDot3);

    const typingRightDotsDiv = document.createElement("div");
    typingRightDotsDiv.classList.add("typing-right-dots-div");
    const typingRightDot1 = document.createElement("p");
    const typingRightDot2 = document.createElement("p");
    const typingRightDot3 = document.createElement("p");
    typingRightDot1.classList.add("typing-dots");
    typingRightDot2.classList.add("typing-right-dots");
    typingRightDot3.classList.add("typing-dots");
    typingRightDot1.classList.add("typing-right-dots");
    typingRightDot2.classList.add("typing-dots");
    typingRightDot3.classList.add("typing-right-dots");
    typingRightDotsDiv.append(typingRightDot1, typingRightDot2, typingRightDot3);

    typingDiv.classList.add("typing-div");
    typingText.classList.add("typing-text");
    typingLeftDot1.textContent = "·";
    typingLeftDot2.textContent = "·";
    typingLeftDot3.textContent = "·";
    typingRightDot1.textContent = "·";
    typingRightDot2.textContent = "·";
    typingRightDot3.textContent = "·";
    typingDiv.append(typingLeftDotsDiv, typingText, typingRightDotsDiv);
    
    return typingDiv;
}

const animateDots = function() { 
    // incSize += 0.01;
    leftDot1Pos -= 1;
    console.log("left:", leftDot1Pos);
    typingLeftDotsArr[0].style.left = `${leftDot1Pos}px`;
    // typingLeftDotsArr[0].style.fontWeight = `bold`;
    // typingLeftDotsArr[0].style.fontSize = `${dotSize+incSize}px`;
    // rightPos -= 1;
    // rightPos += 1;
    rightDot1Pos += 1;
    console.log("right:", rightDot1Pos);
    typingRightDotsArr[0].style.left = `${rightDot1Pos}px`;
    if (rightDot1Pos > 300 && leftDot1Pos < -300) {
        typingLeftDotsArr[0].remove();
        leftDot1Pos = 0;
        leftDot2Pos -= 1;
        console.log("left:", leftDot2Pos);
        typingLeftDotsArr[1].style.left = `${leftDot2Pos}px`;

        typingRightDotsArr[0].remove();
        rightDot1Pos = 0;
        rightDot2Pos += 1;
        console.log("right:", rightDot2Pos);
        typingRightDotsArr[1].style.left = `${rightDot2Pos}px`;
    }
    if (rightDot2Pos > 300 && leftDot2Pos < -300) {
        typingLeftDotsArr[1].remove();
        leftDot2Pos = 0;
        leftDot3Pos -= 1;
        console.log("left:", leftDot3Pos);
        typingLeftDotsArr[2].style.left = `${leftDot3Pos}px`;

        typingRightDotsArr[1].remove();
        rightDot2Pos = 0;
        rightDot3Pos += 1;
        console.log("right:", rightDot3Pos);
        typingRightDotsArr[2].style.left = `${rightDot3Pos}px`;
    }
    if (rightDot3Pos > 300 && leftDot3Pos < -300) {
        typingLeftDotsArr[2].remove();
        typingRightDotsArr[2].remove();
        const prevTypingDiv = document.querySelector(".typing-div");
        prevTypingDiv.remove();
        leftDot3Pos = 0;
        rightDot3Pos = 0;
        console.log("regen");
        const typingDiv = genTypingDiv();
        const chatBox = document.querySelector(".col-1")
        // chatBox.removeChild(chatBox.lastChild);
        chatBox.append(typingDiv);
    }
    // typingRightDotsArr[0].style.fontWeight = `bold`;
    // typingRightDotsArr[0].style.fontSize = `${dotSize+incSize}px`;
    requestAnimationFrame(animateDots);
}
