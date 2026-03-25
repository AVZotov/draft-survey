document.addEventListener('wheel', function () {
    if (document.activeElement.type === 'number') {
        document.activeElement.blur();
    }
});

document.addEventListener('htmx:afterSwap', function (event) {
    if (event.detail.target.id === 'app-modal-content') {
        Alpine.destroyTree(event.detail.target)
        Alpine.initTree(event.detail.target)
    }
})