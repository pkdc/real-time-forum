export const chatSocket = new WebSocket("ws://localhost:8080/chatWs/");
export var targetUserId = null
export const chatForm = document.createElement("form");
const msgArea = document.querySelector(".msgArea")
chatForm.id = "chat-form";
chatForm.addEventListener("submit", function (e) {
    e.preventDefault();
    // add msg
    // send msg to ws
});
const chatInputDiv = document.createElement("div");
chatInputDiv.id = "chat-input-div";
const chatInput = document.createElement("input");
chatInputDiv.append(chatInput);

let animationID = null;
let typingLeftDotsArr, typingRightDotsArr;
let leftDot1Pos = 0;
let rightDot1Pos = 0;
let leftDot2Pos = 0;
let rightDot2Pos = 0;
let leftDot3Pos = 0;
let rightDot3Pos = 0;
let pairsOfDotsRemoved = 0;
const dotSpeed = 2;
const dotDist = 200;
let running = false;
let textSender = "";
let genCount = 0;
// const sendBtn = document.createElement("button");
// sendBtn.textContent = "Send";
// sendBtn.id = "send-btn";
// chatForm.append(chatInputDiv, sendBtn);

document.addEventListener("DOMContentLoaded", function (e) {
    // chatSocket = new WebSocket("ws://localhost:8080/chatWs/");
    console.log("JS attempt to connect to chat");
    chatSocket.onopen = () => console.log("chat connected");
    chatSocket.onclose = () => console.log("Bye chat");
    chatSocket.onerror = (err) => console.log("chat ws Error!");
    chatSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        if (resp.label === "created_room") {
            // console.log(`chat room created between ${resp.sender_id} and ${resp.receiver_id}`);
        } else if (resp.label === "msgIncoming") {
            console.log("recievedChatMsg", resp)
            let msgrow = document.createElement("div")
            let msgtext = document.createElement("p")
            msgrow.className = "msg-row"
            msgtext.className = "msg-text"
            msgtext.textContent = resp.content
            msgrow.append(msgtext)
            msgArea.append(msgrow)
            targetUserId= resp.contactID
        //     var userlist = document.querySelector(".user-list")
        //    var targetUser = document.querySelector(`#li${resp.contactID}`)
        //     console.log("targetuser:", targetUser)
        //     userlist.insertBefore(targetUser, userlist.firstChild)
        } else if (resp.label === "sender-typing") {
            // display typing-in-progress
            // const chatBox = document.querySelector(".chatbox");
            const typingDiv = document.querySelector(".typing-div");
            const typingText = document.querySelector(".typing-text");
            console.log(`${resp.sender} with id ${resp.userID} is typing... to ${resp.contactID}`);
            textSender = resp.sender;
			typingText.textContent = `${textSender} is typing`;
			console.log("fired");
            typingDiv.classList.add("show");
            // typingDotsArr = [...document.querySelectorAll(".typing-dots")];
            typingLeftDotsArr = [...document.querySelectorAll(".typing-left-dots")];
            typingRightDotsArr = [...document.querySelectorAll(".typing-right-dots")];
            if (animationID === null && running === false) {
                animationID = requestAnimationFrame(animateDots);
				running = true;
            }
        } else if (resp.label === "sender-stop-typing") {
			const typingDiv = document.querySelector(".typing-div");
			animationID = null;
			running = false;
			typingDiv.classList.remove("show");
		}
    }
})

export const genTypingDiv = function () {
	console.log("gen");
	const typingDiv = document.createElement("div");
	typingDiv.classList.add("typing-div");
	if (genCount >= 1) typingDiv.classList.add("show"); // show after regen, but don't show at the beginning
    const typingTextDiv = document.createElement("div");
	const typingText = document.createElement("p");
    typingTextDiv.append(typingText);
	genCount += 1;

	const typingLeftDotsDiv = document.createElement("div");
	typingLeftDotsDiv.classList.add("typing-left-dots-div");
	const typingLeftDot1 = document.createElement("p");
	const typingLeftDot2 = document.createElement("p");
	const typingLeftDot3 = document.createElement("p");
	typingLeftDot1.classList.add("typing-dots");
	typingLeftDot1.classList.add("typing-left-dots");
	typingLeftDot2.classList.add("typing-dots");
	typingLeftDot2.classList.add("typing-left-dots");
	typingLeftDot3.classList.add("typing-dots");
	typingLeftDot3.classList.add("typing-left-dots");
	typingLeftDotsDiv.append(typingLeftDot1, typingLeftDot2, typingLeftDot3);

	const typingRightDotsDiv = document.createElement("div");
	typingRightDotsDiv.classList.add("typing-right-dots-div");
	const typingRightDot1 = document.createElement("p");
	const typingRightDot2 = document.createElement("p");
	const typingRightDot3 = document.createElement("p");
    typingRightDot1.classList.add("typing-dots");
	typingRightDot1.classList.add("typing-right-dots");
	typingRightDot2.classList.add("typing-dots");
	typingRightDot2.classList.add("typing-right-dots");
	typingRightDot3.classList.add("typing-dots");
	typingRightDot3.classList.add("typing-right-dots");
	typingRightDotsDiv.append(typingRightDot1, typingRightDot2, typingRightDot3);

	typingText.classList.add("typing-text");
	typingLeftDot1.textContent = "·";
	typingLeftDot2.textContent = "·";
	typingLeftDot3.textContent = "·";
	typingText.textContent = `${textSender} is typing`;
	// typingText.textContent = `I am typing`; 
	typingRightDot1.textContent = "·";
	typingRightDot2.textContent = "·";
	typingRightDot3.textContent = "·";
	typingDiv.append(typingLeftDotsDiv, typingTextDiv, typingRightDotsDiv);

	return typingDiv;
};

