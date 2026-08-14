// ================================================================
// draft-readings.js
// Only vessel indicators and helpers
// ================================================================

const CX = 80, CY = 118, LINE_R = 42;
const MAX_VISUAL_DEG = 45;
const MIN_DOT_Y = 30;
const MAX_DOT_Y = 210;
const BREADTH = 32.26;
const LIST_LOG_SCALE = 8;
const LIST_AMP = 10;
const TRIM_LOG_SCALE = 4;
const TRIM_AMP = 28;

function scaleLog(value, scale) {
    if (value === 0) { return 0 };
    const sign = value > 0 ? -1 : 1
    return sign * Math.log1p(Math.abs(value) * scale)
}

function clamp(val, min, max) {
    return Math.max(min, Math.min(max, val));
}

function g(id) {
    const v = parseFloat(document.getElementById(id)?.value);
    return isNaN(v) ? null : v;
}

function calc(p) {
    const fp = g('fwd_port-' + p), mp = g('mid_port-' + p), ap = g('aft_port-' + p);
    const fs = g('fwd_stbd-' + p), ms = g('mid_stbd-' + p), as_ = g('aft_stbd-' + p);

    // Trim — visual only, Aft mean minus Fwd mean
    const mF = fp !== null && fs !== null ? (fp + fs) / 2 : null;
    const mA = ap !== null && as_ !== null ? (ap + as_) / 2 : null;
    const oTrim = mF !== null && mA !== null ? mA - mF : null;

    // List — visual only, angle from breadth
    let listDeg = null, listSign = 1;
    if (mp !== null && ms !== null) {
        listDeg = Math.atan(Math.abs(mp - ms) / BREADTH) * 180 / Math.PI;
        listSign = (mp !== null && ms !== null && mp < ms) ? -1 : 1;
    }

    // Update vessel indicator only
    updateVessel(p, listDeg !== null ? listDeg * listSign : 0, oTrim);
}

function updateVessel(p, listDeg, trim) {
    const listLine = document.getElementById(p + '-list-line');
    const trimDot = document.getElementById(p + '-trim-dot');
    if (!listLine || !trimDot) return;

    const listVal = listDeg || 0;
    const listScaled = (listVal > 0 ? -1 : 1) * Math.log1p(Math.abs(listVal) * LIST_LOG_SCALE) * LIST_AMP;
    const visualDeg = clamp(listScaled, -MAX_VISUAL_DEG, MAX_VISUAL_DEG);
    const rad = visualDeg * Math.PI / 180;
    listLine.setAttribute('x1', (CX - LINE_R * Math.cos(rad)).toFixed(1));
    listLine.setAttribute('y1', (CY - LINE_R * Math.sin(rad)).toFixed(1));
    listLine.setAttribute('x2', (CX + LINE_R * Math.cos(rad)).toFixed(1));
    listLine.setAttribute('y2', (CY + LINE_R * Math.sin(rad)).toFixed(1));

    const trimVal = trim || 0;
    const trimScaled = (trimVal > 0 ? 1 : -1) * Math.log1p(Math.abs(trimVal) * TRIM_LOG_SCALE) * TRIM_AMP;
    const dotY = clamp(CY + trimScaled, MIN_DOT_Y, MAX_DOT_Y);
    trimDot.setAttribute('cy', dotY.toFixed(1));
}
//Trigger vessel trim list updates on frontend on page load
document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('[id^="fwd_port-"]').forEach(function (el) {
        const idx = el.id.replace('fwd_port-', '');
        calc(idx);
    });
});
//Trigger vessel trim list updates on frontend on page load for HTMX swap support
function initVessels() {
    document.querySelectorAll('[id^="fwd_port-"]').forEach(function (el) {
        const idx = el.id.replace('fwd_port-', '');
        calc(idx);
    });
}

document.addEventListener('DOMContentLoaded', initVessels);
document.addEventListener('htmx:afterSwap', initVessels);

function draftPage(idx) {
    return {
        fwdP: 0, midP: 0, aftP: 0,
        fwdS: 0, midS: 0, aftS: 0,
        mmc: 0, mtcDraftP: 0,
        mtcDraftN: 0,

        init() {
            this.fwdP = g('fwd_port-' + idx) || 0
            this.midP = g('mid_port-' + idx) || 0
            this.aftP = g('aft_port-' + idx) || 0
            this.fwdS = g('fwd_stbd-' + idx) || 0
            this.midS = g('mid_stbd-' + idx) || 0
            this.aftS = g('aft_stbd-' + idx) || 0
            this.mmc = +this.$el.querySelector('.cp-cell--mmc .cp-val').textContent || 0
            this.updtMTCDraft()
        },

        get allFilled() {
            return this.fwdP > 0 && this.midP > 0 && this.aftP > 0 &&
                this.fwdS > 0 && this.midS > 0 && this.aftS > 0
        },

        updtMTCDraft() {
            if (this.allFilled && this.mmc > 0) {
                this.mtcDraftN = (this.mmc - 0.5)
                this.mtcDraftP = (this.mmc + 0.5)
            }
        }
    }
}

// EventDraftCalc — payload is a concatenation of result fragments (calc
// panel, MMC row, LBM, delta MTC, totals) for every draft, each carrying its
// real DOM id. Not routed through htmx's swap engine (this app's SSE
// listeners never are — see 'toast'/'survey-stats' in main.js), so apply
// each fragment by id manually.
document.addEventListener('DOMContentLoaded', function () {
    if (typeof es === 'undefined') return
    es.addEventListener('draft-calc', function (e) {
        const template = document.createElement('template')
        template.innerHTML = e.data
        template.content.querySelectorAll('[id]').forEach(function (fragment) {
            const target = document.getElementById(fragment.id)
            if (!target) return
            target.replaceWith(fragment)
            if (fragment.classList.contains('calc-panel')) {
                const block = fragment.closest('.draft-block')
                if (block) block.dispatchEvent(new CustomEvent('calc-updated'))
            }
        })
    })
})

function delDraftConfirm(id) {
    Alpine.store('confirm').show(
        'Delete Draft',
        'Are you sure you want to delete the last draft?',
        () => htmx.ajax('DELETE', '/api/v1/survey/' + id + '/draft', { swap: 'none' })
    )
}