/**
 * Search Functionality - Sharded Inverted Index
 * Client-side search with Cmd+K / Ctrl+K shortcut
 * Loads manifest once, then fetches shards on-demand
 */

(function() {
    var searchData = null;      // Manifest (pages, sections, shard list)
    var loadedShards = {};      // Cache: prefix -> shard data
    var selectedIndex = -1;
    var manifestLoaded = false;
    var manifestLoading = false;
    var pendingSearch = null;   // Store pending search while loading

    var modal = document.getElementById('search-modal');
    var input = document.getElementById('search-input');
    var results = document.getElementById('search-results');
    var closeBtn = document.getElementById('search-close');
    var searchButton = document.getElementById('search-button');

    if (!modal || !input || !results || !closeBtn) {
        console.warn('Search elements not found');
        return;
    }

    var basePath = document.documentElement.getAttribute('data-base-path') || '';

    // Get shard prefix for a term (first 2 chars)
    function getShardPrefix(term) {
        if (term.length < 2) return term;
        return term.substring(0, 2);
    }

    // Load search manifest (pages + sections + shard list)
    function loadManifest() {
        if (manifestLoaded || manifestLoading) {
            return manifestLoaded ? Promise.resolve() : new Promise(function(resolve) {
                var check = setInterval(function() {
                    if (manifestLoaded) {
                        clearInterval(check);
                        resolve();
                    }
                }, 50);
            });
        }

        manifestLoading = true;
        results.innerHTML = '<div class="search-no-results">Loading...</div>';

        return fetch(basePath + '/search-manifest.json')
            .then(function(response) { return response.json(); })
            .then(function(data) {
                searchData = data;
                searchData.idx = {}; // Initialize empty index, will be populated from shards
                manifestLoaded = true;
                manifestLoading = false;
                results.innerHTML = '';

                // Execute pending search if any
                if (pendingSearch) {
                    var q = pendingSearch;
                    pendingSearch = null;
                    search(q);
                }
            })
            .catch(function(error) {
                console.error('Failed to load search manifest:', error);
                manifestLoading = false;
                results.innerHTML = '<div class="search-no-results">Failed to load search index</div>';
            });
    }

    // Load a specific shard
    function loadShard(prefix) {
        if (loadedShards[prefix]) {
            return Promise.resolve(loadedShards[prefix]);
        }

        // Check if this shard exists
        if (searchData && searchData.shards && searchData.shards.indexOf(prefix) === -1) {
            loadedShards[prefix] = { idx: {} };
            return Promise.resolve(loadedShards[prefix]);
        }

        return fetch(basePath + '/search-shards/' + prefix + '.json')
            .then(function(response) {
                if (!response.ok) {
                    loadedShards[prefix] = { idx: {} };
                    return loadedShards[prefix];
                }
                return response.json();
            })
            .then(function(data) {
                loadedShards[prefix] = data;
                // Merge into searchData.idx for compatibility
                if (data.idx) {
                    Object.keys(data.idx).forEach(function(term) {
                        searchData.idx[term] = data.idx[term];
                    });
                }
                return data;
            })
            .catch(function(error) {
                console.error('Failed to load shard ' + prefix + ':', error);
                loadedShards[prefix] = { idx: {} };
                return loadedShards[prefix];
            });
    }

    // Load multiple shards in parallel
    function loadShards(prefixes) {
        var uniquePrefixes = [];
        prefixes.forEach(function(p) {
            if (uniquePrefixes.indexOf(p) === -1) {
                uniquePrefixes.push(p);
            }
        });
        return Promise.all(uniquePrefixes.map(loadShard));
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
        loadManifest();
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
            .filter(function(w) { return w.length >= 2; });
    }

    // Search using sharded inverted index
    function search(query) {
        if (!query || query.trim() === '') {
            results.innerHTML = '';
            return;
        }

        // Wait for manifest if not loaded
        if (!manifestLoaded) {
            pendingSearch = query;
            if (!manifestLoading) {
                loadManifest();
            }
            return;
        }

        var words = tokenize(query);
        if (words.length === 0) {
            results.innerHTML = '';
            return;
        }

        // Determine which shards we need
        var neededPrefixes = [];
        words.forEach(function(word) {
            var prefix = getShardPrefix(word);
            if (neededPrefixes.indexOf(prefix) === -1) {
                neededPrefixes.push(prefix);
            }

            // For prefix matching, we might need adjacent shards
            // This is a simplification - we only load exact prefix shards
        });

        // Show loading if shards not yet loaded
        var allLoaded = neededPrefixes.every(function(p) { return loadedShards[p]; });
        if (!allLoaded) {
            results.innerHTML = '<div class="search-no-results">Searching...</div>';
        }

        // Load needed shards then execute search
        loadShards(neededPrefixes).then(function() {
            executeSearch(words);
        });
    }

    // Execute search once shards are loaded
    function executeSearch(words) {
        var scores = new Map();

        // Process posting list (compact array format)
        function addPostings(list, multiplier) {
            for (var i = 0; i < list.length; i += 2) {
                var pageID = list[i];
                var score = list[i + 1] * multiplier;
                scores.set(pageID, (scores.get(pageID) || 0) + score);
            }
        }

        // Look up each word in loaded shards
        words.forEach(function(word) {
            var prefix = getShardPrefix(word);
            var shard = loadedShards[prefix];

            if (shard && shard.idx) {
                // Exact match
                if (shard.idx[word]) {
                    addPostings(shard.idx[word], 1);
                }

                // Prefix match for partial words (autocomplete)
                if (word.length >= 3) {
                    Object.keys(shard.idx).forEach(function(indexWord) {
                        if (indexWord !== word && indexWord.indexOf(word) === 0) {
                            addPostings(shard.idx[indexWord], 0.5);
                        }
                    });
                }
            }
        });

        // Convert to array and sort by score
        var matches = Array.from(scores.entries())
            .map(function(entry) {
                return {
                    page: searchData.pages[entry[0]],
                    pageID: entry[0],
                    score: entry[1]
                };
            })
            .filter(function(m) { return m.page; })
            .sort(function(a, b) { return b.score - a.score; })
            .slice(0, 10);

        renderResults(matches);
    }

    // Find matching sections for a page based on search terms
    function findMatchingSections(pageID, words) {
        if (!searchData.sections) return [];

        var matchedSections = [];
        var pageSections = searchData.sections.filter(function(s) { return s.p === pageID; });

        pageSections.forEach(function(section) {
            var sectionWords = tokenize(section.t);
            var matchScore = words.reduce(function(score, word) {
                return score + (sectionWords.some(function(sw) { return sw.indexOf(word) === 0; }) ? 1 : 0);
            }, 0);

            if (matchScore > 0) {
                matchedSections.push({
                    title: section.t,
                    anchor: section.a,
                    score: matchScore
                });
            }
        });

        return matchedSections.sort(function(a, b) { return b.score - a.score; }).slice(0, 3);
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
        var query = input.value;
        var words = tokenize(query);
        var itemIndex = 0;

        matches.forEach(function(match) {
            var page = match.page;
            var matchingSections = findMatchingSections(match.pageID, words);

            // Main page result
            var item = document.createElement('div');
            item.className = 'search-result-item';
            item.setAttribute('role', 'option');
            item.setAttribute('aria-selected', 'false');
            item.setAttribute('tabindex', '-1');
            item.id = 'search-result-' + itemIndex;
            item.dataset.index = itemIndex++;
            item.dataset.url = page.u;

            var title = document.createElement('div');
            title.className = 'search-result-title';
            title.textContent = page.t;

            item.appendChild(title);

            if (page.d) {
                var description = document.createElement('div');
                description.className = 'search-result-description';
                description.textContent = page.d;
                item.appendChild(description);
            }

            // Add matching sections as sub-results
            if (matchingSections.length > 0) {
                var sectionsContainer = document.createElement('div');
                sectionsContainer.className = 'search-result-sections';

                matchingSections.forEach(function(section) {
                    var sectionLink = document.createElement('a');
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
        var items = results.querySelectorAll('.search-result-item');
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
        var items = results.querySelectorAll('.search-result-item');
        if (selectedIndex >= 0 && selectedIndex < items.length) {
            var url = items[selectedIndex].dataset.url;
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
            var items = results.querySelectorAll('.search-result-item');
            if (items.length > 0) {
                selectResult((selectedIndex + 1) % items.length);
            }
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            var downItems = results.querySelectorAll('.search-result-item');
            if (downItems.length > 0) {
                selectResult(selectedIndex <= 0 ? downItems.length - 1 : selectedIndex - 1);
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
