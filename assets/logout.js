const logoutUrl = location.origin + "/logout/";
console.log(logoutUrl);
const logoutHandler = function() {
    fetch(logoutUrl)
    .then(() => {       
        document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 GMT";
        console.log("LogOut");
    })
}


const logoutBtn = document.createElement("div");
logoutBtn.textContent = "Logout";
logoutBtn.addEventListener("click", logoutHandler);

export default logoutBtn;