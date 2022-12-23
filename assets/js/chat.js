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
            setTimeout(() => typingDiv.classList.remove("show"), 5000);
        }
    }
})
