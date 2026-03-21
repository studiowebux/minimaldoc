// Locale Selector
// Handles locale/language switching with localStorage persistence

(function() {
  'use strict';

  var selector = document.getElementById('locale-selector');
  if (!selector) return;

  var dropdown = selector.querySelector('.locale-dropdown');
  var current = selector.querySelector('.locale-current');
  var items = selector.querySelectorAll('.locale-item');
  var isOpen = false;

  // Get version prefixes from version selector if it exists
  var versionPrefixes = [];
  var versionSelector = document.getElementById('version-selector');
  if (versionSelector) {
    versionSelector.querySelectorAll('.version-item').forEach(function(v) {
      var path = v.getAttribute('data-path');
      if (path) versionPrefixes.push(path);
    });
  }

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

  // Extract version prefix from path
  function extractVersionPrefix(path) {
    for (var i = 0; i < versionPrefixes.length; i++) {
      var prefix = '/' + versionPrefixes[i] + '/';
      if (path.startsWith(prefix) || path === '/' + versionPrefixes[i]) {
        return versionPrefixes[i];
      }
    }
    return null;
  }

  // Handle locale selection
  function selectLocale(e) {
    var item = e.target.closest('.locale-item');
    if (!item) return;

    var localeCode = item.getAttribute('data-code');

    // Store preference
    localStorage.setItem('preferredLocale', localeCode);

    // Navigate to the new locale
    var currentPath = window.location.pathname;
    var basePath = document.documentElement.getAttribute('data-base-path') || '';

    // Remove base path from current path
    var relativePath = currentPath;
    if (basePath && relativePath.startsWith(basePath)) {
      relativePath = relativePath.substring(basePath.length);
    }

    // Get locale configuration
    var currentLang = document.documentElement.lang || 'en';
    var hideDefault = document.documentElement.getAttribute('data-hide-default-locale') === 'true';

    // Collect all locale codes
    var localeCodes = [];
    items.forEach(function(i) {
      var code = i.getAttribute('data-code');
      if (code) localeCodes.push(code);
    });

    // Extract current version prefix if any
    var currentVersion = null;

    // First remove any locale prefix
    localeCodes.forEach(function(code) {
      if (relativePath.startsWith('/' + code + '/')) {
        relativePath = relativePath.substring(code.length + 1);
      } else if (relativePath === '/' + code) {
        relativePath = '/';
      }
    });

    // Then check for version prefix
    currentVersion = extractVersionPrefix(relativePath);
    if (currentVersion) {
      // Remove version prefix from path
      var versionPrefix = '/' + currentVersion;
      if (relativePath.startsWith(versionPrefix + '/')) {
        relativePath = relativePath.substring(versionPrefix.length);
      } else if (relativePath === versionPrefix) {
        relativePath = '/';
      }
    }

    // Build new URL
    var newPath = basePath;

    // Get default locale from data attribute
    var defaultLocale = document.documentElement.getAttribute('data-default-locale') || 'en';

    // Add locale prefix only for non-default locales (or if hideDefault is false)
    if (localeCode !== defaultLocale || !hideDefault) {
      newPath += '/' + localeCode;
    }

    // Add version prefix if we were on a versioned page
    if (currentVersion) {
      newPath += '/' + currentVersion;
    }

    // Add the page path
    if (relativePath && !relativePath.startsWith('/')) {
      relativePath = '/' + relativePath;
    }
    newPath += relativePath;

    // Clean up double slashes
    newPath = newPath.replace(/\/+/g, '/');

    // Ensure at least root path
    if (!newPath) {
      newPath = '/';
    }

    // Navigate
    window.location.href = newPath;
  }

  // Event listeners
  current.addEventListener('click', toggleDropdown);

  items.forEach(function(item) {
    item.addEventListener('click', selectLocale);
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
      var first = dropdown.querySelector('.locale-item');
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
    } else if (e.key === 'Enter' && focused.classList.contains('locale-item')) {
      selectLocale(e);
    }
  });
})();
