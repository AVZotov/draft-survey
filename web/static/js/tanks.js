function delTankConfirm(event, surveyID, draftIndex, tankID, kind) {
    // The success toast comes from the server via SSE (TankService.DeleteBW/DeleteFW's
    // Outcome.Toast) — do not also trigger one here, that would show it twice.
    Alpine.store('confirm').show(
        'Delete Tank',
        'Are you sure you want to delete the tank?',
        () => htmx.ajax('DELETE', '/api/v1/survey/' + surveyID + '/tanks/' + draftIndex + '/' + kind + '/' + tankID, {
            target: '#tank-row-' + tankID,
            swap: 'delete swap:250ms'
        })
    )
}

// EventTankCalc — payload is a concatenation of the BW/FW block-header totals
// and the action-bar totals for the draft that changed, each carrying its
// real DOM id. Not routed through htmx's swap engine (this app's SSE
// listeners never are — see 'toast'/'survey-stats' in main.js and 'draft-calc'
// in draft-readings.js), so apply each fragment by id manually.
document.addEventListener('DOMContentLoaded', function () {
    if (typeof es === 'undefined') return
    es.addEventListener('tank-calc', function (e) {
        const template = document.createElement('template')
        template.innerHTML = e.data
        template.content.querySelectorAll('[id]').forEach(function (fragment) {
            const target = document.getElementById(fragment.id)
            if (!target) return
            target.replaceWith(fragment)
        })
    })
})
