/**
 * Claude Assist
 * Opens Claude's web UI with the current page URL
 */

(function() {
    'use strict';

    function init() {
        var btn = document.getElementById('claude-assist-btn');
        if (!btn) return;

        btn.addEventListener('click', function() {
            var pageUrl = window.location.href;
            var customPrompt = btn.dataset.prompt || 'Please help me understand this documentation:';

            // Build the prompt with URL
            var prompt = customPrompt + '\n\n' + pageUrl;

            // Open Claude with the prompt
            var url = 'https://claude.ai/new?q=' + encodeURIComponent(prompt);
            window.open(url, '_blank', 'noopener,noreferrer');
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
