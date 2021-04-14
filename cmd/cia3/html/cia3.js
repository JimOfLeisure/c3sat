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
        /* tried to see if can make alert come to front. Not with Chrome or FF, apparently
            maybe check out web notifications instead
                https://www.w3.org/TR/notifications/
                https://developer.mozilla.org/en-US/docs/Web/API/Notifications_API/Using_the_Notifications_API
                https://developers.google.com/web/fundamentals/push-notifications
                https://www.google.com/search?q=web+push+notifications
            also, would not notify on every refresh; this would only be for alerts
        */
        // window.focus();
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

// Most of the cia3-* elements follow this form, so extend this class
class Cia3Element extends HTMLElement {
    connectedCallback() {
        this.registerAndListen();
    }
    render() {
        this.innerText = 'REPLACE ME';
    }
    registerAndListen() {
        gqlQuery.queryParts.add(this.queryPart);
        window.addEventListener('refresh', () => this.render());
    }
    queryPart = 'REPLACE ME';
}

// TODO: Allow removal of error messages
class Error extends HTMLElement {
    connectedCallback() {
        window.addEventListener('cia3Error', (e) => this.render(e.detail));
    }
    render(errMsg) {
        const p = document.createElement('p');
        p.innerText = errMsg;
        this.appendChild(p);
    }
}

class Filename extends Cia3Element {
    render() {
        this.innerText = data.fileName;
    }
    queryPart = 'fileName';
}

class Fullpath extends Cia3Element {
    render() {
        this.innerText = data.fullPath;
    }
    queryPart = 'fullPath';
}

class Difficulty extends Cia3Element {
    render() {
        this.innerText = data.difficultyNames[data.difficulty[0]].name;
    }
    queryPart = `
        difficulty: int32s(section: "GAME", nth: 2, offset: 20, count: 1)
        difficultyNames: listSection(target: "bic", section: "DIFF", nth: 1) { name: string(offset:0, maxLength: 64) }
    `;
}

class Map extends Cia3Element {
    connectedCallback() {
        let spoilerMask = 0x2;
        this.queryPart = `
            map(playerSpoilerMask: ${spoilerMask}) {
                tileSetWidth
                tileSetHeight
                tiles {
                    hexTerrain
                    chopped
                }
            }
        `;
        this.registerAndListen();
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
        });
    }
    queryPart = '';
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
        let url = location.protocol + "//" + location.host;
        this.innerHTML = `<a href="${url}">${url}</a>`;
    }
}

// TODO: Add controls to customize query and re-query. And remove old query from gqlQuery.
class HexDump extends Cia3Element {
    render() {
        this.innerText = 'Hex dump tool under construction, no controls yet.\n' + data.cia3Hexdump;
    }
    // queryPart = 'cia3Hexdump: hexDump(section: "WRLD", nth: 1, offset: 0, count: 8192)';
    queryPart = 'cia3Hexdump: hexDump(section: "GAME", nth: 2, offset: 0, count: 8192)';
}

class MapX extends Cia3Element {
    render() {
        this.innerText = data.mapx[0];
    }
    queryPart = 'mapx: int32s(section: "WRLD", nth: 2, offset: 8, count: 1)';
}

class MapY extends Cia3Element {
    render() {
        this.innerText = data.mapy[0];
    }
    queryPart = 'mapy: int32s(section: "WRLD", nth: 2, offset: 28, count: 1)';
}

class WorldSize extends Cia3Element {
    render() {
        this.innerText = data.worldSizeNames[data.worldsize.size].name;
    }
    queryPart = `
        worldsize: civ3 { size }
        worldSizeNames: listSection(target: "bic", section: "WSIZ", nth: 1) { name: string(offset:32, maxLength: 32) }
    `;
}

class Barbarians extends Cia3Element {
    render() {
        this.innerText = this.barbariansSettings[data.barbarians.barbariansFinal.toString()];
    }
    queryPart = 'barbarians: civ3 { barbariansFinal }';
    barbariansSettings = {
        '-1': 'No Barbarians',
        '0': 'Sedentary',
        '1': 'Roaming',
        '2': 'Restless',
        '3': 'Raging',
        '4': 'Random'
    };
}

class WorldSeed extends Cia3Element {
    render() {
        this.innerText = data.worldseed.worldSeed;
    }
    queryPart = 'worldseed: civ3 { worldSeed }';
}

