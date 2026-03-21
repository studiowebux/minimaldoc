/**
 * PDF Export
 * Handles PDF export using browser's native print functionality
 */

(function() {
    'use strict';

    function init() {
        var btn = document.getElementById('pdf-export-btn');
        if (!btn) return;

        btn.addEventListener('click', function() {
            btn.classList.add('printing');
            btn.textContent = 'Preparing...';

            // Small delay to show feedback before print dialog
            setTimeout(function() {
                window.print();
                btn.classList.remove('printing');
                btn.textContent = 'Export PDF';
            }, 100);
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
