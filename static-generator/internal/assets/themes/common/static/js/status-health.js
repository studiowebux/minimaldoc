/**
 * Status Health Polling
 * Polls health endpoints for components and updates UI in real-time
 */
(function () {
  'use strict';

  const basePath = document.documentElement.dataset.basePath || '';
  const statusPath = document.documentElement.dataset.statusPath || 'status';
  const statusJsonUrl = `${basePath}/${statusPath}/status.json`;
  const pollingIntervals = new Map();
  const componentState = new Map(); // Track state per component

  /**
   * Fetch status.json and start polling for components with health endpoints
   */
  async function init() {
    try {
      const response = await fetch(statusJsonUrl);
      if (!response.ok) {
        console.warn('[status-health] Failed to fetch status.json:', response.status);
        return;
      }

      const data = await response.json();
      if (!data.components || !Array.isArray(data.components)) {
        console.warn('[status-health] No components found in status.json');
        return;
      }

      for (const component of data.components) {
        if (component.health_endpoint && component.health_interval > 0) {
          // Initialize state for this component
          componentState.set(component.id, {
            interval: component.health_interval,
            lastCheck: null,
            nextCheck: null,
            latency: null,
            status: 'unknown'
          });

          // Add health indicator UI
          addHealthIndicator(component);

          // Start polling
          startPolling(component);
        }
      }

      // Start countdown timer updates
      startCountdownUpdates();
    } catch (err) {
      console.error('[status-health] Error initializing:', err);
    }
  }

  /**
   * Add health check indicator element to component
   */
  function addHealthIndicator(component) {
    const componentEl = document.querySelector(`[data-component-id="${component.id}"]`);
    if (!componentEl) return;

    // Check if indicator already exists
    if (componentEl.querySelector('.health-indicator')) return;

    const indicator = document.createElement('div');
    indicator.className = 'health-indicator';
    indicator.innerHTML = `
      <span class="health-pulse" title="Health check active"></span>
      <span class="health-info">
        <span class="health-latency" title="Response time">--</span>
        <span class="health-next" title="Next check">--</span>
      </span>
    `;
    componentEl.appendChild(indicator);
  }

  /**
   * Start polling a component's health endpoint
   */
  function startPolling(component) {
    const healthUrl = buildHealthUrl(component);
    if (!healthUrl) {
      console.warn(`[status-health] Cannot build health URL for ${component.id}`);
      return;
    }

    // Initial check
    checkHealth(component.id, healthUrl, component.health_interval);

    // Set up interval polling
    const intervalMs = component.health_interval * 1000;
    const intervalId = setInterval(() => {
      checkHealth(component.id, healthUrl, component.health_interval);
    }, intervalMs);

    pollingIntervals.set(component.id, intervalId);
  }

  /**
   * Build the full health endpoint URL
   */
  function buildHealthUrl(component) {
    if (!component.url || !component.health_endpoint) {
      return null;
    }

    try {
      const base = new URL(component.url);
      const endpoint = component.health_endpoint.startsWith('/')
        ? component.health_endpoint
        : '/' + component.health_endpoint;
      return new URL(endpoint, base).toString();
    } catch (err) {
      console.warn(`[status-health] Invalid URL for ${component.id}:`, err);
      return null;
    }
  }

  /**
   * Check a single health endpoint and update UI
   */
  async function checkHealth(componentId, healthUrl, intervalSeconds) {
    const componentEl = document.querySelector(`[data-component-id="${componentId}"]`);
    const statusBadge = componentEl?.querySelector('.component-status');
    const state = componentState.get(componentId);

    const startTime = performance.now();

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 10000);

      // Show pulse animation
      const pulseEl = componentEl?.querySelector('.health-pulse');
      if (pulseEl) pulseEl.classList.add('checking');

      const response = await fetch(healthUrl, {
        method: 'GET',
        mode: 'cors',
        signal: controller.signal,
        headers: {
          'Accept': 'application/json'
        }
      });

      clearTimeout(timeoutId);
      const latency = Math.round(performance.now() - startTime);

      // Update state
      if (state) {
        state.lastCheck = Date.now();
        state.nextCheck = Date.now() + (intervalSeconds * 1000);
        state.latency = latency;
      }

      if (response.ok) {
        updateComponentStatus(componentEl, statusBadge, 'operational', latency);
        if (state) state.status = 'operational';
      } else if (response.status >= 500) {
        updateComponentStatus(componentEl, statusBadge, 'major_outage', latency);
        if (state) state.status = 'major_outage';
      } else {
        updateComponentStatus(componentEl, statusBadge, 'degraded', latency);
        if (state) state.status = 'degraded';
      }

      // Remove pulse animation
      if (pulseEl) pulseEl.classList.remove('checking');

    } catch (err) {
      const latency = Math.round(performance.now() - startTime);

      if (state) {
        state.lastCheck = Date.now();
        state.nextCheck = Date.now() + (intervalSeconds * 1000);
        state.latency = null;
      }

      if (err.name === 'AbortError') {
        updateComponentStatus(componentEl, statusBadge, 'degraded', null);
        if (state) state.status = 'degraded';
      } else {
        updateComponentStatus(componentEl, statusBadge, 'major_outage', null);
        if (state) state.status = 'major_outage';
      }

      const pulseEl = componentEl?.querySelector('.health-pulse');
      if (pulseEl) pulseEl.classList.remove('checking');
    }
  }

  /**
   * Update the component's status in the UI
   */
  function updateComponentStatus(componentEl, statusBadge, status, latency) {
    if (!statusBadge) return;

    // Remove existing status classes
    statusBadge.classList.remove(
      'status-green',
      'status-yellow',
      'status-orange',
      'status-red',
      'status-blue'
    );

    // Add new status class and update text
    const statusConfig = getStatusConfig(status);
    statusBadge.classList.add(`status-${statusConfig.color}`);
    statusBadge.textContent = statusConfig.label;

    // Update component item class for styling
    if (componentEl) {
      componentEl.dataset.liveStatus = status;

      // Update latency display
      const latencyEl = componentEl.querySelector('.health-latency');
      if (latencyEl) {
        if (latency !== null) {
          latencyEl.textContent = `${latency}ms`;
          latencyEl.title = `Response time: ${latency}ms`;
        } else {
          latencyEl.textContent = 'err';
          latencyEl.title = 'Request failed';
        }
      }
    }
  }

  /**
   * Start countdown timer updates
   */
  function startCountdownUpdates() {
    setInterval(() => {
      const now = Date.now();
      for (const [componentId, state] of componentState) {
        const componentEl = document.querySelector(`[data-component-id="${componentId}"]`);
        const nextEl = componentEl?.querySelector('.health-next');
        if (!nextEl || !state.nextCheck) continue;

        const remaining = Math.max(0, Math.ceil((state.nextCheck - now) / 1000));
        nextEl.textContent = `${remaining}s`;
        nextEl.title = `Next check in ${remaining} seconds`;
      }
    }, 1000);
  }

  /**
   * Get status configuration (color and label)
   */
  function getStatusConfig(status) {
    const configs = {
      operational: { color: 'green', label: 'Operational' },
      degraded: { color: 'yellow', label: 'Degraded' },
      partial_outage: { color: 'orange', label: 'Partial Outage' },
      major_outage: { color: 'red', label: 'Major Outage' },
      maintenance: { color: 'blue', label: 'Maintenance' }
    };
    return configs[status] || configs.operational;
  }

  /**
   * Stop all polling (cleanup)
   */
  function stopAllPolling() {
    for (const [id, intervalId] of pollingIntervals) {
      clearInterval(intervalId);
    }
    pollingIntervals.clear();
  }

  // Clean up on page unload
  window.addEventListener('beforeunload', stopAllPolling);

  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
