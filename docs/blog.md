---
title: Blog
description: Latest news and updates
hidden: true
full_width: true
no_header: true
no_widgets: true
---

<style>
.blog-hero {
    padding: 4rem 2rem 3rem;
    text-align: center;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-primary);
}
.blog-hero-content {
    max-width: 800px;
    margin: 0 auto;
}
.blog-hero h1 {
    font-size: 2.5rem;
    font-weight: 700;
    margin: 0 0 1rem;
    color: var(--text-primary);
}
.blog-hero-description {
    font-size: 1.25rem;
    color: var(--text-secondary);
    margin: 0 0 2rem;
}
.blog-search-container {
    position: relative;
    max-width: 500px;
    margin: 0 auto;
}
.blog-search-input {
    width: 100%;
    padding: 1rem 1rem 1rem 3rem;
    font-size: 1rem;
    border: 2px solid var(--border-primary);
    border-radius: 8px;
    background: var(--bg-primary);
    color: var(--text-primary);
    transition: border-color 0.2s, box-shadow 0.2s;
}
.blog-search-input:focus {
    outline: none;
    border-color: var(--link-color);
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}
.blog-search-input::placeholder {
    color: var(--text-muted);
}
.blog-search-icon {
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    width: 1.25rem;
    height: 1.25rem;
    color: var(--text-muted);
    pointer-events: none;
}
.blog-content {
    padding: 3rem 2rem;
    max-width: 1400px;
    margin: 0 auto;
}
.blog-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 2rem;
}
.blog-card {
    background: var(--bg-primary);
    border: 1px solid var(--border-primary);
    border-radius: 12px;
    overflow: hidden;
    transition: border-color 0.2s, box-shadow 0.2s, transform 0.2s;
    cursor: pointer;
}
.blog-card:hover {
    border-color: var(--link-color);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
    transform: translateY(-4px);
}
.blog-card-image {
    width: 100%;
    height: 180px;
    object-fit: cover;
}
.blog-card-body {
    padding: 1.25rem;
}
.blog-card-category {
    display: inline-block;
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--link-color);
    margin-bottom: 0.5rem;
}
.blog-card-title {
    font-size: 1.125rem;
    font-weight: 700;
    margin: 0 0 0.5rem;
    line-height: 1.3;
    color: var(--text-primary);
}
.blog-card-excerpt {
    font-size: 0.875rem;
    color: var(--text-secondary);
    line-height: 1.5;
    margin: 0 0 1rem;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.blog-card-meta {
    display: flex;
    align-items: center;
    gap: 1rem;
    font-size: 0.8rem;
    color: var(--text-muted);
}
.blog-empty, .blog-loading {
    grid-column: 1 / -1;
    text-align: center;
    padding: 4rem 2rem;
    color: var(--text-muted);
}
/* Article View */
.blog-article {
    max-width: 800px;
    margin: 0 auto;
    padding: 3rem 2rem;
}
.blog-article-back {
    display: block;
    color: var(--text-secondary);
    text-decoration: none;
    margin-bottom: 1.5rem;
    font-size: 0.9375rem;
}
.blog-article-back:hover {
    color: var(--link-color);
}
.blog-article-header {
    margin-bottom: 2rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--border-primary);
}
.blog-article-category {
    display: inline-block;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--link-color);
    margin-bottom: 0.75rem;
}
.blog-article-title {
    font-size: 2.25rem;
    font-weight: 700;
    line-height: 1.2;
    margin: 0 0 1rem;
    color: var(--text-primary);
}
.blog-article-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    font-size: 0.9375rem;
    color: var(--text-secondary);
}
.blog-article-featured {
    width: 100%;
    max-height: 400px;
    object-fit: cover;
    border-radius: 12px;
    margin-bottom: 2rem;
}
.blog-article-content {
    font-size: 1.125rem;
    line-height: 1.8;
    color: var(--text-primary);
}
.blog-article-content h2 { font-size: 1.75rem; margin: 2rem 0 1rem; }
.blog-article-content h3 { font-size: 1.5rem; margin: 1.75rem 0 0.75rem; }
.blog-article-content p { margin: 0 0 1.25rem; }
.blog-article-content pre {
    background: var(--bg-tertiary);
    padding: 1rem;
    border-radius: 8px;
    overflow-x: auto;
    margin: 1.5rem 0;
}
.blog-article-content code { font-family: 'Monaco', 'Consolas', monospace; font-size: 0.9em; }
.blog-article-content img { max-width: 100%; border-radius: 8px; margin: 1.5rem 0; }
.blog-article-content a { color: var(--link-color); }
.blog-article-content blockquote {
    border-left: 4px solid var(--link-color);
    margin: 1.5rem 0;
    padding: 0.5rem 0 0.5rem 1.5rem;
    color: var(--text-secondary);
    font-style: italic;
}
@media (max-width: 768px) {
    .blog-hero { padding: 3rem 1rem 2rem; }
    .blog-hero h1 { font-size: 1.75rem; }
    .blog-content { padding: 2rem 1rem; }
    .blog-grid { grid-template-columns: 1fr; gap: 1.5rem; }
    .blog-article { padding: 2rem 1rem; }
    .blog-article-title { font-size: 1.75rem; }
    .blog-article-content { font-size: 1rem; }
}
</style>

<div id="blog-list-view">
<section class="blog-hero">
<div class="blog-hero-content">
<h1>Blog</h1>
<p class="blog-hero-description">Latest news and updates</p>
<div class="blog-search-container">
<input type="text" id="blog-search-input" class="blog-search-input" placeholder="Search posts..." autocomplete="off">
<svg class="blog-search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
<circle cx="11" cy="11" r="8"/>
<path d="m21 21-4.35-4.35"/>
</svg>
</div>
</div>
</section>

