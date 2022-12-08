export const chatSocket = new WebSocket("ws://localhost:8080/chatWs/");

export const chatForm = document.createElement("form");
chatForm.id = "chat-form";
chatForm.addEventListener("submit", function(e) {
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


document.addEventListener("DOMContentLoaded", function(e) {
    // chatSocket = new WebSocket("ws://localhost:8080/chatWs/");
    console.log("JS attempt to connect to chat");
    chatSocket.onopen = () => console.log("chat connected");
    chatSocket.onclose = () => console.log("Bye chat");
    chatSocket.onerror = (err) => console.log("chat ws Error!");
    chatSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        console.log({resp});
        if (resp.label === "created_room") {
            console.log(`chat room created between ${resp.sender_id} and ${resp.receiver_id}`);
        } else if (resp.label === "chat") {
            console.log(resp.content);
        }
        //  else if (resp.label === "chat") {
         
        //     console.log(resp.content);
        // }
    }
})

// const chatBox = document.createElement("form");
// chatBox.id = "chat-form"

// export default chatForm;