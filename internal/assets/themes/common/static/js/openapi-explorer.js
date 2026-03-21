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
    schemas: {},       // All schemas from spec
    currentView: window.OPENAPI_DEFAULT_VIEW || 'path',
    currentEndpoint: null,
    currentSchema: null,
    searchQuery: ''
  };

  // Resolve $ref references in schemas
  // Returns the resolved schema or the original if no ref
  function resolveSchema(schema) {
    if (!schema) return null;

    // If schema has a $ref, resolve it
    if (schema.Ref) {
      const refPath = schema.Ref;
      // Handle #/components/schemas/SchemaName format
      if (refPath.startsWith('#/components/schemas/')) {
        const schemaName = refPath.replace('#/components/schemas/', '');
        const resolved = state.schemas[schemaName];
        if (resolved) {
          // Return a merged object with resolved properties but keep track of original ref
          return { ...resolved, _resolvedFrom: refPath };
        }
      }
      // Return schema with just the ref if we can't resolve
      return schema;
    }

    // Recursively resolve nested schemas
    const resolved = { ...schema };

    if (resolved.Properties) {
      resolved.Properties = {};
      Object.keys(schema.Properties).forEach(key => {
        resolved.Properties[key] = resolveSchema(schema.Properties[key]);
      });
    }

    if (resolved.Items) {
      resolved.Items = resolveSchema(schema.Items);
    }

    if (resolved.AllOf) {
      resolved.AllOf = schema.AllOf.map(s => resolveSchema(s));
    }

    if (resolved.OneOf) {
      resolved.OneOf = schema.OneOf.map(s => resolveSchema(s));
    }

    if (resolved.AnyOf) {
      resolved.AnyOf = schema.AnyOf.map(s => resolveSchema(s));
    }

    return resolved;
  }

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
    setupSidebarCollapse();

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

      // Store schemas
      if (state.spec.schemas) {
        state.schemas = state.spec.schemas;
        populateSchemasNavigation();
      }
    } catch (error) {
      console.error('Failed to load spec data:', error);
    }
  }

  // Populate schemas navigation list
  function populateSchemasNavigation() {
    const navList = document.getElementById('schemas-nav-list');
    if (!navList || !state.schemas) return;

    const schemaNames = Object.keys(state.schemas).sort();
    if (schemaNames.length === 0) {
      navList.innerHTML = '<li class="openapi-nav-empty">No schemas defined</li>';
      return;
    }

    let html = '';
    schemaNames.forEach(name => {
      const schema = state.schemas[name];
      const typeLabel = schema.Type || 'object';
      html += `<li class="openapi-schema-item" data-schema-name="${escapeAttr(name)}">`;
      html += `<a href="#schema-${escapeAttr(name)}" class="openapi-schema-link">`;
      html += `<span class="openapi-schema-type">${escapeHtml(typeLabel)}</span>`;
      html += `<span class="openapi-schema-name">${escapeHtml(name)}</span>`;
      html += `</a></li>`;
    });
    navList.innerHTML = html;

    // Add click handlers
    navList.querySelectorAll('.openapi-schema-link').forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        const name = link.closest('.openapi-schema-item').dataset.schemaName;
        renderSchemaDetail(name);

        // Update active state
        navList.querySelectorAll('.openapi-schema-link').forEach(l => l.classList.remove('active'));
        link.classList.add('active');
      });
    });
  }

  // Load endpoints data (don't render - we show one at a time)
  async function loadEndpointsView(view) {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    try {
      // Load single endpoints.json file (all endpoints in flat array)
      if (state.allEndpoints.length === 0) {
        const response = await fetch('./endpoints.json');
        const data = await response.json();
        state.allEndpoints = data;
        state.endpoints.flat = data;
      }

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

  // Render schema detail view
  function renderSchemaDetail(schemaName) {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    const rawSchema = state.schemas[schemaName];
    if (!rawSchema) {
      container.innerHTML = `<div class="welcome-state"><p>Schema not found: ${escapeHtml(schemaName)}</p></div>`;
      return;
    }

    // Resolve the schema in case it has $refs
    const schema = resolveSchema(rawSchema);

    state.currentSchema = schemaName;
    state.currentEndpoint = null;

    const example = schemaToExample(schema);
    const hasSchema = schema.Properties || schema.Type;
    const hasExample = example !== null && example !== undefined;

    let html = `<div class="schema-detail">`;
    html += `<div class="schema-card">`;
    html += `<div class="schema-card-header">`;
    html += `<h2 class="schema-card-title">${escapeHtml(schemaName)}</h2>`;
    html += `<span class="property-type">${escapeHtml(schema.Type || 'object')}</span>`;
    html += `</div>`;
    html += `<div class="schema-card-body">`;

    // Description
    if (schema.Description) {
      html += `<div class="endpoint-description" style="margin-bottom: 1.5rem;">${schema.Description}</div>`;
    }

    // Only show schema viewer if there's content
    if (hasSchema || hasExample) {
      html += `<div class="schema-viewer">`;
      html += `<div class="schema-viewer-header">`;
      html += `<span class="schema-viewer-title">Schema Definition</span>`;

      // Only show tabs if we have both
      if (hasSchema && hasExample) {
        html += `<div class="schema-viewer-tabs">`;
        html += `<button class="schema-viewer-tab active" data-view="schema">Properties</button>`;
        html += `<button class="schema-viewer-tab" data-view="example">Example</button>`;
        html += `</div>`;
      }
      html += `</div>`;

      html += `<div class="schema-viewer-content">`;

      // Schema view
      if (hasSchema) {
        html += `<div class="schema-view-panel active" data-view="schema">`;
        html += renderSchemaProperties(schema);
        html += `</div>`;
      }

      // Example view
      if (hasExample) {
        const displayStyle = hasSchema ? 'display: none;' : '';
        html += `<div class="schema-view-panel${!hasSchema ? ' active' : ''}" data-view="example" style="${displayStyle}">`;
        html += `<pre class="schema-json-view">${escapeHtml(JSON.stringify(example, null, 2))}</pre>`;
        html += `</div>`;
      }

      html += `</div></div>`;
    }

    html += `</div></div></div>`;

    container.innerHTML = html;

    // Add tab handlers
    container.querySelectorAll('.schema-viewer-tab').forEach(tab => {
      tab.addEventListener('click', () => {
        const view = tab.dataset.view;
        const viewer = tab.closest('.schema-viewer');

        viewer.querySelectorAll('.schema-viewer-tab').forEach(t => {
          t.classList.toggle('active', t.dataset.view === view);
        });

        viewer.querySelectorAll('.schema-view-panel').forEach(p => {
          p.style.display = p.dataset.view === view ? 'block' : 'none';
        });
      });
    });
  }

  // Render all schemas page
  function renderAllSchemas() {
    const container = document.getElementById('endpoints-container');
    if (!container) return;

    state.currentEndpoint = null;
    state.currentSchema = null;

    const schemaNames = Object.keys(state.schemas).sort();

    if (schemaNames.length === 0) {
      container.innerHTML = `<div class="welcome-state"><h2>No Schemas</h2><p>This API specification does not define any reusable schemas.</p></div>`;
      return;
    }

    let html = `<div class="schemas-container">`;
    html += `<header class="openapi-header"><h1 class="openapi-title">Schemas</h1>`;
    html += `<p class="openapi-description">${schemaNames.length} reusable schema definitions</p></header>`;
    html += `<div class="schemas-list">`;

    schemaNames.forEach(name => {
      const rawSchema = state.schemas[name];
      const schema = resolveSchema(rawSchema);
      if (!schema) return;

      html += `<div class="schema-card" id="schema-${escapeAttr(name)}">`;
      html += `<div class="schema-card-header">`;
      html += `<h3 class="schema-card-title">${escapeHtml(name)}</h3>`;
      html += `<span class="property-type">${escapeHtml(schema.Type || 'object')}</span>`;
      html += `</div>`;
      html += `<div class="schema-card-body">`;

      if (schema.Description) {
        html += `<p style="margin-bottom: 1rem; color: var(--openapi-fg-secondary);">${schema.Description}</p>`;
      }

      html += renderSchemaProperties(schema);
      html += `</div></div>`;
    });

    html += `</div></div>`;

    container.innerHTML = html;
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

  // Get server URL
  function getServerURL() {
    const serverSelect = document.getElementById('server-select');
    if (serverSelect) {
      return serverSelect.value;
    }
    return state.spec && state.spec.Servers && state.spec.Servers[0] ? state.spec.Servers[0].URL : '';
  }

  // Build full URL for endpoint
  function buildFullURL(endpoint) {
    const server = getServerURL().replace(/\/$/, '');
    let url = server + endpoint.Path;
    return url;
  }

  // Render Authentication Section
  function renderAuthSection(security, securitySchemes) {
    if (!security || security.length === 0) return '';

    let html = `<div class="endpoint-section auth-section"><h4>Authentication</h4>`;
    html += `<div class="auth-schemes">`;

    security.forEach(req => {
      Object.keys(req).forEach(schemeName => {
        const scheme = securitySchemes && securitySchemes[schemeName];
        const scopes = req[schemeName] || [];

        html += `<div class="auth-scheme">`;
        html += `<div class="auth-scheme-header">`;

        if (scheme) {
          const authType = getAuthTypeDisplay(scheme);
          html += `<span class="auth-scheme-type">${escapeHtml(authType)}</span>`;
          html += `<span class="auth-scheme-name">${escapeHtml(schemeName)}</span>`;
        } else {
          html += `<span class="auth-scheme-type">Auth</span>`;
          html += `<span class="auth-scheme-name">${escapeHtml(schemeName)}</span>`;
        }

        html += `</div>`;

        // Description
        if (scheme && scheme.Description) {
          html += `<div class="auth-scheme-description">${scheme.Description}</div>`;
        }

        // Scheme details table
        if (scheme) {
          html += `<table class="auth-details-table">`;

          if (scheme.Type === 'http') {
            if (scheme.BearerFormat) {
              html += `<tr><td class="auth-detail-label">Format</td><td class="auth-detail-value">${escapeHtml(scheme.BearerFormat)}</td></tr>`;
            }
          } else if (scheme.Type === 'apiKey') {
            html += `<tr><td class="auth-detail-label">Type</td><td class="auth-detail-value">API Key</td></tr>`;
            html += `<tr><td class="auth-detail-label">Parameter</td><td class="auth-detail-value"><code>${escapeHtml(scheme.Name || '')}</code></td></tr>`;
            html += `<tr><td class="auth-detail-label">Location</td><td class="auth-detail-value">${escapeHtml(scheme.In || '')}</td></tr>`;
          } else if (scheme.Type === 'oauth2' && scheme.Flows) {
            html += `<tr><td class="auth-detail-label">Type</td><td class="auth-detail-value">OAuth 2.0</td></tr>`;

            // Show each flow with its URLs
            Object.keys(scheme.Flows).forEach(flowName => {
              const flow = scheme.Flows[flowName];
              if (!flow) return;

              html += `<tr><td class="auth-detail-label">Flow</td><td class="auth-detail-value">${escapeHtml(flowName)}</td></tr>`;
              if (flow.AuthorizationURL) {
                html += `<tr><td class="auth-detail-label">Auth URL</td><td class="auth-detail-value"><code>${escapeHtml(flow.AuthorizationURL)}</code></td></tr>`;
              }
              if (flow.TokenURL) {
                html += `<tr><td class="auth-detail-label">Token URL</td><td class="auth-detail-value"><code>${escapeHtml(flow.TokenURL)}</code></td></tr>`;
              }
              if (flow.RefreshURL) {
                html += `<tr><td class="auth-detail-label">Refresh URL</td><td class="auth-detail-value"><code>${escapeHtml(flow.RefreshURL)}</code></td></tr>`;
              }

              // Show available scopes from flow
              if (flow.Scopes && Object.keys(flow.Scopes).length > 0) {
                html += `<tr><td class="auth-detail-label">Available Scopes</td><td class="auth-detail-value">`;
                Object.keys(flow.Scopes).forEach(scopeName => {
                  const scopeDesc = flow.Scopes[scopeName];
                  html += `<div class="auth-scope-item"><code>${escapeHtml(scopeName)}</code>`;
                  if (scopeDesc) {
                    html += ` - ${escapeHtml(scopeDesc)}`;
                  }
                  html += `</div>`;
                });
                html += `</td></tr>`;
              }
            });
          } else if (scheme.Type === 'openIdConnect') {
            html += `<tr><td class="auth-detail-label">Type</td><td class="auth-detail-value">OpenID Connect</td></tr>`;
            if (scheme.OpenIdConnectUrl) {
              html += `<tr><td class="auth-detail-label">Discovery URL</td><td class="auth-detail-value"><code>${escapeHtml(scheme.OpenIdConnectUrl)}</code></td></tr>`;
            }
          }

          html += `</table>`;
        }

        // Required scopes for this endpoint
        if (scopes.length > 0) {
          html += `<div class="auth-scopes">`;
          html += `<div class="auth-scopes-label">Required Scopes</div>`;
          html += `<div class="auth-scopes-list">`;
          scopes.forEach(scope => {
            html += `<span class="auth-scope">${escapeHtml(scope)}</span>`;
          });
          html += `</div></div>`;
        }

        html += `</div>`;
      });
    });

    html += `</div></div>`;
    return html;
  }

  // Get readable auth type display
  function getAuthTypeDisplay(scheme) {
    if (!scheme) return 'Auth';

    if (scheme.Type === 'http') {
      if (scheme.Scheme === 'bearer') return 'Bearer Token';
      if (scheme.Scheme === 'basic') return 'Basic Auth';
      return scheme.Scheme || 'HTTP';
    }
    if (scheme.Type === 'apiKey') return 'API Key';
    if (scheme.Type === 'oauth2') return 'OAuth 2.0';
    if (scheme.Type === 'openIdConnect') return 'OpenID Connect';

    return scheme.Type || 'Auth';
  }

  // Render Request Body Section
  function renderRequestBodySection(requestBody) {
    if (!requestBody || !requestBody.Content) return '';

    // Check if there's any actual content to show
    const contentTypes = Object.keys(requestBody.Content);
    if (contentTypes.length === 0) return '';

    // Check if at least one content type has schema or example
    const hasContent = contentTypes.some(ct => {
      const mediaType = requestBody.Content[ct];
      if (!mediaType) return false;
      // Resolve schema to check for actual content
      const resolvedSchema = mediaType.Schema ? resolveSchema(mediaType.Schema) : null;
      return resolvedSchema || mediaType.Example;
    });
    if (!hasContent) return '';

    let html = `<div class="endpoint-section"><h4>Request Body</h4>`;

    // Show if required
    if (requestBody.Required) {
      html += `<p style="font-size: var(--openapi-font-sm); color: var(--status-5xx); margin-bottom: 1rem;">Required</p>`;
    }

    // Description
    if (requestBody.Description) {
      html += `<p style="margin-bottom: 1rem; color: var(--openapi-fg-secondary);">${escapeHtml(requestBody.Description)}</p>`;
    }

    // Content types
    contentTypes.forEach(contentType => {
      const mediaType = requestBody.Content[contentType];
      if (!mediaType) return;

      // Resolve schema references
      const resolvedSchema = mediaType.Schema ? resolveSchema(mediaType.Schema) : null;

      // Skip if no schema and no example
      if (!resolvedSchema && !mediaType.Example) return;

      const example = mediaType.Example || (resolvedSchema ? schemaToExample(resolvedSchema) : null);
      const hasSchema = resolvedSchema && (resolvedSchema.Properties || resolvedSchema.Type);
      const hasExample = example !== null && example !== undefined;

      // Skip if nothing to show
      if (!hasSchema && !hasExample) return;

      html += `<div class="schema-viewer">`;
      html += `<div class="schema-viewer-header">`;
      html += `<span class="schema-viewer-title">${escapeHtml(contentType)}</span>`;

      // Only show tabs if we have both schema and example
      if (hasSchema && hasExample) {
        html += `<div class="schema-viewer-tabs">`;
        html += `<button class="schema-viewer-tab active" data-view="schema">Schema</button>`;
        html += `<button class="schema-viewer-tab" data-view="example">Example</button>`;
        html += `</div>`;
      }
      html += `</div>`;

      html += `<div class="schema-viewer-content">`;

      // Schema view
      if (hasSchema) {
        html += `<div class="schema-view-panel active" data-view="schema">`;
        html += renderSchemaProperties(resolvedSchema);
        html += `</div>`;
      }

      // Example view
      if (hasExample) {
        const displayStyle = hasSchema ? 'display: none;' : '';
        html += `<div class="schema-view-panel${!hasSchema ? ' active' : ''}" data-view="example" style="${displayStyle}">`;
        html += `<pre class="schema-json-view">${escapeHtml(JSON.stringify(example, null, 2))}</pre>`;
        html += `</div>`;
      }

      html += `</div></div>`;
    });

    html += `</div>`;
    return html;
  }

  // Render Responses Section
  function renderResponsesSection(responses) {
    if (!responses || Object.keys(responses).length === 0) return '';

    const statusCodes = Object.keys(responses).sort();
    if (statusCodes.length === 0) return '';

    let html = `<div class="endpoint-section"><h4>Responses</h4>`;
    html += `<div class="response-tabs">`;
    statusCodes.forEach((code, i) => {
      html += `<button class="response-tab${i === 0 ? ' active' : ''}" data-code="${code}">${code}</button>`;
    });
    html += `</div>`;

    html += `<div class="response-panels">`;
    statusCodes.forEach((code, i) => {
      const response = responses[code];
      if (!response) return;

      html += `<div class="response-panel${i === 0 ? ' active' : ''}" data-code="${code}">`;

      if (response.Description) {
        html += `<p class="response-desc">${response.Description}</p>`;
      }

      // Response content
      if (response.Content) {
        Object.keys(response.Content).forEach(contentType => {
          const mediaType = response.Content[contentType];
          if (!mediaType) return;

          // Resolve schema references
          const resolvedSchema = mediaType.Schema ? resolveSchema(mediaType.Schema) : null;

          // Skip if no schema and no example
          if (!resolvedSchema && !mediaType.Example) return;

          const example = mediaType.Example || (resolvedSchema ? schemaToExample(resolvedSchema) : null);
          const hasSchema = resolvedSchema && (resolvedSchema.Properties || resolvedSchema.Type);
          const hasExample = example !== null && example !== undefined;

          // Skip if nothing to show
          if (!hasSchema && !hasExample) return;

          html += `<div class="schema-viewer" style="margin-top: 1rem;">`;
          html += `<div class="schema-viewer-header">`;
          html += `<span class="schema-viewer-title">${escapeHtml(contentType)}</span>`;

          // Only show tabs if we have both schema and example
          if (hasSchema && hasExample) {
            html += `<div class="schema-viewer-tabs">`;
            html += `<button class="schema-viewer-tab active" data-view="schema">Schema</button>`;
            html += `<button class="schema-viewer-tab" data-view="example">Example</button>`;
            html += `</div>`;
          }
          html += `</div>`;

          html += `<div class="schema-viewer-content">`;

          // Schema view
          if (hasSchema) {
            html += `<div class="schema-view-panel active" data-view="schema">`;
            html += renderSchemaProperties(resolvedSchema);
            html += `</div>`;
          }

          // Example view
          if (hasExample) {
            const displayStyle = hasSchema ? 'display: none;' : '';
            html += `<div class="schema-view-panel${!hasSchema ? ' active' : ''}" data-view="example" style="${displayStyle}">`;
            html += `<pre class="schema-json-view">${escapeHtml(JSON.stringify(example, null, 2))}</pre>`;
            html += `</div>`;
          }

          html += `</div></div>`;
        });
      }

      html += `</div>`;
    });
    html += `</div></div>`;

    return html;
  }

  // Render schema properties as a list
  function renderSchemaProperties(schema, depth = 0) {
    if (!schema) return '';

    // Resolve $ref if present
    const resolved = resolveSchema(schema);
    if (!resolved) return '';

    let html = '';

    if (resolved.Type === 'object' && resolved.Properties) {
      html += `<ul class="schema-properties-list${depth > 0 ? ' schema-nested' : ''}">`;

      Object.keys(resolved.Properties).forEach(propName => {
        const prop = resolveSchema(resolved.Properties[propName]);
        if (!prop) return;
        const isRequired = resolved.Required && resolved.Required.includes(propName);
        const typeDisplay = getSchemaTypeDisplay(prop);

        html += `<li class="schema-property-row">`;
        html += `<div class="schema-prop-name${isRequired ? ' required' : ''}">${escapeHtml(propName)}</div>`;
        html += `<div class="schema-prop-type">${escapeHtml(typeDisplay)}</div>`;
        html += `<div class="schema-prop-desc">`;
        if (prop.Description) {
          html += escapeHtml(prop.Description);
        }

        // Show constraints
        const constraints = getSchemaConstraints(prop);
        if (constraints.length > 0) {
          html += `<div class="schema-prop-constraints">`;
          constraints.forEach(c => {
            html += `<span class="schema-constraint">${escapeHtml(c)}</span>`;
          });
          html += `</div>`;
        }

        html += `</div>`;
        html += `</li>`;

        // Nested properties
        if (prop.Type === 'object' && prop.Properties && depth < 2) {
          html += renderSchemaProperties(prop, depth + 1);
        } else if (prop.Type === 'array' && prop.Items) {
          const resolvedItems = resolveSchema(prop.Items);
          if (resolvedItems && resolvedItems.Type === 'object' && resolvedItems.Properties && depth < 2) {
            html += renderSchemaProperties(prop.Items, depth + 1);
          }
        }
      });

      html += `</ul>`;
    } else if (resolved.Type === 'array' && resolved.Items) {
      const resolvedItems = resolveSchema(resolved.Items);
      if (resolvedItems && resolvedItems.Type === 'object' && resolvedItems.Properties) {
        html += `<p style="padding: 0.75rem 1rem; color: var(--openapi-fg-muted); font-size: var(--openapi-font-sm);">Array of objects:</p>`;
        html += renderSchemaProperties(resolved.Items, depth);
      } else {
        const itemType = getSchemaTypeDisplay(resolvedItems || resolved.Items);
        html += `<p style="padding: 0.75rem 1rem; font-size: var(--openapi-font-sm);">Array of <code style="background: var(--openapi-hover); padding: 0.2rem 0.4rem; border-radius: 3px;">${escapeHtml(itemType)}</code></p>`;
      }
    } else {
      const typeDisplay = getSchemaTypeDisplay(resolved);
      html += `<p style="padding: 0.75rem 1rem; font-size: var(--openapi-font-sm);">Type: <code style="background: var(--openapi-hover); padding: 0.2rem 0.4rem; border-radius: 3px;">${escapeHtml(typeDisplay)}</code></p>`;
    }

    return html;
  }

  // Get schema constraints
  // Note: Enum is excluded here because it's already shown in getSchemaTypeDisplay()
  function getSchemaConstraints(schema) {
    const constraints = [];
    if (!schema) return constraints;

    if (schema.Minimum != null) constraints.push(`min: ${schema.Minimum}`);
    if (schema.Maximum != null) constraints.push(`max: ${schema.Maximum}`);
    if (schema.MinLength != null) constraints.push(`minLength: ${schema.MinLength}`);
    if (schema.MaxLength != null) constraints.push(`maxLength: ${schema.MaxLength}`);
    if (schema.Pattern) constraints.push(`pattern: ${schema.Pattern}`);
    // Enum is shown in type display, not in constraints to avoid duplication
    if (schema.Default != null) constraints.push(`default: ${JSON.stringify(schema.Default)}`);
    if (schema.Format) constraints.push(`format: ${schema.Format}`);

    return constraints;
  }

  // Create endpoint detail element
  function createEndpointElement(endpoint) {
    const div = document.createElement('div');
    div.className = 'endpoint-detail';
    div.id = generateEndpointID(endpoint);

    const fullUrl = buildFullURL(endpoint);
    const securitySchemes = state.spec ? state.spec.securitySchemes : {};

    let html = '';

    // Main content
    html += `<div class="endpoint-main">`;

    // Header
    html += `<div class="endpoint-header">`;
    html += `<div class="endpoint-title-row">`;
    html += `<span class="endpoint-method endpoint-method-${endpoint.Method.toLowerCase()}">${endpoint.Method}</span>`;
    html += `<code class="endpoint-path">${escapeHtml(endpoint.Path)}</code>`;
    if (endpoint.Deprecated) {
      html += `<span class="endpoint-deprecated-badge">Deprecated</span>`;
    }
    html += `</div>`;
    html += `<div class="endpoint-url-row">`;
    html += `<button class="btn-copy-small" data-action="copy-url">Copy URL</button>`;
    html += `<button class="btn-test" data-action="test">Try it</button>`;
    html += `</div>`;
    html += `</div>`;

    // Summary & Description
    if (endpoint.Summary) {
      html += `<p class="endpoint-summary">${escapeHtml(endpoint.Summary)}</p>`;
    }
    if (endpoint.Description) {
      html += `<div class="endpoint-description">${endpoint.Description}</div>`;
    }

    // Authentication Section
    if (endpoint.Security && endpoint.Security.length > 0) {
      html += renderAuthSection(endpoint.Security, securitySchemes);
    }

    // Parameters
    if (endpoint.Parameters && endpoint.Parameters.length > 0) {
      html += `<div class="endpoint-section"><h4>Parameters</h4>`;
      html += `<table class="params-table"><thead><tr><th>Name</th><th>In</th><th>Type</th><th>Required</th><th>Description</th></tr></thead><tbody>`;
      endpoint.Parameters.forEach(param => {
        const reqClass = param.Required ? ' class="required"' : '';
        // Resolve parameter schema if it has a $ref
        const resolvedSchema = param.Schema ? resolveSchema(param.Schema) : null;
        const typeDisplay = resolvedSchema ? getSchemaTypeDisplay(resolvedSchema) : 'any';
        html += `<tr${reqClass}>`;
        html += `<td><code>${escapeHtml(param.Name)}</code></td>`;
        html += `<td>${escapeHtml(param.In)}</td>`;
        html += `<td>${escapeHtml(typeDisplay)}</td>`;
        html += `<td>${param.Required ? 'Yes' : 'No'}</td>`;
        html += `<td>${param.Description || ''}</td>`;
        html += `</tr>`;
      });
      html += `</tbody></table></div>`;
    }

    // Request Body
    if (endpoint.RequestBody && endpoint.RequestBody.Content) {
      html += renderRequestBodySection(endpoint.RequestBody);
    }

    // Responses
    if (endpoint.Responses && Object.keys(endpoint.Responses).length > 0) {
      html += renderResponsesSection(endpoint.Responses);
    }

    html += `</div>`;

    // Code Samples (right column)
    if (window.OPENAPI_ENABLE_CODE_SAMPLES !== false && window.CodeSamples) {
      html += `<div class="endpoint-samples">`;
      html += window.CodeSamples.createCodeSamplesSection(endpoint, getServerURL(), securitySchemes);
      html += `</div>`;
    }

    div.innerHTML = html;

    // Attach event listeners

    // Test button
    const testBtn = div.querySelector('[data-action="test"]');
    if (testBtn) {
      testBtn.addEventListener('click', () => loadEndpointInTester(endpoint));
    }

    // Copy URL button
    const copyUrlBtn = div.querySelector('[data-action="copy-url"]');
    if (copyUrlBtn) {
      copyUrlBtn.addEventListener('click', (e) => {
        e.preventDefault();
        copyToClipboard(fullUrl, e.currentTarget);
      });
    }

    // Copy JSON buttons
    div.querySelectorAll('[data-copy-json]').forEach(btn => {
      btn.addEventListener('click', (e) => {
        e.preventDefault();
        const json = btn.dataset.copyJson;
        copyToClipboard(json, btn);
      });
    });

    // Response tabs
    div.querySelectorAll('.response-tab').forEach(tab => {
      tab.addEventListener('click', () => {
        const code = tab.dataset.code;
        const section = tab.closest('.endpoint-section');

        // Update tabs
        section.querySelectorAll('.response-tab').forEach(t => {
          t.classList.toggle('active', t.dataset.code === code);
        });

        // Update panels
        section.querySelectorAll('.response-panel').forEach(p => {
          p.classList.toggle('active', p.dataset.code === code);
        });
      });
    });

    // Schema viewer tabs (Schema/Example toggle)
    div.querySelectorAll('.schema-viewer-tab').forEach(tab => {
      tab.addEventListener('click', () => {
        const view = tab.dataset.view;
        const viewer = tab.closest('.schema-viewer');

        // Update tabs
        viewer.querySelectorAll('.schema-viewer-tab').forEach(t => {
          t.classList.toggle('active', t.dataset.view === view);
        });

        // Update panels
        viewer.querySelectorAll('.schema-view-panel').forEach(p => {
          p.style.display = p.dataset.view === view ? 'block' : 'none';
        });
      });
    });

    // Code sample tabs (if CodeSamples loaded)
    if (window.CodeSamples) {
      window.CodeSamples.setupCodeSampleTabs(div);
    }

    return div;
  }

  // Copy to clipboard helper
  async function copyToClipboard(text, button) {
    try {
      await navigator.clipboard.writeText(text);

      // Visual feedback
      const originalText = button.textContent;
      button.textContent = 'Copied!';
      button.classList.add('copied');

      setTimeout(() => {
        button.textContent = originalText;
        button.classList.remove('copied');
      }, 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  }

  // Helper: escape attribute value
  function escapeAttr(text) {
    if (!text) return '';
    return text
      .replace(/&/g, '&amp;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;');
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

    // Resolve $ref if present
    const resolved = resolveSchema(schema);
    if (!resolved) return null;

    if (resolved.Example) return resolved.Example;

    if (resolved.Type === 'object' && resolved.Properties) {
      const obj = {};
      Object.keys(resolved.Properties).forEach(key => {
        obj[key] = schemaToExample(resolved.Properties[key]);
      });
      return obj;
    }

    if (resolved.Type === 'array' && resolved.Items) {
      return [schemaToExample(resolved.Items)];
    }

    // Return default based on type
    switch (resolved.Type) {
      case 'string': return resolved.Default || 'string';
      case 'number': return resolved.Default || 0;
      case 'integer': return resolved.Default || 0;
      case 'boolean': return resolved.Default || false;
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
          if (prop.Deprecated) {
            html += ' <span class="openapi-badge openapi-badge-deprecated">deprecated</span>';
          }
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
          if (prop.Deprecated) {
            html += ' <span class="openapi-badge openapi-badge-deprecated">deprecated</span>';
          }
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

    // Resolve $ref if present
    const resolved = resolveSchema(schema);
    if (!resolved) return 'any';

    let typeStr = resolved.Type || 'object';

    if (resolved.Type === 'array' && resolved.Items) {
      typeStr = getSchemaTypeDisplay(resolved.Items) + '[]';
    }

    if (resolved.Format) {
      typeStr += ` (${resolved.Format})`;
    }

    if (resolved.Enum && resolved.Enum.length > 0) {
      typeStr += ` enum: ${resolved.Enum.join(' | ')}`;
    }

    return typeStr;
  }

  // Setup view toggle
  function setupViewToggle() {
    const buttons = document.querySelectorAll('.view-toggle-btn');

    // Set initial active state from config
    const defaultView = state.currentView || 'path';
    buttons.forEach(btn => {
      const isActive = btn.dataset.view === defaultView;
      btn.classList.toggle('active', isActive);
    });

    // Set initial content visibility
    document.querySelectorAll('[data-view-content]').forEach(el => {
      el.classList.toggle('hidden', el.dataset.viewContent !== defaultView);
    });

    // Setup click handlers
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

        // Handle schemas view separately
        if (view === 'schemas') {
          state.currentView = view;
          renderAllSchemas();
          return;
        }

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

  // Filter endpoints and schemas based on search
  function filterEndpoints(query) {
    // Filter endpoint navigation items
    const navItems = document.querySelectorAll('.openapi-endpoint-item');
    navItems.forEach(item => {
      const text = item.textContent.toLowerCase();
      const matches = query === '' || text.includes(query);
      item.style.display = matches ? '' : 'none';
    });

    // Filter group headers - hide if all children are hidden
    const groups = document.querySelectorAll('.openapi-path-group, .openapi-tag-group');
    groups.forEach(group => {
      const visibleItems = group.querySelectorAll('.openapi-endpoint-item:not([style*="display: none"])');
      group.style.display = visibleItems.length > 0 ? '' : 'none';
    });

    // Filter schema navigation items
    const schemaItems = document.querySelectorAll('.openapi-schema-item');
    schemaItems.forEach(item => {
      const text = item.textContent.toLowerCase();
      const matches = query === '' || text.includes(query);
      item.style.display = matches ? '' : 'none';
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

    // Use event delegation - click anywhere on header to collapse
    navigation.addEventListener('click', (e) => {
      // Check if clicked on header (path or tag)
      const header = e.target.closest('.openapi-path-header, .openapi-tag-header');
      if (!header) return;

      // Don't collapse if clicking on an endpoint link
      if (e.target.closest('.openapi-endpoint-link')) return;

      e.stopPropagation();

      // Find the collapse button and toggle it
      const btn = header.querySelector('.openapi-collapse-btn');
      if (btn) {
        btn.classList.toggle('collapsed');
      }

      // Find the parent item
      const parent = header.closest('.openapi-nav-item');
      if (parent) {
        const isCollapsed = btn && btn.classList.contains('collapsed');
        // Toggle all lists (endpoints and nested paths)
        const lists = parent.querySelectorAll(':scope > .openapi-endpoint-list, :scope > .openapi-nav-nested');
        lists.forEach(list => {
          list.style.display = isCollapsed ? 'none' : '';
        });
      }
    });
  }

  // Setup sidebar collapse toggle
  function setupSidebarCollapse() {
    const toggle = document.getElementById('sidebar-collapse-toggle');
    const container = document.querySelector('.openapi-container');
    if (!toggle || !container) return;

    function updateToggleState(collapsed) {
      toggle.setAttribute('aria-label', collapsed ? 'Expand sidebar' : 'Collapse sidebar');
      toggle.setAttribute('title', collapsed ? 'Expand sidebar' : 'Collapse sidebar');
    }

    // Restore state from localStorage
    const isCollapsed = localStorage.getItem('openapi-sidebar-collapsed') === 'true';
    if (isCollapsed) {
      container.classList.add('sidebar-collapsed');
      updateToggleState(true);
    }

    toggle.addEventListener('click', () => {
      container.classList.toggle('sidebar-collapsed');
      const collapsed = container.classList.contains('sidebar-collapsed');
      localStorage.setItem('openapi-sidebar-collapsed', collapsed);
      updateToggleState(collapsed);
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
