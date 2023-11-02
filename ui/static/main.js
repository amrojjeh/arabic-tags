import { setOffline } from "./base.js"
import { ArabicInput, DeleteErrorsButton, DeleteVowelsButton } from "./arabic-input.js";
import { GrammarTag } from "./grammar-tag.js";

window.save = () => true;

window.addEventListener("offline", () => {
  setOffline(true);
});


window.addEventListener("online", () => {
  if (window.save) {
    if (window.save()) {
      setOffline(false);
    }
  }
});

customElements.define("arabic-input", ArabicInput);
customElements.define("delete-errors", DeleteErrorsButton);
customElements.define("delete-vowels", DeleteVowelsButton);
customElements.define("grammar-tag", GrammarTag);
