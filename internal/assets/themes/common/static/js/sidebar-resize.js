// Sidebar resize functionality
(function() {
  'use strict';

  const STORAGE_KEY_LEFT = 'openapi-sidebar-left-width';
  const STORAGE_KEY_RIGHT = 'openapi-sidebar-right-width';

  // Read constraints from CSS custom properties, falling back to sensible defaults.
  var styles = getComputedStyle(document.documentElement);
  const MIN_WIDTH_LEFT = parseInt(styles.getPropertyValue('--sidebar-min-left'), 10) || 200;
  const MAX_WIDTH_LEFT = parseInt(styles.getPropertyValue('--sidebar-max-left'), 10) || 800;
  const MIN_WIDTH_RIGHT = parseInt(styles.getPropertyValue('--sidebar-min-right'), 10) || 300;
  const MAX_WIDTH_RIGHT = parseInt(styles.getPropertyValue('--sidebar-max-right'), 10) || 800;
  const DEFAULT_WIDTH_LEFT = parseInt(styles.getPropertyValue('--openapi-sidebar-width'), 10) || 350;
  const DEFAULT_WIDTH_RIGHT = parseInt(styles.getPropertyValue('--openapi-tester-width'), 10) || 450;

  function init() {
    initLeftSidebar();
    initRightSidebar();
  }

  // Initialize left sidebar resize
  function initLeftSidebar() {
    const handle = document.getElementById('sidebar-resize-handle');
    const sidebar = document.getElementById('openapi-sidebar-left');

    if (!handle || !sidebar) return;

    // Load saved width from localStorage
    const savedWidth = localStorage.getItem(STORAGE_KEY_LEFT);
    if (savedWidth) {
      setLeftSidebarWidth(parseInt(savedWidth, 10));
    }

    let isResizing = false;
    let startX = 0;
    let startWidth = 0;

    handle.addEventListener('mousedown', (e) => {
      isResizing = true;
      startX = e.clientX;
      startWidth = sidebar.offsetWidth;
      handle.classList.add('resizing');
      document.body.style.cursor = 'col-resize';
      document.body.style.userSelect = 'none';
      e.preventDefault();
    });

    document.addEventListener('mousemove', (e) => {
      if (!isResizing) return;

      const diff = e.clientX - startX;
      const newWidth = Math.min(Math.max(startWidth + diff, MIN_WIDTH_LEFT), MAX_WIDTH_LEFT);
      setLeftSidebarWidth(newWidth);
    });

    document.addEventListener('mouseup', () => {
      if (!isResizing) return;

      isResizing = false;
      handle.classList.remove('resizing');
      document.body.style.cursor = '';
      document.body.style.userSelect = '';

      // Save to localStorage
      const currentWidth = sidebar.offsetWidth;
      localStorage.setItem(STORAGE_KEY_LEFT, currentWidth.toString());
    });

    // Handle double-click to reset
    handle.addEventListener('dblclick', () => {
      setLeftSidebarWidth(DEFAULT_WIDTH_LEFT);
      localStorage.setItem(STORAGE_KEY_LEFT, DEFAULT_WIDTH_LEFT.toString());
    });
  }

  // Initialize right sidebar resize
  function initRightSidebar() {
    const handle = document.getElementById('sidebar-resize-handle-right');
    const sidebar = document.getElementById('openapi-sidebar-right');

    if (!handle || !sidebar) return;

    // Load saved width from localStorage
    const savedWidth = localStorage.getItem(STORAGE_KEY_RIGHT);
    if (savedWidth) {
      setRightSidebarWidth(parseInt(savedWidth, 10));
    }

    let isResizing = false;
    let startX = 0;
    let startWidth = 0;

    handle.addEventListener('mousedown', (e) => {
      isResizing = true;
      startX = e.clientX;
      startWidth = sidebar.offsetWidth;
      handle.classList.add('resizing');
      document.body.style.cursor = 'col-resize';
      document.body.style.userSelect = 'none';
      e.preventDefault();
    });

    document.addEventListener('mousemove', (e) => {
      if (!isResizing) return;

      // For right sidebar, dragging left (negative diff) makes it wider
      const diff = startX - e.clientX;
      const newWidth = Math.min(Math.max(startWidth + diff, MIN_WIDTH_RIGHT), MAX_WIDTH_RIGHT);
      setRightSidebarWidth(newWidth);
    });

    document.addEventListener('mouseup', () => {
      if (!isResizing) return;

      isResizing = false;
      handle.classList.remove('resizing');
      document.body.style.cursor = '';
      document.body.style.userSelect = '';

      // Save to localStorage
      const currentWidth = sidebar.offsetWidth;
      localStorage.setItem(STORAGE_KEY_RIGHT, currentWidth.toString());
    });

    // Handle double-click to reset
    handle.addEventListener('dblclick', () => {
      setRightSidebarWidth(DEFAULT_WIDTH_RIGHT);
      localStorage.setItem(STORAGE_KEY_RIGHT, DEFAULT_WIDTH_RIGHT.toString());
    });
  }

  function setLeftSidebarWidth(width) {
    document.documentElement.style.setProperty('--openapi-sidebar-width', `${width}px`);
  }

  function setRightSidebarWidth(width) {
    document.documentElement.style.setProperty('--openapi-tester-width', `${width}px`);
  }

  // Initialize on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
