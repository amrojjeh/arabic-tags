// TODO(Amr Ojjeh): Validate on the backend as well

// Used for syntax highlighting
const html = String.raw;
let debug = false;

if (debug) {
  var log = console.log;
} else {
  var log = () => { };
}

export class ArabicInput extends HTMLElement {

  constructor() {
    super();

    this.innerHTML = this.initHTML();
    this.HTML = Object.create(null);
    this.HTML.textarea = this.querySelector("textarea");
    this.HTML.highlighted = this.querySelector("div > div > div");

    this.partials = Object.create(null);
    this.partials.highlighted = document.createDocumentFragment();

    this._okays = undefined;
  }

  connectedCallback() {
    this.HTML.textarea.tabindex = "0";
    this.HTML.textarea.addEventListener("keydown", this._filter);
    this.HTML.textarea.addEventListener("input", this._input);
    this.HTML.textarea.addEventListener("scroll", this._scroll);
    this.HTML.textarea.addEventListener("paste", this._paste);
    if (this.getAttribute("autofocus") != null) {
      this.HTML.textarea.focus();
    }
    if (this.getAttribute("value") != null) {
      this.HTML.textarea.value = this.getAttribute("value");
      this.update();
    }
    if (this.getAttribute("id") != null) {
      this.HTML.textarea.setAttribute("hx-put",
        `/excerpt/edit?id=${this.getAttribute("id")}`);
    }
  }

  disconnectedCallback() {
    this.HTML.textarea.tabindex = "-1";
    this.HTML.textarea.removeEventListener("keydown", this._filter);
    this.HTML.textarea.removeEventListener("input", this._input);
    this.HTML.textarea.removeEventListener("paste", this._paste);
    this.HTML.textarea.removeEventListener("scroll", this._scroll);
  }

  initHTML() {
    return html`
      <div dir="rtl" class="h-full py-10 px-2">
        <div class="relative h-full">
          <div class="text-clip overflow-y-auto leading-loose absolute break-words top-0 left-0 h-full w-full text-3xl"></div>
          <textarea spellcheck="false"
            name="content"
            hx-swap="none"
            hx-put="NOT DEFINED"
            hx-trigger="keyup changed delay:500ms"
            hx-indicator=".htmx-indicator"
            class="leading-loose absolute focus:outline-none top-0 left-0 caret-black
            text-transparent bg-transparent h-full w-full text-3xl resize-none"
            placeholder="اكتب..."></textarea>
        </div>
      </div>`;
  }

  forceSave() {
    htmx.ajax("PUT", this.HTML.textarea.getAttribute("hx-put"), {
      swap: "none",
      values: { "content": this.HTML.textarea.value },
    });
  }

  render() {
    const start = Date.now();
    this.HTML.highlighted.innerHTML = "";
    this.HTML.highlighted.appendChild(this.partials.highlighted);
    log("render:", Date.now() - start, "milliseconds");
  }

  createSpan(okay) {
    const span = document.createElement("span");
    if (okay === -1) {
      console.error("Okay should never be -1");
      return;
    }

    if (okay === 1) {
      span.className = "bg-yellow-200 text-yellow-800";
    } else if (okay === 0) {
      span.className = "bg-red-200 text-red-800";
    }
    return span;
  }

  update() {
    const start = Date.now();
    this.deleteDoubleSpaces();
    const text = this.HTML.textarea.value;
    log("Deleted spaces:", Date.now() - start, "milliseconds");

    if (text.length > 0) {
      this._parse(text);
      const frag = this.partials.highlighted = document.createDocumentFragment();
      const okays = this.getOkays();
      let statusQuo = okays[0];
      let span = this.createSpan(statusQuo);
      for (let x = 0; x < text.length; ++x) {
        if (okays[x] === statusQuo) {
          span.innerText += text[x];
        } else {
          frag.appendChild(span);
          statusQuo = okays[x];
          span = this.createSpan(statusQuo);
          span.innerText = text[x];
        }
      }
      frag.appendChild(span)
    }
    log("Fragment updated:", Date.now() - start, "milliseconds");

    this.render();
    const event = new Event("arabic-input-update");
    this.dispatchEvent(event);
    log("Update:", Date.now() - start, "milliseconds");
  }

  hasErrors() {
    const okays = this.getOkays();
    for (let x = 0; x < this.HTML.textarea.value.length; ++x) {
      if (okays[x] === -1) {
        return false;
      }
      if (okays[x] !== 2) {
        return true;
      }
    }
    return false;
  }

  hasTashkeel() {
    const text = this.HTML.textarea.value;
    for (let x = 0; x < text.length; ++x) {
      if (isTashkeel(text[x])) {
        return true;
      }
    }
    return false;
  }

  deleteErrors() {
    let text = "";
    for (let c of this.HTML.textarea.value) {
      if (isValid(c)) {
        text += c;
      }
    }
    this.HTML.textarea.value = text;
    this.forceSave();
    this.update();
  }

  // Called from update, so no need to update
  deleteDoubleSpaces() {
    const text = this.HTML.textarea.value;
    this.HTML.textarea.value = text.replaceAll(/ {2,}/g, " ");
  }

  deleteVowels() {
    let text = "";
    for (let c of this.HTML.textarea.value) {
      if (!isTashkeel(c)) {
        text += c;
      }
    }
    this.HTML.textarea.value = text;
    this.forceSave();
    this.update();
  }

  scaleOkays(size = 1000) {
    // this is done to optimize for a packed SMI array
    if (this._okays == undefined) {
      this._okays = [];
    }
    for (let x = 0; x < size; x++) {
      if (!this._okays.hasOwnProperty(x)) {
        this._okays.push(-1);
      }
    }
    return this._okays;
  }

