function addSentence() {
    len.value++;
    const sentence = document.createElement("div");
    sentence.classList.add("sentence");
    sentence.setAttribute("lang", "ar");
    sentence.innerHTML += `
        <input name=${len.value} type="text" lang="ar">
        <span class="close material-symbols-outlined">cancel</span>`
    const closeButton = sentence.querySelector(".close");
    closeButton.addEventListener("click", delSentence.bind(sentence));

    sentencesEl.appendChild(sentence);
}

function delSentence() {
    if (sentencesEl.children.length > 1) {
        sentencesEl.removeChild(this);
        len.value = 0;
        for (let child of sentencesEl.children) {
            len.value++;
            child.querySelector("input").name = len.value;
        }
    } else {
        this.querySelector("input").value = "";
    }
}

const sentencesEl = document.querySelector("#sentences");
const addButtonEl = document.querySelector("#add-sentences");
const len = document.querySelector("#len");

addButtonEl.addEventListener("click", addSentence);
addSentence()
