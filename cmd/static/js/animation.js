const circles = [];

let animationRunning = false;

const MAX_RADIUS = 50;
const EXPANSION_RATE = 2;
const FADE_RATE = 0.05;

export function drawExpandingCircle(x, y) {
    let radius = 0;
    let alpha = 1.0;
    const canvas = document.getElementById('drawingCanvas');
    const rect = canvas.getBoundingClientRect();
    x = x * (rect.right - rect.left);
    y = y * (rect.bottom - rect.top);

    circles.push({ x, y, radius, alpha });

    if (!animationRunning) {
        animate();
    }
}

function animate() {
    const canvas = document.getElementById('drawingCanvas');
    const ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    circles.forEach((circle, index) => {
        if (circle.radius < MAX_RADIUS) {
            ctx.beginPath();
            ctx.arc(circle.x, circle.y, circle.radius, 0, Math.PI * 2);
            ctx.fillStyle = `rgba(255, 255, 255, ${circle.alpha})`;
            ctx.fill();

            circle.radius += EXPANSION_RATE;
            circle.alpha -= FADE_RATE;
        } else {
            circles.splice(index, 1);
        }
    });

    if (circles.length > 0) {
        animationRunning = true;
        requestAnimationFrame(animate);
    } else {
        animationRunning = false;
    }
}