  getOkays() {
    if (this._okays == undefined) {
      this.scaleOkays();
    }
    return this._okays;
  }

  getOkaysSize(min) {
    if (min > this.getOkays().length) {
      return this.getOkays().length * 2;
    }
    return this.getOkays().length;
  }

  _parse(text) {
    // 0 = NOT OKAY
    // 1 = Tashkeel
    // 2 = OKAY
    const start = Date.now();
    if (this.getOkays().length <= text.length) {
      this.scaleOkays(this.getOkaysSize(text.length));
    }
    let okays = this.getOkays();
    for (let i = 0; i < text.length; ++i) {
      // Check for tashkeel
      const letter = text[i];
      if (!isTashkeel(letter)) {
        const pack = getLetterPack(text, i);
        if (pack.tashkeel.length > 0) {
          okays[i] = 1;
          for (let j = 0; j < pack.tashkeel.length; ++j) {
            okays[i + j + 1] = 1;
          }
          i += pack.tashkeel.length;
          continue;
        }
      }
      // Check for invalid letters
      const valid = isValid(letter);
      if (!valid) {
        okays[i] = 0;
      } else {
        okays[i] = 2;
      }
    }
    log("Parsing:", Date.now() - start, "milliseconds");
  }

  _paste = (e) => {
    e.preventDefault();
    let paste = e.clipboardData.getData("text");
    paste = paste.replaceAll("\n", " ");
    const selectionStart = e.target.selectionStart;
    const value = e.target.value;
    e.target.value = value.substring(0, selectionStart) + paste +
      value.substring(selectionStart + paste.length);
    this.update();
  }

  _filter = (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      return;
    }
    if (e.key === " ") {
      const content = e.target.value;
      const beforeCursor = content[e.target.selectionStart - 1];
      const afterCursor = content[e.target.selectionStart];
      if (beforeCursor === " ") {
        e.preventDefault();
        return;
      } else if (afterCursor === " ") {
        e.target.selectionStart++;
        e.preventDefault();
        return;
      }
    }
  }

  _input = (_e) => {
    this.HTML.highlighted.scrollTop = this.HTML.textarea.scrollTop;
    this.HTML.highlighted.scrollLeft = this.HTML.textarea.scrollLeft;
    this.update(true);
  }

  _scroll = (_e) => {
    this.HTML.highlighted.scrollTop = this.HTML.textarea.scrollTop;
    this.HTML.highlighted.scrollLeft = this.HTML.textarea.scrollLeft;
  }
}

class ArabicInputButton extends HTMLElement {
  constructor(bgColor, fgColor) {
    super();
    this.target = null;

    this.innerHTML = this.initHTML(this.className(bgColor, fgColor));
    this.HTML = Object.create(null);
    this.HTML.root = this.querySelector("button");

  }

  className(bgColor, fgColor) {
    return `${bgColor} capitalize ${fgColor} rounded-lg p-2`;
  }

  connectedCallback() {
    this.HTML.root.addEventListener("click", this._click);
    this.target = document.querySelector("arabic-input");
    this.target.addEventListener("arabic-input-update", this._update);
    this._update();
  }

  disconnectedCallback() {
    this.HTML.root.removeEventListener("click", this._click);
    this.target = null;
  }

  initHTML(className) {
    return html`
        <button type="button" class="${className}">DEFAULT BUTTON</button>
      `;
  }
}

export class DeleteErrorsButton extends ArabicInputButton {
  constructor() {
    super("bg-red-600", "text-white");
    this.HTML.root.innerText = "Delete all errors";
  }

  _click = (_e) => {
    if (confirm("Are you sure you want to delete all errors?")) {
      this.target.deleteErrors();
    }
  }

  _update = (_e) => {
    if (this.target.hasErrors()) {
      this.HTML.root.className = this.className("bg-red-600", "text-white");
      this.HTML.root.removeAttribute("disabled");
    } else {
      this.HTML.root.className = this.className("bg-gray-600", "text-white");
      this.HTML.root.setAttribute("disabled", "");
    }
  }
}

export class DeleteVowelsButton extends ArabicInputButton {
  constructor() {
    super("bg-yellow-600", "text-white");
    this.HTML.root.innerText = "Delete vowels";
  }

  _click = (_e) => {
    if (confirm("Are you sure you want to delete all vowels?")) {
      this.target.deleteVowels();
    }
  }

  _update = (_e) => {
    if (this.target.hasTashkeel()) {
      this.HTML.root.className = this.className("bg-yellow-600", "text-white");
      this.HTML.root.removeAttribute("disabled");
    } else {
      this.HTML.root.className = this.className("bg-gray-600", "text-white");
      this.HTML.root.setAttribute("disabled", "");
    }
  }
}

function getLetterPack(line, indexStart) {
  if (isTashkeel(line[indexStart])) {
    throw new Error("Line should not start with tashkeel!");
  }
  const pack = {
    letter: line[indexStart],
    tashkeel: "",
  };
  for (let i = indexStart + 1; i < line.length; ++i) {
    const char = line[i];
    if (isTashkeel(char)) {
      pack.tashkeel += char;
    } else {
      return pack;
    }
  }
  return pack;
}

function isTashkeel(char) {
  const code = char.codePointAt(0);
  return code >= 0x064B && code <= 0x065F;
}

function isArabicLetter(char) {
  const code = char.codePointAt(0);
  if (code >= 0x0621 && code <= 0x063A) {
    return true;
  }
  if (code >= 0x0641 && code <= 0x064A) {
    return true;
  }
  return false;
}

function isWhitespace(letter) {
  return letter == " ";
}

function isValid(letter) {
  return isArabicLetter(letter) || isWhitespace(letter);
}
