let chatSocket = null;
const uList = document.querySelector(".user-list");

document.addEventListener("DOMContentLoaded", function(e) {
    chatSocket = WebSocket("ws://localhost:8080/chatWs/")
    console.log("JS attempt to connect to chat");
    chatSocket.onopen = () => console.log("chat connected");
    chatSocket.onclose = () => console.log("Bye chat");
    chatSocket.onerror = (err) => console.log("chat ws Error!");
    chatSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        console.log({resp});
        if (resp.label === "userList") {

            console.log(resp.content);
            // remove list item

            // add new list item
        } else if (resp.label === "chat") {
            console.log(resp.content);
        }
    }
})