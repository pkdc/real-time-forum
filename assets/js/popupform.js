let userPopUpPBtns = [...document.querySelectorAll(".button")];
userPopUpPBtns.forEach(function (btn) {
    btn.onclick = function () {
        let userPopUpP = btn.getAttribute("data-userPopUpP");
        document.getElementById(userPopUpP).style.display = "block";
    };
});
let closeBtns = [...document.querySelectorAll(".close")];
closeBtns.forEach(function (btn) {
    btn.onclick = function () {
        let userPopUpP = btn.closest(".userPopUpP");
        userPopUpP.style.display = "none";
    };
});
window.onclick = function (event) {
    if (event.target.className === "userPopUpP") {
        event.target.style.display = "none";
    }
};
