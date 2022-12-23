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
            console.log(`chat room created between ${resp.sender_id} and ${resp.receiver_id}`);
        } else if (resp.label === "msgIncoming") {
            console.log("recievedChatMsg", resp)
            let msgrow = document.createElement("div")
            let msgtext = document.createElement("p")
            let timeOfMsg = document.createElement("p")
            timeOfMsg.classList = "timeofmsg"
            timeOfMsg.textContent=  timenow()
            timeOfMsg.style.fontSize = "9px"
            msgrow.className = "msg-row"
            msgtext.className = "msg-text"
            msgtext.textContent = resp.content
            msgrow.append(msgtext)
            msgrow.append(timeOfMsg)
            msgArea.append(msgrow)
            targetUserId = resp.contactID
            //     var userlist = document.querySelector(".user-list")
            //    var targetUser = document.querySelector(`#li${resp.contactID}`)
            //     console.log("targetuser:", targetUser)
            //     userlist.insertBefore(targetUser, userlist.firstChild)


        }

    }
})
export function timenow() {
    const monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun",
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    // console.log(getMonth(date));
    var currentdate = new Date();
    var datetime =
        monthNames[currentdate.getMonth()] + " "
        + currentdate.getDate() + " "
        + currentdate.getHours() + ":"
        + currentdate.getMinutes() + ":"
        + currentdate.getSeconds();
    return datetime
}