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
      this.log('Tracked:', data.path);
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

      // Send duration on page unload
      window.addEventListener('beforeunload', function() {
        self.sendDuration();
      });

      // Also send duration periodically for long sessions (every 30 seconds)
      setInterval(function() {
        if (!document.hidden && self.pageStartTime > 0) {
          self.sendDuration();
        }
      }, 30000);
    },

    sendDuration: function() {
      if (!this.pageStartTime || !this.currentPath) return;

      var duration = Math.round((Date.now() - this.pageStartTime) / 1000);
      if (duration < 1) return; // Ignore sub-second visits

      var data = {
        site_id: this.config.siteId,
        path: this.currentPath,
        duration: duration,
        session_hash: this.sessionHash
      };

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

    // HTTP
    post: function(path, data, callback) {
      var xhr = new XMLHttpRequest();
      xhr.open('POST', this.config.endpoint + path, true);
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && callback) {
          callback(xhr.status >= 200 && xhr.status < 300);
        }
      };
      xhr.send(JSON.stringify(data));
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
