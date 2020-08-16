
/* eslint-disable no-param-reassign */
function setupTextArtwork(container: HTMLElement, _artwork: Artwork, configuration: ArtworkTextV1) {
  const style = document.createElement('style');
  style.innerHTML = `
  .artwork-container {
    width: 100%;
    height: 100%;
    overflow: auto; /* avoids collapsing margins */
    display: flex;
    justify-content: center;
    align-items: center;
  }

  p {
    text-align: center;
    font-size: 4rem;
    font-family: Avenir;
    font-weight: 400;
  }

  `;
  container.appendChild(style);
  const p = document.createElement('p');
  p.textContent = configuration.text;
  const { 'text-color': textColor = '#FFF', 'background-color': backgroundColor = '#000' } = configuration;
  p.style.color = textColor as string;
  container.style.backgroundColor = backgroundColor as string;
  container.appendChild(p);
}

function setupImageArtwork(
  src: string,
  container: HTMLElement,
  _artwork: Artwork,
  configuration: ArtworkImageV1,
) {
  const imageRelativePath = <string>configuration.image;
  const { 'background-color': backgroundColor = '#000' } = configuration;

  const style = document.createElement('style');
  style.innerHTML = `
  .artwork-container {
    width: 100%;
    height: 100%;
    overflow: auto; /* avoids collapsing margins */
    background: ${backgroundColor as string} url(${src}/${imageRelativePath});
    background-size: cover;
    background-position: center center;
  }
  `;
  container.appendChild(style);
}
/* eslint-enable no-param-reassign */

function setupArtwork(src: string, container: HTMLElement, artwork: Artwork) {
  if (artwork.format === 'text') {
    setupTextArtwork(container, artwork, <ArtworkTextV1>artwork.configuration);
  } else if (artwork.format === 'image') {
    setupImageArtwork(src, container, artwork, <ArtworkImageV1>artwork.configuration);
  }
}

class ArtWorkElement extends HTMLElement {
  constructor() {
    super(); // always call super first in constructor

    // Create a shadow root & container
    this.attachShadow({ mode: 'open' });

    // Update initial artwork (or error message…)
    this.updateArtwork(this.getAttribute('src'));
  }

  static get observedAttributes() { return ['src']; }

  attributeChangedCallback(name: string, oldValue: string, newValue: string) {
    if (name === 'src') this.updateArtwork(newValue);
  }

  showError(message: string) {
    const p = document.createElement('p');
    p.textContent = message;
    p.style.color = '#A00';
    p.classList.add('artwork-error');
    this.shadowRoot.appendChild(p);
  }

  updateArtwork(src: string) {
    const shadow = this.shadowRoot;

    // Cleanup
    const previousContainers = shadow.querySelectorAll('.artwork-container');
    previousContainers.forEach((e) => e.remove());
    const previousErrors = shadow.querySelectorAll('.artwork-error');
    previousErrors.forEach((e) => e.remove());

    if (src === null) {
      this.showError('No artwork selected.');
    } else {
      // Create container
      const container = document.createElement('div');
      container.classList.add('artwork-container');
      shadow.appendChild(container);

      // Fetch artbundle info
      const http = new XMLHttpRequest();
      http.open('GET', `${src}/info.json`);
      http.send();
      http.onreadystatechange = () => {
        if (http.readyState === 4) {
          if (http.status === 200) {
            const artbundle: Artbundle = JSON.parse(http.responseText);
            console.log(artbundle);
            setupArtwork(src, container, artbundle.artwork[0]);
          } else {
            this.showError(`Cannot find artwork '${src}'.`);
          }
        }
      };
    }
  }
}

customElements.define('artwork-element', ArtWorkElement);
