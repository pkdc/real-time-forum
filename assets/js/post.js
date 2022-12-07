let postSocket = null;
const body = document.getElementsByTagName("BODY")[0]
const DisplayPost = document.createElement("div")
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
        console.log({ resp });
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
            console.log("label is now comment----------------------")
            jsonFile = JSON.parse(resp.content)
            console.log("thats json",jsonFile, "and thats resp.content", resp.content)
            if (resp.Content !== null) {
                const comment = document.querySelector(".comment")
              let newCom= CreateComments(jsonFile,parseInt(commentPostId)-1)
            //    comment.append(newCom)
            comment.insertBefore(newCom, comment.children[1])
               console.log("function now over")
            }
        }else if (resp.label=== "showComment"){
            console.log("----------------------------------------- label show")
            jsonFile= JSON.parse(resp.content)
            console.log("and new Json", jsonFile)
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
        const titleForm = document.createElement("form")
        titleForm.addEventListener("submit", showcommentHandler)
        titleButton.setAttribute("value", i)
        titleButton.setAttribute("type", "submit")
        titleButton.addEventListener("click", function (e) {
            const comment = document.querySelector(".comment")
            comment.style.height = "%100"
            let chosenPost = document.querySelector(`#post-${titleButton.value}`)
            console.log("post is choosen")
            let clone = chosenPost.cloneNode(true)
            const closeComments = document.createElement("button")
            closeComments.textContent = String.fromCodePoint(0x274C)
            closeComments.addEventListener("click", function () {
                comment.style.height = "0%"
                PostHandler
                while (comment.firstChild) {
                    comment.removeChild(comment.firstChild)
                }
            })
            console.log("functions working***")
            let comments = CreateComments(jsonFile, i)
            let comForm = CreateCommentForm(titleButton.value)
            comment.append(clone, comments, comForm, closeComments)
            comment.style.height = "100%";
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
        // const titleText = document.createElement("p")
        // titleText.style.fontWeight= "900"
        const contentText = document.createElement("p")
        const categoryText = document.createElement("p")
        categoryText.style.width = "10%"
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
        postDiv.append(titleForm, contentDiv, categoryDiv, userIdDiv)
        allPost.append(postDiv)
    }
    DisplayPost.appendChild(allPost)
}
const PostHandler = function (e) {
    e.preventDefault();
    const formFields = new FormData(e.target);
    console.log("form entries", formFields.entries())
    const payloadObj = Object.fromEntries(formFields.entries());
    payloadObj["label"] = "post";
    console.log("checking target", payloadObj)
    postSocket.send(JSON.stringify(payloadObj));
};

const PostForm = document.createElement("form");
PostForm.addEventListener("submit", PostHandler);


const titleLabelDiv = document.createElement('div');
const titleLabel = document.createElement('label');
titleLabel.textContent = "title";
titleLabel.setAttribute("for", "title");
titleLabelDiv.append(titleLabel);
const titleInputDiv = document.createElement('div');
const titleInput = document.createElement('input');
titleInput.setAttribute("type", "text");
titleInput.setAttribute("name", "title");
titleInput.setAttribute("id", "title");
titleInputDiv.append(titleInput);
//-------------------
const CatDiv = document.createElement('div');
const CatOptionDiv = document.createElement('select');
CatOptionDiv.setAttribute("name", "category_option")
const CatLabel = document.createElement("label");
CatLabel.textContent = "Please choose category";
CatLabel.setAttribute("for", "cat");
CatDiv.append(CatLabel);
const CatInputOpt1 = document.createElement("option");
const CatInputOpt2 = document.createElement("option");
const CatInputOpt3 = document.createElement("option");
const CatInputOpt4 = document.createElement("option");
CatInputOpt1.setAttribute("name", "category_option");
CatInputOpt2.setAttribute("name", "category_option");
CatInputOpt3.setAttribute("name", "category_option");
CatInputOpt4.setAttribute("name", "category_option");
CatInputOpt1.setAttribute("id", "1");
CatInputOpt2.setAttribute("id", "2");
CatInputOpt3.setAttribute("id", "3");
CatInputOpt4.setAttribute("id", "4");
CatInputOpt1.setAttribute("value", "Anthony");
CatInputOpt2.setAttribute("value", "Burak");
CatInputOpt3.setAttribute("value", "David");
CatInputOpt4.setAttribute("value", "Godfrey");
CatInputOpt1.textContent = "Anthony";
CatInputOpt2.textContent = "Burak";
CatInputOpt3.textContent = "David";
CatInputOpt4.textContent = "Godfrey";
CatOptionDiv.setAttribute("id", "category");
CatOptionDiv.append(CatInputOpt1,CatInputOpt2,CatInputOpt3,CatInputOpt4)
const contLabelDiv = document.createElement('div');
const contLabel = document.createElement('label');
contLabel.textContent = "content:";
contLabel.setAttribute("for", "content");
contLabelDiv.append(contLabel);
const contInputDiv = document.createElement('div');
const contInput = document.createElement('input');
contInput.setAttribute("type", "text");
contInput.setAttribute("name", "content");
contInput.setAttribute("id", "content");
contInputDiv.append(contInput);

const PostSubmitDiv = document.createElement('div');
const PostSubmit = document.createElement("button");
PostSubmit.textContent = "Post";
PostSubmit.setAttribute("type", "submit");
PostSubmitDiv.append(PostSubmit);

PostForm.append(titleLabelDiv, titleInputDiv, CatDiv, CatOptionDiv, contLabelDiv, contInputDiv, PostSubmitDiv);
const commentHandler = function (e) {
    e.preventDefault();
    console.log("checking e", e.target[0])
    // const formFields = new FormData(e.target);
    // console.log("checking form", formFields)
    // const payloadObj = Object.fromEntries((e.target[0].value).entries());
    // const payloadObj =Object.create(Object.prototype)
    const payloadObj = Object.create(Object.prototype)
    const payloadObjCom = Object.create(Object.prototype)
    payloadObj["label"] = "Createcomment";
    payloadObj["postID"] = (parseInt(e.submitter.value) + 1) + ""
    payloadObjCom["comment"] = e.target[0].value
    let strCom = JSON.stringify(payloadObjCom)
    payloadObj["commentOfPost"] = strCom
    console.log("checking target", payloadObj)
    commentPostId= payloadObj.postID
    postSocket.send(JSON.stringify(payloadObj));
};
function CreateCommentForm(value) {
    const commentForm = document.createElement("form")
    commentForm.setAttribute("target", "_self")
    commentForm.addEventListener("submit", commentHandler);
    const commentLabelDiv = document.createElement('div');
    const commentLabel = document.createElement('label');
    commentLabel.textContent = "create a comment:";
    commentLabel.setAttribute("for", "comment");
    commentLabelDiv.append(commentLabel);
    const commentInputDiv = document.createElement('div');
    const commentInput = document.createElement('input');
    commentInput.setAttribute("type", "text");
    commentInput.setAttribute("name", "comment");
    commentInput.setAttribute("placeholder", "type here...");
    commentInput.setAttribute("id", "comment");
    commentInputDiv.append(commentInput);
    const commentSubmitDiv = document.createElement('div');
    const commentSubmit = document.createElement("button");
    commentSubmit.textContent = "comment";
    commentSubmit.setAttribute("type", "submit");
    commentSubmit.setAttribute("value", value)
    commentSubmitDiv.append(commentSubmit);
    commentForm.append(commentLabelDiv, commentInputDiv, commentSubmitDiv)
    return commentForm
}

const showcommentHandler = function (e) {
    e.preventDefault();
    // const formFields = new FormData(e.target);
    // const payloadObj = Object.fromEntries(formFields.entries());
    const payloadObj = Object.create(Object.prototype)
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

    console.log(value, "func check", arr[value])
    console.log(arr)
    if (arr[value].postinfo.commentOfPost === "null") {
        console.log("comment of post empty")
        return ""

    } else {
        console.log()
        let comJson = JSON.parse(arr[value].postinfo.commentOfPost)
        const allComments = document.createElement("div")
        allComments.id = "allComments"
        for (let i = 0; i < comJson.length; i++) {
            console.log("*****************************",comJson[i].comInfo.comment)
            const comDiv = document.createElement("div")
            const comContentDiv = document.createElement("div");
            const comUserIdDiv = document.createElement("div");
            comDiv.id = `comment-${i}`;
            comContentDiv.id = `comment-${i}`;
            comUserIdDiv.id = `userId-${i}`;
            const commentText = document.createElement("p")
            const comUserIdText = document.createElement("p")
            let commenTextNode = document.createTextNode(comJson[i].comInfo.comment)
            // let commenTextNode = document.createTextNode(" comment")
            // let comUserIdtextNode = document.createTextNode(comJson[i].comInfo.userID)
            let comUserIdtextNode = document.createTextNode(" userID ")
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
export default {PostForm,DisplayPost};

