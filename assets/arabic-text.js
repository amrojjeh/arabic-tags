// TODO(Amr Ojjeh): 1) Bundle Lit 2) Fix Firefox 3) Move functions outside class
import {LitElement, html, css, map} from 'https://cdn.jsdelivr.net/gh/lit/dist@2/all/lit-all.min.js';

export class ArabicText extends LitElement {
  
  static properties = {
    lines: {},
  }
  
  constructor() {
    super();
    this.lines = [];
  }
  
  render() {
    return html`
      <div dir="rtl" class="border-2 border-yellow-200 h-full py-10 px-2">
        <div class="relative h-full">
          <div class="absolute top-0 left-0 h-full w-full text-2xl">${this._highlighted()}</div>
          <textarea @input=${this._input} spellcheck="false" class="absolute focus:outline-none top-0 left-0 caret-black text-transparent bg-transparent h-full w-full text-2xl resize-none" placeholder="اكتب..."></textarea>
        <div>
      </div>`;
  }

  _getLetterPack(line, indexStart) {
    if (!this._isArabicLetter(line[indexStart])) {
      throw new Error("Line should start with an Arabic letter!");
    }
    const pack = {
      letter: line[indexStart],
      tashkeel: "",
    };
    for (let i = indexStart + 1; i < line.length; ++i) {
      const char = line[i];
      if (this._isTashkeel(char)) {
        pack.tashkeel += char;
      } else {
        return pack;
      }
    }
    return pack;
  }
  
  _input(e) {
    const text = e.target.value;
    let lines = [];
    let currentLine = {ok: true, text: ""};
    for (let i = 0; i < text.length; ++i) {
      const letter = text[i];
      if (this._isArabicLetter(letter)) {
        const pack = this._getLetterPack(text, i);
        if (pack.tashkeel.length > 0) {
          if (currentLine.text) {
            lines.push(currentLine);
            console.log("Pushed", currentLine);
          }
          currentLine = {ok: false, text: letter + pack.tashkeel};
          lines.push(currentLine);
          console.log("Pushed", currentLine);
          i += pack.tashkeel.length;
          currentLine = {};
          continue;
        }
      }
      const valid = this._isValid(letter);
      if (!valid && currentLine.ok) {
        lines.push(currentLine);
        console.log("Pushed", currentLine);
        currentLine = {ok: valid, text: ""};
      } else if (valid && !currentLine.ok) {
        lines.push(currentLine);
        console.log("Pushed", currentLine);
        currentLine = {ok: valid, text: ""};
      }
      currentLine.text = currentLine.text + letter;
    }
    lines.push(currentLine);
    console.log("Final push", currentLine);
    this.lines = lines;
  }
  
  _isValid(letter) {
    return this._isArabicLetter(letter) || this._isPunctuation(letter);
  }
  
  _isArabicLetter(char) {
    const code = char.codePointAt(0);
    if (code >= 0x0621 && code <= 0x063A) {
      return true;
    }
    if (code >= 0x0641 && code <= 0x064A) {
      return true;
    }
    return false;
  }

  _isTashkeel(char) {
    const code = char.codePointAt(0);
    return code >= 0x064B && code <= 0x065F;
  }

  _isEasternNumber(num) {
  }
    
  _isWesternNumber(num) {
    
  }
  
  hasErrors() {
    for (let line of this.lines) {
      if (!line.ok) {
        return true;
      }
    }
    return false;
  }
  
  _isPunctuation(letter) {
    return letter == " ";
  }
  
  _highlighted() {
    return html`
      ${map(this.lines, (line) => line.ok ? html`<span>${line.text}</span>` : html`<span class="bg-red-200 text-red-800">${line.text}</span>`)}`;
  }
  
  // Creating a light DOM for Tailwind
  createRenderRoot() {
    return this;
  }
}

