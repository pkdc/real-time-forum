// import throttle from '/assets/js/node_modules/lodash-es/throttle.js';
import { chatSocket, targetUserId, timenow, genTypingDiv } from "./chat.js";
const userListSocket = new WebSocket("ws://localhost:8080/userListWs/")
const chatBox = document.querySelector(".col-1")
const msgArea = document.querySelector(".msgArea")
const userdiv = document.querySelector(".usernamediv")
let usname = document.createElement("p")
userdiv.appendChild(usname)
let first = 0
let usID
let open = false
let loadMsg = false
let timerId = undefined
var realTargetUser = targetUserId
function throttle(fn, wait) {
    if (timerId) {
        return
    }

    fn()

    timerId = setTimeout(function () {
        timerId = undefined
    }, wait)
}

document.addEventListener("DOMContentLoaded", function (e) {
    realTargetUser = targetUserId
    // userListSocket = new WebSocket("ws://localhost:8080/userListWs/")
    console.log("JS attempt to connect to user list");
    userListSocket.onopen = () => console.log("user list connected");
    userListSocket.onclose = () => console.log("Bye user list");
    userListSocket.onerror = (err) => console.log("user list ws Error!");
    userListSocket.onmessage = (msg) => {

        let arrayOfUsers = []
        let arr
        let numberOfUsers
        const resp = JSON.parse(msg.data);
        if (resp.label === "update") {
            let profileid = document.querySelector(".ProfileID")
            console.log("CHECK THIS RESP", resp)
            if (resp.realUser == parseInt(profileid.textContent)) {

                console.log(resp.online_users);
                const uList = document.querySelector(".user-list");
                // remove list item
                uList.textContent = "";
                // add new list item

                for (const { nickname, status, userID, msgcheck, noti } of resp.online_users) {
                    arrayOfUsers.push(userID)
                    numberOfUsers = (resp.online_users).length
                    const nicknameItem = document.createElement("li");
                    nicknameItem.id = "li" + userID
                    const chatBoxButton = document.createElement("button")
                    chatBoxButton.classList = "nameButtons"
                    const chatBoxForm = document.createElement("form")
                    chatBoxForm.addEventListener("submit", showChatHandler)
                    chatBoxButton.setAttribute("type", "submit")
                    chatBoxButton.value = userID
                    chatBoxButton.id = `ContactID-${userID}`
                    chatBoxButton.addEventListener("click", function (e) {
                        if (open == false) {
                            open = true
                            chatBoxButton.className = "nameButtons"
                            usID = chatBoxButton.value
                            usname.textContent = chatBoxButton.textContent
                            chatBox.id = `chatbox-${usID}`
                            chatBox.style.display = "block"
                            window.onclick = function (event) {
                                // console.log("buttonhere clicked")
                                // scr()
                                open = false
                                if (event.target.className == "closeChat") {
                                    chatBox.style.display = "none"
                                    chatBox.id = "chatbox"
                                    loadMsg = false
                                    first = 0
                                    while (msgArea.firstChild) {
                                        msgArea.removeChild(msgArea.firstChild)
                                    }
                                }
                            }
                        }
                        if (open == true) {
                            usID = chatBoxButton.value
                            chatBox.id = `chatbox-${usID}`
                            usname.textContent = chatBoxButton.textContent
                            loadMsg = false
                            while (msgArea.firstChild) {
                                msgArea.removeChild(msgArea.firstChild)
                            }
                        }

                    })
                    chatBoxForm.append(chatBoxButton)
                    let userNick = document.querySelector(".ProfileNickname") // reg userNick is null
                    chatBoxButton.textContent = `${nickname}`;
                    if (chatBoxButton.textContent == userNick.textContent) {
                        arr = noti.split(",")
                        nicknameItem.style.display = "none"
                    }
                    if (status == false) {
                        nicknameItem.classList = "offline"
                    } else {
                        nicknameItem.classList = "online"
                    }

                    nicknameItem.append(chatBoxForm)
                    uList.append(nicknameItem);


                }
                let not = checkArr(arr, arrayOfUsers)
                if (not.length > 0) {
                    notiDisplay(not)

                }
                // if button already exists do not create another

                if (document.querySelector(".closeChat")) {

                } else {

                    const closeChatBox = document.createElement("button")
                    closeChatBox.textContent = "End Chat"
                    closeChatBox.classList = "closeChat"
                    chatBox.append(closeChatBox)
                }

                console.log("CHECKING IMPORTED USER-1", realTargetUser)
                if (realTargetUser != null) {
                    console.log("CHECKING IMPORTED USER-2", realTargetUser)
                    // let userlist = document.querySelector(".user-list")
                    // let targetUser = document.querySelector(`#li${realTargetUser}`)
                    // userlist.insertBefore(targetUser, userlist.firstChild)
                }
            }
            else {
                console.log("this is all response after logged in as a newUser", resp)
                if ((resp.online_users).length != numberOfUsers) {
                    let inde = regNewUser(resp.online_users)
                    if (typeof (inde) == "number") {
                        console.log("found this index:", inde)
                        let newus = resp.online_users[inde]
                        console.log("newuser:", newus)
                        const uList = document.querySelector(".user-list");
                        let nickname = newus.nickname
                        let status = newus.status
                        let userID = newus.userID
                        let noti = newus.noti
                        arrayOfUsers.push(userID)
                        numberOfUsers = (resp.online_users).length
                        const nicknameItem = document.createElement("li");
                        nicknameItem.id = "li" + userID
                        const chatBoxButton = document.createElement("button")
                        chatBoxButton.className = "nameButtons"
                        const chatBoxForm = document.createElement("form")
                        chatBoxForm.addEventListener("submit", showChatHandler)
                        chatBoxButton.setAttribute("type", "submit")
                        chatBoxButton.value = userID
                        chatBoxButton.id = `ContactID-${userID}`
                        chatBoxButton.addEventListener("click", function (e) {
                            if (open == false) {
                                open = true
                                usID = chatBoxButton.value
                                usname.textContent = chatBoxButton.textContent
                                chatBox.id = `chatbox-${usID}`
                                chatBox.style.display = "block"
                                window.onclick = function (event) {
                                    // console.log("buttonhere clicked")
                                    // scr()
                                    open = false
                                    if (event.target.className == "closeChat") {
                                        chatBox.style.display = "none"
                                        chatBox.id = "chatbox"
                                        loadMsg = false
                                        while (msgArea.firstChild) {
                                            msgArea.removeChild(msgArea.firstChild)
                                        }
                                    }
                                }
                            }
                            if (open == true) {
                                usID = chatBoxButton.value
                                chatBox.id = `chatbox-${usID}`
                                usname.textContent = chatBoxButton.textContent
                                loadMsg = false
                                while (msgArea.firstChild) {
                                    msgArea.removeChild(msgArea.firstChild)
                                }
                            }

                        })
                        chatBoxForm.append(chatBoxButton)
                        let userNick = document.querySelector(".ProfileNickname") // reg userNick is null
                        chatBoxButton.textContent = `${nickname}`;
                        if (chatBoxButton.textContent == userNick.textContent) {
                            arr = noti.split(",")
                            nicknameItem.style.display = "none"
                        }
                        if (status == false) {
                            nicknameItem.classList = "offline"
                        } else {
                            nicknameItem.classList = "online"
                        }

                        nicknameItem.append(chatBoxForm)
                        uList.insertBefore(nicknameItem, uList.children[inde + 1])
                        // if button already exists do not create another

                        if (document.querySelector(".closeChat")) {

                        } else {

                            const closeChatBox = document.createElement("button")
                            closeChatBox.textContent = "End Chat"
                            closeChatBox.classList = "closeChat"
                            chatBox.append(closeChatBox)
                        }

                        console.log("CHECKING IMPORTED USER-1", realTargetUser)
                        if (realTargetUser != null) {
                            console.log("CHECKING IMPORTED USER-2", realTargetUser)
                            // let userlist = document.querySelector(".user-list")
                            // let targetUser = document.querySelector(`#li${realTargetUser}`)
                            // userlist.insertBefore(targetUser, userlist.firstChild)
                        }

                    }
                }
                for (const { nickname, status, userID, msgcheck, noti } of resp.online_users) {
                    arrayOfUsers.push(userID)
                    let userNick = document.querySelector(".ProfileNickname")
                    if (nickname == userNick.textContent) {
                        arr = noti.split(",")
                    }
                    let singleList = document.querySelector("#li" + userID)
                    console.log("SINGLELIST", singleList)
                    if (status == false) {
                        singleList.classList = "offline"
                    } else if (status == true) {
                        singleList.classList = "online"
                    }
                }
                let not = checkArr(arr, arrayOfUsers)
                if (not.length > 0) {
                    notiDisplay(not)

                }
            }
        }
        if (resp.label == "chatBox") {
            // console.log("labelhere")
            // scr()

            if (msgArea.firstElementChild == null) {
                msgArea.addEventListener("scroll", () => { throttle(loadMsgCallback, 200) });
            }
            let js = JSON.parse(resp.content)
            if (js != null) {
                let prevScrollHeight = msgArea.scrollHeight
                for (let i = 0; i < js.length; i++) {
                    let singleMsg = document.createElement("div")
                    let msgContent = document.createElement("p")
                    let timeOfMsg = document.createElement("p")
                    msgContent.classList = "msg-text"
                    timeOfMsg.classList = "timeofmsg"
                    timeOfMsg.textContent = js[i].msgInfo.message_time
                    timeOfMsg.style.fontSize = "9px"
                    msgContent.textContent = js[i].msgInfo.content
                    if (js[i].msgInfo.right_side == true) {
                        singleMsg.classList = "msg-row2"
                    } else {
                        singleMsg.classList = "msg-row"
                    }
                    singleMsg.append(msgContent)
                    singleMsg.append(timeOfMsg)
                    msgArea.insertBefore(singleMsg, msgArea.firstChild)

                    // msgArea.append(singleMsg)
                }
                if (first == 0) {
                    first++
                    msgArea.scrollTop = msgArea.scrollHeight
                } else {
                    msgArea.scrollTop = msgArea.scrollHeight - prevScrollHeight
                }
                console.log(msgArea.scrollTop)
            }
            if (document.querySelector(".chatInput") == null) {
                console.log("creating chat input")
                const chatInput = document.createElement("textarea")
                // chatInput.setAttribute("type", "text")
                const chatFormDiv = document.createElement("div");
                chatFormDiv.classList.add("chat-form-div");
                const chatForm = document.createElement("form")
                const submitChat = document.createElement("button")
                chatForm.addEventListener("submit", SubChatHandler)
                chatForm.id = "chat-form";
                submitChat.setAttribute("type", "submit")
                submitChat.classList = "submitMsg"
                submitChat.textContent = "send"
                chatInput.classList = "chatInput"
                chatForm.append(chatInput, submitChat)
                chatFormDiv.append(chatForm)

                let typingDiv = genTypingDiv();
                
                chatBox.append(chatFormDiv, typingDiv);

                chatInput.addEventListener("input", function(e) {
                    // console.log("input event");
                    const profileid = document.querySelector(".ProfileID");
                    let typingPayloadObj = {};
                    typingPayloadObj["label"] = "typing";
                    typingPayloadObj["sender_id"] = parseInt(profileid.textContent)
                    typingPayloadObj["receiver_id"] = parseInt(usID)
                    console.log(`${typingPayloadObj["sender_id"]} is Typing, and sends to ${typingPayloadObj["receiver_id"]}`);
                    chatSocket.send(JSON.stringify(typingPayloadObj));
                });
                chatInput.addEventListener("blur", function(e) {
                    console.log("blur event");
                    const profileid = document.querySelector(".ProfileID");
                    let stopTypingPayloadObj = {};
                    stopTypingPayloadObj["label"] = "stop-typing";
                    stopTypingPayloadObj["sender_id"] = parseInt(profileid.textContent)
                    stopTypingPayloadObj["receiver_id"] = parseInt(usID)
                    console.log(`${stopTypingPayloadObj["sender_id"]} has stopped Typing, and stop sending to ${stopTypingPayloadObj["receiver_id"]}`);
                    chatSocket.send(JSON.stringify(stopTypingPayloadObj));
                });
            } else {
                console.log("chatinput already exist")
            }
        }
    }
})

