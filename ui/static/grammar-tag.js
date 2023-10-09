const html = String.raw;

// TODO(Amr Ojjeh): Get the list for the tagssssssss
export class GrammarTag extends HTMLElement {
  constructor() {
    super();
    this.data = Object.create(null);
    this.data.value = "";
    this.data.words = [];
    this.data.selectedIndex = undefined;
  }

  connectedCallback() {
    this._initData();
    this._initHTML();
    document.body.addEventListener("keydown", this._keydown);
    this.HTML.input.addEventListener("keydown", this._keydown);
    this.HTML.input.addEventListener("keyup", this._keyup);
    this.render();
  }

  disconnectedCallback() {
    document.body.removeEventListener("keydown", this._keydown);
    this.HTML.input.removeEventListener("keydown", this._keydown);
    this.HTML.input.removeEventListener("keyup", this._keyup);
  }

  selectPrev() {
    this.data.selectedIndex =
      Math.max(this.data.selectedIndex - 1, 0);
    this.render();
  }

  selectNext() {
    this.data.selectedIndex =
      Math.min(this.data.selectedIndex + 1, this.data.words.length - 1);
    this.render();
  }

  shrinkSelection() {
    const wordObj = this.data.words[this.data.selectedIndex];
    if (wordObj.word.length === 1) {
      return;
    }
    const newWord = wordObj.word.substring(0, wordObj.word.length - 1);
    if (!wordObj.shrinked) {
      const nextWord = wordObj.word[wordObj.word.length - 1];
      const nextWordObj = this.genWord(nextWord);
      nextWordObj.leftOver = true;
      this.data.words.splice(this.data.selectedIndex + 1, 0, nextWordObj);
    } else {
      const nextWordObj = this.data.words[this.data.selectedIndex + 1];
      nextWordObj.word =
        wordObj.word[wordObj.word.length - 1] + nextWordObj.word;
    }
    wordObj.word = newWord;
    wordObj.shrinked = true;
    this.render();
  }

  expandSelection() {
    const wordObj = this.data.words[this.data.selectedIndex];
    if (!wordObj.shrinked) {
      return;
    }
    const nextWordObj = this.data.words[this.data.selectedIndex + 1];
    const newWord = wordObj.word + nextWordObj.word[0];
    const nextWord = nextWordObj.word.substring(1);
    if (nextWord === "") {
      this.data.words.splice(this.data.selectedIndex + 1, 1);
      const newNextWord = this.data.words[this.data.selectedIndex + 1];
      if (!newNextWord || (newNextWord && !newNextWord.leftOver)) {
        this.data.words[this.data.selectedIndex].shrinked = false;
      }
    } else {
      nextWordObj.word = nextWord;
    }
    wordObj.word = newWord;
    this.render();
  }

  render() {
    this.HTML.p.innerHTML = "";
    for (let i = 0; i < this.data.words.length; ++i) {
      this.HTML.p.appendChild(this.wordPartial(i));
      if (!this.data.words[i].shrinked) {
        this.HTML.p.append(" ");
      }
    }

    const wordObj = this.data.words[this.data.selectedIndex];
    this.HTML.tagContainer.innerHTML = "";
    for (let i = 0; i < wordObj.tags.length; ++i) {
      const tag = wordObj.tags[i];
      const partial = this.tagPartial(tag);
      this.HTML.tagContainer.appendChild(partial)
    }
  }

  wordPartial(index) {
    const frag = document.createDocumentFragment();
    const word = this.data.words[index].word;
    const span = document.createElement("span");
    span.innerText = word;
    if (index === this.data.selectedIndex) {
      span.classList.add("bg-red-800", "text-white");
    }
    frag.appendChild(span);
    return frag;
  }

  tagPartial(tagValue) {
    const frag = document.createDocumentFragment();
    const li = document.createElement("li");
    li.innerText = tagValue;
    frag.appendChild(li);
    return frag;
  }

  addTag(index, tag) {
    const wordObj = this.data.words[index];
    wordObj.tags.push(tag);
    this.render();
  }

  genWord(word) {
    return {
      word,
      shrinked: false,
      leftOver: false,
      tags: [],
    };
  }

  _initData() {
    if (this.getAttribute("value")) {
      this.data.value = this.getAttribute("value");
    }
    this.data.words = this.data.value.split(" ").map(this.genWord);
    this.data.selectedIndex = 0;
  }

  _initHTML() {
    this.innerHTML = html`
      <div dir="rtl" class="py-10 px-2 h-full">
        <p class="ps-3 pe-3 text-4xl leading-loose">${this.data.value}</p>
        <div class="pt-10 flex flex-col gap-5 mx-auto w-1/2">
          <input autofocus placeholder="اكتب..." type="text" class="text-3xl ps-2 py-2 leading-loose drop-shadow"></input>
          <ul id="tag-container" class="text-3xl list-disc marker:text-red-800 list-inside leading-loose">
          </ul>
        </div>
      </div>`;

    this.HTML = Object.create(null);
    this.HTML.root = this.querySelector("div");
    this.HTML.p = this.querySelector("div > p");
    this.HTML.input = this.querySelector("input");
    this.HTML.tagContainer = this.querySelector("#tag-container");
  }

  _keyup = (e) => {
    if (e.key === "Enter") {
      this.addTag(this.data.selectedIndex, this.HTML.input.value);
      this.HTML.input.value = "";
    }
  }
  _keydown = (e) => {
    const keys = ["Home", "End"];
    if (keys.indexOf(e.key) !== -1) {
      e.preventDefault();
      e.stopPropagation();
      switch (e.key) {
        case "Home":
          if (e.shiftKey) {
            this.expandSelection();
            break;
          }
          this.selectNext();
          break;
        case "End":
          if (e.shiftKey) {
            this.shrinkSelection();
            break;
          }
          this.selectPrev();
          break;
        default:
          console.error("Should not happen");
          break;
      }
    }
  }
}