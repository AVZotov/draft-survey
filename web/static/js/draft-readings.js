// ================================================================
// draft-readings.js
// Only vessel indicator — all calculations are server-side
// ================================================================

const CX = 80, CY = 118, LINE_R = 42;
const LIST_SCALE = 4;
const TRIM_SCALE = 35;
const MAX_VISUAL_DEG = 45;
const MIN_DOT_Y = 30;
const MAX_DOT_Y = 210;
const BREADTH = 32.26;

function clamp(val, min, max) {
    return Math.max(min, Math.min(max, val));
}

function g(id) {
    const v = parseFloat(document.getElementById(id)?.value);
    return isNaN(v) ? null : v;
}

function calc(p) {
    const fp = g(p + '-fp'), mp = g(p + '-mp'), ap = g(p + '-ap');
    const fs = g(p + '-fs'), ms = g(p + '-ms'), as_ = g(p + '-as');

    // Trim — visual only, Aft mean minus Fwd mean
    const mF = fp !== null && fs !== null ? (fp + fs) / 2 : null;
    const mA = ap !== null && as_ !== null ? (ap + as_) / 2 : null;
    const oTrim = mF !== null && mA !== null ? mA - mF : null;

    // List — visual only, angle from breadth
    let listDeg = null, listSign = 1;
    if (mp !== null && ms !== null) {
        listDeg = Math.atan(Math.abs(mp - ms) / BREADTH) * 180 / Math.PI;
        listSign = mp < ms ? 1 : -1;
    }

    // Update vessel indicator only
    updateVessel(p, listDeg !== null ? listDeg * listSign : 0, oTrim);
}

function updateVessel(p, listDeg, trim) {
    const listLine = document.getElementById(p + '-list-line');
    const trimDot = document.getElementById(p + '-trim-dot');
    if (!listLine || !trimDot) return;

    const visualDeg = clamp((listDeg || 0) * LIST_SCALE, -MAX_VISUAL_DEG, MAX_VISUAL_DEG);
    const rad = visualDeg * Math.PI / 180;
    listLine.setAttribute('x1', (CX - LINE_R * Math.cos(rad)).toFixed(1));
    listLine.setAttribute('y1', (CY - LINE_R * Math.sin(rad)).toFixed(1));
    listLine.setAttribute('x2', (CX + LINE_R * Math.cos(rad)).toFixed(1));
    listLine.setAttribute('y2', (CY + LINE_R * Math.sin(rad)).toFixed(1));

    const dotY = clamp(CY + (trim || 0) * TRIM_SCALE, MIN_DOT_Y, MAX_DOT_Y);
    trimDot.setAttribute('cy', dotY.toFixed(1));
}
//Trigger vessel trim list updates on frontend on page load
document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('[id$="-fp"]').forEach(function (el) {
        const idx = el.id.replace('-fp', '');
        calc(idx);
    });
});
//Trigger vessel trim list updates on frontend on page load for HTMX swap support
function initVessels() {
    document.querySelectorAll('[id$="-fp"]').forEach(function (el) {
        const idx = el.id.replace('-fp', '');
        calc(idx);
    });
}

document.addEventListener('DOMContentLoaded', initVessels);
document.addEventListener('htmx:afterSwap', initVessels);