/**
 * Mobile Menu Toggle
 * Handles the mobile hamburger menu for navigation
 */

(function() {
    const toggle = document.getElementById('mobile-menu-toggle');
    const sidebar = document.getElementById('sidebar-left');
    const backdrop = document.getElementById('mobile-menu-backdrop');

    if (!toggle || !sidebar || !backdrop) return;

    function openMenu() {
        sidebar.classList.add('open');
        backdrop.classList.add('open');
        toggle.classList.add('open');
        document.body.style.overflow = 'hidden';
    }

    function closeMenu() {
        sidebar.classList.remove('open');
        backdrop.classList.remove('open');
        toggle.classList.remove('open');
        document.body.style.overflow = '';
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
    const navLinks = sidebar.querySelectorAll('.nav-link');
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
