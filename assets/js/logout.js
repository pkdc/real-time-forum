import userListSocket from "./userList.js";
import { chatSocket } from "./chat.js";
const logoutUrl = location.origin + "/logout/";
console.log(logoutUrl);
const logoutHandler = function() {
    fetch(logoutUrl)
    .then(() => {
        // get cookie val
        const sessionCookie = document.cookie.split(";").find((el) => el.startsWith("session="));
        if (sessionCookie) {
            console.log("session cookie: ", sessionCookie);
            const cookieVal = sessionCookie.split("=")[1];
            console.log("cookie val: ", cookieVal);
    
            // remove cookie
            document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 GMT";
    
            // update user list after a user logout
            let uListPayload = {};
            uListPayload["label"] = "logout-update";
            uListPayload["cookie_value"] = cookieVal;
            console.log("logout UL sending: ", uListPayload);
            userListSocket.send(JSON.stringify(uListPayload));

            // empty user list
            const uList = document.querySelector(".user-list");
            uList.textContent = "";

            // display login and reg
            const navbar = document.querySelector(".navbar")
            navbar.children[0].style.display = "block"
            navbar.children[1].style.display = "block"
            navbar.children[2].style.display = "none"
        }
    })
}


const logoutBtn = document.createElement("button");
logoutBtn.textContent = "Logout";
const logoutDiv = document.querySelector("#logout");
logoutBtn.addEventListener("click", logoutHandler);

export default logoutBtn;