const reset = function() {
    pairsOfDotsRemoved = 0;
    leftDot1Pos = 0;
    rightDot1Pos = 0;
    leftDot2Pos = 0;
    rightDot2Pos = 0;
    leftDot3Pos = 0;
    rightDot3Pos = 0;
	typingLeftDotsArr = [...document.querySelectorAll(".typing-left-dots")];
    typingRightDotsArr = [...document.querySelectorAll(".typing-right-dots")];
}

const animateDots = function () {
	const typingLeftDotsDiv = document.querySelector(".typing-left-dots-div");
	const typingRightDotsDiv = document.querySelector(".typing-right-dots-div");
	if (running === true) {
		if (leftDot1Pos !== null && rightDot1Pos !== null) {
			leftDot1Pos -= dotSpeed;
			console.log("left 1:", leftDot1Pos);
			typingLeftDotsArr[0].style.left = `${leftDot1Pos}px`;

			rightDot1Pos += dotSpeed;
			console.log("right 1:", rightDot1Pos);
			typingRightDotsArr[0].style.left = `${rightDot1Pos}px`;
		}

		if (
			rightDot2Pos !== null &&
			leftDot2Pos !== null &&
			(rightDot1Pos > dotDist || rightDot1Pos === null) &&
			(leftDot1Pos < -dotDist || leftDot1Pos === null)
		) {
			if (pairsOfDotsRemoved <= 0) {
				typingLeftDotsDiv.firstElementChild.remove();
				typingRightDotsDiv.firstElementChild.remove();
				pairsOfDotsRemoved += 1;
			}
			leftDot1Pos = null;

			leftDot2Pos -= dotSpeed;
			console.log("left 2:", leftDot2Pos);
			typingLeftDotsArr[1].style.left = `${leftDot2Pos}px`;

			rightDot1Pos = null;
			rightDot2Pos += dotSpeed;
			console.log("right 2:", rightDot2Pos);
			typingRightDotsArr[1].style.left = `${rightDot2Pos}px`;
		}

		if (
			rightDot3Pos !== null &&
			leftDot3Pos !== null &&
			(rightDot2Pos > dotDist || rightDot2Pos === null) &&
			(leftDot2Pos < -dotDist || leftDot2Pos === null)
		) {
			if (pairsOfDotsRemoved === 1) {
				typingLeftDotsDiv.firstElementChild.remove();
				typingRightDotsDiv.firstElementChild.remove();
				pairsOfDotsRemoved += 1;
			}
			leftDot2Pos = null;
			leftDot3Pos -= dotSpeed;
			console.log("left 3:", leftDot3Pos);
			typingLeftDotsArr[2].style.left = `${leftDot3Pos}px`;

			rightDot2Pos = null;
			rightDot3Pos += dotSpeed;
			console.log("right 3:", rightDot3Pos);
			typingRightDotsArr[2].style.left = `${rightDot3Pos}px`; // why undefined when inputting during the 2nd/3rd loop?
		}
		if (rightDot3Pos > dotDist && leftDot3Pos < -dotDist) {
			if (pairsOfDotsRemoved === 2) {
				typingLeftDotsDiv.firstElementChild.remove();
				typingRightDotsDiv.firstElementChild.remove();
				pairsOfDotsRemoved += 1;
			}
			const prevTypingDiv = document.querySelector(".typing-div");
			// typingLeftDotsDiv.remove();
			// typingRightDotsDiv.remove();
			prevTypingDiv.remove();
			console.log("prevTypingDiv:      ", prevTypingDiv);
			// leftDot3Pos = null;
			// rightDot3Pos = null;
			console.log("regen");
			const typingDiv = genTypingDiv();
			const chatBox = document.querySelector(".col-1");
			chatBox.append(typingDiv);
			reset();
		}

		requestAnimationFrame(animateDots);
	}
};
