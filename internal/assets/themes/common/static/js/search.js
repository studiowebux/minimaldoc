/**
 * Search Functionality - Inverted Index
 * Client-side search with Cmd+K / Ctrl+K shortcut
 */

(function() {
    let searchData = null;
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
        results.innerHTML = '<div class="search-no-results">Loading...</div>';

        return fetch(basePath + '/search-index.json')
            .then(response => response.json())
            .then(data => {
                searchData = data;
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

    // Track the element that triggered the modal for focus restoration
    var triggerElement = null;

    // Open search modal
    function openSearch() {
        triggerElement = document.activeElement;
        modal.classList.add('open');
        modal.setAttribute('aria-hidden', 'false');
        input.value = '';
        input.focus();
        results.innerHTML = '';
        selectedIndex = -1;

        // Trap focus within modal
        document.body.style.overflow = 'hidden';
        loadSearchIndex();
    }

    // Close search modal
    function closeSearch() {
        modal.classList.remove('open');
        modal.setAttribute('aria-hidden', 'true');
        input.value = '';
        results.innerHTML = '';
        selectedIndex = -1;
        document.body.style.overflow = '';

        // Restore focus to trigger element
        if (triggerElement && triggerElement.focus) {
            triggerElement.focus();
        }
    }

    // Tokenize query (same as Go side)
    function tokenize(text) {
        return text.toLowerCase()
            .split(/[^a-z0-9]+/)
            .filter(w => w.length >= 2);
    }

    // Search using inverted index
    // Posting list format: [pageID, score, pageID, score, ...]
    function search(query) {
        if (!query || query.trim() === '' || !searchData) {
            results.innerHTML = '';
            return;
        }

        const words = tokenize(query);
        if (words.length === 0) {
            results.innerHTML = '';
            return;
        }

        // Score accumulator: pageID -> score
        const scores = new Map();

        // Process posting list (compact array format)
        function addPostings(list, multiplier) {
            for (let i = 0; i < list.length; i += 2) {
                const pageID = list[i];
                const score = list[i + 1] * multiplier;
                scores.set(pageID, (scores.get(pageID) || 0) + score);
            }
        }

        // Look up each word in inverted index
        words.forEach(word => {
            // Exact match
            if (searchData.idx[word]) {
                addPostings(searchData.idx[word], 1);
            }

            // Prefix match for partial words (autocomplete)
            if (word.length >= 3) {
                Object.keys(searchData.idx).forEach(indexWord => {
                    if (indexWord !== word && indexWord.startsWith(word)) {
                        addPostings(searchData.idx[indexWord], 0.5);
                    }
                });
            }
        });

        // Convert to array and sort by score
        const matches = Array.from(scores.entries())
            .map(([pageID, score]) => ({
                page: searchData.pages[pageID],
                pageID,
                score
            }))
            .filter(m => m.page)
            .sort((a, b) => b.score - a.score)
            .slice(0, 10);

        renderResults(matches);
    }

    // Find matching sections for a page based on search terms
    function findMatchingSections(pageID, words) {
        if (!searchData.sections) return [];

        const matchedSections = [];
        const pageSections = searchData.sections.filter(s => s.p === pageID);

        pageSections.forEach(section => {
            const sectionWords = tokenize(section.t);
            const matchScore = words.reduce((score, word) => {
                return score + (sectionWords.some(sw => sw.startsWith(word)) ? 1 : 0);
            }, 0);

            if (matchScore > 0) {
                matchedSections.push({
                    title: section.t,
                    anchor: section.a,
                    score: matchScore
                });
            }
        });

        return matchedSections.sort((a, b) => b.score - a.score).slice(0, 3);
    }

    // Render search results
    function renderResults(matches) {
        if (matches.length === 0) {
            results.innerHTML = '<div class="search-no-results" role="status">No results found</div>';
            input.setAttribute('aria-expanded', 'false');
            input.removeAttribute('aria-activedescendant');
            return;
        }

        results.innerHTML = '';
        input.setAttribute('aria-expanded', 'true');
        const query = input.value;
        const words = tokenize(query);
        let itemIndex = 0;

        matches.forEach((match) => {
            const page = match.page;
            const matchingSections = findMatchingSections(match.pageID, words);

            // Main page result
            const item = document.createElement('div');
            item.className = 'search-result-item';
            item.setAttribute('role', 'option');
            item.setAttribute('aria-selected', 'false');
            item.setAttribute('tabindex', '-1');
            item.id = 'search-result-' + itemIndex;
            item.dataset.index = itemIndex++;
            item.dataset.url = page.u;

            const title = document.createElement('div');
            title.className = 'search-result-title';
            title.textContent = page.t;

            item.appendChild(title);

            if (page.d) {
                const description = document.createElement('div');
                description.className = 'search-result-description';
                description.textContent = page.d;
                item.appendChild(description);
            }

            // Add matching sections as sub-results
            if (matchingSections.length > 0) {
                const sectionsContainer = document.createElement('div');
                sectionsContainer.className = 'search-result-sections';

                matchingSections.forEach(section => {
                    const sectionLink = document.createElement('a');
                    sectionLink.className = 'search-result-section';
                    sectionLink.href = page.u + '#' + section.anchor;
                    sectionLink.textContent = '# ' + section.title;
                    sectionLink.addEventListener('click', function(e) {
                        e.stopPropagation();
                        closeSearch();
                        window.location.href = this.href;
                    });
                    sectionsContainer.appendChild(sectionLink);
                });

                item.appendChild(sectionsContainer);
            }

            item.addEventListener('click', function() {
                closeSearch();
                window.location.href = page.u;
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
        items.forEach(function(item) {
            item.classList.remove('selected');
            item.setAttribute('aria-selected', 'false');
        });

        // Add selection to current
        selectedIndex = index;
        items[index].classList.add('selected');
        items[index].setAttribute('aria-selected', 'true');

        // Update aria-activedescendant on input
        input.setAttribute('aria-activedescendant', items[index].id);

        // Scroll into view
        items[index].scrollIntoView({ block: 'nearest', behavior: 'smooth' });
    }

    // Navigate to selected result
    function navigateToSelected() {
        const items = results.querySelectorAll('.search-result-item');
        if (selectedIndex >= 0 && selectedIndex < items.length) {
            const url = items[selectedIndex].dataset.url;
            if (url) {
                closeSearch();
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
