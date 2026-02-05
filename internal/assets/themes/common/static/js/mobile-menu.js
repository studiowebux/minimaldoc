/**
 * Mobile Menu Toggle
 * Handles the mobile hamburger menu for navigation
 */

(function() {
    const toggle = document.getElementById('mobile-menu-toggle');
    const sidebar = document.getElementById('sidebar-left') || document.getElementById('openapi-sidebar-left');
    const backdrop = document.getElementById('mobile-menu-backdrop');

    if (!toggle || !sidebar || !backdrop) return;

    function openMenu() {
        sidebar.classList.add('open');
        backdrop.classList.add('open');
        toggle.classList.add('open');
        toggle.setAttribute('aria-expanded', 'true');
        backdrop.setAttribute('aria-hidden', 'false');
        document.body.style.overflow = 'hidden';

        // Focus first focusable element in sidebar
        var firstFocusable = sidebar.querySelector('a, button, input, [tabindex]:not([tabindex="-1"])');
        if (firstFocusable) {
            firstFocusable.focus();
        }
    }

    function closeMenu() {
        sidebar.classList.remove('open');
        backdrop.classList.remove('open');
        toggle.classList.remove('open');
        toggle.setAttribute('aria-expanded', 'false');
        backdrop.setAttribute('aria-hidden', 'true');
        document.body.style.overflow = '';

        // Return focus to toggle button
        toggle.focus();
    }

    toggle.addEventListener('click', function() {
        if (sidebar.classList.contains('open')) {
            closeMenu();
        } else {
            openMenu();
        }
    });

    backdrop.addEventListener('click', closeMenu);

    // Close menu when clicking a navigation link
    const navLinks = sidebar.querySelectorAll('.nav-link, .openapi-endpoint-link');
    navLinks.forEach(function(link) {
        link.addEventListener('click', closeMenu);
    });

    // Close menu on window resize if viewport becomes wider
    let resizeTimer;
    window.addEventListener('resize', function() {
        clearTimeout(resizeTimer);
        resizeTimer = setTimeout(function() {
            if (window.innerWidth > 768 && sidebar.classList.contains('open')) {
                closeMenu();
            }
        }, 250);
    });
})();
