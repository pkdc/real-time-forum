const userListSocket = new WebSocket("ws://localhost:8080/userListWs/")

document.addEventListener("DOMContentLoaded", function(e) {
    // userListSocket = new WebSocket("ws://localhost:8080/userListWs/")
    console.log("JS attempt to connect to user list");
    userListSocket.onopen = () => console.log("user list connected");
    userListSocket.onclose = () => console.log("Bye user list");
    userListSocket.onerror = (err) => console.log("user list ws Error!");
    userListSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        console.log({resp});
        if (resp.label === "update") {
            console.log(resp.online_users);
            const uList = document.querySelector(".user-list");
            // remove list item
            uList.textContent = "";
            // add new list item
            for (const {nickname, status} of resp.online_users) {
                const nicknameItem = document.createElement("li");
                nicknameItem.textContent = `${nickname} ${status}`;
                uList.append(nicknameItem);
            }
        }
    }
})

// const chatBox = document.createElement("form");
// chatBox.id = "chat-form"

export default userListSocket;