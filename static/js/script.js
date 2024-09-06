
import {setLinePosition, getLinePosition, setTotalHits, getTotalHits} from './state.js'


function updatePostClick(x, y, linePosition, totalHits = null) {
    moveMiddleLine(linePosition);
    drawExpandingCircle(x, y);
    document.getElementById("linePosition").innerText = linePosition;
    if (totalHits != null) {
        document.getElementById("totalHits").innerText = totalHits
    }
}

function moveMiddleLine(linePosition) {
    const percentage = ((linePosition / -2) + 50)
    const a = (100 - percentage) + "%"
    const b = (percentage) + "%"
    left.style.width = a
    right.style.width = b
}

/** socket stuff */
const socket = new WebSocket("ws://localhost:8080/ws");
socket.binaryType = "arraybuffer"

let left = document.getElementById('leftSide')
let right = document.getElementById('rightSide')


/** socket event handlers */
socket.onopen = function(event) {
    console.log("WebSocket is connected.");
};

socket.onmessage = function(event) {
    const arrayBuffer = event.data;
    if (!(arrayBuffer instanceof ArrayBuffer)) {
        console.error("Received data is not an ArrayBuffer");
        return;
    }
    const dataView = new DataView(arrayBuffer);
    if (dataView.getUint8(0) === 0) {
        const x = dataView.getInt32(1, true); // true for little-endian
        const y = dataView.getInt32(5, true);
        const color = dataView.getUint8(9); // 0 for blue, 1 for red
        const linePosition = dataView.getInt32(10, true);
        setLinePosition(linePosition)
        const totalHits = dataView.getInt32(14, true);
        console.log(`Received remote click at: (${x}, ${y}), color: ${color === 0 ? 'blue' : 'red'}, totalHits: ${totalHits}, linePosition ${linePosition}`, );

        updatePostClick(x,y, linePosition, totalHits)

    } else if (dataView.getUint8(0) === 1) {
        const linePosition = dataView.getInt32(1, true);
        setLinePosition(linePosition)
        setTotalHits(linePosition)
        moveMiddleLine(linePosition)
        console.log(`line position ${linePosition}`)
    }
    return
};

socket.onclose = function(event) {
    console.log("WebSocket is closed.");
};


/** click handlers */
const gameFrame = document.getElementById('gameFrame');
const canvas = document.getElementById('drawingCanvas');
// // Adjust canvas size to match the game frame
canvas.width = gameFrame.clientWidth;
canvas.height = gameFrame.clientHeight;
const ctx = canvas.getContext('2d');
const leftSide = document.getElementById("leftSide");
// const rightSide = document.getElementById("rightSide");



gameFrame.addEventListener('click', function(event) {
    // Get the x, y coordinates relative to the frame
    const rect = gameFrame.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    const inLeft = inSideDiv(event, leftSide)
    let color;
    if (inLeft) {
        color = 0;
        setLinePosition(getLinePosition() + 1)
    } else {
        color = 1;
        setLinePosition(getLinePosition() - 1)
    }
    setTotalHits(getTotalHits() + 1)

    updatePostClick(x, y, getLinePosition(), getTotalHits())

    const buffer = new ArrayBuffer(10);
    const dataView = new DataView(buffer);

    dataView.setUint8(0, 0);
    dataView.setInt32(1, x, true); // true for little-endian
    dataView.setInt32(5, y, true); // true for little-endian
    dataView.setUint8(9, color);

    socket.send(buffer);
});

function inSideDiv(event, sideDiv) {
    return event.target === sideDiv || sideDiv.contains(event.target)
}

function drawExpandingCircle(x, y) {
    let radius = 0;
    const maxRadius = 50;
    const expansionRate = 2;
    const fadeRate = 0.05;
    let alpha = 1.0;

    function animate() {
        if (radius < maxRadius) {
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            ctx.beginPath();
            ctx.arc(x, y, radius, 0, Math.PI * 2);
            ctx.fillStyle = `rgba(255, 255, 255, ${alpha})`;
            ctx.fill();

            radius += expansionRate;
            alpha -= fadeRate;

            requestAnimationFrame(animate);
        } else {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
        }
    }

    animate();
}

/** keep alive ping. If this doesn't go off server will terminate websocket */
function sendPing() {
    if (socket.readyState === WebSocket.OPEN) {
        console.log("ping")
        socket.send(new Uint8Array([0xFF]));
    }
}

setInterval(sendPing, 3000);

