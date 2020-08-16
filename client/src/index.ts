import './index.scss';
import './artwork-element.ts';

interface Message {
  artbundle: string;
  rotation: number;
}

function buildWSURL() {
  const loc = window.location;
  const protocol = (loc.protocol === 'https' ? 'wss' : 'ws');
  return `${protocol}://${loc.host}/ws`;
}

const body = document.getElementsByTagName('body')[0];

function updateRotation(rotation: number) {
  body.classList.remove('rotation-0', 'rotation-90', 'rotation-180', 'rotation-270');
  if ([0, 90, 180, 270].includes(rotation)) {
    body.classList.add(`rotation-${rotation}`);
  } else {
    body.classList.add('rotation-0');
  }
}

const artworkElement = document.getElementById('current-artwork');

const ws = new WebSocket(buildWSURL());

ws.onopen = () => {
  // console.log('Connected!');
  // ws.send('Hello, server');
};

ws.onmessage = (evt) => {
  const { artbundle, rotation }: Message = JSON.parse(evt.data);
  updateRotation(rotation);
  artworkElement.setAttribute('src', `/artworks/${artbundle}`);
};
