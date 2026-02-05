/**
 * Status Uptime Display
 * Fetches uptime data from API endpoints (mode: api) and renders calendar grids
 * For mode: incidents, data is pre-rendered by the server
 */
(function () {
  'use strict';

  const basePath = document.documentElement.dataset.basePath || '';

  /**
   * Initialize uptime displays for API mode components
   */
  async function init() {
    const components = document.querySelectorAll('.uptime-component[data-uptime-mode="api"]');

    for (const component of components) {
      const componentId = component.dataset.componentId;
      const endpoint = component.dataset.uptimeEndpoint;

      if (!endpoint) {
        console.warn(`[status-uptime] No endpoint for component ${componentId}`);
        continue;
      }

      await fetchAndRenderUptime(component, componentId, endpoint);
    }

    // Add tooltip functionality to all uptime grids
    initTooltips();
  }

  /**
   * Fetch uptime data from API and render
   */
  async function fetchAndRenderUptime(container, componentId, endpoint) {
    const placeholder = container.querySelector('.uptime-grid-placeholder');

    try {
      // Build full URL from component URL + endpoint
      const componentUrl = getComponentUrl(componentId);
      if (!componentUrl) {
        showError(placeholder, 'No component URL configured');
        return;
      }

      const uptimeUrl = new URL(endpoint, componentUrl).toString();

      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 15000);

      const response = await fetch(uptimeUrl, {
        method: 'GET',
        mode: 'cors',
        signal: controller.signal,
        headers: {
          'Accept': 'application/json'
        }
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        showError(placeholder, `HTTP ${response.status}`);
        return;
      }

      const data = await response.json();
      renderUptimeData(container, placeholder, data);

    } catch (err) {
      if (err.name === 'AbortError') {
        showError(placeholder, 'Request timeout');
      } else {
        showError(placeholder, 'Failed to load');
        console.error(`[status-uptime] Error fetching uptime for ${componentId}:`, err);
      }
    }
  }

  /**
   * Get component URL from status.json data
   */
  function getComponentUrl(componentId) {
    // Try to find URL from data attributes or fetch from status.json
    const statusPath = document.documentElement.dataset.statusPath || 'status';
    const statusJsonUrl = `${basePath}/${statusPath}/status.json`;

    // Synchronous lookup from cached data or return null
    // In practice, we'd cache the status.json on page load
    return null; // Placeholder - component URL should be in data attributes
  }

  /**
   * Render uptime data into the container
   */
  function renderUptimeData(container, placeholder, data) {
    // Update percentage
    const percentEl = container.querySelector('.uptime-percent');
    if (percentEl && data.uptime_percent !== undefined) {
      percentEl.textContent = data.uptime_percent.toFixed(2) + '%';
      percentEl.classList.remove('sla-met', 'sla-missed');
      if (data.sla && data.uptime_percent >= data.sla.target) {
        percentEl.classList.add('sla-met');
      } else {
        percentEl.classList.add('sla-missed');
      }
    }

    // Update SLA stats
    if (data.sla) {
      const statsEl = container.querySelector('.uptime-sla-stats');
      if (statsEl) {
        statsEl.innerHTML = `
          <span class="sla-period" title="Last 7 days">7d: ${data.sla.current_7d.toFixed(2)}%</span>
          <span class="sla-period" title="Last 30 days">30d: ${data.sla.current_30d.toFixed(2)}%</span>
          <span class="sla-period" title="Last 90 days">90d: ${data.sla.current_90d.toFixed(2)}%</span>
          <span class="sla-target" title="SLA Target">Target: ${data.sla.target.toFixed(1)}%</span>
        `;
      }
    }

    // Render grid
    if (data.history && placeholder) {
      const grid = document.createElement('div');
      grid.className = 'uptime-grid';
      grid.dataset.componentId = container.dataset.componentId;

      for (const day of data.history) {
        const dayEl = document.createElement('div');
        dayEl.className = `uptime-day uptime-${day.status}`;
        dayEl.dataset.date = day.date;
        dayEl.dataset.status = day.status;
        dayEl.dataset.downtime = day.downtime_minutes;
        dayEl.title = `${day.date}: ${day.status}${day.downtime_minutes > 0 ? ` (${day.downtime_minutes}m downtime)` : ''}`;
        grid.appendChild(dayEl);
      }

      placeholder.replaceWith(grid);
    }
  }

  /**
   * Show error message in placeholder
   */
  function showError(placeholder, message) {
    if (placeholder) {
      placeholder.innerHTML = `<p class="loading-uptime error">Uptime unavailable: ${message}</p>`;
    }
  }

  /**
   * Initialize custom tooltips for uptime days
   */
  function initTooltips() {
    // Create tooltip element
    const tooltip = document.createElement('div');
    tooltip.className = 'uptime-tooltip';
    tooltip.style.display = 'none';
    document.body.appendChild(tooltip);

    // Status labels
    const statusLabels = {
      'operational': 'Operational',
      'degraded': 'Degraded',
      'partial_outage': 'Partial Outage',
      'major_outage': 'Major Outage',
      'maintenance': 'Maintenance'
    };

    // Event delegation for all uptime days
    document.addEventListener('mouseover', (e) => {
      const day = e.target.closest('.uptime-day');
      if (!day) return;

      const date = day.dataset.date;
      const status = day.dataset.status;
      const downtime = parseInt(day.dataset.downtime, 10) || 0;

      // Format date nicely
      const dateObj = new Date(date + 'T00:00:00');
      const formattedDate = dateObj.toLocaleDateString('en-US', {
        weekday: 'short',
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });

      // Build tooltip content
      let content = `<strong>${formattedDate}</strong><br>`;
      content += `Status: ${statusLabels[status] || status}`;
      if (downtime > 0) {
        const hours = Math.floor(downtime / 60);
        const mins = downtime % 60;
        content += `<br>Downtime: ${hours > 0 ? hours + 'h ' : ''}${mins}m`;
      }

      tooltip.innerHTML = content;
      tooltip.style.display = 'block';

      // Position tooltip
      const rect = day.getBoundingClientRect();
      tooltip.style.left = rect.left + window.scrollX + (rect.width / 2) - (tooltip.offsetWidth / 2) + 'px';
      tooltip.style.top = rect.top + window.scrollY - tooltip.offsetHeight - 8 + 'px';

      // Remove title to prevent native tooltip
      day.removeAttribute('title');
    });

    document.addEventListener('mouseout', (e) => {
      const day = e.target.closest('.uptime-day');
      if (day) {
        tooltip.style.display = 'none';
      }
    });
  }

  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
