import { updatePostClick } from "./state.js";
import { drawExpandingCircle } from "./animation.js";

/** socket stuff */
export const socket = new WebSocket("ws://localhost:8080/ws");
socket.binaryType = "arraybuffer"


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
        console.log(`Received remote click at: (${x}, ${y}), color: ${color === 0 ? 'blue' : 'red'}, totalHits: ${totalHits}, linePosition ${linePosition}`, );
        updatePostClick(linePosition, totalHits)
        drawExpandingCircle(x, y);

    } else if (dataView.getUint8(0) === 1) {
        const linePosition = dataView.getInt32(1, true);
        const totalHits = dataView.getInt32(5, true);
        console.log(`totalHits: ${totalHits}, linePosition ${linePosition}`, );
        updatePostClick(linePosition, totalHits)
    }
    return
};

socket.onclose = function(event) {
    console.log("WebSocket is closed.");
};