class LandMass extends Cia3Element {
    render() {
        this.innerText = this.landMassNames[data.landmass.landMassFinal];
    }
    queryPart = 'landmass: civ3 { landMassFinal }';
    landMassNames = [
        "Archipelago",
        "Continents",
        "Pangea",
        "Random"
    ];
}

class OceanCoverage extends Cia3Element {
    render() {
        this.innerText = this.oceanCoverageNames[data.oceancoverage.oceanCoverageFinal];
    }
    queryPart = 'oceancoverage: civ3 { oceanCoverageFinal }';
    oceanCoverageNames = [
        "80% Water",
        "70% Water",
        "60% Water",
        "Random"
    ];
}

class Climate extends Cia3Element {
    render() {
        this.innerText = this.climateNames[data.climate.climateFinal];
    }
    queryPart = 'climate: civ3 { climateFinal }';
    climateNames = [
        "Arid",
        "Normal",
        "Wet",
        "Random"
    ];
}

class Temperature extends Cia3Element {
    render() {
        this.innerText = this.temperatureNames[data.temperature.temperatureFinal];
    }
    queryPart = 'temperature: civ3 { temperatureFinal }';
    temperatureNames = [
        "Warm",
        "Temperate",
        "Cool",
        "Random"
    ];
}

class Age extends Cia3Element {
    render() {
        this.innerText = this.ageNames[data.age.ageFinal];
    }
    queryPart = 'age: civ3 { ageFinal }';
    ageNames = [
        "3 Billion",
        "4 Billion",
        "5 Billion",
        "Random"
    ];
}

class Civs extends Cia3Element {
    numFields = 111;
    player = 1;
    render() {
        this.innerHTML = '';
        const friendlyTable = document.createElement('table');
        this.appendChild(friendlyTable);
        friendlyTable.innerHTML = `<tr>
            <th>Civ</th>
            <th>Contact</th>
            <th>Relationship</th>
            <th>Will Talk</th>
            <th>Government</th>
            <th>Era</th>
            <th>City Count</th>
        </tr>`;
        data.civs.filter(this.civsFilter, this).forEach((e, i) => {
            const friendlyRow = document.createElement('tr');
            const contactValue = data.civs[this.player].contactWith[e.playerNumber[0]];
            const haveContact = contactValue != 0;
            const warVarlue = data.civs[this.player].atWar[e.playerNumber[0]];
            const atWar = warVarlue != 0;
            // friendlyRow.innerHTML += `<td>${e.playerNumber[0]}</td>`;
            friendlyRow.innerHTML += `<td>${data.race[e.raceId[0]].civName}</td>`;
            friendlyRow.innerHTML += `<td>${this.contactWithName(contactValue)}</td>`;
            friendlyRow.innerHTML += `<td>${haveContact ? this.relationshipName(warVarlue) : '-'}</td>`;
            friendlyRow.innerHTML += `<td>${haveContact ? this.willTalk(e) : '-'}</td>`;
            friendlyRow.innerHTML += `<td>${haveContact ? data.governmentNames[e.governmentType[0]].name : '-'}</td>`;
            friendlyRow.innerHTML += `<td>${haveContact ? data.eraNames[e.era[0]].name : '-'}</td>`;
            // TODO: I don't like using the string value here. Figure out how to use the number without repeating code
            friendlyRow.innerHTML += `<td>${haveContact && this.willTalk(e) == "Yes" ? e.cityCount[0] : '-'}</td>`;
            friendlyTable.appendChild(friendlyRow);
            // if (this.oldCivsData != undefined) {
            //     // put alerts code in here
            // }
        })
        this.oldCivsData = data.civs
    }
    civsFilter (e) {
        return (e.raceId[0] > 0) && (this.player != e.playerNumber[0]); // non-barb, non "player" players
    }
    oldCivsData = undefined;
    queryPart = `
        civs {
            playerNumber: int32s(offset:0, count: 1)
            raceId: int32s(offset:4, count: 1)
            governmentType: int32s(offset:132, count: 1)
            era: int32s(offset:216, count: 1)
            cityCount: int32s(offset:376, count: 1)
            atWar: bytes(offset:3348, count: 32)
            willTalkTo: int32s(offset:2964, count: 32)
            contactWith: int32s(offset:3732, count: 32)
        }
        race {
            civName: string(offset:128, maxLength: 40)
        }
        governmentNames: listSection(target: "bic", section: "GOVT", nth: 1) { name: string(offset:24, maxLength: 64) }
        eraNames: listSection(target: "bic", section: "ERAS", nth: 1) { name: string(offset:0, maxLength: 64) }
    `;
    contactWithName (i) {
        if (i==0) return "No";
        // if (i==1) return "Yes";
        // return "Yes (" + i.toString() + ")"; // apparently other flags exist for units in territory
        return "Yes";
    }
    relationshipName (i) {
        if (i==0) return "Peace";
        // if (i==1) return "WAR";
        // return "war? (" + i.toString() + ")"; // don't know what else there is
        return "WAR";
    }
    willTalkToName (i) {
        if (i==0) return "Yes";
        // if (i==1) return "No";
        // return "No (" + i.toString() + ")"; // don't know what else there is
        return "NO";
    }
    willTalk (e) {
        if (e.atWar[this.player]) {
            return this.willTalkToName(e.willTalkTo[this.player]);
        } else {
            return this.willTalkToName(0);
        }
    }
}

