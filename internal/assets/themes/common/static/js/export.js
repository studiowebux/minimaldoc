// Export Handler - Export to restcli and cURL formats
(function() {
  'use strict';

  // Export to cURL
  function exportCurl(testerState) {
    if (!testerState.endpoint) {
      alert('No endpoint selected');
      return;
    }

    // Build URL
    let url = testerState.server + testerState.endpoint.Path;

    // Replace path parameters
    Object.keys(testerState.parameters).forEach(name => {
      const param = testerState.endpoint.Parameters.find(p => p.Name === name && p.In === 'path');
      if (param) {
        url = url.replace(`{${name}}`, testerState.parameters[name]);
      }
    });

    // Add query parameters
    const queryParams = [];
    Object.keys(testerState.parameters).forEach(name => {
      const param = testerState.endpoint.Parameters.find(p => p.Name === name && p.In === 'query');
      if (param && testerState.parameters[name]) {
        queryParams.push(`${encodeURIComponent(name)}=${encodeURIComponent(testerState.parameters[name])}`);
      }
    });

    if (queryParams.length > 0) {
      url += '?' + queryParams.join('&');
    }

    // Build cURL command
    let cmd = `curl -X ${testerState.endpoint.Method} '${url}'`;

    // Add headers
    Object.keys(testerState.headers).forEach(name => {
      cmd += ` \\\n  -H '${name}: ${testerState.headers[name]}'`;
    });

    // Add auth
    if (testerState.auth.type === 'bearer' && testerState.auth.token) {
      cmd += ` \\\n  -H 'Authorization: ${testerState.auth.token}'`;
    } else if (testerState.auth.type === 'apikey' && testerState.auth.apiKey) {
      if (testerState.auth.apiKeyLocation === 'header') {
        cmd += ` \\\n  -H 'X-API-Key: ${testerState.auth.apiKey}'`;
      }
    } else if (testerState.auth.type === 'oauth2' && testerState.auth.token) {
      cmd += ` \\\n  -H 'Authorization: ${testerState.auth.token}'`;
    }

    // Add content-type and body
    if (['POST', 'PUT', 'PATCH'].includes(testerState.endpoint.Method) && testerState.body) {
      cmd += ` \\\n  -H 'Content-Type: application/json'`;
      cmd += ` \\\n  -d '${testerState.body.replace(/'/g, "'\\''")}'`;
    }

    // Copy to clipboard
    copyToClipboard(cmd);
    showNotification('cURL command copied to clipboard!');

    console.log('Exported cURL:', cmd);
  }

  // Export to restcli format
  function exportRestcli(testerState) {
    if (!testerState.endpoint) {
      alert('No endpoint selected');
      return;
    }

    // Build URL
    let url = testerState.server + testerState.endpoint.Path;

    // Build restcli request object
    const restcliRequest = {
      name: `${testerState.endpoint.Method} ${testerState.endpoint.Path}`,
      description: testerState.endpoint.Summary || testerState.endpoint.Description || '',
      method: testerState.endpoint.Method,
      url: url,
      headers: {},
      pathParams: {},
      queryParams: {},
      body: null
    };

    // Add headers
    Object.keys(testerState.headers).forEach(name => {
      restcliRequest.headers[name] = testerState.headers[name];
    });

    // Add auth headers
    if (testerState.auth.type === 'bearer' && testerState.auth.token) {
      restcliRequest.headers['Authorization'] = testerState.auth.token;
    } else if (testerState.auth.type === 'apikey' && testerState.auth.apiKey) {
      if (testerState.auth.apiKeyLocation === 'header') {
        restcliRequest.headers['X-API-Key'] = testerState.auth.apiKey;
      }
    } else if (testerState.auth.type === 'oauth2' && testerState.auth.token) {
      restcliRequest.headers['Authorization'] = testerState.auth.token;
    }

    // Add parameters
    Object.keys(testerState.parameters).forEach(name => {
      const param = testerState.endpoint.Parameters.find(p => p.Name === name);
      if (param) {
        if (param.In === 'path') {
          restcliRequest.pathParams[name] = testerState.parameters[name];
        } else if (param.In === 'query') {
          restcliRequest.queryParams[name] = testerState.parameters[name];
        } else if (param.In === 'header') {
          restcliRequest.headers[name] = testerState.parameters[name];
        }
      }
    });

    // Add body
    if (testerState.body) {
      try {
        restcliRequest.body = JSON.parse(testerState.body);
        restcliRequest.headers['Content-Type'] = 'application/json';
      } catch (e) {
        restcliRequest.body = testerState.body;
      }
    }

    // Convert to JSON
    const json = JSON.stringify(restcliRequest, null, 2);

    // Copy to clipboard
    copyToClipboard(json);
    showNotification('restcli configuration copied to clipboard!');

    console.log('Exported restcli:', json);
  }

  // Helper: Copy to clipboard
  function copyToClipboard(text) {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(text).catch(err => {
        console.error('Failed to copy:', err);
        fallbackCopy(text);
      });
    } else {
      fallbackCopy(text);
    }
  }

  // Fallback copy method
  function fallbackCopy(text) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
  }

  // Show notification
  function showNotification(message) {
    // Simple alert for now
    // In a real implementation, would use a toast notification
    const notification = document.createElement('div');
    notification.textContent = message;
    notification.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: #49cc90;
      color: white;
      padding: 1rem 1.5rem;
      border-radius: 4px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      z-index: 10000;
      font-weight: 600;
      animation: slideIn 0.3s ease-out;
    `;

    // Add animation
    const style = document.createElement('style');
    style.textContent = `
      @keyframes slideIn {
        from {
          transform: translateX(400px);
          opacity: 0;
        }
        to {
          transform: translateX(0);
          opacity: 1;
        }
      }
    `;
    if (!document.getElementById('export-notification-style')) {
      style.id = 'export-notification-style';
      document.head.appendChild(style);
    }

    document.body.appendChild(notification);

    setTimeout(() => {
      notification.style.animation = 'slideIn 0.3s ease-out reverse';
      setTimeout(() => {
        document.body.removeChild(notification);
      }, 300);
    }, 3000);
  }

  // Export API
  window.ExportHandler = {
    exportCurl,
    exportRestcli
  };
})();
