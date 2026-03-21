// API Tester - Interactive API testing interface
(function() {
  'use strict';

  // State
  const state = {
    endpoint: null,
    server: '',
    auth: {
      type: 'none',
      token: '',
      apiKey: '',
      apiKeyLocation: 'header'
    },
    parameters: {},
    headers: {},
    body: ''
  };

  // Initialize
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  function init() {
    setupServerSelect();
    setupAuthType();
    setupSendButton();
    setupExportButtons();
    setupCloseButton();
    setupAddHeaderButton();

  }

  // Load endpoint into tester
  function loadEndpoint(endpoint) {
    state.endpoint = endpoint;

    // Clear previous state
    state.parameters = {};
    state.headers = {};
    state.body = '';

    // Show tester
    const tester = document.getElementById('openapi-sidebar-right');
    if (tester) {
      tester.classList.add('active');
    }

    // Populate UI
    populateParameters();
    populateRequestBody();

  }

  // Setup server select
  function setupServerSelect() {
    const select = document.getElementById('server-select');
    if (!select) return;

    select.addEventListener('change', (e) => {
      state.server = e.target.value;
    });

    // Set initial value
    if (select.options.length > 0) {
      state.server = select.options[0].value;
    }
  }

  // Setup auth type
  function setupAuthType() {
    const authSelect = document.getElementById('auth-type');
    const authConfig = document.getElementById('auth-config');
    if (!authSelect || !authConfig) return;

    authSelect.addEventListener('change', (e) => {
      state.auth.type = e.target.value;
      renderAuthConfig();
    });
  }

  // Render auth configuration UI
  function renderAuthConfig() {
    const authConfig = document.getElementById('auth-config');
    if (!authConfig) return;

    authConfig.innerHTML = '';

    switch (state.auth.type) {
      case 'bearer':
        authConfig.innerHTML = `
          <input type="text"
                 id="bearer-token-input"
                 placeholder="Enter Bearer token"
                 class="api-tester-select"
                 style="margin-top: 0.5rem"
                 value="${state.auth.token}">
        `;
        const bearerInput = document.getElementById('bearer-token-input');
        if (bearerInput) {
          bearerInput.addEventListener('input', (e) => {
            state.auth.token = e.target.value;
          });
        }
        break;

      case 'apikey':
        authConfig.innerHTML = `
          <input type="text"
                 id="apikey-input"
                 placeholder="Enter API key"
                 class="api-tester-select"
                 style="margin-top: 0.5rem"
                 value="${state.auth.apiKey}">
          <select id="apikey-location" class="api-tester-select" style="margin-top: 0.5rem">
            <option value="header">Header (X-API-Key)</option>
            <option value="query">Query Parameter</option>
          </select>
        `;
        const apikeyInput = document.getElementById('apikey-input');
        const apikeyLocation = document.getElementById('apikey-location');

        if (apikeyInput) {
          apikeyInput.addEventListener('input', (e) => {
            state.auth.apiKey = e.target.value;
          });
        }

        if (apikeyLocation) {
          apikeyLocation.value = state.auth.apiKeyLocation;
          apikeyLocation.addEventListener('change', (e) => {
            state.auth.apiKeyLocation = e.target.value;
          });
        }
        break;

      case 'oauth2':
        authConfig.innerHTML = `
          <button class="btn btn-secondary" id="oauth-authorize-btn" style="margin-top: 0.5rem; width: 100%">
            Authorize with OAuth 2.0
          </button>
          <input type="text"
                 id="oauth-token-display"
                 placeholder="Token will appear here"
                 class="api-tester-select"
                 style="margin-top: 0.5rem"
                 readonly
                 value="${state.auth.token}">
        `;
        const oauthBtn = document.getElementById('oauth-authorize-btn');
        if (oauthBtn) {
          oauthBtn.addEventListener('click', () => {
            if (window.OAuthHandler) {
              window.OAuthHandler.authorize((token) => {
                state.auth.token = token;
                const display = document.getElementById('oauth-token-display');
                if (display) display.value = token;
              });
            }
          });
        }
        break;
    }
  }

  // Populate parameters
  function populateParameters() {
    const container = document.getElementById('parameters-container');
    if (!container || !state.endpoint || !state.endpoint.Parameters) {
      if (container) container.innerHTML = '<p style="color: #666; font-size: 0.875rem">No parameters</p>';
      return;
    }

    container.innerHTML = '';

    state.endpoint.Parameters.forEach(param => {
      const div = document.createElement('div');
      div.style.marginBottom = '0.75rem';

      const label = document.createElement('label');
      label.style.display = 'block';
      label.style.marginBottom = '0.25rem';
      label.style.fontSize = '0.875rem';
      label.style.fontWeight = '600';
      label.innerHTML = `
        ${param.Name}
        <span style="color: #666; font-weight: normal">(${param.In})</span>
        ${param.Required ? '<span style="color: #f93e3e">*</span>' : ''}
      `;

      const input = document.createElement('input');
      input.type = 'text';
      input.className = 'api-tester-select';
      input.placeholder = param.Description || `Enter ${param.Name}`;
      input.value = state.parameters[param.Name] || '';

      input.addEventListener('input', (e) => {
        state.parameters[param.Name] = e.target.value;
      });

      div.appendChild(label);
      div.appendChild(input);
      container.appendChild(div);
    });
  }

  // Populate request body
  function populateRequestBody() {
    const section = document.getElementById('request-body-section');
    const bodyEditor = document.getElementById('request-body');
    const contentTypeSelect = document.getElementById('content-type-select');

    if (!section || !bodyEditor) return;

    if (state.endpoint && state.endpoint.RequestBody) {
      section.style.display = '';

      // Populate content type options
      if (contentTypeSelect && state.endpoint.RequestBody.Content) {
        contentTypeSelect.innerHTML = '';
        Object.keys(state.endpoint.RequestBody.Content).forEach(contentType => {
          const option = document.createElement('option');
          option.value = contentType;
          option.textContent = contentType;
          contentTypeSelect.appendChild(option);
        });

        // Generate example body
        const firstContentType = Object.keys(state.endpoint.RequestBody.Content)[0];
        const mediaType = state.endpoint.RequestBody.Content[firstContentType];
        if (mediaType && mediaType.Schema) {
          bodyEditor.value = JSON.stringify(schemaToExample(mediaType.Schema), null, 2);
          state.body = bodyEditor.value;
        }
      }

      bodyEditor.addEventListener('input', (e) => {
        state.body = e.target.value;
      });
    } else {
      section.style.display = 'none';
    }
  }

  // Setup add header button
  function setupAddHeaderButton() {
    const btn = document.getElementById('add-header-btn');
    if (!btn) return;

    btn.addEventListener('click', () => {
      addHeaderRow();
    });
  }

  // Add header row
  function addHeaderRow(name = '', value = '') {
    const container = document.getElementById('headers-container');
    if (!container) return;

    const div = document.createElement('div');
    div.style.display = 'flex';
    div.style.gap = '0.5rem';
    div.style.marginBottom = '0.5rem';

    const nameInput = document.createElement('input');
    nameInput.type = 'text';
    nameInput.placeholder = 'Header name';
    nameInput.className = 'api-tester-select';
    nameInput.style.flex = '1';
    nameInput.value = name;

    const valueInput = document.createElement('input');
    valueInput.type = 'text';
    valueInput.placeholder = 'Value';
    valueInput.className = 'api-tester-select';
    valueInput.style.flex = '2';
    valueInput.value = value;

    const removeBtn = document.createElement('button');
    removeBtn.textContent = '×';
    removeBtn.className = 'btn-icon';
    removeBtn.style.width = '2.5rem';
    removeBtn.addEventListener('click', () => {
      div.remove();
      delete state.headers[nameInput.value];
    });

    // Track the current header name to handle name changes
    let currentHeaderName = name;

    // Initialize header in state if name and value are provided
    if (name) {
      state.headers[name] = value;
    }

    // Update header when name changes (on blur or Enter)
    const updateHeaderName = () => {
      const newName = nameInput.value.trim();
      if (newName && newName !== currentHeaderName) {
        // Remove old header
        if (currentHeaderName) {
          delete state.headers[currentHeaderName];
        }
        // Add new header
        state.headers[newName] = valueInput.value;
        currentHeaderName = newName;
      }
    };

    nameInput.addEventListener('blur', updateHeaderName);
    nameInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') {
        updateHeaderName();
      }
    });

    // Update header value immediately
    valueInput.addEventListener('input', () => {
      if (currentHeaderName) {
        state.headers[currentHeaderName] = valueInput.value;
      }
    });

    div.appendChild(nameInput);
    div.appendChild(valueInput);
    div.appendChild(removeBtn);

    // Insert before the add button
    const addBtn = document.getElementById('add-header-btn');
    container.insertBefore(div, addBtn);
  }

  // Setup send button
  function setupSendButton() {
    const btn = document.getElementById('send-request-btn');
    if (!btn) return;

    btn.addEventListener('click', async () => {
      await sendRequest();
    });
  }

  // Send API request
  async function sendRequest() {
    if (!state.endpoint) return;

    const responseSection = document.getElementById('response-section');
    const responseStatus = document.getElementById('response-status');
    const responseHeaders = document.getElementById('response-headers');
    const responseBody = document.getElementById('response-body');

    if (!responseSection || !responseStatus || !responseHeaders || !responseBody) return;

    // Show response section
    responseSection.style.display = '';
    responseStatus.innerHTML = 'Sending request...';
    responseHeaders.innerHTML = '';
    responseBody.textContent = '';

    try {
      // Build URL
      let url = state.server + state.endpoint.Path;

      // Replace path parameters
      Object.keys(state.parameters).forEach(name => {
        const param = state.endpoint.Parameters.find(p => p.Name === name && p.In === 'path');
        if (param) {
          url = url.replace(`{${name}}`, encodeURIComponent(state.parameters[name]));
        }
      });

      // Add query parameters
      const queryParams = new URLSearchParams();
      Object.keys(state.parameters).forEach(name => {
        const param = state.endpoint.Parameters.find(p => p.Name === name && p.In === 'query');
        if (param && state.parameters[name]) {
          queryParams.append(name, state.parameters[name]);
        }
      });

      if (state.auth.type === 'apikey' && state.auth.apiKeyLocation === 'query') {
        queryParams.append('api_key', state.auth.apiKey);
      }

      if (queryParams.toString()) {
        url += '?' + queryParams.toString();
      }

      // Build headers
      const headers = { ...state.headers };

      // Add auth headers
      if (state.auth.type === 'bearer' && state.auth.token) {
        headers['Authorization'] = state.auth.token;
      } else if (state.auth.type === 'apikey' && state.auth.apiKey && state.auth.apiKeyLocation === 'header') {
        headers['X-API-Key'] = state.auth.apiKey;
      } else if (state.auth.type === 'oauth2' && state.auth.token) {
        headers['Authorization'] = state.auth.token;
      }

      // Add parameter headers
      Object.keys(state.parameters).forEach(name => {
        const param = state.endpoint.Parameters.find(p => p.Name === name && p.In === 'header');
        if (param && state.parameters[name]) {
          headers[name] = state.parameters[name];
        }
      });

      // Build request options
      const options = {
        method: state.endpoint.Method,
        headers: headers
      };

      // Add body for POST/PUT/PATCH
      if (['POST', 'PUT', 'PATCH'].includes(state.endpoint.Method) && state.body) {
        options.body = state.body;
        if (!headers['Content-Type']) {
          headers['Content-Type'] = 'application/json';
        }
      }

      // Send request
      const response = await fetch(url, options);

      // Display status
      const statusClass = response.ok ? 'status-2xx' : (response.status >= 400 && response.status < 500 ? 'status-4xx' : 'status-5xx');
      responseStatus.innerHTML = `<span class="response-code response-code-${response.status}">${response.status} ${response.statusText}</span>`;

      // Display headers
      let headersText = '';
      response.headers.forEach((value, key) => {
        headersText += `${key}: ${value}\n`;
      });
      responseHeaders.textContent = headersText;

      // Display body
      const contentType = response.headers.get('Content-Type') || '';
      const text = await response.text();

      if (contentType.includes('application/json')) {
        try {
          const json = JSON.parse(text);
          responseBody.textContent = JSON.stringify(json, null, 2);
        } catch {
          responseBody.textContent = text;
        }
      } else {
        responseBody.textContent = text;
      }

    } catch (error) {
      responseStatus.innerHTML = '<span class="response-code" style="background: var(--status-5xx)">Error</span>';
      responseBody.textContent = error.message;
    }
  }

  // Setup export buttons
  function setupExportButtons() {
    const curlBtn = document.getElementById('export-curl-btn');
    const restcliBtn = document.getElementById('export-restcli-btn');

    if (curlBtn) {
      curlBtn.addEventListener('click', () => {
        if (window.ExportHandler) {
          window.ExportHandler.exportCurl(state);
        }
      });
    }

    if (restcliBtn) {
      restcliBtn.addEventListener('click', () => {
        if (window.ExportHandler) {
          window.ExportHandler.exportRestcli(state);
        }
      });
    }
  }

  // Setup close button
  function setupCloseButton() {
    const btn = document.getElementById('api-tester-close');
    if (!btn) return;

    btn.addEventListener('click', () => {
      const tester = document.getElementById('openapi-sidebar-right');
      if (tester) {
        tester.classList.remove('active');
      }
    });
  }

  // Build cURL command (simplified)
  function buildCurlCommand() {
    if (!state.endpoint) return '';

    let url = state.server + state.endpoint.Path;
    let cmd = `curl -X ${state.endpoint.Method} '${url}'`;

    // Add headers
    Object.keys(state.headers).forEach(name => {
      cmd += ` \\\n  -H '${name}: ${state.headers[name]}'`;
    });

    // Add auth
    if (state.auth.type === 'bearer' && state.auth.token) {
      cmd += ` \\\n  -H 'Authorization: ${state.auth.token}'`;
    }

    // Add body
    if (state.body) {
      cmd += ` \\\n  -d '${state.body}'`;
    }

    return cmd;
  }

  // Build restcli config
  function buildRestcliConfig() {
    if (!state.endpoint) return null;

    return {
      name: `${state.endpoint.Method} ${state.endpoint.Path}`,
      method: state.endpoint.Method,
      url: state.server + state.endpoint.Path,
      headers: state.headers,
      body: state.body ? JSON.parse(state.body) : null
    };
  }

  // Helper: Convert schema to example
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

    switch (schema.Type) {
      case 'string': return schema.Default || 'string';
      case 'number': return schema.Default || 0;
      case 'integer': return schema.Default || 0;
      case 'boolean': return schema.Default || false;
      default: return null;
    }
  }

  // Export API
  window.APITester = {
    loadEndpoint,
    state
  };
})();
