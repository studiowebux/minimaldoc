// Shared utility functions for admin UI

/**
 * escapeHtml escapes a string for safe HTML insertion using DOM textContent.
 * @param {string} s - The string to escape
 * @returns {string} HTML-escaped string
 */
function escapeHtml(s) {
    if (!s) return '';
    const div = document.createElement('div');
    div.textContent = s;
    return div.innerHTML;
}
