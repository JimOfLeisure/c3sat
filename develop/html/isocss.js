let xhr = new XMLHttpRequest();
xhr.onload = () => {
	if (xhr.status >= 200 && xhr.status < 300) {
        map = document.getElementById('map');
        map.innerHTML = '';
        mapData = JSON.parse(xhr.responseText);
        let tilesWide = Math.floor(mapData.data.map.tileSetWidth / 2);
        map.style.setProperty('--map-width', tilesWide);
        for (let j = 0; j < mapData.data.map.tileSetHeight; j++) {
            const row = document.createElement('div');
            row.classList += 'row';
            for (let i=0; i < tilesWide; i++) {
                const tile = document.createElement('map-tile');
                // tile.setAttribute('data-tile', JSON.stringify(mapData.data.map.tiles[i]));
                tile.setAttribute('data-terrain', mapData.data.map.tiles[i + j * tilesWide].hexTerrain);
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
    }
}`;

query = `{
    map {
        mapWidth
        mapHeight
        tileSetWidth
        tileSetHeight
        tiles {
            hexTerrain
        }
    } 
}`;

let body = {
    // "operationName":null,
    // "variables":{},
    'query' : query
};

// body = '{"query":"# Write your query or mutation here\n# Trade network ID by civ; nth should be a multiple of 4\n# {\n#   int16s(section: \"TILE\", nth: 4, offset: 26, count: 32)\n# }\n\n# # Starting locations of players\n# {\n#   int32s(section: \"WRLD\", nth: 2, offset: 36, count: 32)\n# }\n\n{\n  hexString(section:\"TILE\", nth: 1, offset:208, count: 4)\n  map(playerSpoilerMask: 6) {\n    mapWidth\n    mapHeight\n    tileSetWidth\n    tileSetHeight\n    tileSetX\n    tileSetY\n    tiles {\n      foo\n      hexTerrain\n    }\n  }\n \tint32s(section: \"WRLD\", nth: 2, offset: 8, count: 6)\n\n}\n"}';
// body = `{"query":"{ hexString(section:\"TILE\", nth: 1, offset:208, count: 4)}"}`;
// body = {"query":"{ map { mapWidth} }"};

xhr.open('POST', 'http://127.0.0.1:8080/graphql');
xhr.setRequestHeader('Content-Type', 'application/json');
xhr.send(JSON.stringify(body));
// xhr.send(body);
// console.log(body);

class MapTile extends HTMLElement {
	connectedCallback () {
		this.render();
	}
	render () {
		const tileDiv = document.createElement('div');
		const textDiv = document.createElement('div');
        this.appendChild(tileDiv);
        this.appendChild(textDiv);
        tileDiv.classList.add('isotile');
        textDiv.classList.add('tiletext');
        if (this.dataset.terrain) {
            textDiv.innerText = this.dataset.terrain;
            switch(this.dataset.terrain[1]) {
                case 'd':
                    this.style.setProperty('--tile-color', 'var(--ocean)');
                    break;
                case 'c':
                    this.style.setProperty('--tile-color', 'var(--sea)');
                    break;
                    case 'b':
                    this.style.setProperty('--tile-color', 'var(--coast)');
                    break;
            }
        } else {
            textDiv.innerText = "🌲⛰️🌴🌳";
        }
	}
}
// TODO: put this in try/catch with friendly output for non-web-component browsers
customElements.define('map-tile', MapTile);
