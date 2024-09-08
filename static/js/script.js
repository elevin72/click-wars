
import {setLinePosition, getLinePosition, setTotalHits, getTotalHits} from './state.js'
import {drawExpandingCircle} from './animation.js'


function updatePostClick(x, y, linePosition, totalHits = null) {
    moveMiddleLine(linePosition);
    drawExpandingCircle(x, y);
    document.getElementById("linePosition").innerText = linePosition;
    if (totalHits != null) {
        setTotalHits(totalHits)
        document.getElementById("totalHits").innerText = totalHits
    }
}

function moveMiddleLine(linePosition) {
    const percentage = ((linePosition / -2) + 50)
    if (percentage > 90 || percentage < 10) {
        return;
    }
    setLinePosition(linePosition)
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
        const x = dataView.getFloat32(1, true); // true for little-endian
        const y = dataView.getFloat32(5, true);
        const color = dataView.getUint8(9); // 0 for blue, 1 for red
        const linePosition = dataView.getInt32(10, true);
        const totalHits = dataView.getInt32(14, true);
        console.log(dataView)
        console.log(`Received remote click at: (${x}, ${y}), color: ${color === 0 ? 'blue' : 'red'}, totalHits: ${totalHits}, linePosition ${linePosition}`, );

        updatePostClick(x,y, linePosition, totalHits)

    } else if (dataView.getUint8(0) === 1) {
        const linePosition = dataView.getInt32(1, true);
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
// const canvas = document.getElementById('drawingCanvas');
// const ctx = canvas.getContext('2d');
// // Adjust canvas size to match the game frame
// canvas.width = gameFrame.clientWidth;
// canvas.height = gameFrame.clientHeight;
const leftSide = document.getElementById("leftSide");
// const rightSide = document.getElementById("rightSide");


addEventListener("resize", onResize)

function onResize(event) {
    console.log("resizing")
    const gameFrame = document.getElementById('gameFrame');
    let canvas = document.getElementById('drawingCanvas');
    canvas.width = gameFrame.clientWidth;
    canvas.height = gameFrame.clientHeight;

}


function normalizeClickLocation(x, y, rect) {
    return {
        x: (x - rect.left) / (rect.right - rect.left),
        y: (y - rect.top) / (rect.bottom - rect.top),
    }
}


gameFrame.addEventListener('click', function(event) {
    // Get the x, y coordinates relative to the frame
    const canvas = document.getElementById('drawingCanvas');
    const rect = canvas.getBoundingClientRect();
    // percentage from the left, top
    const x = (event.clientX - rect.left) / (rect.right - rect.left);
    const y = (event.clientY - rect.top) / (rect.bottom - rect.top);

    const inLeft = insideDiv(event, leftSide)
    let color, newLinePosition;
    const oldLinePosition = getLinePosition()
    if (inLeft) {
        color = 0;
        newLinePosition = oldLinePosition + 1;
    } else {
        color = 1;
        newLinePosition = oldLinePosition - 1;
    }

    updatePostClick(x, y, newLinePosition, getTotalHits() + 1)
    console.log(`Local click at: (${x}, ${y}), color: ${color === 0 ? 'blue' : 'red'}, totalHits: ${getTotalHits()}, linePosition ${getLinePosition()}`, );

    const buffer = new ArrayBuffer(10);
    const dataView = new DataView(buffer);

    dataView.setUint8(0, 0);
    dataView.setFloat32(1, x, true); // true for little-endian
    dataView.setFloat32(5, y, true); // true for little-endian
    dataView.setUint8(9, color);
    console.log(dataView)

    socket.send(buffer);
});

function insideDiv(event, sideDiv) {
    return event.target === sideDiv || sideDiv.contains(event.target)
}

// function drawExpandingCircle(x, y) {
//     const canvas = document.getElementById('drawingCanvas');
//     const ctx = canvas.getContext('2d');
//     let radius = 0;
//     const maxRadius = 50;
//     const expansionRate = 2;
//     const fadeRate = 0.05;
//     let alpha = 1.0;
//     const rect = canvas.getBoundingClientRect();
//     x = x * (rect.right - rect.left)
//     y = y * (rect.bottom - rect.top)

//     function animate() {
//         if (radius < maxRadius) {
//             ctx.clearRect(0, 0, canvas.width, canvas.height);

//             ctx.beginPath();
//             ctx.arc(x, y, radius, 0, Math.PI * 2);
//             ctx.fillStyle = `rgba(255, 255, 255, ${alpha})`;
//             ctx.fill();

//             radius += expansionRate;
//             alpha -= fadeRate;

//             requestAnimationFrame(animate);
//         } else {
//             ctx.clearRect(0, 0, canvas.width, canvas.height);
//         }
//     }

//     animate();
// }

/** keep alive ping. If this doesn't go off server will terminate websocket */
function sendPing() {
    if (socket.readyState === WebSocket.OPEN) {
        console.log("ping")
        socket.send(new Uint8Array([0xFF]));
    }
}

setInterval(sendPing, 3000);
document.getElementById("linePosition").innerText = getLinePosition();
onResize(null)

