const postName = document.getElementById('PostName');
const postDesc = document.getElementById('PostDescription');
const postAuthor = document.getElementById('PostAuthor');
const postDate = document.getElementById('PostDate');

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

function liveUpdate(element, previewClass) {
  var x = element.value;
  var y = document.getElementsByClassName(previewClass);
  Array.prototype.forEach.call(y, el => {
    el.innerText = x;
  });
}

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

const firstHeroSave = firstHero.innerText;
const secondHeroSave = secondHero.innerText;
const uploadButtonSave = uploadButton.innerText;

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

  reader.addEventListener(
    "load",
    () => {
      // convert image file to base64 string
      preview.src = reader.result;
      image.src = reader.result;
    },
    false
  );

  if (file) {
    reader.readAsDataURL(file);
  }
}

function replace(before, after) {
  before.innerHTML = after.innerHTML;
}

function replaceText(before, after) {
  before.innerHTML = after;
}

function deleteFile(el) {
  var temp = (el.closest("div > p"));
  console.log(temp.id);
  if (temp.id === 'uploadButton') {
    c = uploadButtonSave;
    img = 'upload_author-image.svg';
    type = cardAuthorPreview;
  }
  if (temp.id === 'firstHero') {
    c = firstHeroSave;
    img = 'upload_hero-image.svg';
    type = articlePreview;
  }
  if (temp.id === 'secondHero') {
    c = secondHeroSave;
    img = 'upload_hero-image-2.svg';
    type = cardPreview;
  }

  document.getElementById(temp.id).textContent = c;
  var q = (temp.closest("div"));
  q.firstElementChild.src = '../static/images/' + img;
  type.src = '../static/images/default.svg';
}


document.body.addEventListener('click', function (e) {
  if (e.target.className === 'label-disable') {
    e.preventDefault();
  }
});

function validateForm(event) {

  // Get all input elements with the "required" attribute
  const requiredInputs = document.querySelectorAll('input[required]');

  for (let i = 0; i < requiredInputs.length; i++) {
    if (requiredInputs[i].value.trim() === '') {
      alert('Please fill in all required fields.');
      return false; // Prevent form submission
    }
  }

  const warningButton = document.querySelector('.warning-button_good');
  warningButton.style.display = 'block'; // Make the warning button visible
  document.querySelector('.warning-button_wrong').style.display = 'none';

  var formData = new FormData(document.forms.post);

  var object = {};
  var imagePromises = [];

  formData.forEach(function (value, key) {
    if (value instanceof File) {
      var reader = new FileReader();
      var imagePromise = new Promise(function (resolve, reject) {
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
      var json = JSON.stringify(object);

      console.log(json);
    })
    .catch(function (error) {
      console.error('An error occurred while reading the image:', error);
    });


  return true; 
}

document.addEventListener("DOMContentLoaded", function () { //Check if fields are not empty
  var elements = document.getElementsByTagName("INPUT");
  for (var i = 0; i < elements.length; i++) {
    elements[i].oninvalid = function (e) {
      e.target.setCustomValidity("");
      if (!e.target.validity.valid) {
        e.target.setCustomValidity("This field cannot be left blank");
        document.querySelector('.warning-button_wrong').style.display = 'block';
      }
    };
    elements[i].oninput = function (e) {
      e.target.setCustomValidity("");
    };
  }
})


