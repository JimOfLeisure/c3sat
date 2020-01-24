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
        console.log(xhr);
        let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Received non-2xx response on data query. Live updates will continue, but the latest save file is not shown here.`});
        dispatchEvent(cia3Error);
    }
}

xhr.onerror = e => {
    console.log(e);
    let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Data query request failed. Live updates will continue, but the latest save file is not shown here.`});
    dispatchEvent(cia3Error);
}

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
        console.log("failed xhr request:", pollXhr);
        let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Received non-2xx response on polling query. Live updates have stopped.`});
        dispatchEvent(cia3Error);
    }
}

pollXhr.onerror = e => {
    console.error("Long poll returned error");
    console.log(e);
    let cia3Error = new CustomEvent("cia3Error", { 'detail' : `Polling error. Live updates have stopped. Correct and refresh page.`});
    dispatchEvent(cia3Error);
}

class GqlQuery {
    // Using Set to deduplicate queries
    queryParts = new Set();
    query() {
        return '{' + Array.from(this.queryParts).join('\n') + '}';
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
        gqlQuery.queryParts.add('fileName');
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerText = data.fileName;
    }
}

class Difficulty extends HTMLElement {
    connectedCallback() {
        gqlQuery.queryParts.add('difficulty: int32s(section: "GAME", nth: 2, offset: 20, count: 1)');
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerText = this.difficultyNames[data.difficulty[0]];
    }
    difficultyNames = [
        "Chieftan",
        "Warlord",
        "Regent",
        "Monarch",
        "Emperor",
        "Demigod",
        "Deity",
        "Sid"
    ]
    
}

class Map extends HTMLElement {
    connectedCallback() {
        let spoilerMask = 0x2;
        gqlQuery.queryParts.add(`
            map(playerSpoilerMask: ${spoilerMask}) {
                tileSetWidth
                tileSetHeight
                tiles {
                    hexTerrain
                    chopped
                }
            }
        `);
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerHTML = '';
        let tilesWide = Math.floor(data.map.tileSetWidth / 2);
        this.style.setProperty('--map-width', tilesWide);
        data.map.tiles.forEach( (e, i) => {
            const tile = document.createElement('cia3-tile');
            if (e.hexTerrain) {
                tile.setAttribute('data-terrain', e.hexTerrain);
            }
            if (e.chopped) {
                tile.setAttribute('data-chopped', 'true');
            }
            if ((i + tilesWide) % data.map.tileSetWidth == 0) {
                tile.classList.add('odd-row');
            }
            this.appendChild(tile);
        })
    }
}

class Tile extends HTMLElement {
	connectedCallback () {
		this.render();
	}
    baseTerrainCss = {
        '0': 'desert',
        '1': 'plains',
        '2': 'grassland',
        '3': 'tundra',
        'b': 'coast',
        'c': 'sea',
        'd': 'ocean'
    }
    overlayTerrain = {
        '4': 'fp',
        '5': 'hill',
        '6': '⛰️',
        '7': '🌲',
        '8': '🌴',
        '9': 'marsh',
        'a': '🌋'
    }
    render () {
        const tileDiv = document.createElement('div');
        this.appendChild(tileDiv);
        tileDiv.classList.add('isotile');
        if (this.dataset.chopped == 'true') {
            const chopDiv = document.createElement('div');
            chopDiv.classList.add('chopped');
            this.appendChild(chopDiv);
        }
        let terr = this.dataset.terrain;
        if (terr) {
            if (this.baseTerrainCss[terr[1]]) {
                this.style.setProperty('--tile-color', `var(--${this.baseTerrainCss[terr[1]]})`);
            }
            if (this.overlayTerrain[terr[0]]) {
                const terrOverlayDiv = document.createElement('div');
                this.appendChild(terrOverlayDiv);
                terrOverlayDiv.className = 'terrain-overlay';
                terrOverlayDiv.innerText = this.overlayTerrain[terr[0]];
            }
        }
        let text = this.dataset.text;
        if (text) {
            const textDiv = document.createElement('div');
            textDiv.classList.add('tiletext');
            this.appendChild(textDiv);
        }
    }
}

class Url extends HTMLElement {
    connectedCallback() {
        this.render();
    }
    render() {
        this.innerHTML = `<a href="${location.href}" target="_blank">${location.href}</a>`;
    }
}

// TODO: Add controls to customize query and re-query. And remove old query from gqlQuery.
class HexDump extends HTMLElement {
    connectedCallback() {
        gqlQuery.queryParts.add(this.queryPart);
        window.addEventListener('refresh', () => this.render());
    }
    render() {
        this.innerText = data.cia3Hexdump;
    }
    queryPart = 'cia3Hexdump: hexDump(section: "GAME", nth: 2, offset: -4, count: 256)'
}


window.customElements.define('cia3-error', Error);
window.customElements.define('cia3-filename', Filename);
window.customElements.define('cia3-difficulty', Difficulty);
window.customElements.define('cia3-map', Map);
window.customElements.define('cia3-tile', Tile);
window.customElements.define('cia3-url', Url);
window.customElements.define('cia3-hexdump', HexDump);
pollNow();
