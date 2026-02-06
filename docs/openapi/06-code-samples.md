---
title: Code Samples
description: Auto-generated code examples for API endpoints
tags:
  - openapi
  - code-samples
---

# Code Samples

MinimalDoc generates code samples for every endpoint.

## Supported Languages

| Language   | Library    |
| ---------- | ---------- |
| cURL       | Native     |
| JavaScript | Fetch API  |
| Python     | requests   |
| Go         | net/http   |
| Swift      | URLSession |

## Configuration

Enable code samples:

```yaml
openapi:
  enable_code_samples: true
```

Disable:

```yaml
openapi:
  enable_code_samples: false
```

## Sample Output

### cURL

```bash
curl -X POST 'https://api.example.com/v1/users' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### JavaScript

```javascript
const response = await fetch("https://api.example.com/v1/users", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    Authorization: "Bearer YOUR_TOKEN",
  },
  body: JSON.stringify({
    name: "John Doe",
    email: "john@example.com",
  }),
});

const data = await response.json();
```

### Python

```python
import requests

response = requests.post(
    'https://api.example.com/v1/users',
    headers={
        'Content-Type': 'application/json',
        'Authorization': 'Bearer YOUR_TOKEN'
    },
    json={
        'name': 'John Doe',
        'email': 'john@example.com'
    }
)

data = response.json()
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func main() {
    body := map[string]string{
        "name":  "John Doe",
        "email": "john@example.com",
    }
    jsonBody, _ := json.Marshal(body)

    req, _ := http.NewRequest("POST",
        "https://api.example.com/v1/users",
        bytes.NewBuffer(jsonBody))

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer YOUR_TOKEN")

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()
}
```

### Swift

```swift
import Foundation

let url = URL(string: "https://api.example.com/v1/users")!
var request = URLRequest(url: url)
request.httpMethod = "POST"
request.setValue("application/json", forHTTPHeaderField: "Content-Type")
request.setValue("Bearer YOUR_TOKEN", forHTTPHeaderField: "Authorization")

let body: [String: Any] = [
    "name": "John Doe",
    "email": "john@example.com"
]
request.httpBody = try? JSONSerialization.data(withJSONObject: body)

let task = URLSession.shared.dataTask(with: request) { data, response, error in
    // Handle response
}
task.resume()
```

## Features

### Path Parameters

Substituted in URL:

```bash
# Spec: GET /users/{id}
curl -X GET 'https://api.example.com/v1/users/123'
```

### Query Parameters

Appended to URL:

```bash
# Spec: GET /users?limit=10&offset=0
curl -X GET 'https://api.example.com/v1/users?limit=10&offset=0'
```

### Request Body

Generated from schema:

```yaml
# Spec schema
User:
  type: object
  properties:
    name:
      type: string
    email:
      type: string
      format: email
```

Produces:

```json
{
  "name": "string",
  "email": "user@example.com"
}
```

### Authentication

Included based on security scheme:

```bash
# Bearer
-H 'Authorization: Bearer YOUR_TOKEN'

# API Key
-H 'X-API-Key: YOUR_API_KEY'

# Basic
-H 'Authorization: Basic base64(user:pass)'
```

## Example Generation

### From Schema

```yaml
components:
  schemas:
    CreateUserRequest:
      type: object
      required:
        - name
        - email
      properties:
        name:
          type: string
          example: "John Doe"
        email:
          type: string
          format: email
          example: "john@example.com"
        age:
          type: integer
          minimum: 0
          maximum: 150
```

Generated example:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 25
}
```

### Precedence

1. `example` field (exact value)
2. `examples` field (first example)
3. Generated from constraints

## Copy to Clipboard

Each sample has a copy button:

```
[curl] [JavaScript] [Python] [Go] [Swift]

┌─────────────────────────────────────┐
│ curl -X GET 'https://...'           │
│                              [Copy] │
└─────────────────────────────────────┘
```

## Customization

### Server URL

Samples use selected server:

```yaml
servers:
  - url: https://api.example.com/v1
    description: Production
  - url: http://localhost:3000/v1
    description: Local
```

### Content Type

Based on request body:

```yaml
requestBody:
  content:
    application/json: # Uses JSON
      schema:
        $ref: "#/components/schemas/User"
    multipart/form-data: # Uses form data
      schema:
        type: object
        properties:
          file:
            type: string
            format: binary
```
