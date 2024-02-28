up.link.config.instantSelectors.push('a[href]')
up.link.config.followSelectors.push('a[href]')
up.form.config.submitSelectors.push(['form'])
up.compiler("[select]", function(el) {
  const input = el.querySelector("input")
  input.select()
  input.addEventListener("keydown", (e) => {
    if (e.key == "Escape") {
      e.preventDefault()
      up.follow(el.querySelector("a"))
    }
  })
})

let defaultVisible = up.viewport.config.autoFocusVisible
up.viewport.config.autoFocusVisible = (options) =>
  defaultVisible(options) && !up.fragment.matches(options.element, ':main')

// Used for syntax highlighting
const html = String.raw
let debug = false

if (debug) {
  var log = console.log
} else {
  var log = () => { }
}

class ArabicInput extends HTMLElement {

  constructor() {
    super()

    this.partials = Object.create(null)
    this.partials.highlighted = document.createDocumentFragment()

    this._okays = undefined
  }

  connectedCallback() {
    this.innerHTML = this.initHTML()
    this.HTML = {
      textarea: this.querySelector("textarea"),
      highlighted: this.querySelector("div > div > div"),
      form: this.querySelector("form")
    }

    this.HTML.textarea.tabindex = "0"
    if (typeof this.getAttribute("readonly") !== "string") {
      this.HTML.textarea.addEventListener("keydown", this._filter)
      this.HTML.textarea.addEventListener("input", this._input)
      this.HTML.textarea.addEventListener("paste", this._paste)
    } else {
      this.HTML.textarea.setAttribute("readonly", "")
      this.HTML.highlighted.classList.add("bg-gray-100/50")
    }
    if (this.getAttribute("value") != null) {
      this.HTML.textarea.value = this.getAttribute("value")
    }
    if (this.getAttribute("url") != null) {
      this.HTML.form.action = this.getAttribute("url");
    } else {
      console.error("There is no save url!")
    }
    if (this.getAttribute("punctuation") != null) {
      this.punctuation = new RegExp(this.getAttribute("punctuation"))
    } else {
      this.punctuation = ""
    }

    this.HTML.textarea.addEventListener("scroll", this._scroll)
    this.HTML.textarea.focus()
    this.update()
  }

  disconnectedCallback() {
    this.HTML.textarea.tabindex = "-1"
    this.HTML.textarea.removeEventListener("keydown", this._filter)
    this.HTML.textarea.removeEventListener("input", this._input)
    this.HTML.textarea.removeEventListener("paste", this._paste)
    this.HTML.textarea.removeEventListener("scroll", this._scroll)
  }

  initHTML() {
    return html`
      <div dir="rtl" class="h-full">
        <div class="relative h-full">
          <div class="ps-5 text-clip overflow-y-auto leading-loose absolute break-words top-0 left-0 h-full w-full text-3xl"></div>
          <form action="" method="post" up-autosubmit>
            <textarea spellcheck="false"
              name="content"
              autofocus
              class="ps-5 leading-loose absolute focus:outline-none top-0 left-0 caret-black
              text-transparent bg-transparent h-full w-full text-3xl resize-none"
              placeholder="اكتب..."></textarea>
          </form>
        </div>
      </div>`
  }

  save() {
    up.submit("form", {})
  }

  render() {
    const start = Date.now()
    this.HTML.highlighted.innerHTML = ""
    this.HTML.highlighted.appendChild(this.partials.highlighted)
    log("render:", Date.now() - start, "milliseconds")
  }

  createSpan(okay) {
    const span = document.createElement("span")
    if (okay === -1) {
      console.error("Okay should never be -1")
      return
    }

    if (okay === 1) {
      span.className = "bg-yellow-200 text-yellow-800"
    } else if (okay === 0) {
      span.className = "bg-red-200 text-red-800"
    }
    return span
  }

