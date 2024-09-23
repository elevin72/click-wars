let linePosition = parseInt(document.getElementById("linePosition").innerText)
let totalHits = parseInt(document.getElementById("totalHits").innerText)

export function getTotalHits() {
    return totalHits;
}

export function getLinePosition() {
    return linePosition;
}

export function updateTotalHits(newTotalHits) {
    totalHits = newTotalHits;
    document.getElementById("totalHits").innerText = newTotalHits;
}

export function updateLinePostion(newLinePosition) {
    linePosition = newLinePosition;
    const percentage = ((linePosition * 2) + 50)
    let leftPerecentage, rightPercentage;
    if (percentage > 90 || percentage < 10) {
        leftPerecentage = "10%" ? percentage <= 10 : "90%"
        rightPercentage = "90%" ? percentage >= 90 : "10%"
    } else {
        leftPerecentage = percentage.toString() + "%"
        rightPercentage = (100 - percentage).toString() + "%"
    }

    if (percentage > 90 || percentage < 10) {
        return;
    }

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
    document.getElementById("leftCount").innerText = (totalHits + linePosition) / 2;
    document.getElementById("rightCount").innerText = (totalHits - linePosition) / 2;
}



