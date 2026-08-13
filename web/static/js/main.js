// SSE FOUNDATION
const es = new EventSource('/events')

es.addEventListener('toast', (e) => {
    document.getElementById('toast-container').insertAdjacentHTML('beforeend', e.data)
})

es.addEventListener('alert', (e) => {
    const d = JSON.parse(e.data)
    const show = () => {
        const header = document.getElementById('app-alert-header')
        const message = document.getElementById('app-alert-message')
        const dialog = document.getElementById('app-alert')
        if (!header || !message || !dialog) {
            requestAnimationFrame(show)
            return
        }
        header.textContent = d.header
        message.textContent = d.message
        dialog.showModal()
    }
    show()
})

es.addEventListener('survey-stats', (e) => {
    const el = document.getElementById('survey-stats')
    if (el) el.innerHTML = e.data
})

// NEW SSE FOUNDATION — add above
// --------------------------------

// DELETE CANDIDATES below

document.addEventListener('alpine:init', () => {
    Alpine.store('toast', {
        header: '',
        message: '',
        visible: false,

        show(header, message) {
            this.header = header;
            this.message = message;
            this.visible = true;
            setTimeout(() => {
                this.header = '';
                this.message = '';
                this.visible = false
            }, 2000)
        }
    })

    Alpine.store('alert', {
        header: '',
        message: '',

        show(header, message) {
            this.header = header
            this.message = message
            document.getElementById('app-alert').showModal()
        },

        ok() {
            this.header = ''
            this.message = ''
            document.getElementById('app-alert').close()
        }
    })

    Alpine.store('confirm', {
        header: '',
        message: '',
        visible: false,
        onConfirm: null,

        show(header, message, callback) {
            this.header = header
            this.message = message
            this.onConfirm = callback
            document.getElementById('app-confirm').showModal()
        },


        ok() {
            if (this.onConfirm) {
                this.onConfirm()
            }
            this.reset()
        },

        cancel() {
            this.reset()
        },

        reset() {
            this.header = ''
            this.message = ''
            this.visible = false
            this.onConfirm = null
            document.getElementById('app-confirm').close()
        }
    })
})


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

//Inline usage in templ components do not remove
function requireFields(ids) {
    for (const id of ids) {
        const el = document.getElementById(id)
        if (!el || el.value === '') {
            event.preventDefault();
            return;
        }
    }
}