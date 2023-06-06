function validateForm() {
    const emailValue = document.getElementById('email');
    const passwordValue = document.getElementById('password');
    const warningMessage = document.querySelector('.warning-button_wrong')
    let emailWarningVisibility = document.getElementById('email-req');
    let passwordWarningVisibility = document.getElementById('password-req');

    if (emailValue.value === '') {
        emailValue.style.borderBottomColor = "#E86961";
        emailWarningVisibility.innerHTML = 'Email is required';
        emailWarningVisibility.style.visibility = 'visible';
    }
    else {
        if (!isEmailValid(emailValue.value)) {
            emailWarningVisibility.innerHTML = 'Email is not correct';
            emailWarningVisibility.style.visibility = 'visible';
            emailValue.style.borderBottomColor = "#E86961";
        }
        else {
            emailWarningVisibility.style.visibility = 'hidden';
        }
    }

    if (passwordValue.value === '') {
        passwordWarningVisibility.style.visibility = 'visible';
        passwordValue.style.borderBottomColor = "#E86961";
    }

    if (emailValue.value === '' || passwordValue.value === '' || !isEmailValid(emailValue.value)) {
        warningMessage.classList.add('warning-button_active');
        return false;
    } else {
        warningMessage.classList.remove('warning-button_active');
    }

    const formData = new FormData(document.forms.login);
    // let object = {};
    // formData.forEach(function (value, key) {
    //     object[key] = value;
    // });
    const json = JSON.stringify(formData);
    console.log(json);
    return false
}

function checkField(element, elementID) {
    let elementId = document.getElementById(elementID);
    element.onchange = function (e) {
        if (!e.value) {
            e.preventDefault();
            elementId.style.visibility = 'hidden';
            element.style.backgroundColor = '#F7F7F7';
            element.style.borderBottomColor = "#2E2E2E"
        }
        if (element.value === '') {
            element.style.backgroundColor = '#FFFFFF';
        }
    }
}

function isEmailValid(email) {
    const emailRegexp = new RegExp(
        /^[a-zA-Z0-9][\-_\.\+\!\#\$\%\&\'\*\/\=\?\^\`\{\|]{0,1}([a-zA-Z0-9][\-_\.\+\!\#\$\%\&\'\*\/\=\?\^\`\{\|]{0,1})*[a-zA-Z0-9]@[a-zA-Z0-9][-\.]{0,1}([a-zA-Z][-\.]{0,1})*[a-zA-Z0-9]\.[a-zA-Z0-9]{1,}([\.\-]{0,1}[a-zA-Z]){0,}[a-zA-Z0-9]{0,}$/i
    )
    return (emailRegexp.test(email));
}


