/**
 * ScrollSpy
 * Highlights the active section in the table of contents as the user scrolls
 */

class ScrollSpy {
    constructor() {
        this.sections = [];
        this.tocLinks = [];
        this.currentActiveLink = null;
        this.ticking = false;
        this.offset = 80; // Offset from top of viewport
        this.init();
    }

    init() {
        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.setup());
        } else {
            this.setup();
        }
    }

    setup() {
        // Find all headings with IDs
        this.sections = Array.from(
            document.querySelectorAll('h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]')
        );

        // Find all TOC links
        this.tocLinks = Array.from(
            document.querySelectorAll('.toc-link')
        );

        if (this.sections.length === 0 || this.tocLinks.length === 0) {
            return; // No sections or TOC, nothing to do
        }

        // Add click handlers to TOC links to prevent horizontal scrolling
        this.tocLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const href = link.getAttribute('href');
                if (href && href.startsWith('#')) {
                    const id = href.substring(1);
                    const target = document.getElementById(id);
                    if (target) {
                        // Update URL hash
                        window.history.pushState({}, '', href);
                        // Scroll to target without horizontal realignment
                        target.scrollIntoView({ behavior: 'smooth', block: 'start', inline: 'nearest' });
                    }
                }
            });
        });

        // Bind scroll event with throttling
        window.addEventListener('scroll', () => this.onScroll(), { passive: true });

        // Initial update
        this.updateActiveLink();
    }

    onScroll() {
        // Use requestAnimationFrame for smooth performance
        if (!this.ticking) {
            window.requestAnimationFrame(() => {
                this.updateActiveLink();
                this.ticking = false;
            });
            this.ticking = true;
        }
    }

    updateActiveLink() {
        // Find the current section
        let currentSectionId = null;
        const scrollPosition = window.pageYOffset || document.documentElement.scrollTop;

        // Check each section to find which one is currently in view
        for (let i = this.sections.length - 1; i >= 0; i--) {
            const section = this.sections[i];
            const sectionTop = section.offsetTop - this.offset;

            if (scrollPosition >= sectionTop) {
                currentSectionId = section.getAttribute('id');
                break;
            }
        }

        // If we haven't scrolled past any section, activate the first one
        if (!currentSectionId && this.sections.length > 0) {
            currentSectionId = this.sections[0].getAttribute('id');
        }

        // Update active state on TOC links
        this.tocLinks.forEach(link => {
            const href = link.getAttribute('href');
            const linkId = href ? href.substring(1) : null; // Remove '#' from href

            if (linkId === currentSectionId) {
                if (this.currentActiveLink !== link) {
                    // Remove active class from previous link
                    if (this.currentActiveLink) {
                        this.currentActiveLink.classList.remove('active');
                    }

                    // Add active class to current link
                    link.classList.add('active');
                    this.currentActiveLink = link;

                    // Scroll TOC to show active item (if needed)
                    this.scrollTocToActiveItem(link);
                }
            }
        });
    }

    scrollTocToActiveItem(link) {
        const toc = link.closest('.sidebar-right');
        if (!toc) return;

        const tocRect = toc.getBoundingClientRect();
        const linkRect = link.getBoundingClientRect();

        // Check if link is out of view in the TOC
        if (linkRect.top < tocRect.top || linkRect.bottom > tocRect.bottom) {
            // Calculate proper scroll position to center the link
            // Use getBoundingClientRect for accurate positioning
            const linkOffsetInToc = linkRect.top - tocRect.top + toc.scrollTop;
            const targetScroll = linkOffsetInToc - (tocRect.height / 2) + (linkRect.height / 2);

            toc.scrollTo({
                top: Math.max(0, targetScroll),
                behavior: 'smooth'
            });
        }
    }
}

// Initialize scrollspy
const scrollspy = new ScrollSpy();
