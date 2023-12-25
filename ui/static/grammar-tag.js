import { setOffline } from "./base.js";
const html = String.raw;

export class GrammarTag extends HTMLElement {
  constructor() {
    super();
    this.tags = [];
  }

  connectedCallback() {
    this._initData();
    this._initHTML();
    window.save = this.save.bind(this);
    if (!this.data.disabled) {
      this.HTML.input.addEventListener("keydown", this._keydownInput);
      this.HTML.input.addEventListener("input", this._input);
      this.HTML.btn_expand.addEventListener("click", this.expandSelection.bind(this));
      this.HTML.btn_reduce.addEventListener("click", this.shrinkSelection.bind(this));
    }
    document.body.addEventListener("keydown", this._keydownDoc);
    this.HTML.input.addEventListener("keydown", this._keydownDoc);
    this.HTML.btn_next.addEventListener("click", this.selectNext.bind(this));
    this.HTML.btn_prev.addEventListener("click", this.selectPrev.bind(this));
    this.render();
  }

  disconnectedCallback() {
    window.save = () => true;
    document.body.removeEventListener("keydown", this._keydownDoc);
    this.HTML.input.removeEventListener("keydown", this._keydownDoc);
    this.HTML.input.removeEventListener("keydown", this._keydownInput);
    this.HTML.input.removeEventListener("input", this._input);
  }

  selectPrev() {
    this.data.selectedIndex =
      Math.max(this.data.selectedIndex - 1, 0);
    if (this.data.words[this.data.selectedIndex].punctuation) {
      if (this.data.selectedIndex === 0) {
        return this.selectNext();
      }
      return this.selectPrev();
    }
    this.render();
  }

