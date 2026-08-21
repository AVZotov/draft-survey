function delTankConfirm(event) {
    const url = event.currentTarget.dataset.deleteUrl
    const tankID = event.currentTarget.dataset.tankId
    // The success toast comes from the server via SSE (TankService.DeleteBW/DeleteFW's
    // Outcome.Toast) — do not also trigger one here, that would show it twice.
    Alpine.store('confirm').show(
        'Delete Tank',
        'Are you sure you want to delete the tank?',
        () => htmx.ajax('DELETE', url, {
            target: '#tank-row-' + tankID,
            swap: 'delete swap:250ms'
        })
    )
}

// EventTankCalc — payload is a concatenation of the BW/FW block-header totals,
// the action-bar totals, and the full BW/FW <tbody> rows for the draft that
// changed, each carrying its real DOM id. Not routed through htmx's swap
// engine (this app's SSE listeners never are — see 'toast'/'survey-stats' in
// main.js and 'draft-calc' in draft-readings.js), so apply each fragment by
// id manually. The payload must be parsed inside a <table> wrapper — the
// HTML parser silently drops <tbody>/<tr> tags outside a table ancestor
// (including inside a bare <template>), which otherwise makes the tank rows
// vanish while the plain <div> fragments (headers, totals) still apply.
//
// htmx.process() is required after each replaceWith(): replacing a node with
// raw parsed HTML happens outside htmx's own request/swap cycle, so htmx
// never binds hx-put/hx-get/hx-trigger on the new elements — the row's
// move-up/down buttons, autosave, and corrections-modal trigger would only
// work once (on the original server-rendered row) and go dead on every
// subsequent SSE-driven replace. See htmx docs: "Processes new content,
// enabling htmx behavior. Useful when content is added to the DOM outside of
// the normal htmx request cycle."
document.addEventListener('DOMContentLoaded', function () {
    if (typeof es === 'undefined') return
    es.addEventListener('tank-calc', function (e) {
        const container = document.createElement('div')
        container.innerHTML = '<table>' + e.data + '</table>'
        container.querySelectorAll('[id]').forEach(function (fragment) {
            const target = document.getElementById(fragment.id)
            if (!target) return
            target.replaceWith(fragment)
            htmx.process(fragment)
        })
    })
})
