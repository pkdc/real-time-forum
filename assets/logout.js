const logoutUrl = location.origin + "/logout/";
console.log(logoutUrl);
const logoutHandler = function() {
    fetch(logoutUrl)
    .then(() => {       
        // document.cookie = "";
        // browser.cookies.remove
        console.log("LogOut");
    })
}


const logoutBtn = document.createElement("div");
logoutBtn.textContent = "Logout";
logoutBtn.addEventListener("click", logoutHandler);

export default logoutBtn;