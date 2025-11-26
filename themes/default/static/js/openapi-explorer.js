// OpenAPI Explorer - Main UI controller
(function() {
  'use strict';

  // State
  const state = {
    spec: null,
    endpoints: {
      byPath: [],
      byTag: [],
      flat: []
    },
    allEndpoints: [],  // Flat list of all endpoints for easy lookup
    currentView: window.OPENAPI_DEFAULT_VIEW || 'path',
    currentEndpoint: null,
    searchQuery: ''
  };

  // Initialize on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  async function init() {
    console.log('OpenAPI Explorer initializing...');

    // Load spec data
    await loadSpecData();

    // Setup UI
    setupViewToggle();
    setupSearch();
    setupNavigation();
    setupCollapseButtons();
    setupCopyButtons();

    // Load initial view
    await loadEndpointsView(state.currentView);

    // Handle hash navigation
    handleHashNavigation();
    window.addEventListener('hashchange', handleHashNavigation);

    console.log('OpenAPI Explorer ready');
  }

  // Load spec metadata
  async function loadSpecData() {
    try {
      const response = await fetch('./spec-data.json');
      state.spec = await response.json();
      console.log('Loaded spec:', state.spec.title, 'v' + state.spec.version);
    } catch (error) {
      console.error('Failed to load spec data:', error);
    }
  }

  // Load endpoints data (don't render - we show one at a time)
  async function loadEndpointsView(view) {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    try {
      let filename;
      switch (view) {
        case 'path':
          filename = 'endpoints-by-path-0.json';
          break;
        case 'tag':
          filename = 'endpoints-by-tag-0.json';
          break;
        case 'flat':
          filename = 'endpoints-flat-0.json';
          break;
        default:
          filename = 'endpoints-by-path-0.json';
      }

      const response = await fetch('./' + filename);
      const data = await response.json();

      // Store in state
      state.endpoints[viewToKey(view)] = data;

      // Build flat endpoint list for easy lookup
      state.allEndpoints = extractAllEndpoints(data, view);

      // Show welcome message if no endpoint is selected
      if (!state.currentEndpoint) {
        showWelcomeMessage();
      }

    } catch (error) {
      console.error('Failed to load endpoints:', error);
      const container = document.getElementById('endpoints-container');
      if (container) {
        container.innerHTML = '<div class="loading-state"><p>Failed to load endpoints</p></div>';
      }
    }
  }

  // Extract all endpoints into a flat list for easy lookup
  function extractAllEndpoints(data, view) {
    const endpoints = [];

    if (view === 'flat') {
      return data;
    } else if (view === 'tag') {
      data.forEach(group => {
        if (group.Endpoints) {
          endpoints.push(...group.Endpoints);
        }
      });
    } else if (view === 'path') {
      const extractFromGroups = (groups) => {
        groups.forEach(group => {
          if (group.Endpoints) {
            endpoints.push(...group.Endpoints);
          }
          if (group.Children) {
            extractFromGroups(group.Children);
          }
        });
      };
      extractFromGroups(data);
    }

    return endpoints;
  }

  // Show welcome message
  function showWelcomeMessage() {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    container.innerHTML = `
      <div class="welcome-state">
        <h2>Welcome to ${state.spec ? state.spec.title : 'API'} Documentation</h2>
        <p>Select an endpoint from the navigation to view its documentation.</p>
      </div>
    `;
  }

  // Render a single endpoint (Spotify-style)
  function renderSingleEndpoint(endpoint) {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    // Clear and render just this endpoint
    container.innerHTML = '';
    const el = createEndpointElement(endpoint);
    container.appendChild(el);

    // Update state
    state.currentEndpoint = endpoint;

    // Scroll to top
    container.scrollTop = 0;
  }

  // Create endpoint detail element
  function createEndpointElement(endpoint) {
    const div = document.createElement('div');
    div.className = 'endpoint-detail';
    div.id = generateEndpointID(endpoint);

    let html = `
      <div class="endpoint-detail-header">
        <div class="endpoint-method-path">
          <span class="endpoint-method endpoint-method-${endpoint.Method.toLowerCase()}">${endpoint.Method}</span>
          <span class="endpoint-path">${endpoint.Path}</span>
        </div>
        <div class="endpoint-actions">
          <button class="btn-icon" data-action="test" title="Test this endpoint">Test</button>
          <button class="btn-icon" data-action="copy-link" title="Copy link">Link</button>
        </div>
      </div>
    `;

    if (endpoint.Summary) {
      html += `<h2 class="endpoint-summary">${escapeHtml(endpoint.Summary)}</h2>`;
    }

    if (endpoint.Description) {
      html += `<div class="endpoint-description">${endpoint.Description}</div>`;
    }

    if (endpoint.Deprecated) {
      html += `<div class="endpoint-deprecated-notice"><strong>DEPRECATED:</strong> This endpoint is deprecated</div>`;
    }

    // Parameters
    if (endpoint.Parameters && endpoint.Parameters.length > 0) {
      html += `
        <div class="endpoint-section">
          <h3>Parameters</h3>
          <table class="parameters-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>In</th>
                <th>Required</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
      `;
      endpoint.Parameters.forEach(param => {
        html += `
          <tr>
            <td><code>${escapeHtml(param.Name)}</code></td>
            <td><code>${param.Schema ? escapeHtml(param.Schema.Type || 'any') : 'any'}</code></td>
            <td><span class="param-in">${escapeHtml(param.In)}</span></td>
            <td>${param.Required ? 'Yes' : 'No'}</td>
            <td>${param.Description || ''}</td>
          </tr>
        `;
      });
      html += `
            </tbody>
          </table>
        </div>
      `;
    }

    // Request Body
    if (endpoint.RequestBody) {
      html += `
        <div class="endpoint-section">
          <h3>Request Body${endpoint.RequestBody.Required ? ' (required)' : ''}</h3>
          ${endpoint.RequestBody.Description ? `<p>${endpoint.RequestBody.Description}</p>` : ''}
          <div class="schema-viewer">
            ${renderSchema(endpoint.RequestBody.Content)}
          </div>
        </div>
      `;
    }

    // Responses
    if (endpoint.Responses && Object.keys(endpoint.Responses).length > 0) {
      html += `<div class="endpoint-section"><h3>Responses</h3><div class="responses-container">`;

      Object.keys(endpoint.Responses).forEach(statusCode => {
        const response = endpoint.Responses[statusCode];
        html += `
          <div class="response-item">
            <div class="response-header">
              <span class="response-code response-code-${statusCode}">${statusCode}</span>
              <span class="response-description">${response.Description || ''}</span>
            </div>
            ${response.Content ? `<div class="response-content">${renderSchema(response.Content)}</div>` : ''}
          </div>
        `;
      });

      html += `</div></div>`;
    }

    // Security
    if (endpoint.Security && endpoint.Security.length > 0) {
      html += `<div class="endpoint-section"><h3>Security</h3><div class="security-requirements">`;

      endpoint.Security.forEach(req => {
        Object.keys(req).forEach(schemeName => {
          const scopes = req[schemeName];
          html += `
            <div class="security-requirement">
              <code>${escapeHtml(schemeName)}</code>
              ${scopes && scopes.length > 0 ? ` (scopes: ${scopes.map(s => escapeHtml(s)).join(', ')})` : ''}
            </div>
          `;
        });
      });

      html += `</div></div>`;
    }

    div.innerHTML = html;

    // Attach event listeners
    const testBtn = div.querySelector('[data-action="test"]');
    if (testBtn) {
      testBtn.addEventListener('click', () => loadEndpointInTester(endpoint));
    }

    const copyBtn = div.querySelector('[data-action="copy-link"]');
    if (copyBtn) {
      copyBtn.addEventListener('click', (e) => {
        e.preventDefault();
        copyEndpointLink(div.id, e.currentTarget);
      });
    }

    return div;
  }

  // Render schema with accordion (Spotify-style)
  function renderSchema(content) {
    if (!content) return '';

    let html = '<div class="schema-content">';
    Object.keys(content).forEach(contentType => {
      const mediaType = content[contentType];
      html += `<div class="schema-media-type">`;
      html += `<div class="schema-header">`;
      html += `<strong>${escapeHtml(contentType)}</strong>`;

      // Add copy button for schema
      const schemaId = `schema-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      if (mediaType.Schema) {
        html += `<button class="schema-copy-btn" data-schema-id="${schemaId}" aria-label="Copy JSON">Copy JSON</button>`;

        // Store schema in a global map for copy button access
        if (!window.__schemaMap) window.__schemaMap = {};
        window.__schemaMap[schemaId] = mediaType.Schema;
      }
      html += `</div>`;

      if (mediaType.Schema) {
        html += renderSchemaObject(mediaType.Schema, 0);
      } else if (mediaType.Example) {
        const jsonCode = JSON.stringify(mediaType.Example, null, 2);
        html += `
          <div class="code-block-wrapper">
            <button class="code-copy-btn" aria-label="Copy code">Copy</button>
            <pre class="schema-json">${jsonCode}</pre>
          </div>
        `;
      }

      html += '</div>';
    });
    html += '</div>';

    return html;
  }

  // Convert schema to example JSON
  function schemaToExample(schema) {
    if (!schema) return null;

    if (schema.Example) return schema.Example;

    if (schema.Type === 'object' && schema.Properties) {
      const obj = {};
      Object.keys(schema.Properties).forEach(key => {
        obj[key] = schemaToExample(schema.Properties[key]);
      });
      return obj;
    }

    if (schema.Type === 'array' && schema.Items) {
      return [schemaToExample(schema.Items)];
    }

    // Return default based on type
    switch (schema.Type) {
      case 'string': return schema.Default || 'string';
      case 'number': return schema.Default || 0;
      case 'integer': return schema.Default || 0;
      case 'boolean': return schema.Default || false;
      default: return null;
    }
  }

  // Render schema object with nested properties
  function renderSchemaObject(schema, depth = 0) {
    if (!schema) return '';

    let html = '';

    // For object types with properties, show as collapsible list
    if (schema.Type === 'object' && schema.Properties) {
      html += '<div class="schema-properties">';

      Object.keys(schema.Properties).forEach(propName => {
        const prop = schema.Properties[propName];
        const isRequired = schema.Required && schema.Required.includes(propName);
        const hasNested = (prop.Type === 'object' && prop.Properties) ||
                         (prop.Type === 'array' && prop.Items);

        html += '<div class="schema-property">';

        if (hasNested) {
          // Use details/summary for collapsible nested objects
          html += '<details class="schema-details">';
          html += '<summary class="schema-summary">';
          html += `<span class="property-name">${escapeHtml(propName)}</span>`;
          html += `<span class="property-type">${escapeHtml(getSchemaTypeDisplay(prop))}</span>`;
          if (isRequired) {
            html += '<span class="property-required">required</span>';
          }
          html += '</summary>';

          if (prop.Description) {
            html += `<div class="property-description">${prop.Description}</div>`;
          }

          // Render nested content
          if (prop.Type === 'object' && prop.Properties) {
            html += renderSchemaObject(prop, depth + 1);
          } else if (prop.Type === 'array' && prop.Items) {
            html += '<div class="schema-array-items">';
            html += '<span class="array-items-label">Array items:</span>';
            html += renderSchemaObject(prop.Items, depth + 1);
            html += '</div>';
          }

          html += '</details>';
        } else {
          // Simple property - no nesting
          html += '<div class="schema-property-simple">';
          html += `<span class="property-name">${escapeHtml(propName)}</span>`;
          html += `<span class="property-type">${escapeHtml(getSchemaTypeDisplay(prop))}</span>`;
          if (isRequired) {
            html += '<span class="property-required">required</span>';
          }
          if (prop.Description) {
            html += `<span class="property-description">${prop.Description}</span>`;
          }
          html += '</div>';
        }

        html += '</div>';
      });

      html += '</div>';
    }
    // For array types, show items schema
    else if (schema.Type === 'array' && schema.Items) {
      html += '<div class="schema-array-items">';
      html += '<span class="array-items-label">Array of:</span>';
      html += renderSchemaObject(schema.Items, depth);
      html += '</div>';
    }
    // For primitive types or no properties
    else {
      html += '<div class="schema-primitive">';
      html += `<span class="property-type">${escapeHtml(getSchemaTypeDisplay(schema))}</span>`;
      if (schema.Description) {
        html += `<span class="property-description">${schema.Description}</span>`;
      }
      html += '</div>';
    }

    return html;
  }

  // Get display string for schema type
  function getSchemaTypeDisplay(schema) {
    if (!schema) return 'any';

    let typeStr = schema.Type || 'object';

    if (schema.Type === 'array' && schema.Items) {
      typeStr = getSchemaTypeDisplay(schema.Items) + '[]';
    }

    if (schema.Format) {
      typeStr += ` (${schema.Format})`;
    }

    if (schema.Enum && schema.Enum.length > 0) {
      typeStr += ` enum: ${schema.Enum.join(' | ')}`;
    }

    return typeStr;
  }

  // Setup view toggle
  function setupViewToggle() {
    const buttons = document.querySelectorAll('.view-toggle-btn');
    buttons.forEach(btn => {
      btn.addEventListener('click', async () => {
        const view = btn.dataset.view;

        // Update active state
        buttons.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');

        // Update view content visibility
        document.querySelectorAll('[data-view-content]').forEach(el => {
          el.classList.toggle('hidden', el.dataset.viewContent !== view);
        });

        // Load and render endpoints
        state.currentView = view;
        await loadEndpointsView(view);
      });
    });
  }

  // Setup search
  function setupSearch() {
    const searchInput = document.getElementById('openapi-search');
    if (!searchInput) return;

    searchInput.addEventListener('input', debounce((e) => {
      state.searchQuery = e.target.value.toLowerCase();
      filterEndpoints(state.searchQuery);
    }, 300));
  }

  // Filter endpoints based on search
  function filterEndpoints(query) {
    // Filter navigation items in the sidebar, not main content
    const navItems = document.querySelectorAll('.openapi-endpoint-item');

    navItems.forEach(item => {
      const text = item.textContent.toLowerCase();
      const matches = query === '' || text.includes(query);
      item.style.display = matches ? '' : 'none';
    });

    // Also filter group headers - hide if all children are hidden
    const groups = document.querySelectorAll('.openapi-path-group, .openapi-tag-group');
    groups.forEach(group => {
      const visibleItems = group.querySelectorAll('.openapi-endpoint-item:not([style*="display: none"])');
      group.style.display = visibleItems.length > 0 ? '' : 'none';
    });
  }

  // Setup navigation item highlighting
  function setupNavigation() {
    const navigation = document.querySelector('.openapi-navigation');
    if (!navigation) return;

    // Use event delegation for nav links
    navigation.addEventListener('click', (e) => {
      const link = e.target.closest('.openapi-endpoint-link');
      if (!link) return;

      e.preventDefault();
      const href = link.getAttribute('href');
      if (href) {
        window.location.hash = href;
      }
    });
  }

  // Setup collapse buttons
  function setupCollapseButtons() {
    const navigation = document.querySelector('.openapi-navigation');
    if (!navigation) return;

    // Use event delegation for collapse buttons
    navigation.addEventListener('click', (e) => {
      const btn = e.target.closest('.openapi-collapse-btn');
      if (!btn) return;

      e.stopPropagation();
      btn.classList.toggle('collapsed');

      // Find the parent item
      const parent = btn.closest('.openapi-nav-item');
      if (parent) {
        // Toggle all lists (endpoints and nested paths)
        const lists = parent.querySelectorAll('.openapi-endpoint-list, .openapi-nav-nested');
        lists.forEach(list => {
          list.style.display = btn.classList.contains('collapsed') ? 'none' : '';
        });
      }
    });
  }

  // Setup copy buttons for code blocks
  function setupCopyButtons() {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    // Use event delegation for copy buttons
    container.addEventListener('click', async (e) => {
      const btn = e.target.closest('.code-copy-btn, .schema-copy-btn');
      if (!btn) return;

      e.preventDefault();

      // Get code from different sources
      let code;

      // Schema copy button
      if (btn.classList.contains('schema-copy-btn')) {
        const schemaId = btn.dataset.schemaId;
        if (schemaId && window.__schemaMap && window.__schemaMap[schemaId]) {
          const schema = window.__schemaMap[schemaId];
          code = JSON.stringify(schemaToExample(schema), null, 2);
        }
      }
      // Code copy button
      else if (btn.classList.contains('code-copy-btn')) {
        const wrapper = btn.closest('.code-block-wrapper');
        if (wrapper) {
          const pre = wrapper.querySelector('pre');
          if (pre) {
            code = pre.textContent;
          }
        }
      }

      if (!code) return;

      try {
        await navigator.clipboard.writeText(code);

        // Visual feedback
        const originalText = btn.textContent;
        btn.textContent = 'Copied!';
        btn.classList.add('copied');

        setTimeout(() => {
          btn.textContent = originalText;
          btn.classList.remove('copied');
        }, 2000);
      } catch (err) {
        console.error('Failed to copy code:', err);
        btn.textContent = 'Failed';
        setTimeout(() => {
          btn.textContent = originalText || 'Copy';
        }, 2000);
      }
    });
  }

  // Handle hash navigation
  function handleHashNavigation() {
    const hash = window.location.hash.slice(1);

    // If no hash, show welcome
    if (!hash) {
      showWelcomeMessage();
      return;
    }

    // Find endpoint by ID
    if (!state.allEndpoints || state.allEndpoints.length === 0) {
      console.warn('No endpoints loaded yet');
      return;
    }

    const endpoint = state.allEndpoints.find(ep => {
      const id = generateEndpointID(ep);
      return id === hash;
    });

    if (endpoint) {
      // Render this endpoint
      renderSingleEndpoint(endpoint);

      // Highlight active nav item
      document.querySelectorAll('.openapi-endpoint-link').forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('href') === '#' + hash) {
          link.classList.add('active');
        }
      });
    } else {
      console.warn('Endpoint not found:', hash);
      showWelcomeMessage();
    }
  }

  // Load endpoint in tester
  function loadEndpointInTester(endpoint) {
    state.currentEndpoint = endpoint;

    // Trigger API tester to show and populate
    if (window.APITester) {
      window.APITester.loadEndpoint(endpoint);
    }

    console.log('Loading endpoint in tester:', endpoint.Method, endpoint.Path);
  }

  // Copy endpoint link
  async function copyEndpointLink(id, button) {
    const url = window.location.origin + window.location.pathname + '#' + id;

    try {
      await navigator.clipboard.writeText(url);
      console.log('Link copied:', url);

      // Show visual feedback
      if (button) {
        const originalText = button.textContent;
        button.textContent = '✓ Copied';
        button.style.color = 'var(--method-post)';
        setTimeout(() => {
          button.textContent = originalText;
          button.style.color = '';
        }, 2000);
      }
    } catch (err) {
      console.error('Failed to copy link:', err);
      // Fallback for browsers that don't support clipboard API
      const textarea = document.createElement('textarea');
      textarea.value = url;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.select();
      try {
        document.execCommand('copy');
        if (button) {
          button.textContent = '✓ Copied';
          button.style.color = 'var(--method-post)';
          setTimeout(() => {
            button.textContent = 'Link';
            button.style.color = '';
          }, 2000);
        }
      } catch (e) {
        alert('Failed to copy link. URL: ' + url);
      }
      document.body.removeChild(textarea);
    }
  }

  // Helper: Convert view name to state key
  function viewToKey(view) {
    switch (view) {
      case 'path': return 'byPath';
      case 'tag': return 'byTag';
      case 'flat': return 'flat';
      default: return 'byPath';
    }
  }

  // Helper: Escape HTML
  function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  // Helper: Generate endpoint ID (matches Go template logic)
  function generateEndpointID(endpoint) {
    if (endpoint.OperationID) {
      return endpoint.OperationID;
    }
    // Fallback: use Method-Path with non-alphanumeric chars replaced
    let id = endpoint.Method + '-' + endpoint.Path;
    // Replace non-alphanumeric characters with dashes, avoiding consecutive dashes
    let result = '';
    let lastWasDash = false;
    for (let i = 0; i < id.length; i++) {
      const char = id[i];
      if (/[a-zA-Z0-9]/.test(char)) {
        result += char;
        lastWasDash = false;
      } else if (!lastWasDash) {
        result += '-';
        lastWasDash = true;
      }
    }
    // Trim trailing dash
    return result.replace(/-$/, '');
  }

  // Helper: Debounce
  function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  }

  // Export for use by other modules
  window.OpenAPIExplorer = {
    state,
    loadEndpointInTester
  };
})();
