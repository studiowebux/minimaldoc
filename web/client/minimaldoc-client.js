/**
 * MinimalDoc Client - Lightweight tracking and interaction library
 * Cookie-free, privacy-first analytics and widgets
 * < 3KB minified
 */
(function() {
  'use strict';

  var MinimalDoc = {
    config: {
      endpoint: '',
      siteId: '',
      features: [],
      debug: false
    },

    // Time tracking state
    pageStartTime: 0,
    currentPath: '',
    sessionHash: '',
    pageCount: 0, // Track pages viewed in session for bounce rate

    init: function(options) {
      if (!options.endpoint || !options.siteId) {
        console.error('[MinimalDoc] endpoint and siteId are required');
        return;
      }

      this.config.endpoint = options.endpoint.replace(/\/$/, '');
      this.config.siteId = options.siteId;
      this.config.features = options.features || ['analytics'];
      this.config.debug = options.debug || false;

      // Generate session hash for this browser session
      this.sessionHash = this.generateSessionHash();

      if (this.hasFeature('analytics')) {
        this.trackPageView();
        this.setupSPATracking();
        this.setupTimeTracking();
      }

      if (this.hasFeature('feedback')) {
        this.initFeedbackWidget();
      }

      if (this.hasFeature('newsletter')) {
        this.initNewsletterForms();
      }

      if (this.hasFeature('private-docs')) {
        this.checkPrivateAccess();
      }

      if (this.hasFeature('blog')) {
        this.initBlog();
      }

      if (this.hasFeature('forum')) {
        this.initForum();
      }

      this.log('Initialized with features:', this.config.features);
    },

    generateSessionHash: function() {
      // Generate a simple anonymous session hash
      // This is NOT tracking users across sessions, just grouping page views in a single session
      var stored = sessionStorage.getItem('minimaldoc_session');
      if (stored) return stored;

      var hash = Math.random().toString(36).substring(2, 15);
      sessionStorage.setItem('minimaldoc_session', hash);
      return hash;
    },

    getPageCount: function() {
      return parseInt(sessionStorage.getItem('minimaldoc_page_count') || '0', 10);
    },

    incrementPageCount: function() {
      var count = this.getPageCount() + 1;
      sessionStorage.setItem('minimaldoc_page_count', count.toString());
      this.pageCount = count;
      return count;
    },

    hasFeature: function(name) {
      return this.config.features.indexOf(name) !== -1;
    },

    log: function() {
      if (this.config.debug) {
        console.log.apply(console, ['[MinimalDoc]'].concat(Array.prototype.slice.call(arguments)));
      }
    },

    // Analytics
    trackPageView: function(path) {
      var trackPath = path || window.location.pathname;

      // Increment page count for bounce tracking
      this.incrementPageCount();

      var data = {
        site_id: this.config.siteId,
        path: trackPath,
        referrer: document.referrer || '',
        screen_width: window.innerWidth,
        session_hash: this.sessionHash
      };

      // Record page start time for duration tracking
      this.pageStartTime = Date.now();
      this.currentPath = trackPath;

      this.post('/api/analytics/track', data);
      this.log('Tracked:', data.path, '(page', this.pageCount, 'in session)');
    },

    setupSPATracking: function() {
      var self = this;
      var lastPath = window.location.pathname;

      // Track pushState/replaceState
      var originalPush = history.pushState;
      var originalReplace = history.replaceState;

      history.pushState = function() {
        self.sendDuration(); // Send duration before navigating
        originalPush.apply(this, arguments);
        if (window.location.pathname !== lastPath) {
          lastPath = window.location.pathname;
          self.trackPageView();
        }
      };

      history.replaceState = function() {
        self.sendDuration(); // Send duration before navigating
        originalReplace.apply(this, arguments);
        if (window.location.pathname !== lastPath) {
          lastPath = window.location.pathname;
          self.trackPageView();
        }
      };

      // Track popstate (back/forward)
      window.addEventListener('popstate', function() {
        self.sendDuration();
        if (window.location.pathname !== lastPath) {
          lastPath = window.location.pathname;
          self.trackPageView();
        }
      });
    },

    setupTimeTracking: function() {
      var self = this;

      // Record when page becomes visible/hidden
      document.addEventListener('visibilitychange', function() {
        if (document.hidden) {
          self.sendDuration();
        } else {
          // Reset start time when page becomes visible again
          self.pageStartTime = Date.now();
        }
      });

      // Send duration on page unload (mark as final for bounce detection)
      window.addEventListener('beforeunload', function() {
        self.sendDuration(true);
      });

      // Also send duration periodically for long sessions (every 30 seconds)
      setInterval(function() {
        if (!document.hidden && self.pageStartTime > 0) {
          self.sendDuration();
        }
      }, 30000);
    },

    sendDuration: function(isFinal) {
      if (!this.pageStartTime || !this.currentPath) return;

      var duration = Math.round((Date.now() - this.pageStartTime) / 1000);
      if (duration < 1) return; // Ignore sub-second visits

      var data = {
        site_id: this.config.siteId,
        path: this.currentPath,
        duration: duration,
        session_hash: this.sessionHash
      };

      // Only mark as bounce on final send (page unload) and if only 1 page viewed
      if (isFinal && this.getPageCount() === 1) {
        data.is_bounce = true;
        this.log('Marking as bounce (single page session)');
      }

      // Use sendBeacon for reliability on page unload
      if (navigator.sendBeacon) {
        navigator.sendBeacon(
          this.config.endpoint + '/api/analytics/duration',
          JSON.stringify(data)
        );
      } else {
        this.post('/api/analytics/duration', data);
      }

      this.log('Duration sent:', duration, 'seconds for', this.currentPath);
      this.pageStartTime = Date.now(); // Reset for next measurement
    },

    // Feedback Widget
    initFeedbackWidget: function() {
      var self = this;
      var widgets = document.querySelectorAll('[data-minimaldoc-feedback]');

      this.log('Found', widgets.length, 'feedback widget(s)');

      widgets.forEach(function(widget) {
        self.log('Rendering feedback widget for:', widget);
        self.renderFeedbackWidget(widget);
      });
    },

    renderFeedbackWidget: function(container) {
      var self = this;
      var path = container.getAttribute('data-path') || window.location.pathname;

      container.innerHTML = [
        '<div class="minimaldoc-feedback">',
        '  <span class="minimaldoc-feedback-label">Was this helpful?</span>',
        '  <div class="minimaldoc-feedback-stars" data-rating="0">',
        '    <button data-star="1" aria-label="1 star">&#9734;</button>',
        '    <button data-star="2" aria-label="2 stars">&#9734;</button>',
        '    <button data-star="3" aria-label="3 stars">&#9734;</button>',
        '    <button data-star="4" aria-label="4 stars">&#9734;</button>',
        '    <button data-star="5" aria-label="5 stars">&#9734;</button>',
        '  </div>',
        '  <div class="minimaldoc-feedback-form" style="display:none;">',
        '    <textarea placeholder="Any feedback? (optional)" rows="2"></textarea>',
        '    <button type="submit">Submit</button>',
        '  </div>',
        '  <div class="minimaldoc-feedback-thanks" style="display:none;">',
        '    Thanks for your feedback!',
        '  </div>',
        '</div>'
      ].join('\n');

      var stars = container.querySelector('.minimaldoc-feedback-stars');
      var form = container.querySelector('.minimaldoc-feedback-form');
      var thanks = container.querySelector('.minimaldoc-feedback-thanks');
      var textarea = container.querySelector('textarea');
      var submitBtn = container.querySelector('button[type="submit"]');
      var rating = 0;

      // Star click handler
      stars.addEventListener('click', function(e) {
        var btn = e.target.closest('[data-star]');
        if (!btn) return;

        rating = parseInt(btn.getAttribute('data-star'), 10);
        self.updateStars(stars, rating);
        form.style.display = 'block';
      });

      // Submit handler
      submitBtn.addEventListener('click', function() {
        if (rating === 0) return;

        self.submitFeedback(path, rating, textarea.value.trim());
        stars.style.display = 'none';
        form.style.display = 'none';
        thanks.style.display = 'block';
      });
    },

    updateStars: function(container, rating) {
      var buttons = container.querySelectorAll('[data-star]');
      buttons.forEach(function(btn) {
        var star = parseInt(btn.getAttribute('data-star'), 10);
        btn.innerHTML = star <= rating ? '&#9733;' : '&#9734;';
        btn.classList.toggle('active', star <= rating);
      });
    },

    submitFeedback: function(path, rating, feedback) {
      var data = {
        site_id: this.config.siteId,
        path: path,
        rating: rating
      };
      if (feedback) {
        data.feedback = feedback;
      }

      this.post('/api/feedback', data);
      this.log('Feedback submitted:', rating, 'stars');
    },

    // Newsletter
    initNewsletterForms: function() {
      var self = this;
      var forms = document.querySelectorAll('[data-minimaldoc-newsletter]');

      forms.forEach(function(form) {
        form.addEventListener('submit', function(e) {
          e.preventDefault();
          var emailInput = form.querySelector('input[type="email"]');
          if (!emailInput || !emailInput.value) return;

          self.subscribe(emailInput.value.trim(), function(success) {
            if (success) {
              var msg = form.querySelector('.minimaldoc-newsletter-msg') ||
                        document.createElement('div');
              msg.className = 'minimaldoc-newsletter-msg';
              msg.textContent = 'Check your email to confirm subscription!';
              if (!msg.parentNode) form.appendChild(msg);
              emailInput.value = '';
            }
          });
        });
      });
    },

    subscribe: function(email, callback) {
      var data = {
        site_id: this.config.siteId,
        email: email
      };

      this.post('/api/newsletter/subscribe', data, callback);
      this.log('Newsletter subscription:', email);
    },

    // Private Docs Access
    checkPrivateAccess: function() {
      var self = this;
      var path = window.location.pathname;

      this.log('Checking access for:', path);

      this.get('/api/docs/check?path=' + encodeURIComponent(path), function(response) {
        if (!response) {
          self.log('Failed to check access');
          return;
        }

        self.log('Access check result:', response);

        if (response.is_protected && !response.has_access) {
          self.handleAccessDenied(response);
        }
      });
    },

    handleAccessDenied: function(response) {
      var self = this;
      var reason = response.reason;
      var requiredRole = response.required_role;

      this.log('Access denied:', reason, 'required role:', requiredRole);

      // Find content element to overlay
      var content = document.querySelector('main') ||
                    document.querySelector('article') ||
                    document.querySelector('.content') ||
                    document.body;

      // Create overlay
      var overlay = document.createElement('div');
      overlay.className = 'minimaldoc-access-overlay';
      overlay.innerHTML = [
        '<div class="minimaldoc-access-modal">',
        '  <h2>Access Required</h2>',
        reason === 'authentication_required'
          ? '  <p>This content requires authentication. Please log in to continue.</p>'
          : '  <p>Your account does not have permission to view this content. Required role: <strong>' + requiredRole + '</strong></p>',
        '  <div class="minimaldoc-access-actions">',
        reason === 'authentication_required'
          ? '    <a href="' + self.config.endpoint + '/login?site_id=' + encodeURIComponent(self.config.siteId) + '&redirect=' + encodeURIComponent(window.location.href) + '" class="btn btn-primary">Log In</a>'
          : '    <a href="/" class="btn">Go to Home</a>',
        '  </div>',
        '</div>'
      ].join('\n');

      // Apply styles
      overlay.style.cssText = 'position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.85);z-index:9999;display:flex;align-items:center;justify-content:center;';

      var modal = overlay.querySelector('.minimaldoc-access-modal');
      modal.style.cssText = 'background:var(--bg-color,#fff);padding:2rem;border-radius:8px;max-width:400px;text-align:center;color:var(--text-color,#333);';

      // Hide original content
      content.style.filter = 'blur(10px)';
      content.style.pointerEvents = 'none';

      document.body.appendChild(overlay);
    },

    // Blog API
    blog: {
      parent: null,

      list: function(options) {
        var self = this;
        options = options || {};
        var params = [];

        if (options.limit) params.push('limit=' + options.limit);
        if (options.offset) params.push('offset=' + options.offset);
        if (options.category) params.push('category=' + encodeURIComponent(options.category));
        if (options.tag) params.push('tag=' + encodeURIComponent(options.tag));
        if (options.search) params.push('q=' + encodeURIComponent(options.search));

        var query = params.length ? '?' + params.join('&') : '';

        return new Promise(function(resolve, reject) {
          self.parent.get('/api/blog/posts' + query, function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch posts'));
            }
          });
        });
      },

      get: function(slug) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/blog/posts/' + encodeURIComponent(slug), function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Post not found'));
            }
          });
        });
      },

      related: function(slug, limit) {
        var self = this;
        var query = limit ? '?limit=' + limit : '';
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/blog/posts/' + encodeURIComponent(slug) + '/related' + query, function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch related posts'));
            }
          });
        });
      }
    },

    initBlog: function() {
      var self = this;
      this.blog.parent = this;

      // Auto-render blog containers
      var containers = document.querySelectorAll('[data-minimaldoc-blog]');
      this.log('Found', containers.length, 'blog container(s)');

      containers.forEach(function(container) {
        self.renderBlogContainer(container);
      });
    },

    renderBlogContainer: function(container) {
      var self = this;
      var type = container.getAttribute('data-minimaldoc-blog') || 'list';
      var limit = parseInt(container.getAttribute('data-limit'), 10) || 10;
      var category = container.getAttribute('data-category');
      var tag = container.getAttribute('data-tag');
      var slug = container.getAttribute('data-slug');
      var template = container.getAttribute('data-template') || 'default';
      var compact = container.hasAttribute('data-compact');

      container.innerHTML = '<div class="minimaldoc-blog-loading">Loading...</div>';

      if (type === 'single' && slug) {
        this.blog.get(slug).then(function(post) {
          container.innerHTML = self.renderPost(post, template, true, false);
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-blog-error">Post not found</div>';
        });
      } else if (type === 'related' && slug) {
        this.blog.related(slug, limit).then(function(response) {
          container.innerHTML = self.renderPostList(response.posts || [], template, false, compact);
        }).catch(function() {
          container.innerHTML = '';
        });
      } else if (type === 'categories') {
        this.blog.list({ limit: 100 }).then(function(response) {
          container.innerHTML = self.renderCategories(response.posts || []);
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-blog-error">Failed to load categories</div>';
        });
      } else {
        this.blog.list({ limit: limit, category: category, tag: tag }).then(function(response) {
          container.innerHTML = self.renderPostList(response.posts || [], template, false, compact);
          if (response.total > limit && !compact) {
            container.innerHTML += self.renderPagination(response.total, limit, 0);
          }
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-blog-error">Failed to load posts</div>';
        });
      }
    },

    renderCategories: function(posts) {
      var categories = {};
      posts.forEach(function(post) {
        if (post.category) {
          categories[post.category] = (categories[post.category] || 0) + 1;
        }
      });

      var categoryNames = Object.keys(categories).sort();
      if (!categoryNames.length) {
        return '<div class="minimaldoc-blog-empty">No categories yet</div>';
      }

      var html = '<div class="minimaldoc-blog-categories">';
      categoryNames.forEach(function(name) {
        html += '<a href="/blog?category=' + encodeURIComponent(name) + '" class="minimaldoc-blog-category-item">';
        html += '<span class="minimaldoc-blog-category-name">' + this.escapeHtml(name) + '</span>';
        html += '<span class="minimaldoc-blog-category-count">' + categories[name] + '</span>';
        html += '</a>';
      }, this);
      html += '</div>';
      return html;
    },

    renderPostList: function(posts, template, full, compact) {
      var self = this;
      if (!posts.length) {
        return '<div class="minimaldoc-blog-empty">No posts found</div>';
      }
      var listClass = compact ? 'minimaldoc-blog-list minimaldoc-blog-list--compact' : 'minimaldoc-blog-list';
      return '<div class="' + listClass + '">' +
        posts.map(function(post) { return self.renderPost(post, template, full, compact); }).join('') +
        '</div>';
    },

    renderPost: function(post, template, full, compact) {
      var postClass = 'minimaldoc-blog-post';
      if (full) postClass += ' minimaldoc-blog-post--full';

      var html = '<article class="' + postClass + '">';

      if (post.featured_image && !compact) {
        html += '<img src="' + this.escapeHtml(post.featured_image) + '" alt="" class="minimaldoc-blog-image" loading="lazy">';
      }

      html += '<div class="minimaldoc-blog-content">';
      html += '<h2 class="minimaldoc-blog-title">';
      html += full ? this.escapeHtml(post.title) : '<a href="/blog/' + this.escapeHtml(post.slug) + '">' + this.escapeHtml(post.title) + '</a>';
      html += '</h2>';

      html += '<div class="minimaldoc-blog-meta">';
      if (post.published_at) {
        html += '<time datetime="' + post.published_at + '">' + this.formatDate(post.published_at) + '</time>';
      }
      if (post.reading_time && !compact) {
        html += '<span class="minimaldoc-blog-reading-time">' + post.reading_time + ' min read</span>';
      }
      if (post.category && !compact) {
        html += '<span class="minimaldoc-blog-category">' + this.escapeHtml(post.category) + '</span>';
      }
      html += '</div>';

      if (full && post.content_html) {
        html += '<div class="minimaldoc-blog-body">' + post.content_html + '</div>';
      } else if (post.description && !compact) {
        html += '<p class="minimaldoc-blog-excerpt">' + this.escapeHtml(post.description) + '</p>';
      }

      if (post.tags && post.tags.length && !compact) {
        html += '<div class="minimaldoc-blog-tags">';
        post.tags.forEach(function(tag) {
          html += '<span class="minimaldoc-blog-tag">' + this.escapeHtml(tag) + '</span>';
        }, this);
        html += '</div>';
      }

      if (!full && !compact) {
        html += '<a href="/blog/' + this.escapeHtml(post.slug) + '" class="minimaldoc-blog-read-more">Read more</a>';
      }

      html += '</div></article>';
      return html;
    },

    renderPagination: function(total, limit, offset) {
      var pages = Math.ceil(total / limit);
      var current = Math.floor(offset / limit);
      var html = '<nav class="minimaldoc-blog-pagination">';
      for (var i = 0; i < pages && i < 10; i++) {
        html += '<button data-page="' + i + '"' + (i === current ? ' class="active"' : '') + '>' + (i + 1) + '</button>';
      }
      html += '</nav>';
      return html;
    },

    formatDate: function(dateStr) {
      var date = new Date(dateStr);
      return date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
    },

    escapeHtml: function(text) {
      if (!text) return '';
      var div = document.createElement('div');
      div.textContent = text;
      return div.innerHTML;
    },

    // Forum API
    forum: {
      parent: null,

      categories: function() {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/categories', function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch categories'));
            }
          });
        });
      },

      topics: function(options) {
        var self = this;
        options = options || {};
        var params = [];

        if (options.limit) params.push('limit=' + options.limit);
        if (options.offset) params.push('offset=' + options.offset);
        if (options.category) params.push('category_id=' + encodeURIComponent(options.category));
        if (options.search) params.push('q=' + encodeURIComponent(options.search));
        if (options.status) params.push('status=' + options.status);

        var query = params.length ? '?' + params.join('&') : '';

        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/topics' + query, function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch topics'));
            }
          });
        });
      },

      getTopic: function(slug) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/topics/by-slug/' + encodeURIComponent(slug), function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Topic not found'));
            }
          });
        });
      },

      getPosts: function(topicSlug, options) {
        var self = this;
        options = options || {};
        var params = [];

        if (options.limit) params.push('limit=' + options.limit);
        if (options.offset) params.push('offset=' + options.offset);

        var query = params.length ? '?' + params.join('&') : '';

        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/topics/by-slug/' + encodeURIComponent(topicSlug) + '/posts' + query, function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch posts'));
            }
          });
        });
      },

      createTopic: function(data) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.post('/api/forum/topics', data, function(success) {
            if (success) {
              resolve();
            } else {
              reject(new Error('Failed to create topic'));
            }
          });
        });
      },

      createPost: function(topicSlug, data) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.post('/api/forum/topics/by-slug/' + encodeURIComponent(topicSlug) + '/posts', data, function(success) {
            if (success) {
              resolve();
            } else {
              reject(new Error('Failed to create post'));
            }
          });
        });
      },

      likeTopic: function(id) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.post('/api/forum/topics/' + id + '/like', {}, function(success) {
            if (success) {
              resolve();
            } else {
              reject(new Error('Failed to like topic'));
            }
          });
        });
      },

      likePost: function(id) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.post('/api/forum/posts/' + id + '/like', {}, function(success) {
            if (success) {
              resolve();
            } else {
              reject(new Error('Failed to like post'));
            }
          });
        });
      },

      bookmark: function(topicId) {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.post('/api/forum/topics/' + topicId + '/bookmark', {}, function(success) {
            if (success) {
              resolve();
            } else {
              reject(new Error('Failed to bookmark'));
            }
          });
        });
      },

      search: function(query, limit) {
        var self = this;
        var params = ['q=' + encodeURIComponent(query)];
        if (limit) params.push('limit=' + limit);

        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/search?' + params.join('&'), function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Search failed'));
            }
          });
        });
      },

      notifications: function(options) {
        var self = this;
        options = options || {};
        var params = [];

        if (options.unread) params.push('unread=true');
        if (options.limit) params.push('limit=' + options.limit);

        var query = params.length ? '?' + params.join('&') : '';

        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/notifications' + query, function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch notifications'));
            }
          });
        });
      },

      tags: function() {
        var self = this;
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/tags', function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch tags'));
            }
          });
        });
      },

      leaderboard: function(limit) {
        var self = this;
        var query = limit ? '?limit=' + limit : '';
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/forum/leaderboard' + query, function(response) {
            if (response) {
              resolve(response);
            } else {
              reject(new Error('Failed to fetch leaderboard'));
            }
          });
        });
      }
    },

    initForum: function() {
      var self = this;
      this.forum.parent = this;

      var containers = document.querySelectorAll('[data-minimaldoc-forum]');
      this.log('Found', containers.length, 'forum container(s)');

      containers.forEach(function(container) {
        self.renderForumContainer(container);
      });
    },

    renderForumContainer: function(container) {
      var self = this;
      var type = container.getAttribute('data-minimaldoc-forum') || 'latest';
      var limit = parseInt(container.getAttribute('data-limit'), 10) || 10;
      var category = container.getAttribute('data-category');
      var topicSlug = container.getAttribute('data-topic');

      container.innerHTML = '<div class="minimaldoc-forum-loading">Loading...</div>';

      if (type === 'categories') {
        this.forum.categories().then(function(response) {
          container.innerHTML = self.renderForumCategories(response.categories || []);
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-forum-error">Failed to load categories</div>';
        });
      } else if (type === 'topics') {
        this.forum.topics({ category: category, limit: limit }).then(function(response) {
          container.innerHTML = self.renderForumTopics(response.topics || []);
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-forum-error">Failed to load topics</div>';
        });
      } else if (type === 'topic' && topicSlug) {
        Promise.all([
          this.forum.getTopic(topicSlug),
          this.forum.getPosts(topicSlug)
        ]).then(function(results) {
          container.innerHTML = self.renderForumTopic(results[0].topic, results[1].posts || []);
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-forum-error">Topic not found</div>';
        });
      } else if (type === 'latest') {
        this.forum.topics({ limit: limit }).then(function(response) {
          container.innerHTML = self.renderForumTopics(response.topics || []);
        }).catch(function() {
          container.innerHTML = '<div class="minimaldoc-forum-error">Failed to load topics</div>';
        });
      } else if (type === 'search') {
        container.innerHTML = self.renderForumSearch();
      }
    },

    renderForumCategories: function(categories) {
      if (!categories.length) {
        return '<div class="minimaldoc-forum-empty">No categories</div>';
      }

      var html = '<div class="minimaldoc-forum-categories">';
      categories.forEach(function(cat) {
        html += '<div class="minimaldoc-forum-category">';
        html += '<div class="minimaldoc-forum-category-color" style="background:' + (cat.Color || '#3b82f6') + '"></div>';
        html += '<div class="minimaldoc-forum-category-info">';
        html += '<h3><a href="/forum/category/' + this.escapeHtml(cat.Slug) + '">' + this.escapeHtml(cat.Name) + '</a></h3>';
        if (cat.Description) {
          html += '<p>' + this.escapeHtml(cat.Description) + '</p>';
        }
        html += '</div>';
        html += '<div class="minimaldoc-forum-category-stats">';
        html += '<span>' + (cat.TopicCount || 0) + ' topics</span>';
        html += '</div>';
        html += '</div>';
      }, this);
      html += '</div>';
      return html;
    },

    renderForumTopics: function(topics) {
      if (!topics.length) {
        return '<div class="minimaldoc-forum-empty">No topics yet</div>';
      }

      var html = '<div class="minimaldoc-forum-topics">';
      topics.forEach(function(topic) {
        var badges = '';
        if (topic.is_pinned) badges += '<span class="minimaldoc-forum-badge minimaldoc-forum-badge--pinned">Pinned</span>';
        if (topic.is_solved) badges += '<span class="minimaldoc-forum-badge minimaldoc-forum-badge--solved">Solved</span>';

        html += '<div class="minimaldoc-forum-topic">';
        html += '<div class="minimaldoc-forum-topic-main">';
        html += '<h4><a href="/forum/topic/' + this.escapeHtml(topic.slug) + '">' + this.escapeHtml(topic.title) + '</a></h4>';
        html += badges;
        html += '<div class="minimaldoc-forum-topic-meta">';
        html += '<span>' + this.escapeHtml(topic.author_name || 'Anonymous') + '</span>';
        if (topic.category_name) {
          html += '<span class="minimaldoc-forum-topic-category">' + this.escapeHtml(topic.category_name) + '</span>';
        }
        html += '<span>' + this.formatDate(topic.created_at) + '</span>';
        html += '</div>';
        html += '</div>';
        html += '<div class="minimaldoc-forum-topic-stats">';
        html += '<span>' + topic.post_count + ' replies</span>';
        html += '<span>' + topic.view_count + ' views</span>';
        html += '<span>' + topic.like_count + ' likes</span>';
        html += '</div>';
        html += '</div>';
      }, this);
      html += '</div>';
      return html;
    },

    renderForumTopic: function(topic, posts) {
      var html = '<div class="minimaldoc-forum-topic-full">';

      // Topic header
      html += '<div class="minimaldoc-forum-topic-header">';
      html += '<h1>' + this.escapeHtml(topic.title) + '</h1>';
      html += '<div class="minimaldoc-forum-topic-meta">';
      html += '<span>by ' + this.escapeHtml(topic.author_name || 'Anonymous') + '</span>';
      html += '<span>' + this.formatDate(topic.created_at) + '</span>';
      html += '<span>' + topic.view_count + ' views</span>';
      html += '</div>';
      html += '</div>';

      // Topic content
      html += '<div class="minimaldoc-forum-topic-content">';
      html += topic.content_html || this.escapeHtml(topic.content);
      html += '</div>';

      // Posts
      html += '<div class="minimaldoc-forum-posts">';
      html += '<h3>' + posts.length + ' Replies</h3>';
      posts.forEach(function(post) {
        html += '<div class="minimaldoc-forum-post' + (post.is_solution ? ' minimaldoc-forum-post--solution' : '') + '">';
        if (post.is_solution) {
          html += '<div class="minimaldoc-forum-solution-badge">Solution</div>';
        }
        html += '<div class="minimaldoc-forum-post-author">';
        html += '<strong>' + this.escapeHtml(post.author_name || 'Anonymous') + '</strong>';
        html += '<span>' + this.formatDate(post.created_at) + '</span>';
        html += '</div>';
        html += '<div class="minimaldoc-forum-post-content">';
        html += post.content_html || this.escapeHtml(post.content);
        html += '</div>';
        html += '<div class="minimaldoc-forum-post-actions">';
        html += '<button onclick="MinimalDoc.forum.likePost(\'' + post.id + '\')">' + post.like_count + ' Likes</button>';
        html += '</div>';
        html += '</div>';
      }, this);
      html += '</div>';

      html += '</div>';
      return html;
    },

    renderForumSearch: function() {
      var html = '<div class="minimaldoc-forum-search">';
      html += '<form onsubmit="MinimalDoc.handleForumSearch(event)">';
      html += '<input type="search" name="q" placeholder="Search forum..." class="minimaldoc-forum-search-input">';
      html += '<button type="submit">Search</button>';
      html += '</form>';
      html += '<div id="minimaldoc-forum-search-results"></div>';
      html += '</div>';
      return html;
    },

    handleForumSearch: function(e) {
      e.preventDefault();
      var self = this;
      var form = e.target;
      var query = form.querySelector('input[name="q"]').value.trim();
      if (!query) return;

      var resultsDiv = document.getElementById('minimaldoc-forum-search-results');
      resultsDiv.innerHTML = '<div class="minimaldoc-forum-loading">Searching...</div>';

      this.forum.search(query, 20).then(function(response) {
        resultsDiv.innerHTML = self.renderForumTopics(response.topics || []);
      }).catch(function() {
        resultsDiv.innerHTML = '<div class="minimaldoc-forum-error">Search failed</div>';
      });
    },

    getAuthToken: function() {
      // Check for session cookie or stored token
      var cookies = document.cookie.split(';');
      for (var i = 0; i < cookies.length; i++) {
        var cookie = cookies[i].trim();
        if (cookie.indexOf('minimaldoc_session=') === 0) {
          return cookie.substring('minimaldoc_session='.length);
        }
      }
      return sessionStorage.getItem('minimaldoc_token') || null;
    },

    // HTTP
    post: function(path, data, callback) {
      var xhr = new XMLHttpRequest();
      xhr.open('POST', this.config.endpoint + path, true);
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.setRequestHeader('X-API-Key', this.config.siteId);
      var token = this.getAuthToken();
      if (token) {
        xhr.setRequestHeader('Authorization', 'Bearer ' + token);
      }
      xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && callback) {
          callback(xhr.status >= 200 && xhr.status < 300);
        }
      };
      xhr.send(JSON.stringify(data));
    },

    get: function(path, callback) {
      var self = this;
      var xhr = new XMLHttpRequest();
      xhr.open('GET', this.config.endpoint + path, true);
      xhr.setRequestHeader('X-API-Key', this.config.siteId);
      var token = this.getAuthToken();
      if (token) {
        xhr.setRequestHeader('Authorization', 'Bearer ' + token);
      }
      xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
          if (xhr.status >= 200 && xhr.status < 300) {
            try {
              callback(JSON.parse(xhr.responseText));
            } catch (e) {
              self.log('Failed to parse response:', e);
              callback(null);
            }
          } else {
            callback(null);
          }
        }
      };
      xhr.send();
    }
  };

  // Auto-init from script tag attributes
  // Note: document.currentScript may be null for deferred scripts in some browsers
  var script = document.currentScript || document.querySelector('script[data-endpoint][data-site-id]');
  if (script) {
    var endpoint = script.getAttribute('data-endpoint');
    var siteId = script.getAttribute('data-site-id');
    var features = script.getAttribute('data-features');
    var debug = script.hasAttribute('data-debug');

    if (endpoint && siteId) {
      var doInit = function() {
        MinimalDoc.init({
          endpoint: endpoint,
          siteId: siteId,
          features: features ? features.split(',') : ['analytics'],
          debug: debug
        });
      };

      // For deferred scripts, DOM may already be ready
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', doInit);
      } else {
        doInit();
      }
    }
  }

  // Export
  window.MinimalDoc = MinimalDoc;
})();
