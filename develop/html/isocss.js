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
            for (let i=0; i < tilesWide; i++) {
                const tile = document.createElement('map-tile');
                const index = i + j * tilesWide;
                tile.setAttribute('data-terrain', mapData.data.map.tiles[index].hexTerrain);
                tile.setAttribute('data-chopped', mapData.data.map.tiles[index].chopped);
                row.appendChild(tile);
            }
            map.appendChild(row);
        }
	} else {
		console.error(xhr.status, 'Data fetch failed. Response text follows.');
		console.log(xhr.responseText);
	}
}
let query = `{
    map(playerSpoilerMask: 6) {
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
		const textDiv = document.createElement('div');
        this.appendChild(tileDiv);
        this.appendChild(textDiv);
        tileDiv.classList.add('isotile');
        textDiv.classList.add('tiletext');
        if (this.dataset.chopped == 'true') {
            const chopDiv = document.createElement('div');
            chopDiv.classList.add('chopped');
            this.appendChild(chopDiv);
            // textDiv.innerText += "C";
        }
        let terr = this.dataset.terrain;
        if (terr) {
            if (baseTerrainCss[terr[1]]) {
                this.style.setProperty('--tile-color', `var(--${baseTerrainCss[terr[1]]})`);
            }

            if (overlayTerrain[terr[0]]) {
                textDiv.innerText = overlayTerrain[terr[0]];
            }
        } else {
            textDiv.innerText = "🌲⛰️🌴🌳";
        }
	}
}
// TODO: put this in try/catch with friendly output for non-web-component browsers
customElements.define('map-tile', MapTile);
