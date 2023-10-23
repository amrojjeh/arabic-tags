const html = String.raw;

export class GrammarTag extends HTMLElement {
  constructor() {
    super();
    this.data = {
      words: [],
      id: undefined,
      selectedIndex: undefined,
      autocomplete: [],
      autocomplete_selected: -1,
    }

    this.tags = [];
  }

  connectedCallback() {
    this._initData();
    this._initHTML();
    document.body.addEventListener("keydown", this._keydown);
    this.HTML.input.addEventListener("keydown", this._keydown);
    this.HTML.input.addEventListener("keyup", this._keyup);
    this.HTML.input.addEventListener("input", this._input);
    this.render();
  }

  disconnectedCallback() {
    document.body.removeEventListener("keydown", this._keydown);
    this.HTML.input.removeEventListener("keydown", this._keydown);
    this.HTML.input.removeEventListener("keyup", this._keyup);
    this.HTML.input.removeEventListener("input", this._input);
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
    this.save();
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
    this.save();
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
      const partial = this.tagPartial(i, tag);
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

  tagPartial(i, tagValue) {
    const frag = document.createDocumentFragment();
    const li = document.createElement("li");
    li.innerText = tagValue;
    li.classList.add("hover:line-through", "hover:cursor-pointer", "max-w-fit");
    li.addEventListener("click", this._clickTag);
    li.setAttribute("data-i", i);
    frag.appendChild(li);
    return frag;
  }

  autocompletePartial(text, boldStart, boldEnd) {
    const p = document.createElement("p");
    p.classList.add("group", "hover:bg-sky-500", "hover:text-white",
      "cursor-pointer", "bg-white", "text-3xl", "ps-2", "leading-loose", "border");
    p.append(text.substring(0, boldStart));
    const strong = document.createElement("strong");
    strong.classList.add("text-red-800", "group-hover:text-white");
    strong.innerText = text.substring(boldStart, boldEnd);
    p.append(strong);
    p.append(text.substring(boldEnd));
    p.addEventListener("click", this._clickAutocomplete);
    return p;
  }

  selectAutocomplete(index) {
    const item = this.HTML.autocomplete.children[index];
    item.classList.replace("bg-white", "bg-sky-500");
    item.classList.add("text-white");
    const strong = item.querySelector("strong");
    strong.classList.replace("text-red-800", "text-white");
  }

  unselectAutocomplete(index) {
    const item = this.HTML.autocomplete.children[index];
    item.classList.replace("bg-sky-500", "bg-white");
    item.classList.remove("text-white");
    const strong = item.querySelector("strong");
    strong.classList.replace("text-white", "text-red-800");
  }

  addTag(index, tag) {
    if (this.tags.indexOf(tag) === -1) {
      alert(`Tag ${tag} is not supported!`);
      return;
    }
    const wordObj = this.data.words[index];
    if (wordObj.tags.indexOf(tag) !== -1) {
      console.log("Duplicate found");
      return;
    }
    wordObj.tags.push(tag);
    this.save();
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
      this.data.words = JSON.parse(this.getAttribute("value")).words;
    } else {
      console.error("There's no value!");
    }

    if (this.getAttribute("id") != null) {
      this.data.id = this.getAttribute("id");
    } else {
      console.error("There's no id!");
    }
    this.data.selectedIndex = 0;

    const t_tags = document.body.querySelector("#template-tags");
    for (let p of t_tags.content.children) {
      this.tags.push(p.innerText);
    }
  }

  _initHTML() {
    this.innerHTML = html`
      <div dir="rtl" class="py-10 px-2 h-full">
        <p class="ps-3 pe-3 text-4xl leading-loose"></p>
        <div class="pt-10 flex flex-col gap-5 mx-auto w-1/2">
          <div class="flex flex-col">
            <input autofocus placeholder="اكتب..." type="text" class="w-full text-3xl ps-2 py-2 leading-loose drop-shadow"></input>
            <div class="relative w-full">
              <div id="autocomplete" class="absolute left-0 right-0 top-0 select-none max-h-72 overflow-y-auto">
              </div>
            </div>
          </div>
          <ul id="tag-container" class="text-3xl list-disc marker:text-red-800 list-inside leading-loose">
          </ul>
        </div>
      </div>`;

    this.HTML = {
      root: this.querySelector("div"),
      p: this.querySelector("div > p"),
      input: this.querySelector("input"),
      tagContainer: this.querySelector("#tag-container"),
      autocomplete: this.querySelector("#autocomplete"),
    }
  }

  save() {
    if (!this.save.timeout || Date.now() - this.save.date > 500) {
      this.save.timeout = setTimeout(htmx.ajax.bind(this), 500, "PUT",
        `/excerpt/grammar?id=${this.data.id}`, {
        swap: "none",
        values: { "content": { "words": this.data.words } }
      });
      this.save.date = Date.now();
    } else if (Date.now() - this.save.date < 500) {
      clearTimeout(this.save.timeout);
      this.save.timeout = undefined;
      this.save.date = undefined;
      this.save();
    }
  }

  renderAutocomplete() {
    this.data.autocomplete_selected = 0;
    this.HTML.autocomplete.innerHTML = "";
    if (this.HTML.input.value.length > 1) {
      this.data.autocomplete = this.tags.filter(x => x.indexOf(this.HTML.input.value) !== -1);
      for (let i = 0; i < this.data.autocomplete.length; ++i) {
        const l = this.data.autocomplete[i];
        const start = l.indexOf(this.HTML.input.value);
        const p = this.autocompletePartial(l, start,
          start + this.HTML.input.value.length,
          i === this.data.autocomplete_selected);
        this.HTML.autocomplete.append(p);
      }
      this.selectAutocomplete(this.data.autocomplete_selected);
    }
  }

  autoCompleteSelectBelow() {
    if (this.data.autocomplete_selected !== -1) {
      this.unselectAutocomplete(this.data.autocomplete_selected);
    }
    this.data.autocomplete_selected =
      (this.data.autocomplete_selected + 1) % this.data.autocomplete.length;
    this.selectAutocomplete(this.data.autocomplete_selected);
  }

  autoCompleteSelectAbove() {
    if (this.data.autocomplete_selected !== -1) {
      this.unselectAutocomplete(this.data.autocomplete_selected);
    }
    this.data.autocomplete_selected--;
    if (this.data.autocomplete_selected < 0) {
      this.data.autocomplete_selected = this.data.autocomplete.length - 1;
    }
    this.selectAutocomplete(this.data.autocomplete_selected);
  }

  _keyup = (e) => {
    if (e.key === "Enter") {
      if (this.data.autocomplete.length > 0) {
        if (this.addTag(this.data.selectedIndex, this.data.autocomplete[0])) {
          this.HTML.input.value = "";
        }
      }
    } else if (e.key === "ArrowDown") {
      this.autoCompleteSelectBelow();
    } else if (e.key === "ArrowUp") {
      this.autoCompleteSelectAbove();
    }
  }

  _input = (_e) => {
    this.renderAutocomplete();
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

  _clickTag = (e) => {
    const index = e.target.getAttribute("data-i");
    const wordObj = this.data.words[this.data.selectedIndex];
    wordObj.tags.splice(index, 1);
    this.save();
    this.render();
  }

  _clickAutocomplete = (e) => {
    this.HTML.input.value = "";
    this.addTag(this.data.selectedIndex, e.target.innerText);
    this.renderAutocomplete();
  }
}