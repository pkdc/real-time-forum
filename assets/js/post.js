let postSocket = null;
// const body = document.getElementsByTagName("BODY")[0]
const DisplayPost = document.createElement("div")
DisplayPost.className = "ContainAllPost"
var commentPostId
let jsonFile
document.addEventListener("DOMContentLoaded", function () {
    postSocket = new WebSocket("ws://localhost:8080/postWs/");
    console.log("JS attempt to connect post");
    postSocket.onopen = () => console.log("connected-postCreate");
    postSocket.onclose = () => console.log("Bye-postCreate");
    postSocket.onerror = (err) => console.log("Error!-postCreate", err);
    postSocket.onmessage = (msg) => {
        const resp = JSON.parse(msg.data);
        if (resp.label === "Greet") {
            jsonFile = JSON.parse(resp.content)
            if (jsonFile !== null) {
                createPost(jsonFile)
            }
        } else if (resp.label === "post") {
            jsonFile = JSON.parse(resp.content)
            createPost(jsonFile)
        }
        else if (resp.label === "Createcomment") {
            jsonFile = JSON.parse(resp.content)
            if (resp.Content !== null) {
                const comment = document.querySelector(".comment")
                let newCom = CreateComments(jsonFile, parseInt(commentPostId) - 1)
                //    comment.append(newCom)
                console.log("NEWCOM", newCom)
                comment.insertBefore(newCom, comment.children[1])
            }
        } else if (resp.label === "showComment") {
            jsonFile = JSON.parse(resp.content)
        }
    }
});
function createPost(arr) {

    document.querySelectorAll("#allPost").forEach(e => {
        e.remove();
    });
    const allPost = document.createElement("div")
    allPost.id = "allPost"
    for (let i = 0; i < arr.length; i++) {
        const postDiv = document.createElement("div")
        const titleDiv = document.createElement("div");
        const titleButton = document.createElement("button")
        titleButton.className = "titleButtonForm"
        const titleForm = document.createElement("form")
        titleForm.addEventListener("submit", showcommentHandler)
        titleButton.setAttribute("value", i)
        titleButton.setAttribute("type", "submit")
        titleButton.addEventListener("click", function (e) {
            const comment = document.querySelector(".comment")
            comment.style.left = "20%"
            let chosenPost = document.querySelector(`#post-${titleButton.value}`)
            console.log("post is choosen")
            let clone = chosenPost.cloneNode(true)
            const closeComments = document.createElement("button")
            closeComments.className = "closeCommentsButton"
            // closeComments.textContent = String.fromCodePoint(0x274C)
            closeComments.textContent = "Close Comment"
            closeComments.addEventListener("click", function () {
                comment.style.left = "-70%"
                PostHandler
                while (comment.firstChild) {
                    comment.removeChild(comment.firstChild)
                }
            })
            let comments = CreateComments(jsonFile, i)
            let comForm = CreateCommentForm(titleButton.value)
            comment.append(clone, comments, comForm, closeComments)
            // comment.style.height = "100%";
        })
        titleButton.innerText = (arr[i].postinfo.title)
        titleButton.style.padding = "0"
        titleButton.style.border = "none"
        titleButton.style.background = "none"
        const contentDiv = document.createElement("div");
        const categoryDiv = document.createElement("div");
        const userIdDiv = document.createElement("div");
        postDiv.id = `post-${i}`;
        titleDiv.id = `title-${i}`;
        contentDiv.id = `content-${i}`;
        categoryDiv.id = `category-${i}`;
        userIdDiv.id = `id-${i}`;
        postDiv.className = `singlePostDiv`;
        titleDiv.className = `postTitle`;
        contentDiv.className = `postContent`;
        categoryDiv.className = `postCategoryId`;
        userIdDiv.className = `postAuthorId`;
        // const titleText = document.createElement("p")
        // titleText.style.fontWeight= "900"
        const contentText = document.createElement("p")
        const categoryText = document.createElement("p")
        // categoryText.style.width = "10%"
        const userIdText = document.createElement("p")
        // const titletextNode = document.createTextNode(arr[i].postinfo.title) 
        // titleText.appendChild(titletextNode)
        const contenttextNode = document.createTextNode(arr[i].postinfo.Content)
        contentText.appendChild(contenttextNode)
        const categorytextNode = document.createTextNode(arr[i].postinfo.category_option)
        categoryText.appendChild(categorytextNode)
        const userIdtextNode = document.createTextNode(arr[i].postinfo.userID)
        userIdText.appendChild(userIdtextNode)
        // titleDiv.append(titleText)
        titleDiv.append(titleButton)
        titleForm.append(titleDiv)
        contentDiv.append(contentText)
        categoryDiv.append(categoryText)
        userIdDiv.append(userIdText)
        categoryDiv.style.display = "inline-block"
        userIdDiv.style.display = "inline-block"
        userIdDiv.style.float = "right"
        const bottomline = document.createElement("div");
        bottomline.className = "bottomline";
        bottomline.append(categoryDiv, userIdDiv);
        postDiv.append(titleForm, contentDiv, bottomline)
        allPost.append(postDiv)
    }
    DisplayPost.appendChild(allPost)
}
const PostHandler = function (e) {
    e.preventDefault();
    const formFields = new FormData(e.target);
    const payloadObj = Object.fromEntries(formFields.entries());
    payloadObj["label"] = "post";
    let profid = document.querySelector(".ProfileID").textContent
    payloadObj["userID"] = document.querySelector(".ProfileID").textContent
    console.log("checking target", payloadObj)
    postSocket.send(JSON.stringify(payloadObj));
};