  update() {
    const start = Date.now()
    this.deleteDoubleSpaces()
    const text = this.HTML.textarea.value
    log("Deleted spaces:", Date.now() - start, "milliseconds")

    if (text.length > 0) {
      this._parse(text)
      const frag = this.partials.highlighted = document.createDocumentFragment()
      const okays = this.getOkays()
      let statusQuo = okays[0]
      let span = this.createSpan(statusQuo)
      for (let x = 0; x < text.length; ++x) {
        if (okays[x] === statusQuo) {
          span.innerText += text[x]
        } else {
          frag.appendChild(span)
          statusQuo = okays[x]
          span = this.createSpan(statusQuo)
          span.innerText = text[x]
        }
      }
      frag.appendChild(span)
    }
    log("Fragment updated:", Date.now() - start, "milliseconds")

    this.render()
    const event = new Event("arabic-input-update")
    this.dispatchEvent(event)
    log("Update:", Date.now() - start, "milliseconds")
  }

  hasErrors() {
    const okays = this.getOkays()
    for (let x = 0; x < this.HTML.textarea.value.length; ++x) {
      if (okays[x] === -1) {
        return false
      }
      if (okays[x] !== 2) {
        return true
      }
    }
    return false
  }

  hasTashkeel() {
    const text = this.HTML.textarea.value
    for (let x = 0; x < text.length; ++x) {
      if (isTashkeel(text[x])) {
        return true
      }
    }
    return false
  }

  deleteErrors() {
    let text = ""
    for (let c of this.HTML.textarea.value) {
      if (this.isValid(c)) {
        text += c
      }
    }
    this.HTML.textarea.value = text
    this.save()
    this.update()
  }

  // Called from update, so no need to update
  deleteDoubleSpaces() {
    const text = this.HTML.textarea.value
    this.HTML.textarea.value = text.replaceAll(/ {2,}/g, " ")
  }

  deleteVowels() {
    let text = ""
    for (let c of this.HTML.textarea.value) {
      if (!isTashkeel(c)) {
        text += c
      }
    }
    this.HTML.textarea.value = text
    this.save()
    this.update()
  }

  scaleOkays(size = 1000) {
    // this is done to optimize for a packed SMI array
    if (this._okays == undefined) {
      this._okays = []
    }
    for (let x = 0; x < size; x++) {
      if (!this._okays.hasOwnProperty(x)) {
        this._okays.push(-1)
      }
    }
    return this._okays
  }

  getOkays() {
    if (this._okays == undefined) {
      this.scaleOkays()
    }
    return this._okays
  }

  getOkaysSize(min) {
    if (min > this.getOkays().length) {
      return this.getOkays().length * 2
    }
    return this.getOkays().length
  }

  isValid(letter) {
    return isArabicLetter(letter) || isWhitespace(letter) || this.punctuation.test(letter)
  }

  _parse(text) {
    // 0 = NOT OKAY
    // 1 = Tashkeel
    // 2 = OKAY
    const start = Date.now()
    if (this.getOkays().length <= text.length) {
      this.scaleOkays(this.getOkaysSize(text.length))
    }
    let okays = this.getOkays()
    for (let i = 0; i < text.length; ++i) {
      // Check for tashkeel
      const letter = text[i]
      if (!isTashkeel(letter)) {
        const pack = getLetterPack(text, i)
        if (pack.tashkeel.length > 0) {
          okays[i] = 1
          for (let j = 0; j < pack.tashkeel.length; ++j) {
            okays[i + j + 1] = 1
          }
          i += pack.tashkeel.length
          continue
        }
      }
      // Check for invalid letters
      const valid = this.isValid(letter)
      if (!valid) {
        okays[i] = 0
      } else {
        okays[i] = 2
      }
    }
    log("Parsing:", Date.now() - start, "milliseconds")
  }

  _paste = (e) => {
    e.preventDefault()
    let paste = e.clipboardData.getData("text")
    paste = paste.replaceAll("\n", " ")
    paste = paste.trimLeft()
    paste = paste.trimRight()
    const selectionStart = e.target.selectionStart
    const value = e.target.value
    e.target.value = value.substring(0, selectionStart) + paste +
      value.substring(selectionStart + paste.length)
    this.update()
  }

  _filter = (e) => {
    if (e.key === "Enter") {
      e.preventDefault()
      return
    }
    if (e.key === " ") {
      const content = e.target.value
      const beforeCursor = content[e.target.selectionStart - 1]
      const afterCursor = content[e.target.selectionStart]
      if (beforeCursor === " ") {
        e.preventDefault()
        return
      } else if (afterCursor === " ") {
        e.target.selectionStart++
        e.preventDefault()
        return
      }
    }
  }

