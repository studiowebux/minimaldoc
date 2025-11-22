/**
 * External Links
 * Opens external links in a new tab with target="_blank"
 */

(function() {
    function init() {
        const markdownBody = document.querySelector('.markdown-body');
        if (!markdownBody) return;

        const links = markdownBody.querySelectorAll('a[href]');

        links.forEach(function(link) {
            const href = link.getAttribute('href');

            // Check if link is external
            if (href && (href.startsWith('http://') || href.startsWith('https://'))) {
                const currentDomain = window.location.hostname;
                try {
                    const linkUrl = new URL(href);
                    if (linkUrl.hostname !== currentDomain) {
                        link.setAttribute('target', '_blank');
                        link.setAttribute('rel', 'noopener noreferrer');
                    }
                } catch (e) {
                    // Invalid URL, skip
                }
            }
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