  selectNext() {
    this.data.selectedIndex =
      Math.min(this.data.selectedIndex + 1, this.data.words.length - 1);
    if (this.data.words[this.data.selectedIndex].punctuation) {
      if (this.data.selectedIndex === this.data.words.length - 1) {
        return this.selectPrev();
      }
      return this.selectNext();
    }
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
      nextWordObj.preceding = wordObj.preceding;
      wordObj.preceding = false;
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
      wordObj.preceding = nextWordObj.preceding;
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
      if (!this.data.words[i].shrinked && !this.data.words[i].preceding) {
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
    span.classList.add("decoration-blue-400", "underline-offset-[.7em]", "pe-0.5");
    span.innerText = word;
    if (!this.data.words[index].punctuation) {
      if (index === this.data.selectedIndex) {
        span.classList.add("bg-red-800", "text-white");
      } else {
        span.classList.add("hover:text-red-800", "cursor-pointer");
        span.addEventListener("click", () => {
          this.data.selectedIndex = index;
          this.render();
        })
      }
    }
    if (this.data.words[index].tags.length) {
      span.classList.add("underline");
    }
    frag.appendChild(span);
    return frag;
  }

  tagPartial(i, tagValue) {
    const frag = document.createDocumentFragment();
    const li = document.createElement("li");
    li.innerText = tagValue;
    if (!this.data.disabled) {
      li.classList.add("hover:line-through", "cursor-pointer", "max-w-fit",
        "decoration-red-500");
      li.addEventListener("click", this._clickTag);
    }
    li.setAttribute("data-i", i);
    frag.appendChild(li);
    return frag;
  }

  autocompletePartial(text, boldStart, boldEnd) {
    const p = document.createElement("p");
    p.classList.add("group", "hover:bg-sky-400", "hover:text-white",
      "cursor-pointer", "bg-white", "text-3xl", "ps-2", "leading-loose", "border");
    p.setAttribute("data-autocomplete", "");
    p.append(text.substring(0, boldStart));
    const strong = document.createElement("strong");
    strong.classList.add("text-sky-400", "group-hover:text-white");
    strong.innerText = text.substring(boldStart, boldEnd);
    p.append(strong);
    p.append(text.substring(boldEnd));
    p.addEventListener("click", this._clickAutocomplete);
    return p;
  }

  selectAutocomplete(index, smooth) {
    const item = this.HTML.autocomplete.children[index];
    item.classList.replace("bg-white", "bg-sky-500");
    item.classList.add("text-white");
    const strong = item.querySelector("strong");
    strong.classList.replace("text-sky-400", "text-white");
    item.scrollIntoView({
      behavior: smooth ? "smooth" : "instant",
      block: "center",
    });
  }

  unselectAutocomplete(index) {
    const item = this.HTML.autocomplete.children[index];
    item.classList.replace("bg-sky-500", "bg-white");
    item.classList.remove("text-white");
    const strong = item.querySelector("strong");
    strong.classList.replace("text-white", "text-sky-400");
  }

  addTag(index, tag) {
    if (this.tags.indexOf(tag) === -1) {
      alert(`Tag ${tag} is not supported!`);
      return;
    }
    const wordObj = this.data.words[index];
    console.log("Duplicate found");
    if (wordObj.tags.indexOf(tag) !== -1) {
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
      punctuation: false,
      preceding: false,
    };
  }

  _initData() {
    this.data = {
      words: [],
      id: undefined,
      selectedIndex: undefined,
      autocomplete: [],
      autocomplete_selected: -1,
      disabled: false,
    }

    if (this.getAttribute("disabled") != null) {
      this.data.disabled = true;
    }

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
    if (this.getAttribute("shared") != null) {
      this.data.shared = true;
    } else {
      this.data.shared = false;
    }
    this.data.selectedIndex = 0;

    const t_tags = document.body.querySelector("#template-tags");
    for (let p of t_tags.content.children) {
      this.tags.push(p.innerText);
    }
  }

  _initHTML() {
    this.innerHTML = html`
      <div dir="rtl" class="py-10 px-2 h-full grid grid-rows-1 grid-cols-[400px_auto]">
        <div class="border-e-2">
          <h2 class="text-3xl text-center font-sans">Navigation</h2>
          <div dir="ltr" class="font-sans grid grid-rows-2 grid-cols-2 gap-2 p-2">
            <button type="button" id="btn_next" class="bg-sky-600 text-white rounded-lg p-2 group">
              <p class="material-symbols-outlined">arrow_back</p>
              <p class="text-sm text-center">Shortcut: <kbd class="bg-slate-100 shadow-key text-black font-black border-1 ps-1 pe-1">-</kbd></p>
            </button>
            <button type="button" id="btn_prev" class="bg-sky-600 text-white rounded-lg p-2">
              <p class="material-symbols-outlined">arrow_forward</p>
              <p class="text-sm text-center">Shortcut: <kbd class="bg-slate-100 shadow-key text-black font-black border-1 ps-1 pe-1">=</kbd></p>
            </button>
            <button type="button" id="btn_expand" class="bg-sky-600 text-white rounded-lg p-2">
              <p class="material-symbols-outlined">text_select_move_back_word</p>
              <p class="text-sm text-center">Shortcut: <kbd class="bg-slate-100 shadow-key text-black font-black border-1 ps-1 pe-1">_</kbd></p>
            </button>
            <button type="button" id="btn_reduce" class="bg-sky-600 text-white rounded-lg p-2">
              <p class="material-symbols-outlined">text_select_move_forward_word</p>
              <p class="text-sm text-center">Shortcut: <kbd class="bg-slate-100 shadow-key text-black font-black border-1 ps-1 pe-1">+</kbd></p>
            </button>
          </div>
          <h2 dir="ltr" class="text-3xl text-center font-sans">Grammatical Tags</h2>
          <ul id="tag-container" class="text-3xl list-disc marker:text-green-600 list-inside leading-loose">
          </ul>
        </div>
        <div class="ms-8">
          <p class="ps-3 pe-3 text-5xl leading-loose"></p>
          <div class="pt-10 flex flex-col gap-5 w-1/2">
            <div class="flex flex-col">
              <input autofocus placeholder="اكتب..." type="text" class="w-full text-3xl ps-2 py-2 leading-loose drop-shadow rounded-lg"></input>
              <div class="relative w-full">
                <div id="autocomplete" class="absolute left-0 right-0 top-0 select-none max-h-72 overflow-y-auto">
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>`;

    this.HTML = {
      root: this.querySelector("div"),
      p: this.querySelector("div > p"),
      input: this.querySelector("input"),
      tagContainer: this.querySelector("#tag-container"),
      autocomplete: this.querySelector("#autocomplete"),
      btn_next: this.querySelector("#btn_next"),
      btn_prev: this.querySelector("#btn_prev"),
      btn_expand: this.querySelector("#btn_expand"),
      btn_reduce: this.querySelector("#btn_reduce"),
    }

    if (this.data.disabled) {
      this.HTML.input.style.display = "none";
      this.HTML.btn_expand.style.display = "none";
      this.HTML.btn_reduce.style.display = "none";
    }
  }

  save() {
    function forceSave() {
      htmx.ajax("PUT", `/excerpt/grammar?id=${this.data.id}`, {
        swap: "none",
        values: {
          content: { words: this.data.words },
          share: this.data.shared,
        }
      }).then(() => setOffline(false), () => setOffline(true));
    }
    if (!this.save.timeout || Date.now() - this.save.date > 500) {
      this.save.timeout = setTimeout(forceSave.bind(this));
      this.save.date = Date.now();
    } else if (Date.now() - this.save.date < 500) {
      clearTimeout(this.save.timeout);
      this.save.timeout = undefined;
      this.save.date = undefined;
      this.save();
    }
    return false;
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
      if (this.data.autocomplete.length > 0) {
        this.selectAutocomplete(this.data.autocomplete_selected, false);
      }
    }
  }

