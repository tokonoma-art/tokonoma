import './index.scss';
import './artwork-element.ts';

interface Message {
  artworkKey: string;
}

function buildWSURL() {
  const loc = window.location;
  const protocol = (loc.protocol === 'https' ? 'wss' : 'ws');
  return `${protocol}://${loc.host}/ws`;
}

const artworkElement = document.getElementById('current-artwork');

const ws = new WebSocket(buildWSURL());

ws.onopen = () => {
  console.log('Connected!');
  ws.send('Hello, server');
};

ws.onmessage = (evt) => {
  const { artworkKey }: Message = JSON.parse(evt.data);
  artworkElement.setAttribute('src', `/artworks/${artworkKey}.artbundle`);
};
