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
        visible: false,

        show(header, message) {
            this.header = header;
            this.message = message;
            this.visible = true;
        },

        ok() {
            this.visible = false;
            this.header = '';
            this.message = '';
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