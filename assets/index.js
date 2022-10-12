// import RegisterForm from "./reg.js";
import loginForm from "./login.js";

const formDiv = document.querySelector(".form-div");

formDiv.append(loginForm);

document.addEventListener("DOMContentLoaded", function() {
    formDiv.classList.add("show");brew
});

let socket = null; 
document.addEventListener("DOMContentLoaded", function() {
    socket = new WebSocket("ws://localhost:8080/ws");
    console.log("JS attempt to connect");

    socket.onopen = () => console.log("connected");
    socket.onclose = () => console.log("Bye");
    socket.onerror = (err) => console.log("Error!");
    socket.onmessage = (msg) => {
        console.log(msg);
        console.log(JSON.parse(msg));
        const resp = JSON.parse(msg);
        console.log(resp.Label);
        if (resp.Label === "login")
    }
});