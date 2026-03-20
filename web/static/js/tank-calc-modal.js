(() => {
    const TYPE_LABELS = [
        { trim: 'Trim Rows', volume: null },
        { trim: 'Sounding Correction by Trim', volume: 'Volume Table (trim = 0)' },
        { trim: 'Volume Correction by Trim', volume: 'Base Volume Table (trim = 0)' },
    ];

    const TYPE_VALUES = [
        'Standard (Volume by Trim)',
        'Sounding Correction',
        'Volume Correction',
    ];

    let currentType = 0;

    window.setType = function (idx) {
        currentType = idx;

        document.querySelectorAll('.tcm-type-tab').forEach((t, i) => {
            t.classList.toggle('active', i === idx);
        });

        document.getElementById('calib-type-val').value = TYPE_VALUES[idx];
        document.getElementById('trim-rows-title').textContent = TYPE_LABELS[idx].trim;

        const secVol = document.getElementById('section-volume');
        if (idx > 0) {
            secVol.classList.add('visible');
            document.getElementById('volume-rows-title').textContent = TYPE_LABELS[idx].volume;
        } else {
            secVol.classList.remove('visible');
        }
    };

    window.toggleList = function (checkbox) {
        document.getElementById('section-list').classList.toggle('visible', checkbox.checked);
    };
})();