const loadMsgCallback = function () {
    console.log("scrolling", msgArea.scrollTop)
    if (msgArea.scrollTop <= 200) {
        // console.log(`Loading msg ... msgArea.scrollTop = ${msgArea.scrollTop} load msg when value <= 20 `);
        loadMsg = true
        loadPrevMsgsHandler();

    }
}

const loadPrevMsgsHandler = function () {
    let payloadObj = {};
    let profileid = document.querySelector(".ProfileID")
    payloadObj["label"] = "createChat";
    payloadObj["userID"] = parseInt(profileid.textContent) /* after login change to loggedUserID */
    payloadObj["contactID"] = parseInt(usID)
    payloadObj["loadMsg"] = loadMsg
    userListSocket.send(JSON.stringify(payloadObj));
};

const showChatHandler = function (e) {
    e.preventDefault();
    let payloadObj = {}
    let profileid = document.querySelector(".ProfileID")
    payloadObj["label"] = "createChat";
    payloadObj["userID"] = parseInt(profileid.textContent) /* after login change to loggedUserID */
    payloadObj["contactID"] = parseInt(usID)
    payloadObj["loadMsg"] = loadMsg
    // payloadObj["loadMsg"] = true
    userListSocket.send(JSON.stringify(payloadObj));
    let chatPayloadObj = {};
    chatPayloadObj["label"] = "createChat";
    chatPayloadObj["sender_id"] = parseInt(profileid.textContent)/* after login change to loggedUserID */
    chatPayloadObj["receiver_id"] = parseInt(usID)
    chatSocket.send(JSON.stringify(chatPayloadObj));
};
const SubChatHandler = function (e) {
    e.preventDefault();
    let chatPayloadObj = {};
    let profileid = document.querySelector(".ProfileID")
    let chatInput = document.querySelector(".chatInput")
    let msgrow = document.createElement("div")
    let msgtext = document.createElement("p")
    msgrow.className = "msg-row2"
    msgtext.className = "msg-text"
    msgtext.textContent = chatInput.value
    let userlist = document.querySelector(".user-list")
    realTargetUser = usID
    let targetUser = document.querySelector(`#li${usID}`)
    userlist.insertBefore(targetUser, userlist.firstChild)
    let timeOfMsg = document.createElement("p")
    timeOfMsg.classList = "timeofmsg"
    timeOfMsg.textContent = timenow()
    timeOfMsg.style.fontSize = "9px"
    msgrow.append(msgtext)
    msgArea.append(msgrow)
    // ***********************NEED TO UPDATE USERLIST *********************
    // let userlist = document.querySelector(".user-list")
    // while (userlist.firstChild) {
    //     console.log("userlist,", userlist.length)
    //     userlist.removeChild(userlist.firstChild)
    // }
    chatPayloadObj["label"] = "chat";
    chatPayloadObj["sender_id"] = parseInt(profileid.textContent)/* after login change to loggedUserID */
    chatPayloadObj["receiver_id"] = parseInt(usID)
    chatPayloadObj["content"] = chatInput.value
    chatInput.value = ""
    chatSocket.send(JSON.stringify(chatPayloadObj));

}
function checkArr(arr, userArr) {
    for (let i = 0; i < arr.length; i++) {
        arr[i] = parseInt(arr[i])
    }
    const intersection = arr.filter(element => userArr.includes(element));
    return intersection
}
function notiDisplay(arr) {

    for (let i = 0; i < arr.length; i++) {
        console.log("notification:", arr)
        let x = document.querySelector("#ContactID-" + arr[i])
        if (x != null) {
            let k = document.querySelector(`#chatbox-${arr[i]}`)
            console.log("K:", k)
            if (k == null) {
                x.classList.add("notif")
            }
        }
    }
}
// function scr() {
//     console.log("AUTO SCROLLLLL")
//     console.log("first", msgArea.scrollTop)
//     console.log("height", msgArea.scrollHeight)
//     console.log("width", msgArea.scrollWidth)
//     msgArea.scrollTop = msgArea.scrollHeight;
//     console.log("second", msgArea.scrollTop)
// }

function regNewUser(newuserlist) {
    let userlist = document.querySelector(".user-list")
    let lists = userlist.children
    // let intercept
    let array1 = []
    let array2 = []
    for (let i = 0; i < lists.length; i++) {

        // console.log("ALLCHILDS", lists[i])
        // console.log("ALLOFNEWCHILDS", newuserlist[i])
        let idofchild = lists[i].id
        idofchild = idofchild.replace("li", "")
        array1.push(parseInt(idofchild))
        // console.log("idofchild,", idofchild, "::::",idofchild.length)
        // console.log(parseInt(idofchild) == newuserlist[i].userID)
        // console.log(parseInt(idofchild) ,"::::", newuserlist[i].userID)
        // if (parseInt(idofchild) != newuserlist[i].userID) {
        //     intercept = i
        //     return intercept
        // }
    }
    for (let i = 0; i < newuserlist.length; i++) {
        array2.push(newuserlist[i].userID)
    }
    var array3 = array2.filter((o) => array1.indexOf(o) === -1);
    let final = array3[0]
    for (let i = 0; i < newuserlist.length; i++) {
        if (newuserlist[i].userID == final) {
            return i
        }
    }
}
export default userListSocket;