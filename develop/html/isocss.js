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
        textDiv.innerText = "🌲🌲🌲";
	}
}
// TODO: put this in try/catch with friendly output for non-web-component browsers
customElements.define('map-tile', MapTile);