class CivsDev extends Cia3Element {
    numFields = 111;
    render() {
        const player = 1;
        // this.innerHTML = JSON.stringify(data.civs, null, '  ');
        this.innerHTML = '';
        const table = document.createElement('table');
        const friendlyTable = document.createElement('table');
        const diplomaticTable = document.createElement('table');
        const hexDumps = document.createElement('div');
        hexDumps.classList += "dump";
        this.appendChild(friendlyTable);
        this.appendChild(table);
        this.appendChild(diplomaticTable);
        // table.innerHTML = '<tr><th>Player #</th><th>RACE ID</th>' + '<th>?</th>'.repeat(this.numFields - 2) + '</tr>';
        let headers = "";
        for (let i = 2; i < this.numFields; i++) {
            headers += `<th>${i} 0x${(i*4).toString(16)} ${i*4}</th>`
        }
        table.innerHTML = `<tr><th>First ${this.numFields} int32s / Player</th><th>RACE ID</th>` +  headers + '</tr>';
        const attitudeFieldCount = 23;
        headers = "";
        for (let i = 0; i < attitudeFieldCount; i++) {
            headers += `<th>${i} 0x${(i*4).toString(16)} ${i*4}</th>`
        }
        diplomaticTable.innerHTML = '<tr><th>attitude/diplomacy? Player - Opponent</th>' +  headers + '</tr>';
        friendlyTable.innerHTML = `<tr>
            <th>Friendly Table / Player #</th>
            <th>Civ Name</th>
            <th>Contact with player ${player.toString()}</th>
            <th>Player ${player.toString()} Contact With Them</th>
            <th>Relationship with player ${player.toString()}</th>
            <th>Player ${player.toString()} Relationship With Them</th>
            <th>Will Talk to player ${player.toString()}</th>
            <th>Government</th>
            <th>Mobilization</th>
            <th>Tiles Discovered</th>
            <th>Era</th>
            <th>Research Beakers</th>
            <th>Current Research Tech</th>
            <th>Current Research Turns</th>
            <th># Future Techs</th>
            <th># Armies</th>
            <th># Units</th>
            <th># Miltary Units</th>
            <th># Cities</th>
        </tr>`;
        data.civs.filter(this.civsFilter, this).forEach((e, i) => {
            const friendlyRow = document.createElement('tr');
            friendlyRow.innerHTML += `<td>${e.playerNumber[0]}</td>`;
            friendlyRow.innerHTML += `<td>${data.race[e.raceId[0]].civName}</td>`;
            friendlyRow.innerHTML += `<td>${this.contactWithName(e.contactWith[player])}</td>`;
            friendlyRow.innerHTML += `<td>${this.contactWithName(data.civs[player].contactWith[e.playerNumber[0]])}</td>`;
            friendlyRow.innerHTML += `<td>${this.relationshipName(e.atWar[player])}</td>`;
            friendlyRow.innerHTML += `<td>${this.relationshipName(data.civs[player].atWar[e.playerNumber[0]])}</td>`;
            // Unsure of willTalkTo data location, and unsure if it's an int32[32] array. Only see it for player 1 so far
            // Seems to be turns until they talk without battle impacts, but also counts down a few turns after making peace
            //   Maybe this prevents them from redeclaring war for a few turns, too?
            friendlyRow.innerHTML += `<td>${this.willTalkToName(e.willTalkTo[player])}</td>`;
            friendlyRow.innerHTML += `<td>${data.governmentNames[e.governmentType[0]].name}</td>`;
            friendlyRow.innerHTML += `<td>${e.mobilizationLevel[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.tilesDiscovered[0]}</td>`;
            friendlyRow.innerHTML += `<td>${data.eraNames[e.era[0]].name}</td>`;
            friendlyRow.innerHTML += `<td>${e.researchBeakers[0]}</td>`;
            let currentResearchId = e.currentResearchId[0];
            let currentResearchName;
            if (currentResearchId < 0 ) {
                currentResearchName = "-"
            } else {
                currentResearchName = data.techNames[currentResearchId].name
            }
            friendlyRow.innerHTML += `<td>${currentResearchName}</td>`;
            friendlyRow.innerHTML += `<td>${e.currentResearchTurns[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.futureTechsCount[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.armiesCount[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.unitCount[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.militaryUnitCount[0]}</td>`;
            friendlyRow.innerHTML += `<td>${e.cityCount[0]}</td>`;
            // friendlyRow.innerHTML += `<td>${}</td>`;
            // friendlyRow.innerHTML += `<td>${}</td>`;
            // friendlyRow.innerHTML += `<td>${}</td>`;
            // friendlyRow.innerHTML += `<td>${}</td>`;
            friendlyTable.appendChild(friendlyRow);
            const row = document.createElement('tr');
            e.int32s.forEach((ee, ii) => {
                const td = document.createElement('td');
                switch (ii) {
                    case 1:
                        if (ee != -1) {
                            td.innerText = data.race[ee].civName;
                            break;
                        }
                    default:
                        td.innerText = ee;
                }
                row.appendChild(td);
            });
            table.appendChild(row);

            for (let i = 1; i < 9; i++) {
                const intOffset = 23 * i;
                const diploRow = document.createElement('tr');
                diploRow.style = "color: lightgray;";
                diploRow.innerHTML += `<td>${e.playerNumber[0].toString()} - ${i}</td>`;
                for (let i = 0; i < 23; i++) {
                    let myValue = e.attitudes[intOffset + i]
                    diploRow.innerHTML += `<td${myValue !=0 ? ' style="color: black; font-weight: bold;"': ''}>${myValue.toString()}</td>`;
                }
                diplomaticTable.appendChild(diploRow);
            }

            if (this.oldCivsData != undefined) {
                const hexDiff = document.createElement('div');
                // let foo = Diff.diffWordsWithSpace(this.oldCivsData[e.int32s[0]].hexDump, data.civs[e.int32s[0]].hexDump);
                let foo = Diff.createTwoFilesPatch("old", "new" ,this.oldCivsData[e.int32s[0]].hexDump, data.civs[e.int32s[0]].hexDump);
                // foo.forEach(function(part){
                //     // green for additions, red for deletions
                //     // grey for common parts
                //     let color = part.added ? 'green' :
                //       part.removed ? 'red' : 'grey';
                //     let span = document.createElement('span');
                //     span.style.color = color;
                //     span.appendChild(document
                //       .createTextNode(part.value));
                //     hexDiff.appendChild(span);
                // });
                let diff2Html = Diff2Html.html(foo, {
                    drawFileList: true,
                    matching: 'words',
                    outputFormat: 'side-by-side',
                });
                hexDiff.innerHTML = '*****\n#' + e.int32s[0] + ' : ' + data.race[e.int32s[1]].civName + ' Diff\n*****\n' + diff2Html;
                this.appendChild(hexDiff);
            }
            hexDumps.innerHTML += '*****\n#' + e.int32s[0] + ' : ' + data.race[e.int32s[1]].civName + '\n*****\n' +
                e.hexDump.replace(/((?: 00)+)/g, '<span class="dim">$1</span>')
                .replace(/(\.+)/g, '<span class="dim">$1</span>')
                .replace(/((?: ff)+)/g, '<span class="medium">$1</span>')
                ;
        })
        this.appendChild(hexDumps);
        this.oldCivsData = data.civs
    }
    civsFilter (e) {
        return e.raceId[0] > 0; // non-barb players
        // return e.playerNumber[0] == 1; // first human player only
    }
    oldCivsData = undefined;
    queryPart = `
        civs {
            int32s(offset:0, count: ${this.numFields})
            hexDump
            playerNumber: int32s(offset:0, count: 1)
            raceId: int32s(offset:4, count: 1)
            governmentType: int32s(offset:132, count: 1)
            mobilizationLevel: int32s(offset:136, count: 1)
            tilesDiscovered: int32s(offset:140, count: 1)
            era: int32s(offset:252, count: 1)
            researchBeakers: int32s(offset:220, count: 1)
            currentResearchId: int32s(offset:224, count: 1)
            currentResearchTurns: int32s(offset:228, count: 1)
            futureTechsCount: int32s(offset:232, count: 1)
            armiesCount: int32s(offset:364, count: 1)
            unitCount: int32s(offset:368, count: 1)
            militaryUnitCount: int32s(offset:372, count: 1)
            cityCount: int32s(offset:376, count: 1)
            attitudes: int32s(offset:404, count: 736)
            atWar: bytes(offset:3348, count: 32)
            willTalkTo: int32s(offset:2964, count: 32)
            contactWith: int32s(offset:3732, count: 32)
        }
        race {
            leaderName: string(offset:0, maxLength: 32)
            leaderTitle: string(offset:32, maxLength: 24)
            adjective:  string(offset:88, maxLength: 40)
            civName: string(offset:128, maxLength: 40)
            objectNoun: string(offset:168, maxLength: 40)
        }
        governmentNames: listSection(target: "bic", section: "GOVT", nth: 1) { name: string(offset:24, maxLength: 64) }
        eraNames: listSection(target: "bic", section: "ERAS", nth: 1) { name: string(offset:0, maxLength: 64) }
        techNames: listSection(target: "bic", section: "TECH", nth: 1) { name: string(offset:0, maxLength: 32) }
    `;
    contactWithName (i) {
        if (i==0) return "No";
        // if (i==1) return "Yes";
        // return i.toString(); // don't know what else there is
        return "Yes";
    }
    relationshipName (i) {
        if (i==0) return "Peace";
        // if (i==1) return "WAR";
        // return i.toString(); // don't know what else there is
        return "WAR";
    }
    willTalkToName (i) {
        if (i==0) return "Yes";
        // if (i==1) return "No";
        // return "No: " + i.toString(); // don't know what else there is
        return "No";
    }
}

