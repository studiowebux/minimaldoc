---
title: OpenAPI Views
description: Organize API endpoints by path, tag, or flat list
tags:
  - openapi
  - navigation
---

# OpenAPI Views

MinimalDoc provides three ways to organize API endpoints.

## View Modes

### Path View

Group endpoints by URL path hierarchy:

```
/users
  GET    List users
  POST   Create user
/users/{id}
  GET    Get user
  PUT    Update user
  DELETE Delete user
/orders
  GET    List orders
  POST   Create order
```

Set as default:

```yaml
openapi:
  default_view: "path"
```

### Tag View

Group endpoints by OpenAPI tags:

```
Users
  GET    /users          List users
  POST   /users          Create user
  GET    /users/{id}     Get user
  PUT    /users/{id}     Update user
  DELETE /users/{id}     Delete user

Orders
  GET    /orders         List orders
  POST   /orders         Create order
  GET    /orders/{id}    Get order
```

Set as default:

```yaml
openapi:
  default_view: "tag"
```

Requires tags in your spec:

```yaml
paths:
  /users:
    get:
      tags:
        - Users
      summary: List users
```

### Flat View

All endpoints in a single list:

```
GET    /users          List users
POST   /users          Create user
GET    /users/{id}     Get user
PUT    /users/{id}     Update user
DELETE /users/{id}     Delete user
GET    /orders         List orders
POST   /orders         Create order
```

Set as default:

```yaml
openapi:
  default_view: "flat"
```

## Switching Views

Users can switch views via tabs in the UI:

```
[Path] [Tag] [Flat] [Schemas]
```

The default view is shown initially.

## Single Endpoint View

Click any endpoint to see focused documentation:

```
GET /users/{id}

Description: Get a specific user by ID

Parameters:
  id (path, required): User ID

Responses:
  200: User found
  404: User not found
```

Includes:
- Full description
- Parameters table
- Request body schema
- Response schemas
- Code samples
- Try It interface

## Navigation

### Collapsible Groups

Path and tag groups collapse/expand:

```
▼ /users
    GET    List users
    POST   Create user
▶ /orders (collapsed)
```

Large APIs start collapsed for performance.

### Search

Filter endpoints by typing:

```
Search: [user          ]

Results:
  GET  /users         List users
  POST /users         Create user
  GET  /users/{id}    Get user
```

### Method Colors

HTTP methods are color-coded:

| Method | Color |
|--------|-------|
| GET | Green |
| POST | Blue |
| PUT | Orange |
| PATCH | Yellow |
| DELETE | Red |

## Schemas Tab

Dedicated view for data models:

```
[Path] [Tag] [Flat] [Schemas]

Schemas:
  User
  Order
  Product
  Error
```

Click to expand schema details:

```
User (object)
├── id (integer, required)
├── name (string, required)
├── email (string, format: email)
└── created_at (string, format: date-time)
```

