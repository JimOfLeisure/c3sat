/*
    loading this module starts long polling and xhr queries (if appropriate)
    pollXhr implements long polling
        if event received, launch xhr query
    Each custom element adds its query part to gqlQuery (type GqlQuery)
        use GraphQL aliases for each custom element for clarity, especially if arguments differ
    GqlQuery dedups query parts and builds the xhr request body
    Upon successful xhr query, "refresh" event dispatched
        each custom element listens to refresh and renders itself when received
    event "cia3Error" dispatched on errors, with detail being a string which the Error custom element will add to itself
*/

let xhr = new XMLHttpRequest();
let pollXhr = new XMLHttpRequest();
let pollSince = Date.now() - 86400000;
const longPollTimeout = 30;
let data = {};

xhr.onload = () => {
	if (xhr.status >= 200 && xhr.status < 300) {
        data = JSON.parse(xhr.responseText).data;
        const refreshData = new CustomEvent("refresh");
        dispatchEvent(refreshData);
	} else {
        // TODO: Handle non-2xx results
        console.log(xhr);
    }
}

// TODO: Handle xhr errors
xhr.onerror = e => console.log(e);

let pollNow = () => {
    pollXhr.open('GET', `/events?timeout=${longPollTimeout}&category=refresh&since_time=${pollSince}`);
    pollXhr.send();
}

pollXhr.onload = () => {
    if (pollXhr.status >= 200 & pollXhr.status < 300) {
        let pollData = JSON.parse(pollXhr.responseText);
        if (typeof pollData.events != 'undefined') {
            pollSince = pollData.events[0].timestamp;
            xhr.open('POST', '/graphql');
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.send(JSON.stringify(gqlQuery.body()));
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

pollXhr.onerror = e => {
    console.error("Long poll returned error");
    console.log(e);
    let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Polling error. Live updates have stopped. Correct and refresh page.`});
    dispatchEvent(cia3Error);
}

class GqlQuery {
    queryParts = new Array();
    query() {
        // Using Set to deduplicate array
        return '{' + [...new Set(this.queryParts)].join('\n') + '}';
    }
    body() {
        return {
            'query' : this.query()
        }
    }
}
let gqlQuery = new GqlQuery();

class Error extends HTMLElement {
    connectedCallback() {
        window.addEventListener('cia3Error', (e) => this.render(e.detail));
    }
    render(errMsg) {
        const p = document.createElement('p');
        p.innerText = errMsg;
        this.appendChild(p);
        // this.innerText = errMsg;
    }
}

class Filename extends HTMLElement {
    connectedCallback() {
        gqlQuery.queryParts.push('fileName');
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerText = data.fileName;
    }
}

class Difficulty extends HTMLElement {
    connectedCallback() {
        gqlQuery.queryParts.push('difficulty: int32s(section: "GAME", nth: 2, offset: 20, count: 1)');
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerText = difficultyNames[data.difficulty[0]];
    }
}

class Map extends HTMLElement {
    connectedCallback() {
        gqlQuery.queryParts.push('testQueryNotReal: int32s(section: "GAME", nth: 2, offset: 20, count: 1)');
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerText = difficultyNames[data.difficulty[0]];
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
window.customElements.define('cia3-map', Map);
pollNow();
