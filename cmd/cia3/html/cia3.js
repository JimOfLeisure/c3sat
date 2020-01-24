let pollXhr = new XMLHttpRequest();
let pollSince = Date.now() - 86400000
const longPollTimeout = 30

let pollNow = () => {
    pollXhr.open('GET', `/events?timeout=${longPollTimeout}&category=refresh&since_time=${pollSince}`);
    pollXhr.send();
}

pollXhr.onload = () => {
    if (pollXhr.status >= 200 & pollXhr.status < 300) {
        let pollData = JSON.parse(pollXhr.responseText);
        if (typeof pollData.events != 'undefined') {
            pollSince = pollData.events[0].timestamp;
            console.log('poll event received at ' + pollSince);
        }
        if (pollData.timeout != undefined) {
            pollSince = pollData.timestamp;
        }
        pollNow();
    } else {
        // TODO: Better error handling. For now just passing the xhr object to a function which usually gets a ProgressEvent
        // pollError(pollXhr);
        console.log(pollXhr);
    }
}

pollXhr.onerror = e => pollError(e);

let pollError = (e) => {
    console.error("Long poll returned error");
    console.log(e);
    const errorDiv = document.getElementsByTagName('cia3-error');
    const errMsg = document.createElement('p');
    errMsg.innerText = `Polling error. Live updates have stopped. Correct and refresh page.`
    errorDiv[0].appendChild(errMsg);

}

class Error extends HTMLElement {
    connectedCallback() {
        this.render();
    }
    render() {
        this.innerText = "Errors go here";
    }
}

class Filename extends HTMLElement {
    connectedCallback() {
        this.render()
    }
    render() {
        this.innerText = "Filename goes here";
    }
}

class Difficulty extends HTMLElement {
    connectedCallback() {
        this.render()
    }
    render() {
        this.innerText = "Difficulty goes here; warlord/emperor/etc";
    }
}

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

window.customElements.define('cia3-error', Error);
window.customElements.define('cia3-filename', Filename);
window.customElements.define('cia3-difficulty', Difficulty);
pollNow();
