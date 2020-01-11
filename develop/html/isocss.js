let xhr = new XMLHttpRequest();
xhr.onload = () => {
	if (xhr.status >= 200 && xhr.status < 300) {
        map = document.getElementById('map');
        map.innerHTML = '';
        mapData = JSON.parse(xhr.responseText);
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
	} else {
		console.error(xhr.status, 'Data fetch failed. Response text follows.');
		console.log(xhr.responseText);
	}
}

let spoilerMask = 0x2;
let query = `{
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

class MapTile extends HTMLElement {
	connectedCallback () {
		this.render();
	}
	render () {
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
        this.appendChild(tileDiv);
        tileDiv.classList.add('isotile');
        if (this.dataset.chopped == 'true') {
            const chopDiv = document.createElement('div');
            chopDiv.classList.add('chopped');
            this.appendChild(chopDiv);
        }
        let terr = this.dataset.terrain;
        if (terr) {
            if (baseTerrainCss[terr[1]]) {
                this.style.setProperty('--tile-color', `var(--${baseTerrainCss[terr[1]]})`);
            }
            if (overlayTerrain[terr[0]]) {
                const terrOverlayDiv = document.createElement('div');
                this.appendChild(terrOverlayDiv);
                terrOverlayDiv.className = 'terrain-overlay';
                terrOverlayDiv.innerText = overlayTerrain[terr[0]];
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
// TODO: put this in try/catch with friendly output for non-web-component browsers
customElements.define('map-tile', MapTile);
