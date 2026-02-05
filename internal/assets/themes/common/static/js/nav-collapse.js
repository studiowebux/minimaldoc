// Navigation collapse/expand functionality
document.addEventListener('DOMContentLoaded', function() {
    const navToggles = document.querySelectorAll('.nav-toggle-btn');

    navToggles.forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();

            const navItem = this.closest('.nav-item');
            if (navItem) {
                navItem.classList.toggle('collapsed');
                this.setAttribute('aria-expanded', !navItem.classList.contains('collapsed'));
                saveNavState();
            }
        });
    });
});

function saveNavState() {
    const collapsedSections = [];
    document.querySelectorAll('.nav-item.collapsed .nav-item-with-children').forEach(container => {
        const titleEl = container.querySelector('.nav-section, .nav-link-parent');
        if (titleEl) {
            collapsedSections.push(titleEl.textContent.trim());
        }
    });
    localStorage.setItem('navCollapsedSections', JSON.stringify(collapsedSections));
}