  _input = (_e) => {
    this.HTML.highlighted.scrollTop = this.HTML.textarea.scrollTop
    this.HTML.highlighted.scrollLeft = this.HTML.textarea.scrollLeft
    this.update(true)
  }

  _scroll = (_e) => {
    this.HTML.highlighted.scrollTop = this.HTML.textarea.scrollTop
    this.HTML.highlighted.scrollLeft = this.HTML.textarea.scrollLeft
  }
}

class ArabicInputButton extends HTMLElement {
  constructor(bgColor, fgColor) {
    super()
    this.target = null
    this.bgColor = bgColor
    this.fgColor = fgColor
  }

  className(bgColor, fgColor) {
    return `${bgColor} capitalize ${fgColor} rounded-lg p-2`
  }

  connectedCallback() {
    this.innerHTML = this.initHTML(this.className(this.bgColor, this.fgColor))
    this.HTML = Object.create(null)
    this.HTML.root = this.querySelector("button")

    this.HTML.root.addEventListener("click", this._click)
    this.target = document.querySelector("arabic-input")
    this.target.addEventListener("arabic-input-update", this._update)
    this._update()
  }

  disconnectedCallback() {
    this.HTML.root.removeEventListener("click", this._click)
    this.target = null
  }

  initHTML(className) {
    return html`
        <button type="button" class="${className}">DEFAULT BUTTON</button>
      `
  }
}

class DeleteErrorsButton extends ArabicInputButton {
  constructor() {
    super("bg-red-600", "text-white")
  }

  connectedCallback() {
    super.connectedCallback()
    this.HTML.root.innerText = "Delete all errors"
  }

  _click = (_e) => {
    if (confirm("Are you sure you want to delete all errors?")) {
      this.target.deleteErrors()
    }
  }

  _update = (_e) => {
    if (this.target.hasErrors()) {
      this.HTML.root.className = this.className("bg-red-600", "text-white")
      this.HTML.root.removeAttribute("disabled")
    } else {
      this.HTML.root.className = this.className("bg-gray-600", "text-white")
      this.HTML.root.setAttribute("disabled", "")
    }
  }
}

class DeleteVowelsButton extends ArabicInputButton {
  constructor() {
    super("bg-yellow-600", "text-white")
  }

  connectedCallback() {
    super.connectedCallback()
    this.HTML.root.innerText = "Delete vowels"
  }

  _click = (_e) => {
    if (confirm("Are you sure you want to delete all vowels?")) {
      this.target.deleteVowels()
    }
  }

  _update = (_e) => {
    if (this.target.hasTashkeel()) {
      this.HTML.root.className = this.className("bg-yellow-600", "text-white")
      this.HTML.root.removeAttribute("disabled")
    } else {
      this.HTML.root.className = this.className("bg-gray-600", "text-white")
      this.HTML.root.setAttribute("disabled", "")
    }
  }
}

function getLetterPack(line, indexStart) {
  if (isTashkeel(line[indexStart])) {
    throw new Error("Line should not start with tashkeel!")
  }
  const pack = {
    letter: line[indexStart],
    tashkeel: "",
  }
  for (let i = indexStart + 1; i < line.length; ++i) {
    const char = line[i]
    if (isTashkeel(char)) {
      pack.tashkeel += char
    } else {
      return pack
    }
  }
  return pack
}

function isTashkeel(char) {
  const code = char.codePointAt(0)
  return code >= 0x064B && code <= 0x065F
}

function isArabicLetter(char) {
  const code = char.codePointAt(0)
  if (code >= 0x0621 && code <= 0x063A) {
    return true
  }
  if (code >= 0x0641 && code <= 0x064A) {
    return true
  }
  return false
}

function isWhitespace(letter) {
  return letter == " "
}

customElements.define("arabic-input", ArabicInput)
customElements.define("delete-errors", DeleteErrorsButton)
customElements.define("delete-vowels", DeleteVowelsButton)