const PostForm = document.createElement("form");
PostForm.addEventListener("submit", PostHandler);
PostForm.className = "newPostFormSection"

const titleLabelDiv = document.createElement('div');
titleLabelDiv.className = "newPostTitlesArea"
const titleLabel = document.createElement('label');
titleLabel.className = "newPostLabelArea"
titleLabel.textContent = "title";
titleLabel.setAttribute("for", "title");
titleLabelDiv.append(titleLabel);
const titleInputDiv = document.createElement('div');
titleInputDiv.className = "newPostTitleInputDiv"
const titleInput = document.createElement('textarea');
titleInput.className = "newPostInput"
titleInput.setAttribute("maxlength", 60);
titleInput.setAttribute("name", "title");
titleInput.setAttribute("id", "title");
titleInputDiv.append(titleInput);
//-------------------
const CatDiv = document.createElement('div');
CatDiv.className = "newPostCatDiv"
const CatOptionDiv = document.createElement('select');
CatOptionDiv.setAttribute("name", "category_option")
const CatLabel = document.createElement("label");
CatLabel.textContent = "Please choose category";
CatLabel.className = "categoryLableTitle"
CatLabel.setAttribute("for", "cat");
CatDiv.append(CatLabel);
const CatInputOpt1 = document.createElement("option");
const CatInputOpt2 = document.createElement("option");
const CatInputOpt3 = document.createElement("option");
const CatInputOpt4 = document.createElement("option");
const CatInputOpt5 = document.createElement("option");
const CatInputOpt6 = document.createElement("option");
const CatInputOpt7 = document.createElement("option");
const CatInputOpt8 = document.createElement("option");
const CatInputOpt9 = document.createElement("option");
const CatInputOpt10 = document.createElement("option");
CatInputOpt1.setAttribute("name", "category_option");
CatInputOpt2.setAttribute("name", "category_option");
CatInputOpt3.setAttribute("name", "category_option");
CatInputOpt4.setAttribute("name", "category_option");
CatInputOpt5.setAttribute("name", "category_option");
CatInputOpt6.setAttribute("name", "category_option");
CatInputOpt7.setAttribute("name", "category_option");
CatInputOpt8.setAttribute("name", "category_option");
CatInputOpt9.setAttribute("name", "category_option");
CatInputOpt10.setAttribute("name", "category_option");
CatInputOpt1.setAttribute("id", "1");
CatInputOpt2.setAttribute("id", "2");
CatInputOpt3.setAttribute("id", "3");
CatInputOpt4.setAttribute("id", "4");
CatInputOpt5.setAttribute("id", "5");
CatInputOpt6.setAttribute("id", "6");
CatInputOpt7.setAttribute("id", "7");
CatInputOpt8.setAttribute("id", "8");
CatInputOpt9.setAttribute("id", "9");
CatInputOpt10.setAttribute("id", "10");
CatInputOpt1.setAttribute("value", "HTML");
CatInputOpt2.setAttribute("value", "CSS");
CatInputOpt3.setAttribute("value", "JAVASCRIPT");
CatInputOpt4.setAttribute("value", "GOLANG");
CatInputOpt5.setAttribute("value", "DOCKER");
CatInputOpt6.setAttribute("value", "PHP");
CatInputOpt7.setAttribute("value", "WEBSOCKET");
CatInputOpt8.setAttribute("value", "SQL");
CatInputOpt9.setAttribute("value", "PYTHON");
CatInputOpt10.setAttribute("value", "MISCELLANEOUS");
CatInputOpt1.textContent = "HTML";
CatInputOpt2.textContent = "CSS";
CatInputOpt3.textContent = "JAVASCRIPT";
CatInputOpt4.textContent = "GOLANG";
CatInputOpt5.textContent = "DOCKER";
CatInputOpt6.textContent = "PHP";
CatInputOpt7.textContent = "WEBSOCKET";
CatInputOpt8.textContent = "SQL";
CatInputOpt9.textContent = "PYTHON";
CatInputOpt10.textContent = "MISCELLANEOUS";
CatOptionDiv.setAttribute("id", "category");
CatOptionDiv.append(CatInputOpt1, CatInputOpt2, CatInputOpt3, CatInputOpt4, CatInputOpt5, CatInputOpt6, CatInputOpt7, CatInputOpt8, CatInputOpt9, CatInputOpt10)

const contLabelDiv = document.createElement('div');
contLabelDiv.className = "newPostContentTitleDiv"
const contLabel = document.createElement('label');
contLabel.className = "newPostContentTitle"
contLabel.textContent = "content:";
contLabel.setAttribute("for", "content");
contLabelDiv.append(contLabel);
const contInputDiv = document.createElement('div');
contInputDiv.className = "newPostContentInputArea"
const contInput = document.createElement('textarea');
contInput.className = "newPostContentInputBox"
contInput.setAttribute("name", "content");
contInput.setAttribute("id", "content");
contInputDiv.append(contInput);

