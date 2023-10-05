const html = Strings.raw;

export class GrammarTag extends HTMLElement {
  constructor() {
    super();
    this.data = Object.create(null);
    this.data.value = "";
  }

  connectedCallback() {
    if (this.getAttribute("value")) {
      this.data.value = this.getAttribute("value");
    }

    this.innerHTML = this.initHTML();
    this.HTML = Object.create(null);
  }

  initHTML() {
    return html`
      <div>
      </div>`;
  }
}