function delTankConfirm(event, surveyID, draftIndex, tankID) {
    const block = event.target.closest('[hx-headers]')
    const headers = block ? JSON.parse(block.getAttribute('hx-headers')) : {}

    Alpine.store('confirm').show(
        'Delete Tank',
        'Are you sure you want to delete the tank?',
        () => htmx.ajax('DELETE', '/api/v1/survey/' + surveyID + '/tanks/' + draftIndex + '/bw-tank/' + tankID, {
            target: '#tank-row-' + tankID,
            swap: 'delete swap:250ms',
            headers: headers
        }).then(() => {
            Alpine.store('toast').show('Done', 'Tank removed successfully')
        })
    )
}