const PostSubmitDiv = document.createElement('div');
PostSubmitDiv.className = "newPostSubmitButtonDiv"
const PostSubmit = document.createElement("button");
PostSubmit.className = "newPostSubmitButton"
PostSubmit.textContent = "Post";
PostSubmit.setAttribute("type", "submit");
PostSubmitDiv.append(PostSubmit);

PostForm.append(titleLabelDiv, titleInputDiv, CatDiv, CatOptionDiv, contLabelDiv, contInputDiv, PostSubmitDiv);
const commentHandler = function (e) {
    e.preventDefault();
    const payloadObj = {}
    const payloadObjCom = {}
    payloadObj["label"] = "Createcomment";
    payloadObj["postID"] = (parseInt(e.submitter.value) + 1) + ""
    payloadObj["userID"] = document.querySelector(".ProfileID").textContent
    payloadObjCom["comment"] = e.target[0].value
    let strCom = JSON.stringify(payloadObjCom)
    payloadObj["commentOfPost"] = strCom
    console.log("checking target", payloadObj)
    commentPostId = payloadObj.postID
    postSocket.send(JSON.stringify(payloadObj));
};
function CreateCommentForm(value) {
    const commentForm = document.createElement("form")
    commentForm.className = "creatingCommentForm"
    commentForm.setAttribute("target", "_self")
    commentForm.addEventListener("submit", commentHandler);
    const commentLabelDiv = document.createElement('div');
    commentLabelDiv.className = "creatingCommentRespondDiv"
    const commentLabel = document.createElement('label');
    commentLabel.className = "creatingCommentRespondLabel"
    commentLabel.textContent = "create a comment:";
    commentLabel.setAttribute("for", "comment");
    commentLabelDiv.append(commentLabel);
    const commentInputDiv = document.createElement('div');
    commentInputDiv.className = "commentInputDiv"
    const commentInput = document.createElement('textarea');
    commentInput.className = "commentInput"
    // commentInput.setAttribute("type", "text");
    commentInput.setAttribute("name", "comment");
    commentInput.setAttribute("placeholder", "type here...");
    commentInput.setAttribute("id", "comment");
    commentInputDiv.append(commentInput);
    const commentSubmitDiv = document.createElement('div');
    commentSubmitDiv.className = "commentSubmitButtonDiv"
    const commentSubmit = document.createElement("button");
    commentSubmit.className = "commentSubmitButton"
    commentSubmit.textContent = "Submit Comment";
    commentSubmit.setAttribute("type", "submit");
    commentSubmit.setAttribute("value", value)
    commentSubmitDiv.append(commentSubmit);
    commentForm.append(commentLabelDiv, commentInputDiv, commentSubmitDiv)
    return commentForm
}

const showcommentHandler = function (e) {
    e.preventDefault();
    const payloadObj = {}
    payloadObj["label"] = "showComment";
    payloadObj["postID"] = (parseInt(e.submitter.value) + 1) + ""
    payloadObj["commentOfPost"] = jsonFile[e.submitter.value].postinfo.commentOfPost
    postSocket.send(JSON.stringify(payloadObj));

};
function CreateComments(arr, value) {
    console.log("CREATING COMMENTS")
    document.querySelectorAll("#allComments").forEach(e => {
        e.remove();
    });

    if (arr[value].postinfo.commentOfPost === "null") {
        console.log("comment of post empty")
        return ""

    } else {
        let comJson = JSON.parse(arr[value].postinfo.commentOfPost)
        const allComments = document.createElement("div")
        allComments.id = "allComments"
        for (let i = 0; i < comJson.length; i++) {
            console.log("cominfo", comJson[i].comInfo)
            const comDiv = document.createElement("div")
            const comContentDiv = document.createElement("div");
            const comUserIdDiv = document.createElement("div");
            comDiv.id = `comment-${i}`;
            comContentDiv.id = `comment-${i}`;
            comUserIdDiv.id = `userId-${i}`;
            comDiv.className = `singleComment`;
            comContentDiv.className = `singleCommentContent`;
            comUserIdDiv.className = `singleCommentAuthor`;
            const commentText = document.createElement("p")
            const comUserIdText = document.createElement("p")
            let commenTextNode = document.createTextNode(comJson[i].comInfo.comment)
            // let commenTextNode = document.createTextNode(" comment")
            // let comUserIdtextNode = document.createTextNode(comJson[i].comInfo.userID)
            let comUserIdtextNode = document.createTextNode(comJson[i].comInfo.userID)
            commentText.appendChild(commenTextNode)
            comUserIdText.appendChild(comUserIdtextNode)
            comContentDiv.append(commentText)
            comUserIdDiv.append(comUserIdText)
            comDiv.append(comContentDiv, comUserIdDiv)
            allComments.append(comDiv)
        }
        return allComments
    }
}

// }
export default { PostForm, DisplayPost };

