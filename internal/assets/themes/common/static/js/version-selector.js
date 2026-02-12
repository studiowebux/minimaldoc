// Version Selector
// Handles version switching with localStorage persistence

(function() {
  'use strict';

  var selector = document.getElementById('version-selector');
  if (!selector) return;

  var dropdown = selector.querySelector('.version-dropdown');
  var current = selector.querySelector('.version-current');
  var items = selector.querySelectorAll('.version-item');
  var isOpen = false;

  // Get stored version preference
  var storedVersion = localStorage.getItem('preferredVersion');

  // Toggle dropdown
  function toggleDropdown(e) {
    e.preventDefault();
    e.stopPropagation();
    isOpen = !isOpen;
    dropdown.classList.toggle('open', isOpen);
    current.setAttribute('aria-expanded', isOpen);
  }

  // Close dropdown
  function closeDropdown() {
    isOpen = false;
    dropdown.classList.remove('open');
    current.setAttribute('aria-expanded', 'false');
  }

  // Handle version selection
  function selectVersion(e) {
    var item = e.target.closest('.version-item');
    if (!item) return;

    var versionPath = item.getAttribute('data-path');
    var versionName = item.getAttribute('data-name');

    // Store preference
    localStorage.setItem('preferredVersion', versionName);

    // Navigate to the new version
    var currentPath = window.location.pathname;
    var basePath = document.documentElement.getAttribute('data-base-path') || '';

    // Remove base path and current version prefix from path
    var relativePath = currentPath;
    if (basePath && relativePath.startsWith(basePath)) {
      relativePath = relativePath.substring(basePath.length);
    }

    // Remove any existing version prefix
    var versions = [];
    items.forEach(function(i) {
      var p = i.getAttribute('data-path');
      if (p) versions.push(p);
    });

    versions.forEach(function(v) {
      if (relativePath.startsWith('/' + v + '/')) {
        relativePath = relativePath.substring(v.length + 1);
      }
    });

    // Build new URL
    var newPath = basePath;
    if (versionPath) {
      newPath += '/' + versionPath;
    }
    newPath += relativePath;

    // Navigate
    window.location.href = newPath;
  }

  // Event listeners
  current.addEventListener('click', toggleDropdown);

  items.forEach(function(item) {
    item.addEventListener('click', selectVersion);
  });

  // Close on outside click
  document.addEventListener('click', function(e) {
    if (!selector.contains(e.target)) {
      closeDropdown();
    }
  });

  // Close on escape
  document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
      closeDropdown();
    }
  });

  // Keyboard navigation
  current.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' || e.key === ' ') {
      toggleDropdown(e);
    } else if (e.key === 'ArrowDown' && isOpen) {
      e.preventDefault();
      var first = dropdown.querySelector('.version-item');
      if (first) first.focus();
    }
  });

  dropdown.addEventListener('keydown', function(e) {
    var focused = document.activeElement;
    var allItems = Array.from(items);
    var index = allItems.indexOf(focused);

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      var next = allItems[index + 1] || allItems[0];
      next.focus();
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      var prev = allItems[index - 1] || allItems[allItems.length - 1];
      prev.focus();
    } else if (e.key === 'Enter' && focused.classList.contains('version-item')) {
      selectVersion(e);
    }
  });
})();
