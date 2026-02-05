---
title: Authentication
description: Configure API authentication for testing
tags:
  - openapi
  - authentication
  - security
---

# Authentication

MinimalDoc supports OpenAPI security schemes for API testing.

## Supported Methods

| Method | Type | Description |
|--------|------|-------------|
| Bearer Token | `http: bearer` | JWT or opaque tokens |
| API Key | `apiKey` | Header or query parameter |
| OAuth 2.0 | `oauth2` | Authorization Code flow |
| Basic Auth | `http: basic` | Username/password |

## Bearer Token

### Spec Definition

```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - bearerAuth: []
```

### Usage in UI

```
Authentication: Bearer Token
Token: [eyJhbGciOiJIUzI1NiIs...]
```

Generated header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

## API Key

### Header-Based

```yaml
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key

security:
  - apiKeyAuth: []
```

Usage:

```
Authentication: API Key (Header)
X-API-Key: [sk_live_abc123...]
```

### Query Parameter

```yaml
components:
  securitySchemes:
    apiKeyQuery:
      type: apiKey
      in: query
      name: api_key
```

Usage:

```
Authentication: API Key (Query)
api_key: [sk_live_abc123...]
```

Added to URL:

```
GET /users?api_key=sk_live_abc123
```

## OAuth 2.0

### Authorization Code Flow

```yaml
components:
  securitySchemes:
    oauth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://auth.example.com/authorize
          tokenUrl: https://auth.example.com/token
          scopes:
            read: Read access
            write: Write access
```

### Usage in UI

1. Click "Authorize"
2. Redirect to auth provider
3. Grant permissions
4. Token stored automatically

```
Authentication: OAuth 2.0
Status: Authorized
Scopes: read, write
[Revoke]
```

### Token Refresh

OAuth tokens refresh automatically when expired.

## Basic Auth

```yaml
components:
  securitySchemes:
    basicAuth:
      type: http
      scheme: basic
```

Usage:

```
Authentication: Basic
Username: [admin]
Password: [••••••••]
```

Generated header:

```
Authorization: Basic YWRtaW46cGFzc3dvcmQ=
```

## Per-Endpoint Security

Override global security:

```yaml
paths:
  /public:
    get:
      security: []  # No auth required
      summary: Public endpoint

  /admin:
    get:
      security:
        - bearerAuth: []
        - apiKeyAuth: []  # Either method works
      summary: Admin endpoint
```

## Multiple Security Options

```yaml
security:
  - bearerAuth: []
  - apiKeyAuth: []
```

UI shows dropdown:

```
Authentication: [Bearer Token ▼]
  Bearer Token
  API Key
```

## Combined Security

Require multiple auth methods:

```yaml
security:
  - bearerAuth: []
    apiKeyAuth: []  # Both required
```

## Token Persistence

### Session Storage

Tokens stored in browser session:
- Cleared on tab close
- Isolated per origin

### Local Storage

Enable persistent tokens:

```
[x] Remember authentication
```

Tokens persist across sessions.

## Security Best Practices

### In Documentation

```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: |
        Obtain token from POST /auth/login.
        Include in Authorization header.
        Tokens expire after 1 hour.
```

### For Users

- Use test/sandbox credentials
- Never commit real tokens
- Revoke tokens after testing
- Use short-lived tokens

## Configuration

### Disable Testing

Hides auth configuration:

```yaml
openapi:
  enable_testing: false
```

### Documentation Only

Show security requirements without testing:

```yaml
openapi:
  enable_testing: false
  enable_code_samples: true  # Shows auth in samples
```

## Troubleshooting

### 401 Unauthorized

- Check token format
- Verify token not expired
- Confirm correct auth method

### CORS Issues

- Ensure server allows credentials
- Check preflight response headers

### OAuth Redirect

- Verify callback URL allowed
- Check scopes match requirements