class CivTech extends Cia3Element {
    render () {
        this.innerHTML = '';
        // tech list offset 852 + 4*numContinents? Grab a bunch more "techs" to ensure have all of them
        // SEEMS TO WORK! TODO: better handle the tech count and offset; right now I'm just picking enough to hopefully cover likely continents and default # of techs
        // Also TODO: I'm grabbing numContinents from *after* the tech data, so the real authoritative count must be elsewhere.
        let intOffset = data.numContinents[0];
        console.log(intOffset);
        data.techList.forEach((e, i) => {
            if (data.techCivMask[i+intOffset] != 0) {
                this.innerHTML += `${e.name} -`;
                data.civs.forEach((ee, ii) => {
                    if ((2**ii & data.techCivMask[i+intOffset]) !=0) {
                        console.log(i, i+intOffset, ii);
                        this.innerHTML += ` ${data.race[ee.raceId[0]].civName}`;
                    }
                });
                this.innerHTML += `<br>`;
            }
        });
    }
    queryPart = `civs {
        playerNumber: int32s(offset:0, count: 1)
        raceId: int32s(offset:4, count: 1)
    }
    race {
        civName: string(offset:128, maxLength: 40)
    }
    techCivMask: int32s(section: "GAME", nth: 2, offset: 852, count: 140)
    techList: listSection(target: "bic", section: "TECH", nth: 1) {
        name: string(offset:0, maxLength: 32)
    }
    numContinents: int16s(section:"WRLD", nth: 1, offset: 4, count: 1)
    `;
}

