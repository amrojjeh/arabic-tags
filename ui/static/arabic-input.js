// TODO(Amr Ojjeh): Validate on the backend as well

// Used for syntax highlighting
const html = String.raw;

export class ArabicInput extends HTMLElement {

  constructor() {
    super();

    this.innerHTML = this.initHTML();
    this.HTML = Object.create(null);
    this.HTML.textarea = this.querySelector("textarea");
    this.HTML.highlighted = this.querySelector("div > div > div");

    this.partials = Object.create(null);
    this.partials.highlighted = document.createDocumentFragment();

    this.lines = [];
  }

  connectedCallback() {
    this.HTML.textarea.tabindex = "0";
    this.HTML.textarea.addEventListener("keydown", this._filter);
    this.HTML.textarea.addEventListener("input", this._input);
    this.HTML.textarea.addEventListener("paste", this._paste);
    if (this.getAttribute("autofocus") != null) {
      this.HTML.textarea.focus();
    }
  }

  disconnectedCallback() {
    this.HTML.textarea.tabindex = "-1";
    this.HTML.textarea.removeEventListener("keydown", this._filter);
    this.HTML.textarea.removeEventListener("input", this._input);
    this.HTML.textarea.removeEventListener("paste", this._paste);
  }

  initHTML() {
    return html`
      <div dir="rtl" class="h-full py-10 px-2">
        <div class="relative h-full">
          <div class="absolute break-words top-0 left-0 h-full w-full text-2xl"></div>
          <textarea spellcheck="false" class="absolute focus:outline-none 
            top-0 left-0 caret-black text-transparent bg-transparent h-full
            w-full text-2xl resize-none" placeholder="اكتب..."></textarea>
        </div>
      </div>`;
  }

  render() {
    this.HTML.highlighted.innerHTML = "";
    this.HTML.highlighted.appendChild(this.partials.highlighted);
  }

  update() {
    this.lines = parse(this.HTML.textarea.value);
    const frag = this.partials.highlighted = document.createDocumentFragment();
    for (const line of this.lines) {
      const span = document.createElement("span");
      span.innerText = line.text;
      if (!line.ok) {
        span.className = "bg-red-200 text-red-800";
      }
      frag.appendChild(span)
    }
    this.render();
  }

  hasErrors() {
    for (let line of this.lines) {
      if (!line.ok) {
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
    this.HTML.textarea.value = text.trim().replaceAll(/ +/g, " ");
    this.update();
  }

  deleteVowels() {
    let text = "";
    for (let c of this.HTML.textarea.value) {
      if (!isTashkeel(c)) {
        text += c;
      }
    }
    this.HTML.textarea.value = text;
    this.update();
  }

  _paste = (e) => {
    e.preventDefault();
    let paste = e.clipboardData.getData("text");
    paste = paste.trim().replaceAll("\n", " ");
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
    this.update();
  }
}

class ArabicInputButton extends HTMLElement {
  constructor() {
    super();
    this.innerHTML = this.initHTML();
    this.HTML = Object.create(null);
    this.HTML.root = this.querySelector("button");

    this.target = null;
  }

  connectedCallback() {
    this.HTML.root.addEventListener("click", this._click);
    this.target = document.querySelector("arabic-input");
  }

  disconnectedCallback() {
    this.HTML.root.removeEventListener("click", this._click);
    this.target = null;
  }

  initHTML() {
    return html`
        <button type="button" class="bg-red-600 capitalize text-white rounded-lg p-2">DEFAULT BUTTON</button>
      `;
  }
}

// TODO(Amr Ojjeh): Gray out if there are no errors
export class DeleteErrorsButton extends ArabicInputButton {
  constructor() {
    super();
    this.HTML.root.innerText = "Delete all errors";
  }

  _click = (_e) => {
    if (confirm("Are you sure you want to delete all errors?")) {
      this.target.deleteErrors();
    }
  }
}

export class DeleteVowelsButton extends ArabicInputButton  {
  constructor() {
    super();
    this.HTML.root.innerText = "Delete vowels";
  }

  _click = (_e) => {
    if (confirm("Are you sure you want to delete all vowels?")) {
      this.target.deleteVowels();
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

function getSpaces(line, indexStart) {
  if (line[indexStart] !== " ") {
    throw new Error("Line should start with a space!");
  }
  const pack = {
    first: line[indexStart],
    extra: "",
  };
  for (let i = indexStart + 1; i < line.length; ++i) {
    const char = line[i];
    if (char === " ") {
      pack.extra += char;
    } else {
      return pack;
    }
  }
  return pack;
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

function isPunctuation(letter) {
  return letter == " ";
}

function isValid(letter) {
  return isArabicLetter(letter) || isPunctuation(letter);
}

function parse(text, debug=true) {
  if (debug) {
    var log = console.log;
  } else {
    var log = () => {};
  }

  let lines = [];
  let currentLine = {ok: true, text: ""};
  for (let i = 0; i < text.length; ++i) {
    // Check for tashkeel
    const letter = text[i];
    if (!isTashkeel(letter)) {
      const pack = getLetterPack(text, i);
      if (pack.tashkeel.length > 0) {
        if (currentLine.text) {
          lines.push(currentLine);
          log("Pushed", currentLine);
        }
        currentLine = {ok: false, text: letter + pack.tashkeel};
        lines.push(currentLine);
        log("Pushed", currentLine);
        i += pack.tashkeel.length;
        currentLine = {ok: true, text: ""};
        continue;
      }
    }
    // Check for double space
    if (letter === " ") {
      currentLine.text += "\u00A0";
      const spaces = getSpaces(text, i);
      if (spaces.extra.length > 0) {
        if (currentLine.text) {
          lines.push(currentLine);
          log("Pushed", currentLine);
        }
        currentLine = {ok: false, text: "\u00A0".repeat(spaces.extra.length)}
        lines.push(currentLine);
        log("Pushed", currentLine);
        i += spaces.extra.length;
        currentLine = {ok: true, text: ""};
      }
      continue;
    }
    // Check for invalid letters
    const valid = isValid(letter);
    if (!valid && currentLine.ok) {
      lines.push(currentLine);
      log("Pushed", currentLine);
      currentLine = {ok: valid, text: ""};
    } else if (valid && !currentLine.ok) {
      lines.push(currentLine);
      log("Pushed", currentLine);
      currentLine = {ok: valid, text: ""};
    }
    currentLine.text = currentLine.text + letter;
  }
  lines.push(currentLine);
  log("Final push", currentLine);
  return lines;
}