<section class="blog-content">
<div class="blog-grid" id="blog-posts"></div>
</section>
</div>

<div id="blog-article-view" style="display:none;">
<article class="blog-article" id="blog-article"></article>
</div>

<script>
(function() {
    var listView = document.getElementById('blog-list-view');
    var articleView = document.getElementById('blog-article-view');
    var articleContainer = document.getElementById('blog-article');
    var postsContainer = document.getElementById('blog-posts');
    var searchInput = document.getElementById('blog-search-input');

    function getSlug() {
        var params = new URLSearchParams(window.location.search);
        return params.get('article');
    }

    function showList() {
        listView.style.display = '';
        articleView.style.display = 'none';
        document.title = 'Blog';
    }

    function showArticle() {
        listView.style.display = 'none';
        articleView.style.display = '';
    }

    function esc(t) {
        if (!t) return '';
        var d = document.createElement('div');
        d.textContent = t;
        return d.innerHTML;
    }

    function formatDate(dateStr) {
        var d = new Date(dateStr);
        return d.toLocaleDateString(undefined, {year:'numeric', month:'long', day:'numeric'});
    }

    function formatDateShort(dateStr) {
        var d = new Date(dateStr);
        return d.toLocaleDateString(undefined, {year:'numeric', month:'short', day:'numeric'});
    }

    function renderCards(posts) {
        if (!posts.length) {
            postsContainer.innerHTML = '<div class="blog-empty">No posts yet</div>';
            return;
        }
        var html = '';
        posts.forEach(function(post) {
            html += '<article class="blog-card" onclick="window.location.href=\'?article=' + encodeURIComponent(post.slug) + '\'">';
            if (post.featured_image) {
                html += '<img src="' + esc(post.featured_image) + '" alt="" class="blog-card-image" loading="lazy">';
            }
            html += '<div class="blog-card-body">';
            if (post.category) {
                html += '<span class="blog-card-category">' + esc(post.category) + '</span>';
            }
            html += '<h2 class="blog-card-title">' + esc(post.title) + '</h2>';
            if (post.description) {
                html += '<p class="blog-card-excerpt">' + esc(post.description) + '</p>';
            }
            html += '<div class="blog-card-meta">';
            if (post.published_at) {
                html += '<span>' + formatDateShort(post.published_at) + '</span>';
            }
            if (post.reading_time) {
                html += '<span>' + post.reading_time + ' min read</span>';
            }
            html += '</div></div></article>';
        });
        postsContainer.innerHTML = html;
    }

    function renderArticle(post) {
        var html = '<header class="blog-article-header">';
        html += '<a href="?" class="blog-article-back">&larr; Back to Blog</a>';
        if (post.category) {
            html += '<span class="blog-article-category">' + esc(post.category) + '</span>';
        }
        html += '<h1 class="blog-article-title">' + esc(post.title) + '</h1>';
        html += '<div class="blog-article-meta">';
        if (post.published_at) {
            html += '<time datetime="' + post.published_at + '">' + formatDate(post.published_at) + '</time>';
        }
        if (post.reading_time) {
            html += '<span>' + post.reading_time + ' min read</span>';
        }
        html += '</div></header>';
        if (post.featured_image) {
            html += '<img src="' + esc(post.featured_image) + '" alt="" class="blog-article-featured">';
        }
        html += '<div class="blog-article-content">';
        html += post.content_html || '<p>' + esc(post.description || '') + '</p>';
        html += '</div>';
        articleContainer.innerHTML = html;
        document.title = post.title + ' | Blog';
    }

    function loadPosts(query) {
        postsContainer.innerHTML = '<div class="blog-loading">Loading...</div>';
        if (window.MinimalDoc && window.MinimalDoc.blog) {
            var opts = { limit: 12 };
            if (query) opts.search = query;
            window.MinimalDoc.blog.list(opts).then(function(r) {
                renderCards(r.posts || []);
            }).catch(function() {
                postsContainer.innerHTML = '<div class="blog-empty">Failed to load posts</div>';
            });
        } else {
            setTimeout(function() { loadPosts(query); }, 100);
        }
    }

    function loadArticle(slug) {
        articleContainer.innerHTML = '<div class="blog-loading" style="padding:4rem 2rem;text-align:center;">Loading...</div>';
        if (window.MinimalDoc && window.MinimalDoc.blog) {
            window.MinimalDoc.blog.get(slug).then(function(r) {
                if (r.post) {
                    renderArticle(r.post);
                } else {
                    articleContainer.innerHTML = '<div class="blog-empty" style="padding:4rem 2rem;text-align:center;">Article not found</div>';
                }
            }).catch(function() {
                articleContainer.innerHTML = '<div class="blog-empty" style="padding:4rem 2rem;text-align:center;">Failed to load article</div>';
            });
        } else {
            setTimeout(function() { loadArticle(slug); }, 100);
        }
    }

    function init() {
        var slug = getSlug();
        if (slug) {
            showArticle();
            loadArticle(slug);
        } else {
            showList();
            loadPosts();
        }
    }

    if (searchInput) {
        var timeout;
        searchInput.addEventListener('input', function() {
            clearTimeout(timeout);
            timeout = setTimeout(function() {
                loadPosts(searchInput.value.trim());
            }, 300);
        });
    }

    window.renderBlogCards = renderCards;
    init();
})();
</script>
