import userListSocket from "./userList.js";

const logoutUrl = location.origin + "/logout/";
console.log(logoutUrl);
const logoutHandler = function() {
    fetch(logoutUrl)
    .then(() => {
        // get cookie val
        const sessionCookie = document.cookie.split(";").find((el) => el.startsWith("session="));
        console.log("session cookie: ", sessionCookie);
        const cookieVal = sessionCookie.split("=")[1];
        console.log("cookie val: ", cookieVal);

        // remove cookie
        document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 GMT";

        // update user list after a user logout
        let uListPayload = {};
        uListPayload["label"] = "update";
        uListPayload["cookie_value"] = cookieVal;
        console.log("logout UL sending: ", uListPayload);
        userListSocket.send(JSON.stringify(uListPayload));
    })
}


const logoutBtn = document.createElement("div");
logoutBtn.textContent = "Logout";
logoutBtn.addEventListener("click", logoutHandler);

export default logoutBtn;