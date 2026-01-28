// Code Samples Generator for OpenAPI Explorer
(function() {
  'use strict';

  // Generate example value from schema
  function generateExampleBody(schema) {
    if (!schema) return null;

    if (schema.Example !== undefined && schema.Example !== null) {
      return schema.Example;
    }

    if (schema.Type === 'object' && schema.Properties) {
      const obj = {};
      Object.keys(schema.Properties).forEach(function(key) {
        obj[key] = generateExampleBody(schema.Properties[key]);
      });
      return obj;
    }

    if (schema.Type === 'array' && schema.Items) {
      return [generateExampleBody(schema.Items)];
    }

    // Return defaults based on type
    switch (schema.Type) {
      case 'string':
        if (schema.Format === 'date') return '2024-01-15';
        if (schema.Format === 'date-time') return '2024-01-15T10:30:00Z';
        if (schema.Format === 'email') return 'user@example.com';
        if (schema.Format === 'uuid') return '550e8400-e29b-41d4-a716-446655440000';
        if (schema.Format === 'uri') return 'https://example.com';
        if (schema.Enum && schema.Enum.length > 0) return schema.Enum[0];
        return schema.Default || 'string';
      case 'number':
      case 'integer':
        return schema.Default !== undefined ? schema.Default : 0;
      case 'boolean':
        return schema.Default !== undefined ? schema.Default : false;
      default:
        return null;
    }
  }

  // Get request body example from endpoint
  function getRequestBodyExample(endpoint) {
    if (!endpoint.RequestBody || !endpoint.RequestBody.Content) {
      return null;
    }

    // Prefer application/json
    const jsonContent = endpoint.RequestBody.Content['application/json'];
    if (jsonContent) {
      if (jsonContent.Example) return jsonContent.Example;
      if (jsonContent.Schema) return generateExampleBody(jsonContent.Schema);
    }

    // Fall back to first content type
    const firstKey = Object.keys(endpoint.RequestBody.Content)[0];
    if (firstKey) {
      const content = endpoint.RequestBody.Content[firstKey];
      if (content.Example) return content.Example;
      if (content.Schema) return generateExampleBody(content.Schema);
    }

    return null;
  }

  // Build URL with path parameters replaced
  function buildURL(server, path, parameters) {
    let url = server + path;

    // Replace path parameters
    if (parameters) {
      parameters.forEach(function(param) {
        if (param.In === 'path') {
          const placeholder = '{' + param.Name + '}';
          const value = param.Example || param.Schema && param.Schema.Example || '<' + param.Name + '>';
          url = url.replace(placeholder, value);
        }
      });
    }

    return url;
  }

  // Get query parameters as string
  function getQueryString(parameters) {
    if (!parameters) return '';

    const queryParams = parameters.filter(function(p) {
      return p.In === 'query';
    });

    if (queryParams.length === 0) return '';

    const parts = queryParams.map(function(p) {
      const value = p.Example || p.Schema && p.Schema.Example || '<' + p.Name + '>';
      return encodeURIComponent(p.Name) + '=' + encodeURIComponent(value);
    });

    return '?' + parts.join('&');
  }

  // Get security headers
  function getSecurityHeaders(endpoint, securitySchemes) {
    const headers = [];

    if (!endpoint.Security || endpoint.Security.length === 0) {
      return headers;
    }

    endpoint.Security.forEach(function(req) {
      Object.keys(req).forEach(function(schemeName) {
        const scheme = securitySchemes && securitySchemes[schemeName];
        if (!scheme) return;

        if (scheme.Type === 'http' && scheme.Scheme === 'bearer') {
          headers.push({ name: 'Authorization', value: 'Bearer <token>' });
        } else if (scheme.Type === 'http' && scheme.Scheme === 'basic') {
          headers.push({ name: 'Authorization', value: 'Basic <credentials>' });
        } else if (scheme.Type === 'apiKey' && scheme.In === 'header') {
          headers.push({ name: scheme.Name, value: '<api-key>' });
        } else if (scheme.Type === 'oauth2') {
          headers.push({ name: 'Authorization', value: 'Bearer <access_token>' });
        } else if (scheme.Type === 'openIdConnect') {
          headers.push({ name: 'Authorization', value: 'Bearer <id_token>' });
        }
      });
    });

    return headers;
  }

  // Generate curl command
  function generateCurl(endpoint, server, securitySchemes) {
    const method = endpoint.Method.toUpperCase();
    const url = buildURL(server, endpoint.Path, endpoint.Parameters);
    const queryString = getQueryString(endpoint.Parameters);
    const fullUrl = url + queryString;
    const body = getRequestBodyExample(endpoint);
    const secHeaders = getSecurityHeaders(endpoint, securitySchemes);

    let lines = [];
    lines.push("curl -X " + method + " '" + fullUrl + "'");

    // Content-Type for body
    if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
      lines.push("  -H 'Content-Type: application/json'");
    }

    // Security headers
    secHeaders.forEach(function(h) {
      lines.push("  -H '" + h.name + ": " + h.value + "'");
    });

    // Body
    if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
      lines.push("  -d '" + JSON.stringify(body, null, 2).replace(/'/g, "'\\''") + "'");
    }

    return lines.join(' \\\n');
  }

  // Generate JavaScript fetch
  function generateFetch(endpoint, server, securitySchemes) {
    const method = endpoint.Method.toUpperCase();
    const url = buildURL(server, endpoint.Path, endpoint.Parameters);
    const queryString = getQueryString(endpoint.Parameters);
    const fullUrl = url + queryString;
    const body = getRequestBodyExample(endpoint);
    const secHeaders = getSecurityHeaders(endpoint, securitySchemes);

    let lines = [];
    lines.push("const response = await fetch('" + fullUrl + "', {");
    lines.push("  method: '" + method + "',");

    // Headers
    const headers = {};
    if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
      headers['Content-Type'] = 'application/json';
    }
    secHeaders.forEach(function(h) {
      headers[h.name] = h.value;
    });

    if (Object.keys(headers).length > 0) {
      lines.push("  headers: {");
      Object.keys(headers).forEach(function(key, i, arr) {
        const comma = i < arr.length - 1 ? ',' : '';
        lines.push("    '" + key + "': '" + headers[key] + "'" + comma);
      });
      lines.push("  },");
    }

    // Body
    if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
      lines.push("  body: JSON.stringify(" + JSON.stringify(body, null, 4).replace(/\n/g, '\n  ') + ")");
    }

    lines.push("});");
    lines.push("const data = await response.json();");

    return lines.join('\n');
  }

  // Generate Go http code
  function generateGoHttp(endpoint, server, securitySchemes) {
    const method = endpoint.Method.toUpperCase();
    const url = buildURL(server, endpoint.Path, endpoint.Parameters);
    const queryString = getQueryString(endpoint.Parameters);
    const fullUrl = url + queryString;
    const body = getRequestBodyExample(endpoint);
    const secHeaders = getSecurityHeaders(endpoint, securitySchemes);
    const hasBody = body && (method === 'POST' || method === 'PUT' || method === 'PATCH');

    let lines = [];
    lines.push('package main');
    lines.push('');
    lines.push('import (');
    if (hasBody) {
      lines.push('\t"bytes"');
      lines.push('\t"encoding/json"');
    }
    lines.push('\t"net/http"');
    lines.push(')');
    lines.push('');
    lines.push('func main() {');

    if (hasBody) {
      lines.push('\tbody, _ := json.Marshal(map[string]interface{}{');
      Object.keys(body).forEach(function(key) {
        const val = body[key];
        if (typeof val === 'string') {
          lines.push('\t\t"' + key + '": "' + val + '",');
        } else if (typeof val === 'number' || typeof val === 'boolean') {
          lines.push('\t\t"' + key + '": ' + val + ',');
        } else {
          lines.push('\t\t"' + key + '": ' + JSON.stringify(val) + ',');
        }
      });
      lines.push('\t})');
      lines.push('\treq, _ := http.NewRequest("' + method + '", "' + fullUrl + '", bytes.NewBuffer(body))');
    } else {
      lines.push('\treq, _ := http.NewRequest("' + method + '", "' + fullUrl + '", nil)');
    }

    if (hasBody) {
      lines.push('\treq.Header.Set("Content-Type", "application/json")');
    }
    secHeaders.forEach(function(h) {
      lines.push('\treq.Header.Set("' + h.name + '", "' + h.value + '")');
    });

    lines.push('');
    lines.push('\tclient := &http.Client{}');
    lines.push('\tresp, _ := client.Do(req)');
    lines.push('\tdefer resp.Body.Close()');
    lines.push('}');

    return lines.join('\n');
  }

  // Generate Python requests code
  function generatePython(endpoint, server, securitySchemes) {
    const method = endpoint.Method.toLowerCase();
    const url = buildURL(server, endpoint.Path, endpoint.Parameters);
    const queryString = getQueryString(endpoint.Parameters);
    const fullUrl = url + queryString;
    const body = getRequestBodyExample(endpoint);
    const secHeaders = getSecurityHeaders(endpoint, securitySchemes);
    const hasBody = body && (method === 'post' || method === 'put' || method === 'patch');

    let lines = [];
    lines.push('import requests');
    lines.push('');

    lines.push('response = requests.' + method + '(');
    lines.push("    '" + fullUrl + "',");

    // Headers
    const headers = {};
    if (hasBody) {
      headers['Content-Type'] = 'application/json';
    }
    secHeaders.forEach(function(h) {
      headers[h.name] = h.value;
    });

    if (Object.keys(headers).length > 0) {
      lines.push('    headers={');
      Object.keys(headers).forEach(function(key) {
        lines.push("        '" + key + "': '" + headers[key] + "',");
      });
      lines.push('    },');
    }

    if (hasBody) {
      lines.push('    json=' + JSON.stringify(body, null, 4).replace(/\n/g, '\n    '));
    }

    lines.push(')');
    lines.push('data = response.json()');

    return lines.join('\n');
  }

  // Generate Swift URLSession code
  function generateSwift(endpoint, server, securitySchemes) {
    const method = endpoint.Method.toUpperCase();
    const url = buildURL(server, endpoint.Path, endpoint.Parameters);
    const queryString = getQueryString(endpoint.Parameters);
    const fullUrl = url + queryString;
    const body = getRequestBodyExample(endpoint);
    const secHeaders = getSecurityHeaders(endpoint, securitySchemes);
    const hasBody = body && (method === 'POST' || method === 'PUT' || method === 'PATCH');

    let lines = [];
    lines.push('import Foundation');
    lines.push('');
    lines.push('let url = URL(string: "' + fullUrl + '")!');
    lines.push('var request = URLRequest(url: url)');
    lines.push('request.httpMethod = "' + method + '"');

    if (hasBody) {
      lines.push('request.setValue("application/json", forHTTPHeaderField: "Content-Type")');
    }

    secHeaders.forEach(function(h) {
      lines.push('request.setValue("' + h.value + '", forHTTPHeaderField: "' + h.name + '")');
    });

    if (hasBody) {
      lines.push('');
      lines.push('let body: [String: Any] = ' + JSON.stringify(body, null, 4).replace(/"/g, '\\"').replace(/\n/g, '\n'));
      lines.push('request.httpBody = try? JSONSerialization.data(withJSONObject: body)');
    }

    lines.push('');
    lines.push('let task = URLSession.shared.dataTask(with: request) { data, response, error in');
    lines.push('    guard let data = data else { return }');
    lines.push('    // Handle response');
    lines.push('}');
    lines.push('task.resume()');

    return lines.join('\n');
  }

  // Create code samples section HTML
  function createCodeSamplesSection(endpoint, server, securitySchemes) {
    const languages = [
      { id: 'curl', name: 'curl', generator: generateCurl },
      { id: 'fetch', name: 'JavaScript', generator: generateFetch },
      { id: 'swift', name: 'Swift', generator: generateSwift },
      { id: 'go', name: 'Go', generator: generateGoHttp },
      { id: 'python', name: 'Python', generator: generatePython }
    ];

    let html = '<div class="endpoint-section code-samples-section">';
    html += '<h3 class="code-samples-title">Code Samples</h3>';
    html += '<div class="code-samples-content">';

    // Tabs
    html += '<div class="code-sample-tabs">';
    languages.forEach(function(lang, i) {
      const activeClass = i === 0 ? ' active' : '';
      html += '<button class="code-sample-tab' + activeClass + '" data-lang="' + lang.id + '">' + lang.name + '</button>';
    });
    html += '</div>';

    // Panels
    html += '<div class="code-sample-panels">';
    languages.forEach(function(lang, i) {
      const activeClass = i === 0 ? ' active' : '';
      const code = lang.generator(endpoint, server, securitySchemes);
      html += '<div class="code-sample-panel' + activeClass + '" data-lang="' + lang.id + '">';
      html += '<div class="code-sample-header">';
      html += '<button class="btn-copy" data-code="' + escapeAttr(code) + '">Copy</button>';
      html += '</div>';
      html += '<pre class="code-sample-code"><code>' + escapeHtml(code) + '</code></pre>';
      html += '</div>';
    });
    html += '</div>';

    html += '</div>';
    html += '</div>';

    return html;
  }

  // Setup tab switching for code samples
  function setupCodeSampleTabs(container) {
    container.addEventListener('click', function(e) {
      const tab = e.target.closest('.code-sample-tab');
      if (!tab) return;

      const section = tab.closest('.code-samples-section');
      if (!section) return;

      const lang = tab.dataset.lang;

      // Update tabs
      section.querySelectorAll('.code-sample-tab').forEach(function(t) {
        t.classList.toggle('active', t.dataset.lang === lang);
      });

      // Update panels
      section.querySelectorAll('.code-sample-panel').forEach(function(p) {
        p.classList.toggle('active', p.dataset.lang === lang);
      });
    });
  }

  // Helper: escape HTML
  function escapeHtml(text) {
    if (!text) return '';
    var div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
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

  // Export
  window.CodeSamples = {
    generateCurl: generateCurl,
    generateFetch: generateFetch,
    generateSwift: generateSwift,
    generateGoHttp: generateGoHttp,
    generatePython: generatePython,
    generateExampleBody: generateExampleBody,
    createCodeSamplesSection: createCodeSamplesSection,
    setupCodeSampleTabs: setupCodeSampleTabs
  };
})();
