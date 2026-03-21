(function () {
  'use strict';

  var btn = document.getElementById('copy-md-btn');
  if (!btn) return;

  var basePath = document.documentElement.dataset.basePath || '';
  var slug = btn.dataset.slug;
  if (!slug) return;

  var mdPath = basePath + '/' + slug + '.html.md';
  var originalText = btn.textContent;

  function fallbackCopy(text) {
    var ta = document.createElement('textarea');
    ta.value = text;
    ta.className = 'copy-fallback-textarea';
    document.body.appendChild(ta);
    ta.select();
    try {
      document.execCommand('copy');
    } finally {
      document.body.removeChild(ta);
    }
  }

  btn.addEventListener('click', function () {
    var textPromise = fetch(mdPath).then(function (res) {
      if (!res.ok) throw new Error('Not found');
      return res.text();
    });

    var done;
    if (navigator.clipboard && typeof ClipboardItem !== 'undefined') {
      done = navigator.clipboard.write([
        new ClipboardItem({
          'text/plain': textPromise.then(function (text) {
            return new Blob([text], { type: 'text/plain' });
          })
        })
      ]);
    } else {
      done = textPromise.then(function (text) {
        fallbackCopy(text);
      });
    }

    done
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
