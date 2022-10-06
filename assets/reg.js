const RegisterForm = document.createElement("form");

// login form
// name label
const RnameLabelDiv = document.createElement('div');
const RnameLabel = document.createElement('label');
RnameLabel.textContent = "Please Enter Your First Name :";
RnameLabel.setAttribute("for", "name");
RnameLabelDiv.append(RnameLabel);
// name input
const RnameInputDiv = document.createElement('div');
const RnameInput = document.createElement('input');
RnameInput.setAttribute("type", "text");
RnameInput.setAttribute("name", "name");
RnameInput.setAttribute("id", "name");
RnameInput.setAttribute("placeholder", "eg: Nick")
RnameInputDiv.append(RnameInput);
// Nickname
//  Age
//  Gender
//  First Name
//  Last Name
//  E-mail
//  Password
// pw label
const RpwLabelDiv = document.createElement('div');
const RpwLabel = document.createElement('label');
RpwLabel.textContent = "Please Enter Your Password:";
RpwLabel.setAttribute("for", "pw");
RpwLabelDiv.append(RpwLabel);
// password input
const RpwInputDiv = document.createElement('div');
const RpwInput = document.createElement('input');
RpwInput.setAttribute("type", "password");
RpwInput.setAttribute("name", "pw");
RpwInput.setAttribute("id", "pw");
RpwInputDiv.append(RpwInput);

RegisterForm.append(RnameLabelDiv, RnameInputDiv, RpwLabelDiv, RpwInputDiv);
export default RegisterForm;