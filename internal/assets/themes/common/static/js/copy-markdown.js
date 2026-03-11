(function () {
  'use strict';

  var btn = document.getElementById('copy-md-btn');
  if (!btn) return;

  var basePath = document.documentElement.dataset.basePath || '';
  var slug = btn.dataset.slug;
  if (!slug) return;

  var mdPath = basePath + '/' + slug + '.html.md';
  var originalText = btn.textContent;

  btn.addEventListener('click', function () {
    fetch(mdPath)
      .then(function (res) {
        if (!res.ok) throw new Error('Not found');
        return res.text();
      })
      .then(function (text) {
        return navigator.clipboard.writeText(text);
      })
      .then(function () {
        btn.textContent = 'Copied';
        setTimeout(function () {
          btn.textContent = originalText;
        }, 2000);
      })
      .catch(function () {
        btn.textContent = 'Failed';
        setTimeout(function () {
          btn.textContent = originalText;
        }, 2000);
      });
  });
})();
