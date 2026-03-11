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

    // Custom Events
    trackEvent: function(name, category, value) {
      if (!name) {
        this.log('Event name is required');
        return;
      }

      var data = {
        site_id: this.config.siteId,
        name: name,
        path: window.location.pathname,
        session_hash: this.sessionHash
      };

      if (category) {
        data.category = category;
      }
      if (value !== undefined && value !== null) {
        data.value = String(value);
      }

      this.post('/api/analytics/event', data);
      this.log('Event tracked:', name, category || '', value || '');
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
          ? '    <a href="' + self.config.endpoint + '/admin/login?redirect=' + encodeURIComponent(window.location.href) + '" class="btn btn-primary">Log In</a>'
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
        var params = ['site_id=' + encodeURIComponent(self.parent.config.siteId)];

        if (options.limit) params.push('limit=' + options.limit);
        if (options.offset) params.push('offset=' + options.offset);
        if (options.category) params.push('category=' + encodeURIComponent(options.category));
        if (options.tag) params.push('tag=' + encodeURIComponent(options.tag));
        if (options.search) params.push('q=' + encodeURIComponent(options.search));

        var query = '?' + params.join('&');

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
        var query = '?site_id=' + encodeURIComponent(self.parent.config.siteId);
        return new Promise(function(resolve, reject) {
          self.parent.get('/api/blog/posts/' + encodeURIComponent(slug) + query, function(response) {
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
        var params = ['site_id=' + encodeURIComponent(self.parent.config.siteId)];
        if (limit) params.push('limit=' + limit);
        var query = '?' + params.join('&');
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
          // Use custom render function if available
          if (window.renderBlogCards) {
            window.renderBlogCards(container, response.posts || []);
          } else {
            container.innerHTML = self.renderPostList(response.posts || [], template, false, compact);
            if (response.total > limit && !compact) {
              container.innerHTML += self.renderPagination(response.total, limit, 0);
            }
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
