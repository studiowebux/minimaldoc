// Knowledge Base Search
(function() {
    'use strict';

    var searchInput = document.getElementById('kb-search-input');
    var searchResults = document.getElementById('kb-search-results');

    if (!searchInput || !searchResults) return;

    var searchData = null;
    var debounceTimer = null;

    // Get base path from HTML element
    var basePath = document.documentElement.dataset.basePath || '';

    // Load search index
    function loadSearchIndex() {
        if (searchData) return Promise.resolve(searchData);

        var kbPath = window.location.pathname.split('/').filter(Boolean)[0] || 'kb';
        var indexUrl = basePath + '/' + kbPath + '/kb-search.json';

        return fetch(indexUrl)
            .then(function(response) {
                if (!response.ok) throw new Error('Failed to load search index');
                return response.json();
            })
            .then(function(data) {
                searchData = data;
                return data;
            })
            .catch(function(err) {
                console.error('KB search index load error:', err);
                return [];
            });
    }

    // Search function
    function search(query) {
        if (!searchData || !query.trim()) {
            hideResults();
            return;
        }

        var terms = query.toLowerCase().split(/\s+/).filter(Boolean);
        if (terms.length === 0) {
            hideResults();
            return;
        }

        var results = searchData.filter(function(item) {
            var text = (item.title + ' ' + item.description + ' ' + item.category + ' ' + (item.tags || []).join(' ') + ' ' + item.content).toLowerCase();
            return terms.every(function(term) {
                return text.indexOf(term) !== -1;
            });
        });

        // Score and sort results
        results = results.map(function(item) {
            var score = 0;
            var titleLower = item.title.toLowerCase();
            var descLower = (item.description || '').toLowerCase();

            terms.forEach(function(term) {
                if (titleLower.indexOf(term) !== -1) score += 10;
                if (descLower.indexOf(term) !== -1) score += 5;
                if ((item.tags || []).some(function(tag) { return tag.toLowerCase().indexOf(term) !== -1; })) score += 3;
            });

            return { item: item, score: score };
        }).sort(function(a, b) {
            return b.score - a.score;
        }).slice(0, 8);

        showResults(results.map(function(r) { return r.item; }));
    }

    // Show results
    function showResults(results) {
        if (results.length === 0) {
            searchResults.innerHTML = '<div class="kb-search-no-results">No results found</div>';
            searchResults.classList.add('active');
            return;
        }

        var html = results.map(function(item) {
            return '<a href="' + escapeHtml(item.url) + '" class="kb-search-result">' +
                '<div class="kb-search-result-category">' + escapeHtml(item.category) + '</div>' +
                '<div class="kb-search-result-title">' + escapeHtml(item.title) + '</div>' +
                (item.description ? '<div class="kb-search-result-description">' + escapeHtml(item.description) + '</div>' : '') +
                '</a>';
        }).join('');

        searchResults.innerHTML = html;
        searchResults.classList.add('active');
    }

    // Hide results
    function hideResults() {
        searchResults.classList.remove('active');
        searchResults.innerHTML = '';
    }

    // Escape HTML
    function escapeHtml(text) {
        if (!text) return '';
        var div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // Input handler with debounce
    searchInput.addEventListener('input', function() {
        var query = searchInput.value;

        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(function() {
            loadSearchIndex().then(function() {
                search(query);
            });
        }, 150);
    });

    // Focus handler - load index early
    searchInput.addEventListener('focus', function() {
        loadSearchIndex();
        if (searchInput.value.trim()) {
            search(searchInput.value);
        }
    });

    // Close on click outside
    document.addEventListener('click', function(e) {
        if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
            hideResults();
        }
    });

    // Keyboard navigation
    searchInput.addEventListener('keydown', function(e) {
        if (!searchResults.classList.contains('active')) return;

        var items = searchResults.querySelectorAll('.kb-search-result');
        var current = searchResults.querySelector('.kb-search-result:focus');
        var index = current ? Array.prototype.indexOf.call(items, current) : -1;

        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                if (index < items.length - 1) {
                    items[index + 1].focus();
                } else if (index === -1 && items.length > 0) {
                    items[0].focus();
                }
                break;
            case 'ArrowUp':
                e.preventDefault();
                if (index > 0) {
                    items[index - 1].focus();
                } else if (index === 0) {
                    searchInput.focus();
                }
                break;
            case 'Escape':
                hideResults();
                searchInput.blur();
                break;
        }
    });

    // Allow keyboard nav from results back to input
    searchResults.addEventListener('keydown', function(e) {
        if (e.key === 'ArrowUp') {
            var items = searchResults.querySelectorAll('.kb-search-result');
            var current = document.activeElement;
            if (items[0] === current) {
                e.preventDefault();
                searchInput.focus();
            }
        }
    });
})();
