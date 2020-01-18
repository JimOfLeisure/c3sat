let xhr = new XMLHttpRequest();
xhr.onload = () => {
	if (xhr.status >= 200 && xhr.status < 300) {
        map = document.getElementById('map');
        map.innerHTML = '';
        mapData = JSON.parse(xhr.responseText);
        const fileName = document.getElementById("fileName");
        fileName.innerText = mapData.data.fileName;
        let tilesWide = Math.floor(mapData.data.map.tileSetWidth / 2);
        map.style.setProperty('--map-width', tilesWide);
        for (let j = 0, l = mapData.data.map.tileSetHeight; j < l; j++) {
            const row = document.createElement('div');
            row.classList += 'row';
            map.appendChild(row);
            for (let i=0; i < tilesWide; i++) {
                const tile = document.createElement('map-tile');
                const index = i + j * tilesWide;
                if (mapData.data.map.tiles[index].hexTerrain) {
                    tile.setAttribute('data-terrain', mapData.data.map.tiles[index].hexTerrain);
                }
                if (mapData.data.map.tiles[index].chopped) {
                    tile.setAttribute('data-chopped', 'true');
                }
                row.appendChild(tile);
            }
        }
    // If non-WebComponent browser, manually render each tile
    if (typeof customElements == 'undefined') {
        nonWebComponentRender();
    }
	} else {
		console.error(xhr.status, 'Data fetch failed. Response text follows.');
		console.log(xhr.responseText);
    }
}

let spoilerMask = 0x2;
let query = `{
    fileName
    map(playerSpoilerMask: ${spoilerMask}) {
        tileSetWidth
        tileSetHeight
        tiles {
            hexTerrain
            chopped
        }
    }
}`;

let body = {
    // "operationName":null,
    // "variables":{},
    'query' : query
};

xhr.open('POST', 'http://127.0.0.1:8080/graphql');
xhr.setRequestHeader('Content-Type', 'application/json');
xhr.send(JSON.stringify(body));

function renderMapTile (e) {
    const baseTerrainCss = {
        '0': 'desert',
        '1': 'plains',
        '2': 'grassland',
        '3': 'tundra',
        'b': 'coast',
        'c': 'sea',
        'd': 'ocean'
    }
    const overlayTerrain = {
        '4': 'fp',
        '5': 'hill',
        '6': '⛰️',
        '7': '🌲',
        '8': '🌴',
        '9': 'marsh',
        'a': '🌋'
    }
    const tileDiv = document.createElement('div');
    e.appendChild(tileDiv);
    tileDiv.classList.add('isotile');
    if (e.dataset.chopped == 'true') {
        const chopDiv = document.createElement('div');
        chopDiv.classList.add('chopped');
        e.appendChild(chopDiv);
    }
    let terr = e.dataset.terrain;
    if (terr) {
        if (baseTerrainCss[terr[1]]) {
            e.style.setProperty('--tile-color', `var(--${baseTerrainCss[terr[1]]})`);
        }
        if (overlayTerrain[terr[0]]) {
            const terrOverlayDiv = document.createElement('div');
            e.appendChild(terrOverlayDiv);
            terrOverlayDiv.className = 'terrain-overlay';
            terrOverlayDiv.innerText = overlayTerrain[terr[0]];
        }
    }
    let text = e.dataset.text;
    if (text) {
        const textDiv = document.createElement('div');
        textDiv.classList.add('tiletext');
        e.appendChild(textDiv);
    }
}

class MapTile extends HTMLElement {
	connectedCallback () {
        renderMapTile(this);
		// this.render();
	}
	// render () {
	// }
}

function nonWebComponentRender() {
    foo = document.getElementsByTagName('map-tile');
    for (i=0; i< foo.length; i++) {
        renderMapTile(foo[i]);
    }
}

// If non-WebComponent browser, manually render each tile
if (typeof customElements == 'undefined') {
    document.addEventListener('DOMContentLoaded',() => {
        nonWebComponentRender();
    });
} else {
    customElements.define('map-tile', MapTile);
}

