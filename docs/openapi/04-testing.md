---
title: API Testing
description: Interactive API testing with Try It interface
tags:
  - openapi
  - testing
---

# API Testing

MinimalDoc includes an interactive API tester for trying endpoints.

## Try It Interface

Each endpoint has a "Try It" button:

```
GET /users/{id}

[Try It]

Parameters:
  id: [123        ]

Headers:
  Authorization: [Bearer token...]

[Send Request]
```

## Features

| Feature | Description |
|---------|-------------|
| Parameter Input | Fill path, query, header params |
| Request Body | Edit JSON/form body |
| Authentication | Configure auth headers |
| Response Viewer | See status, headers, body |
| History | Recent requests saved |

## Using Try It

### 1. Select Endpoint

Navigate to an endpoint in the API docs.

### 2. Configure Parameters

Fill required parameters:

```
Path Parameters:
  id: [123]

Query Parameters:
  limit: [10]
  offset: [0]
```

### 3. Set Authentication

If endpoint requires auth:

```
Authentication:
  [Bearer Token ▼]
  Token: [eyJhbG...]
```

### 4. Edit Request Body

For POST/PUT/PATCH:

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

### 5. Send Request

Click "Send Request" to execute.

### 6. View Response

```
Status: 200 OK
Time: 145ms

Headers:
  Content-Type: application/json
  X-Request-Id: abc123

Body:
{
  "id": 123,
  "name": "John Doe",
  "email": "john@example.com"
}
```

## Server Selection

If spec defines multiple servers:

```yaml
servers:
  - url: https://api.example.com/v1
    description: Production
  - url: https://staging-api.example.com/v1
    description: Staging
  - url: http://localhost:3000/v1
    description: Local
```

Select in the UI:

```
Server: [Production ▼]
```

## Request Body

### JSON Editor

Syntax-highlighted JSON editor:

```json
{
  "name": "New User",
  "email": "user@example.com",
  "role": "admin"
}
```

### Schema Hints

Editor shows expected schema:

```
Expected: User object
  name (string, required)
  email (string, required)
  role (string, enum: user, admin)
```

### Example Generation

Click "Generate Example" for sample body:

```json
{
  "name": "string",
  "email": "user@example.com",
  "role": "user"
}
```

## Response Handling

### Success

```
Status: 201 Created

{
  "id": 456,
  "name": "New User",
  ...
}
```

### Error

```
Status: 400 Bad Request

{
  "error": "validation_error",
  "message": "Invalid email format",
  "field": "email"
}
```

### Headers

View response headers:

```
Content-Type: application/json
X-RateLimit-Remaining: 99
X-Request-Id: req_abc123
```

## CORS Considerations

Browser-based testing requires CORS headers:

```yaml
# Your API must return:
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE
Access-Control-Allow-Headers: Authorization, Content-Type
```

If CORS blocked:
- Test against local server
- Use browser extension
- Test via cURL export

## Configuration

### Disable Testing

For read-only docs:

```yaml
openapi:
  enable_testing: false
```

### Enable Export Only

Show cURL export without Try It:

```yaml
openapi:
  enable_testing: false
  enable_export: true
```

## Export to cURL

Generate cURL command for any request:

```bash
curl -X GET 'https://api.example.com/v1/users/123' \
  -H 'Authorization: Bearer eyJhbG...' \
  -H 'Accept: application/json'
```

Copy and run in terminal.

## Export to restcli

Generate restcli format:

```http
GET https://api.example.com/v1/users/123
Authorization: Bearer eyJhbG...
Accept: application/json
```

Compatible with VS Code REST Client extension.

## Request History

Recent requests are saved in browser storage:

```
Recent Requests:
  GET /users/123 (2 min ago)
  POST /users (5 min ago)
  GET /orders (10 min ago)
```

Click to replay with same parameters.
