// Navigation collapse/expand functionality
document.addEventListener('DOMContentLoaded', function() {
    const navToggles = document.querySelectorAll('.nav-toggle-btn');

    console.log('Nav collapse initialized, found', navToggles.length, 'toggles');

    navToggles.forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();

            const navItem = this.closest('.nav-item');

            console.log('Toggle clicked, navItem:', navItem);

            if (navItem) {
                // Toggle collapsed state
                const wasCollapsed = navItem.classList.contains('collapsed');
                navItem.classList.toggle('collapsed');

                console.log('Toggled from', wasCollapsed, 'to', !wasCollapsed);

                // Save state to localStorage
                saveNavState();
            }
        });
    });

    // Restore navigation state from localStorage
    restoreNavState();
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

function restoreNavState() {
    try {
        const collapsedSections = JSON.parse(localStorage.getItem('navCollapsedSections') || '[]');

        document.querySelectorAll('.nav-item-with-children').forEach(container => {
            const titleEl = container.querySelector('.nav-section, .nav-link-parent');
            const navItem = container.closest('.nav-item');

            if (titleEl && navItem && collapsedSections.includes(titleEl.textContent.trim())) {
                // Don't collapse if it's the active section
                if (!navItem.classList.contains('active')) {
                    navItem.classList.add('collapsed');
                }
            }
        });
    } catch (e) {
        // Ignore errors in localStorage
    }
}
