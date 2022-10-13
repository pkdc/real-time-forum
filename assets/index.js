// import RegisterForm from "./reg.js";
import loginForm from "./login.js";

const formDiv = document.querySelector(".form-div");

formDiv.append(loginForm);

document.addEventListener("DOMContentLoaded", function() {
    formDiv.classList.add("show");
});

let socket = null; 
document.addEventListener("DOMContentLoaded", function() {
    socket = new WebSocket("ws://localhost:8080/ws/");
    console.log("JS attempt to connect");

    socket.onopen = () => console.log("connected");
    socket.onclose = () => console.log("Bye");
    socket.onerror = (err) => console.log("Error!");
    socket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        console.log({resp});
        if (resp.label === "Greet") {
            console.log(resp.content);
        }
    }

    // can send smth back
    // const loginForm = document.querySelector("#login-form");
    loginForm.addEventListener("submit", function(e) {
        e.preventDefault();
        const formFields = new FormData(e.target);
        const payloadObj = Object.fromEntries(formFields.entries());
        payloadObj["label"] = "login";
        console.log({payloadObj});
        socket.send(JSON.stringify(payloadObj));
    })
});