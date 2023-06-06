window.addEventListener('load', () => {
    namesUpdate();
    picturesUpdate();
    resetMenu();
    formValidation();
})

function namesUpdate () {
    const postName = document.getElementById('PostName');
    const postDesc = document.getElementById('PostDescription');
    const postAuthor = document.getElementById('PostAuthor');
    const postDate = document.getElementById('PostDate');

    function liveUpdate(element, previewClass) {
        const elFrom = element.value;
        const elTo = document.querySelectorAll('.' + previewClass);
        elTo.forEach(el => {
            el.innerText = elFrom;
        });
    }

    postName.addEventListener('keyup', function () {
        liveUpdate(postName, "preview__h3");
    });

    postDesc.addEventListener('keyup', function () {
        liveUpdate(postDesc, "preview__description");
    });

    postAuthor.addEventListener('keyup', function () {
        liveUpdate(postAuthor, "preview__author");
    });

    postDate.addEventListener('change', function () {
        liveUpdate(postDate, "preview__date");
    });
}

function picturesUpdate(){
    const firstUploadPreview = document.getElementById("firstHeroPreview");
    const secondUploadPreview = document.getElementById("secondHeroPreview");
    const thirdUploadPreview = document.getElementById("authorPreview");

    const articlePreview = document.getElementById("articlePreview");
    const cardPreview = document.getElementById("cardPreview");
    const cardAuthorPreview = document.getElementById("cardAuthorPreview");

    const firstFileInput = document.getElementById("firstHeroInput");
    const secondFileInput = document.getElementById("secondHeroInput");
    const thirdFileInput = document.getElementById("authorInput");

    const uploadButton = document.getElementById('uploadButton');
    const firstHero = document.getElementById('firstHero');
    const secondHero = document.getElementById('secondHero');

    const menuButtons = document.getElementById('menuButtons');

    firstFileInput.addEventListener("change", function () {
        previewFile(firstUploadPreview, articlePreview, this.files[0]);
        replace(firstHero, menuButtons);
    });

    secondFileInput.addEventListener("change", function () {
        previewFile(secondUploadPreview, cardPreview, this.files[0]);
        replace(secondHero, menuButtons);
    });

    thirdFileInput.addEventListener("change", function () {
        previewFile(thirdUploadPreview, cardAuthorPreview, this.files[0]);
        replace(uploadButton, menuButtons);
    });

    function previewFile(preview, image, file) {
        const reader = new FileReader();

        const loadHandler = () => {
            preview.src = reader.result;
            image.src = reader.result;

            reader.removeEventListener("load", loadHandler);
        };

        reader.addEventListener("load", loadHandler, false);

        if (file) {
            reader.readAsDataURL(file);
        }
    }

    function replace(before, after) {
        before.innerHTML = after.innerHTML;
    }
}

function resetMenu(){
    const firstHeroSave = document.getElementById("firstHeroInput").innerText;
    const secondHeroSave = document.getElementById("secondHeroInput").innerText;
    const uploadButtonSave = document.getElementById("authorInput").innerText;

    document.body.addEventListener('click', function (a) {
        if (a.target.className === 'label-disable') {
            a.preventDefault();
            deleteFile(a.target, firstHeroSave, secondHeroSave, uploadButtonSave);
        }
    });

    function deleteFile(el, firstHeroSave, secondHeroSave, uploadButtonSave) {

        const articlePreview = document.getElementById("articlePreview");
        const cardPreview = document.getElementById("cardPreview");
        const cardAuthorPreview = document.getElementById("cardAuthorPreview");

        let temp = (el.closest("div > p"));
        let text, img, type;
        if (temp.id === 'uploadButton') {
            text = uploadButtonSave;
            img = 'upload_author-image.svg';
            type = cardAuthorPreview;
        }
        if (temp.id === 'firstHero') {
            text = firstHeroSave;
            img = 'upload_hero-image.svg';
            type = articlePreview;
        }
        if (temp.id === 'secondHero') {
            text = secondHeroSave;
            img = 'upload_hero-image-2.svg';
            type = cardPreview;
        }

        document.getElementById(temp.id).textContent = text;
        const find = (temp.closest("div"));
        find.firstElementChild.src = '/static/images/' + img;
        type.src = '/static/images/default.svg';
    }
}

function formValidation(){
    let elements = document.getElementsByTagName("INPUT");
    let textValidate = document.getElementsByClassName("main-box__input-warning");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("invalid", function (e) {
            if (!e.target.validity.valid) {
                e.preventDefault();
                document.querySelector('.warning-button_wrong').classList.add("warning-button_active");
                textValidate[i - 1].classList.add("main-box__input-warning_active");
                elements[i].classList.add("main-box__input__warning");
            }
        });

        elements[i].addEventListener("change", function (a) {
            if (!a.value) {
                textValidate[i - 1].classList.remove("main-box__input-warning_active");
                elements[i].classList.remove("main-box__input__warning");
            }
        });
    }

    addEventListener("click", function (a) {
        if (a.target.id === "FormButton" || a.target.id === "PostContent") {
            checkContent();
        }
    })

    addEventListener("submit", function (e){
        e.preventDefault()
        if (checkContent()){
            submitForm()
        }
    })

    function checkContent() {
        const textArea = document.getElementById("PostContent").value;
        const textWarning = document.getElementById("content-warning");
        if (textArea === '') {
            textWarning.classList.add("main-box__input-warning_active");
            return false
        } else {
            textWarning.classList.remove("main-box__input-warning_active");
            return true
        }
    }

    function submitForm() {

        const formData = new FormData(document.forms.form);

        const textarea = document.getElementById('PostContent');
        formData.append('content', textarea.value);

        const heroFirst = document.getElementById('firstHeroInput');
        let heroExt = heroFirst.files[0].name.split('.').pop();

        const authorImg = document.getElementById('authorInput');
        let authorExt = authorImg.files[0].name.split('.').pop();

        formData.append('hero_ext', '.' + heroExt);
        formData.append('author_ext', '.' + authorExt);


        let object = {};
        let imagePromises = [];

        formData.forEach(function (value, key) {
            if (value instanceof File) {
                let reader = new FileReader();
                let imagePromise = new Promise(function (resolve, reject) {
                    reader.onload = function (event) {
                        object[key] = event.target.result;
                        resolve();
                    };
                    reader.onerror = function (event) {
                        reject(event.error);
                    };
                });

                reader.readAsDataURL(value);
                imagePromises.push(imagePromise);
            } else {
                object[key] = value;
            }
        });

        Promise.all(imagePromises)
            .then(function () {
                const json = JSON.stringify(object);
                fetch("/api/post", {
                    method: 'post',
                    body: json,
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                    .then(function (response) {
                        console.log('Data sent to the server:', response);
                    })
                    .catch(function (error) {
                        console.error('An error occurred while sending the data:', error);
                    });

                console.log(json);
                document.querySelector('.warning-button_wrong').classList.remove("warning-button_active");
                document.querySelector('.warning-button_good').classList.add("warning-button_active");
            })
            .catch(function (error) {
                console.error('An error occurred while reading the image:', error);
            });


        return false;
    }
}


