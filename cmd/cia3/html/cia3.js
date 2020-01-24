export class Filename extends HTMLElement {
    connectedCallback() {
        this.render()
    }
    render() {
        this.innerText = "Filename goes here";
    }
}
window.customElements.define('cia3-filename', Filename);

export class Difficulty extends HTMLElement {
    connectedCallback() {
        this.render()
    }
    render() {
        this.innerText = "Difficulty goes here; warlord/emperor/etc";
    }
}
window.customElements.define('cia3-difficulty', Difficulty);

const difficultyNames = [
    "Chieftan",
    "Warlord",
    "Regent",
    "Monarch",
    "Emperor",
    "Demigod",
    "Deity",
    "Sid"
]