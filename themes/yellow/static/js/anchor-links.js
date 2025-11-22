/**
 * Anchor Links
 * Adds clickable anchor links to all headings
 */

(function() {
    function init() {
        const markdownBody = document.querySelector('.markdown-body');
        if (!markdownBody) return;

        const headings = markdownBody.querySelectorAll('h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]');

        headings.forEach(function(heading) {
            const id = heading.getAttribute('id');
            if (!id) return;

            heading.style.position = 'relative';
            heading.style.cursor = 'pointer';

            const anchor = document.createElement('a');
            anchor.className = 'header-anchor';
            anchor.href = '#' + id;
            anchor.textContent = '#';
            anchor.setAttribute('aria-label', 'Link to ' + heading.textContent);

            anchor.addEventListener('click', function(e) {
                e.preventDefault();
                const url = new URL(window.location);
                url.hash = '#' + id;
                window.history.pushState({}, '', url);
                heading.scrollIntoView({ behavior: 'smooth', block: 'start', inline: 'nearest' });
            });

            heading.appendChild(anchor);
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