  autoCompleteSelectBelow(repeated) {
    if (this.data.autocomplete_selected !== -1) {
      this.unselectAutocomplete(this.data.autocomplete_selected);
    }
    this.data.autocomplete_selected++;
    if (this.data.autocomplete_selected >= this.data.autocomplete.length) {
      this.data.autocomplete_selected = 0;
      this.selectAutocomplete(this.data.autocomplete_selected, false);
      return;
    }
    this.selectAutocomplete(this.data.autocomplete_selected, !repeated);
  }

  autoCompleteSelectAbove(repeated) {
    if (this.data.autocomplete_selected !== -1) {
      this.unselectAutocomplete(this.data.autocomplete_selected);
    }
    this.data.autocomplete_selected--;
    if (this.data.autocomplete_selected < 0) {
      this.data.autocomplete_selected = this.data.autocomplete.length - 1;
      this.selectAutocomplete(this.data.autocomplete_selected, false);
      return;
    }
    this.selectAutocomplete(this.data.autocomplete_selected, !repeated);
  }

  _input = (_e) => {
    this.renderAutocomplete();
  }

  _keydownInput = (e) => {
    const keys = ["ArrowDown", "ArrowUp", "Enter"];
    if (keys.indexOf(e.key) !== -1) {
      e.preventDefault();
      e.stopPropagation();
      switch (e.key) {
        case "ArrowDown":
          if (this.data.autocomplete.length > 0) {
            this.autoCompleteSelectBelow(e.repeat);
          }
          break;
        case "ArrowUp":
          if (this.data.autocomplete.length > 0) {
            this.autoCompleteSelectAbove(e.repeat);
          }
          break;
        case "Enter":
          if (this.data.autocomplete.length > 0) {
            this.addTag(this.data.selectedIndex,
              this.data.autocomplete[this.data.autocomplete_selected]);
            this.HTML.input.value = "";
            this.renderAutocomplete();
          }
          break;
        default:
          console.error("Should not happen");
          break;
      }
    }
  }

  _keydownDoc = (e) => {
    const keys = ["-", "=", "_", "+"];
    if (keys.indexOf(e.key) !== -1) {
      e.preventDefault();
      e.stopPropagation();
      switch (e.key) {
        case "-":
          this.selectNext();
          break;
        case "=":
          this.selectPrev();
          break;
        case "_":
          if (!this.data.disabled) {
            this.expandSelection();
          }
          break;
        case "+":
          if (!this.data.disabled) {
            this.shrinkSelection();
          }
          break;
        default:
          console.error("Should not happen");
          break;
      }
    }
  }

  _clickTag = (e) => {
    console.log("tag clicked");
    const index = e.target.getAttribute("data-i");
    const wordObj = this.data.words[this.data.selectedIndex];
    wordObj.tags.splice(index, 1);
    this.save();
    this.render();
  }

  _clickAutocomplete = (e) => {
    this.HTML.input.value = "";
    this.addTag(this.data.selectedIndex, e.currentTarget.innerText);
    this.renderAutocomplete();
  }
}
