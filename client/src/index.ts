import './index.scss';

interface Message {
  artworkKey: string;
}

function buildWSURL() {
  const loc = window.location;
  const protocol = (loc.protocol === 'https' ? 'wss' : 'ws');
  return `${protocol}://${loc.host}/ws`;
}

const ws = new WebSocket(buildWSURL());

const current = <HTMLElement>document.getElementsByClassName('current')[0];
const currentArtworkKeySpan = document.getElementsByClassName('current-artwork-key')[0];

ws.onopen = () => {
  console.log('Connected!');
  ws.send('Hello, server');
};

ws.onmessage = (evt) => {
  const { artworkKey }: Message = JSON.parse(evt.data);
  const http = new XMLHttpRequest();
  const url = `/artworks/${artworkKey}.json`;
  http.open('GET', url);
  http.send();
  http.onreadystatechange = () => {
    if (http.readyState === 4 && http.status === 200) {
      const artwork = JSON.parse(http.responseText);
      console.log(artwork);
      currentArtworkKeySpan.textContent = artwork.key;
      current.style.backgroundColor = artwork.color;
    }
  };
};
