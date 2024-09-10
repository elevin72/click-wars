
import { getLinePosition, getTotalHits, updatePostClick } from './state.js'
import { socket } from './socket.js'
import { drawExpandingCircle } from './animation.js';


function onClick(event) {
    const canvas = document.getElementById('drawingCanvas');
    const rect = canvas.getBoundingClientRect();
    const x = (event.clientX - rect.left) / (rect.right - rect.left);
    const y = (event.clientY - rect.top) / (rect.bottom - rect.top);

    let color, newLinePosition;
    const oldLinePosition = getLinePosition()
    if (event.target === leftSide || leftSide.contains(event.target)) {
        color = 0;
        newLinePosition = oldLinePosition + 1;
    } else {
        color = 1;
        newLinePosition = oldLinePosition - 1;
    }

    updatePostClick(newLinePosition, getTotalHits() + 1)
    drawExpandingCircle(x, y)

    console.log(`Local click at: (${x}, ${y}), color: ${color === 0 ? 'blue' : 'red'}, totalHits: ${getTotalHits()}, linePosition ${getLinePosition()}`, );

    const buffer = new ArrayBuffer(10);
    const dataView = new DataView(buffer);
    dataView.setUint8(0, 0);
    dataView.setFloat32(1, x, true); // true for little-endian
    dataView.setFloat32(5, y, true); // true for little-endian
    dataView.setUint8(9, color);

    socket.send(buffer);
};


function sendPing() {
    if (socket.readyState === WebSocket.OPEN) {
        console.log("ping")
        socket.send(new Uint8Array([0xFF]));
    }
}

function onResize(event) {
    console.log("resizing")
    const gameFrame = document.getElementById('gameFrame');
    let canvas = document.getElementById('drawingCanvas');
    canvas.width = gameFrame.clientWidth;
    canvas.height = gameFrame.clientHeight;

}
onResize(null)

setInterval(sendPing, 3000);
addEventListener("resize", onResize)
document.getElementById('gameFrame').addEventListener('click', onClick)