class HexDiffAll extends Cia3Element {
    render () {
        this.innerHTML = '';
        
        this.innerHTML += `${data.fileName}<br>`;
        if (this.oldData != undefined) {
            const hexDiff = document.createElement('div');
            let foo = Diff.createTwoFilesPatch(this.oldData.fileName, data.fileName ,this.oldData.hexDumpAll, data.hexDumpAll);
            let diff2Html = Diff2Html.html(foo, {
                drawFileList: true,
                matching: 'none',
                outputFormat: 'side-by-side',
            });
            hexDiff.innerHTML = diff2Html;
            this.appendChild(hexDiff);
        }
        this.oldData = data;
    }
    oldData = undefined;
    queryPart = `fileName hexDumpAll`;
}

class Trade extends Cia3Element {
    player = 1;
    render () {
        let playerGold = data.tradeCivs[this.player].gold[0] + data.tradeCivs[this.player].gold[1];
        this.innerHTML = '';
        const friendlyTable = document.createElement('table');
        this.appendChild(friendlyTable);
        friendlyTable.innerHTML = `<tr>
            <th>Civ</th>
            <th>Gold</th>
            <th>Tech to Buy</th>
            <th>Tech to Sell</th>
            <th>Gold</th>
        </tr>`;
        data.tradeCivs.filter(this.willTalk, this).forEach((e, i) => {
            const friendlyRow = document.createElement('tr');
            friendlyRow.innerHTML += `<td>${data.tradeRace[e.raceId[0]].civName}</td>`;
            let techsToBuy = new Array();
            let techsToSell = new Array();
            data.techList.forEach((ee, ii) => {
                const thisTech = data.techCivMask[ii];
                // If mismatched truthiness, one player has it and the other doesn't
                if (!!(thisTech & 2**e.playerNumber[0]) != !!(thisTech & 2**this.player)) {
                    if (!(thisTech & 2**this.player)) {
                        if (this.hasPreReq(ii, this.player)) {
                            techsToBuy.push(ee.name);
                        }
                    } else {
                        if (this.hasPreReq(ii, e.playerNumber[0])) {
                            techsToSell.push(ee.name);
                        }
                    }
                }
            });
            friendlyRow.innerHTML += `<td>${e.gold[0] + e.gold[1]}</td>`;
            friendlyRow.innerHTML += `<td>${techsToBuy.join(", ")}</td>`;
            friendlyRow.innerHTML += `<td>${techsToSell.join(", ")}</td>`;
            friendlyRow.innerHTML += `<td>${playerGold}</td>`;
            friendlyTable.appendChild(friendlyRow);
        });
    }
    willTalk(e) {
        // Non-barb, non-player civs that exist, has contact with player, and will talk (not atWar OR willTalkTo == 0)
        return e.raceId[0] > 0
            && this.player != e.playerNumber[0]
            && data.tradeCivs[this.player].contactWith[e.playerNumber[0]]
            && (
                (! e.atWar[this.player])
                || (e.willTalkTo[this.player] == 0)
            )
            ;
    }
    // Check tech prereqs for player, and also era requirement. return boolean
    hasPreReq(techIndex, playerNumber) {
        return (data.tradeCivs[playerNumber].era[0] >= data.techList[techIndex].era[0]) &&
            data.techList[techIndex].prereqs.every(e => {
                if (e < 0) {
                    return true;
                }
                if (!!(data.techCivMask[e] & 2**playerNumber)) {
                    return true;
                }
                return false;
            });
    }
    queryPart = `
    tradeRace: race {
        civName: string(offset:128, maxLength: 40)
    }
    techCivMask
    techList: listSection(target: "bic", section: "TECH", nth: 1) {
        name: string(offset:0, maxLength: 32)
        era: int32s(offset: 68, count: 1)
        prereqs: int32s(offset: 84, count: 4)
    }
    tradeCivs: civs {
        playerNumber: int32s(offset:0, count: 1)
        raceId: int32s(offset:4, count: 1)
        gold: int32s(offset:40, count: 2)
        atWar: bytes(offset:3348, count: 32)
        willTalkTo: int32s(offset:2964, count: 32)
        contactWith: int32s(offset:3732, count: 32)
        era: int32s(offset:216, count: 1)
    }
`;
}


window.customElements.define('cia3-error', Error);
window.customElements.define('cia3-filename', Filename);
window.customElements.define('cia3-fullpath', Fullpath);
window.customElements.define('cia3-difficulty', Difficulty);
window.customElements.define('cia3-map', Map);
window.customElements.define('cia3-tile', Tile);
window.customElements.define('cia3-url', Url);
window.customElements.define('cia3-hexdump', HexDump);
window.customElements.define('cia3-mapx', MapX);
window.customElements.define('cia3-mapy', MapY);
window.customElements.define('cia3-worldsize', WorldSize);
window.customElements.define('cia3-barbarians', Barbarians);
window.customElements.define('cia3-worldseed', WorldSeed);
window.customElements.define('cia3-landmass', LandMass);
window.customElements.define('cia3-oceancoverage', OceanCoverage);
window.customElements.define('cia3-climate', Climate);
window.customElements.define('cia3-temperature', Temperature);
window.customElements.define('cia3-age', Age);
window.customElements.define('cia3-civs', Civs);
window.customElements.define('cia3-civsdev', CivsDev);
window.customElements.define('cia3-civstech', CivTech);
window.customElements.define('cia3-hexdiffall', HexDiffAll);
window.customElements.define('cia3-trade', Trade);
pollNow();
