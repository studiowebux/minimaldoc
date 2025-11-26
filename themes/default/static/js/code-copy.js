// Code Copy Button
// Adds copy buttons to all code blocks in documentation
(function() {
  'use strict';

  // Initialize on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  function init() {
    addCopyButtonsToCodeBlocks();
  }

  function addCopyButtonsToCodeBlocks() {
    // Find all pre elements that contain code
    const preElements = document.querySelectorAll('pre:has(code), pre');

    preElements.forEach(pre => {
      // Skip if already has a copy button
      if (pre.parentElement.classList.contains('code-block-wrapper')) {
        return;
      }

      // Create wrapper
      const wrapper = document.createElement('div');
      wrapper.className = 'code-block-wrapper';

      // Create copy button
      const button = document.createElement('button');
      button.className = 'code-copy-btn';
      button.setAttribute('aria-label', 'Copy code');
      button.textContent = 'Copy';

      // Add click handler
      button.addEventListener('click', async () => {
        const code = pre.textContent;

        try {
          await navigator.clipboard.writeText(code);

          // Visual feedback
          const originalText = button.textContent;
          button.textContent = 'Copied!';
          button.classList.add('copied');

          setTimeout(() => {
            button.textContent = originalText;
            button.classList.remove('copied');
          }, 2000);
        } catch (err) {
          console.error('Failed to copy code:', err);
          button.textContent = 'Failed';
          setTimeout(() => {
            button.textContent = 'Copy';
          }, 2000);
        }
      });

      // Wrap the pre element
      pre.parentNode.insertBefore(wrapper, pre);
      wrapper.appendChild(button);
      wrapper.appendChild(pre);
    });
  }
})();
