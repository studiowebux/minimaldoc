/**
 * Search Functionality
 * Client-side search with Cmd+K / Ctrl+K shortcut and search button
 */

(function() {
    let searchIndex = [];
    let selectedIndex = -1;
    let indexLoaded = false;
    let indexLoading = false;

    const modal = document.getElementById('search-modal');
    const input = document.getElementById('search-input');
    const results = document.getElementById('search-results');
    const closeBtn = document.getElementById('search-close');
    const searchButton = document.getElementById('search-button');

    if (!modal || !input || !results || !closeBtn) {
        console.warn('Search elements not found');
        return;
    }

    const basePath = document.documentElement.getAttribute('data-base-path') || '';

    // Load search index lazily
    function loadSearchIndex() {
        if (indexLoaded || indexLoading) return Promise.resolve();

        indexLoading = true;
        results.innerHTML = '<div class="search-no-results">Loading search index...</div>';

        return fetch(basePath + '/search-index.json')
            .then(response => response.json())
            .then(data => {
                searchIndex = data.entries || [];
                indexLoaded = true;
                indexLoading = false;
                results.innerHTML = '';
            })
            .catch(error => {
                console.error('Failed to load search index:', error);
                indexLoading = false;
                results.innerHTML = '<div class="search-no-results">Failed to load search index</div>';
            });
    }

    // Open search modal
    function openSearch() {
        modal.classList.add('open');
        input.value = '';
        input.focus();
        results.innerHTML = '';
        selectedIndex = -1;

        // Load index when modal opens
        loadSearchIndex();
    }

    // Close search modal
    function closeSearch() {
        modal.classList.remove('open');
        input.value = '';
        results.innerHTML = '';
        selectedIndex = -1;
    }

    // Calculate similarity between two strings (Levenshtein distance based)
    function similarity(s1, s2) {
        const longer = s1.length > s2.length ? s1 : s2;
        const shorter = s1.length > s2.length ? s2 : s1;

        if (longer.length === 0) return 1.0;

        const editDistance = levenshteinDistance(longer, shorter);
        return (longer.length - editDistance) / longer.length;
    }

    // Levenshtein distance algorithm
    function levenshteinDistance(str1, str2) {
        const matrix = [];

        for (let i = 0; i <= str2.length; i++) {
            matrix[i] = [i];
        }

        for (let j = 0; j <= str1.length; j++) {
            matrix[0][j] = j;
        }

        for (let i = 1; i <= str2.length; i++) {
            for (let j = 1; j <= str1.length; j++) {
                if (str2.charAt(i - 1) === str1.charAt(j - 1)) {
                    matrix[i][j] = matrix[i - 1][j - 1];
                } else {
                    matrix[i][j] = Math.min(
                        matrix[i - 1][j - 1] + 1,
                        matrix[i][j - 1] + 1,
                        matrix[i - 1][j] + 1
                    );
                }
            }
        }

        return matrix[str2.length][str1.length];
    }

    // Perform search with multi-word support and fuzzy matching
    function search(query) {
        if (!query || query.trim() === '') {
            results.innerHTML = '';
            return;
        }

        // Split query into words and convert to lowercase
        const queryWords = query.toLowerCase().trim().split(/\s+/);
        const matches = [];

        searchIndex.forEach(entry => {
            let totalScore = 0;
            const titleLower = entry.title.toLowerCase();
            const descriptionLower = (entry.description || '').toLowerCase();
            const contentLower = (entry.content || '').toLowerCase();
            const tagsLower = (entry.tags || []).join(' ').toLowerCase();

            // Combine all searchable text
            const allText = `${titleLower} ${descriptionLower} ${contentLower} ${tagsLower}`;

            // Score each query word
            queryWords.forEach(word => {
                let wordScore = 0;

                // Exact word match in title (highest score)
                if (titleLower === word) {
                    wordScore += 100;
                } else if (titleLower.includes(word)) {
                    wordScore += 50;
                } else {
                    // Fuzzy match in title
                    const titleWords = titleLower.split(/\s+/);
                    titleWords.forEach(titleWord => {
                        const sim = similarity(word, titleWord);
                        if (sim > 0.7) {
                            wordScore += sim * 40;
                        }
                    });
                }

                // Exact word match in description
                if (descriptionLower.includes(word)) {
                    wordScore += 30;
                } else {
                    // Fuzzy match in description
                    const descWords = descriptionLower.split(/\s+/);
                    descWords.forEach(descWord => {
                        const sim = similarity(word, descWord);
                        if (sim > 0.7) {
                            wordScore += sim * 20;
                        }
                    });
                }

                // Exact word match in content
                if (contentLower.includes(word)) {
                    wordScore += 15;
                } else {
                    // Fuzzy match in content
                    const contentWords = contentLower.split(/\s+/).slice(0, 50); // Limit for performance
                    contentWords.forEach(contentWord => {
                        const sim = similarity(word, contentWord);
                        if (sim > 0.8) {
                            wordScore += sim * 10;
                        }
                    });
                }

                // Exact word match in tags
                if (tagsLower.includes(word)) {
                    wordScore += 35;
                }

                totalScore += wordScore;
            });

            // Bonus for matching all query words
            const allWordsFound = queryWords.every(word =>
                allText.includes(word) ||
                allText.split(/\s+/).some(textWord => similarity(word, textWord) > 0.7)
            );

            if (allWordsFound && queryWords.length > 1) {
                totalScore *= 1.5; // 50% bonus for matching all words
            }

            if (totalScore > 0) {
                matches.push({ entry, score: totalScore });
            }
        });

        // Sort by score descending
        matches.sort((a, b) => b.score - a.score);

        // Render results
        renderResults(matches.slice(0, 10)); // Show top 10 results
    }

    // Render search results
    function renderResults(matches) {
        if (matches.length === 0) {
            results.innerHTML = '<div class="search-no-results">No results found</div>';
            return;
        }

        results.innerHTML = '';

        matches.forEach((match, index) => {
            const item = document.createElement('div');
            item.className = 'search-result-item';
            item.dataset.index = index;
            item.dataset.url = match.entry.url;

            const title = document.createElement('div');
            title.className = 'search-result-title';
            title.textContent = match.entry.title;

            const description = document.createElement('div');
            description.className = 'search-result-description';
            description.textContent = match.entry.description || '';

            const content = document.createElement('div');
            content.className = 'search-result-content';
            content.textContent = match.entry.content || '';

            item.appendChild(title);
            if (match.entry.description) {
                item.appendChild(description);
            }
            if (match.entry.content) {
                item.appendChild(content);
            }

            item.addEventListener('click', function() {
                window.location.href = match.entry.url;
            });

            results.appendChild(item);
        });

        selectedIndex = -1;
    }

    // Select result item
    function selectResult(index) {
        const items = results.querySelectorAll('.search-result-item');
        if (index < 0 || index >= items.length) return;

        // Remove previous selection
        items.forEach(item => item.classList.remove('selected'));

        // Add selection to current
        selectedIndex = index;
        items[index].classList.add('selected');

        // Scroll into view
        items[index].scrollIntoView({ block: 'nearest', behavior: 'smooth' });
    }

    // Navigate to selected result
    function navigateToSelected() {
        const items = results.querySelectorAll('.search-result-item');
        if (selectedIndex >= 0 && selectedIndex < items.length) {
            const url = items[selectedIndex].dataset.url;
            if (url) {
                window.location.href = url;
            }
        }
    }

    // Keyboard shortcuts
    document.addEventListener('keydown', function(e) {
        // Cmd+K / Ctrl+K to open search
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
            e.preventDefault();
            if (!modal.classList.contains('open')) {
                openSearch();
            }
            return;
        }

        // Escape to close
        if (e.key === 'Escape' && modal.classList.contains('open')) {
            closeSearch();
            return;
        }

        // Arrow keys for navigation (only when modal is open)
        if (!modal.classList.contains('open')) return;

        if (e.key === 'ArrowDown') {
            e.preventDefault();
            const items = results.querySelectorAll('.search-result-item');
            if (items.length > 0) {
                selectResult((selectedIndex + 1) % items.length);
            }
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            const items = results.querySelectorAll('.search-result-item');
            if (items.length > 0) {
                selectResult(selectedIndex <= 0 ? items.length - 1 : selectedIndex - 1);
            }
        } else if (e.key === 'Enter') {
            e.preventDefault();
            navigateToSelected();
        }
    });

    // Input event for search
    input.addEventListener('input', function() {
        search(this.value);
    });

    // Search button click
    if (searchButton) {
        searchButton.addEventListener('click', openSearch);
    }

    // Close button
    closeBtn.addEventListener('click', closeSearch);

    // Click outside modal to close
    modal.addEventListener('click', function(e) {
        if (e.target === modal) {
            closeSearch();
        }
    });
})();
