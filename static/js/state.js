import { drawExpandingCircle } from "./animation.js";
let linePosition = 0;
let totalHits = 0;

export function setTotalHits(newTotalHits) {
    totalHits = newTotalHits;
}

export function getTotalHits() {
    return totalHits;
}

export function setLinePosition(newLinePosition) {
    linePosition = newLinePosition;
}

export function getLinePosition() {
    return linePosition;
}

export function updateTotalHits(newTotalHits) {
    setTotalHits(newTotalHits);
    document.getElementById("totalHits").innerText = newTotalHits;
}

export function updateLinePostion(linePosition) {
    const percentage = ((linePosition / -2) + 50)
    if (percentage > 90 || percentage < 10) {
        return;
    }

    setLinePosition(linePosition)

    const leftPerecentage = (100 - percentage) + "%"
    const rightPercentage = (percentage) + "%"

    // update screen elements
    document.getElementById("leftSide").style.width = leftPerecentage;
    document.getElementById("rightSide").style.width = rightPercentage;
    document.getElementById("linePosition").innerText = linePosition;
}

export function updatePostClick(linePosition, totalHits) {
    if (linePosition != null) {
        updateLinePostion(linePosition)
    }
    if (totalHits != null) {
        updateTotalHits(totalHits)
    }
}



