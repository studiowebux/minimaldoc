---
title: Schemas
description: Browse and explore API data models
tags:
  - openapi
  - schemas
---

# Schemas

The Schemas tab provides a dedicated view for exploring API data models.

## Accessing Schemas

Click the "Schemas" tab in API documentation:

```
[Path] [Tag] [Flat] [Schemas]
```

## Schema List

All schemas from `components/schemas` displayed alphabetically:

```
Schemas (12)
├── CreateUserRequest
├── Error
├── Order
├── OrderItem
├── PaginatedResponse
├── Product
├── UpdateUserRequest
├── User
└── ...
```

## Schema Details

Click a schema to expand:

```
User (object)
│
├── id
│   type: integer
│   format: int64
│   readOnly: true
│
├── name (required)
│   type: string
│   minLength: 1
│   maxLength: 100
│
├── email (required)
│   type: string
│   format: email
│
├── role
│   type: string
│   enum: [user, admin, moderator]
│   default: user
│
├── created_at
│   type: string
│   format: date-time
│   readOnly: true
│
└── profile
    $ref: #/components/schemas/UserProfile
```

## Type Indicators

| Type      | Display                     |
| --------- | --------------------------- |
| `object`  | Object icon, expandable     |
| `array`   | Array icon, shows item type |
| `string`  | String icon                 |
| `integer` | Number icon                 |
| `boolean` | Boolean icon                |
| `$ref`    | Reference link              |

## Property Information

### Constraints

```
age
  type: integer
  minimum: 0
  maximum: 150
  description: User's age in years
```

### Formats

```
email
  type: string
  format: email

created_at
  type: string
  format: date-time

id
  type: string
  format: uuid
```

### Enums

```
status
  type: string
  enum:
    - pending
    - active
    - suspended
  default: pending
```

### Required Fields

Required fields marked with indicator:

```
├── name (required)
├── email (required)
├── age
```

## Nested Schemas

Expand nested objects inline:

```
User (object)
└── profile
    └── UserProfile (object)
        ├── bio
        │   type: string
        ├── avatar_url
        │   type: string
        │   format: uri
        └── social_links
            └── array of SocialLink
                ├── platform
                └── url
```

## Arrays

Array schemas show item type:

```
tags
  type: array
  items:
    type: string
  minItems: 1
  maxItems: 10
  uniqueItems: true
```

## References

Click `$ref` to navigate:

```
order_items
  type: array
  items:
    $ref: #/components/schemas/OrderItem  [→]
```

## JSON Examples

Auto-generated examples from schema:

```json
// User
{
  "id": 12345,
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "created_at": "2025-01-15T10:30:00Z",
  "profile": {
    "bio": "Software developer",
    "avatar_url": "https://example.com/avatar.jpg",
    "social_links": [
      {
        "platform": "twitter",
        "url": "https://twitter.com/johndoe"
      }
    ]
  }
}
```

## Search

Filter schemas by name:

```
Search: [user    ]

Results:
  CreateUserRequest
  UpdateUserRequest
  User
  UserProfile
```

## Composition

### allOf (Inheritance)

```yaml
AdminUser:
  allOf:
    - $ref: "#/components/schemas/User"
    - type: object
      properties:
        permissions:
          type: array
          items:
            type: string
```

Displayed as:

```
AdminUser (object)
  extends: User
  └── permissions
      type: array
```

### oneOf (Union)

```yaml
Response:
  oneOf:
    - $ref: "#/components/schemas/SuccessResponse"
    - $ref: "#/components/schemas/ErrorResponse"
```

Displayed as:

```
Response (oneOf)
  ├── SuccessResponse
  └── ErrorResponse
```

### anyOf

```yaml
Value:
  anyOf:
    - type: string
    - type: integer
    - type: boolean
```

Displayed as:

```
Value (anyOf)
  ├── string
  ├── integer
  └── boolean
```

## Discriminator

```yaml
Pet:
  discriminator:
    propertyName: petType
    mapping:
      dog: "#/components/schemas/Dog"
      cat: "#/components/schemas/Cat"
```

Displayed with discriminator indicator:

```
Pet (object)
  discriminator: petType
  ├── Dog (petType: "dog")
  └── Cat (petType: "cat")
```

