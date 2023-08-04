function addSentence() {
    sentencesEl.innerHTML += `
    <div class="sentence">
        <span class="material-symbols-outlined">cancel</span>
        <input type="text" lang="ar">
    </div>`
}

const sentencesEl = document.querySelector("#sentences");
const addButtonEl = document.querySelector("#add-sentences");

addButtonEl.addEventListener("click", addSentence);
