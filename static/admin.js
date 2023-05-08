const postName = document.getElementById('PostName');
const postDesc = document.getElementById('PostDescription');
const postAuthor = document.getElementById('PostAuthor');
const postDate = document.getElementById('PostDate');

postName.addEventListener('keyup', function () {
  liveUpdate(postName, "PreviewPostName");
});

postDesc.addEventListener('keyup', function () {
  liveUpdate(postDesc, "PreviewPostDescription");
});

postAuthor.addEventListener('keyup', function () {
  liveUpdate(postAuthor, "PreviewPostAuthor");
});

postDate.addEventListener('keyup', function () {
  liveUpdate(postDate, "PreviewPostDate");
});

function liveUpdate(element, previewClass) {
  var x = element.value;
  var y = document.getElementsByClassName(previewClass);
  Array.prototype.forEach.call(y, el => {
    el.innerText = x;
  });
}

const uploadpreview1 = document.getElementById("hero-preview");
const uploadpreview2 = document.getElementById("hero-preview2");
const uploadpreview3 = document.getElementById("author-preview");

const articlepreview = document.getElementById("article-preview");
const cardpreview = document.getElementById("card-preview");
const cardauthorpreview = document.getElementById("card-author-preview");

const fileInput1 = document.getElementById("hero-input");
const fileInput2 = document.getElementById("hero-input2");
const fileInput3 = document.getElementById("author-input");

const uploadButton = document.getElementById('upload-button');
const firstHero = document.getElementById('first-hero');
const secondHero = document.getElementById('second-hero');

const firstHeroSave = firstHero.innerText;
const secondHeroSave = secondHero.innerText;
const uploadButtonSave = uploadButton.innerText;

const menuButtons = document.getElementById('menu-buttons');

fileInput1.addEventListener("change", function () {
  previewFile(uploadpreview1, articlepreview, this.files[0]);
  replace(firstHero, menuButtons);
});

fileInput2.addEventListener("change", function () {
  previewFile(uploadpreview2, cardpreview, this.files[0]);
  replace(secondHero, menuButtons);
});

fileInput3.addEventListener("change", function () {
  previewFile(uploadpreview3, cardauthorpreview, this.files[0]);
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
  if (temp.id === 'upload-button'){
    c = uploadButtonSave;
    img = 'upload_author-image.svg';
    type = cardauthorpreview;
  }
  if (temp.id === 'first-hero'){
    c = firstHeroSave;
    img = 'upload_hero-image.svg';
    type = articlepreview;
  }
  if (temp.id === 'second-hero'){
    c = secondHeroSave;
    img = 'upload_hero-image-2.svg';
    type = cardpreview;
  }

  document.getElementById(temp.id).textContent = c;
  var q = (temp.closest("div"));
  q.firstElementChild.src = '../static/images/' + img;
  type.src = '../static/images/default.svg';
}


document.body.addEventListener('click', function(e) {
  if(e.target.className === 'label-disable') {
    e.preventDefault();
  }
});

function validateForm(event) {
  // Prevent the form from submitting
  event.preventDefault();

  // Get all input elements with the "required" attribute
  const requiredInputs = document.querySelectorAll('input[required]');

  // Check if any required fields are empty
  for (let i = 0; i < requiredInputs.length; i++) {
    if (requiredInputs[i].value.trim() === '') {
      alert('Please fill in all required fields.');
      return false; // Prevent form submission
    }
  }

  // All fields are completed
  const warningButton = document.querySelector('.warning-button_good');
  warningButton.style.display = 'block'; // Make the warning button visible
  document.querySelector('.warning-button_wrong').style.display = 'none';

  // Optionally, you can reset the form after displaying the warning button
  // event.target.reset();

  return true; // Allow form submission
}

document.addEventListener("DOMContentLoaded", function() {
  var elements = document.getElementsByTagName("INPUT");
  for (var i = 0; i < elements.length; i++) {
      elements[i].oninvalid = function(e) {
          e.target.setCustomValidity("");
          if (!e.target.validity.valid) {
              e.target.setCustomValidity("This field cannot be left blank");
              document.querySelector('.warning-button_wrong').style.display = 'block';
          }
      };
      elements[i].oninput = function(e) {
          e.target.setCustomValidity("");
      };
  }
}) 

