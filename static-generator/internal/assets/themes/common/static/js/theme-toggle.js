/**
 * Theme Toggle
 * Handles light/dark theme switching with localStorage persistence
 */

class ThemeToggle {
    constructor() {
        // Get theme from localStorage or default to light
        this.theme = localStorage.getItem('theme') || 'light';
        this.toggleButton = null;
        this.init();
    }

    init() {
        // Set initial theme
        document.documentElement.setAttribute('data-theme', this.theme);

        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.bindEvents());
        } else {
            this.bindEvents();
        }
    }

    bindEvents() {
        this.toggleButton = document.getElementById('theme-toggle');

        if (this.toggleButton) {
            this.toggleButton.addEventListener('click', () => this.toggle());
        }
    }

    toggle() {
        // Switch theme
        this.theme = this.theme === 'light' ? 'dark' : 'light';

        // Update DOM
        document.documentElement.setAttribute('data-theme', this.theme);

        // Persist to localStorage
        localStorage.setItem('theme', this.theme);

        // Dispatch custom event for other components
        window.dispatchEvent(new CustomEvent('themechange', {
            detail: { theme: this.theme }
        }));
    }

    getTheme() {
        return this.theme;
    }
}

// Initialize theme toggle
const themeToggle = new ThemeToggle();
