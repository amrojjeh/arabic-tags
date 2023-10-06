const html = String.raw;

export class GrammarTag extends HTMLElement {
  constructor() {
    super();
    this.data = Object.create(null);
    this.data.value = "";
    this.data.words = [];
  }

  connectedCallback() {
    this._initData();
    this._initHTML();
  }

  disconnectedCallback() {

  }

  selectNext() {

  }

  _initData() {
    if (this.getAttribute("value")) {
      this.data.value = this.getAttribute("value");
    }
    this.data.words = this.data.value.split(" ");
  }

  _initHTML() {
    // div > input because we want to add suggestions as well
    this.innerHTML = html`
      <div dir="rtl" class="py-10 px-2 h-full">
        <p class="ps-3 pe-3 text-3xl leading-loose">${this.data.value}</p>
        <div class="pt-10 flex flex-col gap-5">
          <input autofocus placeholder="اكتب..." type="text" class="text-2xl ps-2 py-2 leading-loose drop-shadow mx-auto"></input>
        </div>
      </div>`;

    this.HTML = Object.create(null);
    this.HTML.root = this.querySelector("div");
  }
}