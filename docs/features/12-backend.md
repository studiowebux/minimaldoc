---
title: Backend Server (Optional)
description: Self-hosted backend for analytics, feedback, and newsletter features
menu_order: 12
---

# Backend Server

MinimalDoc includes an optional backend server that adds dynamic features to your static documentation site. The backend is 100% opt-in - your static sites work without it.

## Overview

The backend provides:

| Feature | Description |
|---------|-------------|
| Analytics | Cookie-free, privacy-first page view tracking |
| Feedback | Page rating widget with optional comments |
| Newsletter | Email subscription with verification |
| Admin Portal | Web-based dashboard for all features |

## Architecture

```
┌─────────────────┐     ┌──────────────────┐
│  Static Site    │────▶│ minimaldoc-server│
│  (HTML/JS)      │     │  (Go binary)     │
└─────────────────┘     └────────┬─────────┘
                                 │
                        ┌────────▼─────────┐
                        │    Database      │
                        │ (SQLite/Postgres)│
                        └──────────────────┘
```

The CLI generator remains unchanged. The backend is a separate binary that your static site calls via JavaScript.

## Quick Start

### 1. Build the Server

```bash
git clone https://github.com/studiowebux/minimaldoc
cd minimaldoc
make build-server
```

### 2. Start with SQLite

```bash
export DB_DRIVER=sqlite
export DB_URL=./minimaldoc.db
export AUTH_JWT_SECRET=your-secret-key-at-least-32-chars
export SERVER_PORT=8090

./minimaldoc-server
```

### 3. Bootstrap

Create your first site and admin user:

```bash
curl -X POST http://localhost:8090/api/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "your-password",
    "site_name": "My Docs",
    "domain": "docs.example.com"
  }'
```

Save the returned `site_id` and `api_key`.

### 4. Enable in Static Site

Add to your `config.yaml`:

```yaml
analytics:
  enabled: true
  providers:
    - type: minimaldoc
      enabled: true
      config:
        endpoint: "http://localhost:8090"
        site_id: "your-site-id-from-bootstrap"
        features: "analytics,feedback,newsletter"
```

Rebuild your site:

```bash
minimaldoc build
```

## Features

### Analytics

Cookie-free page view tracking with:

- No personal data collection
- Hashed session IDs (non-reversible)
- Country detection from IP (IP discarded immediately)
- Device type detection
- Referrer tracking

Data available in admin dashboard:
- Total page views
- Unique visitors (estimated via session hashes)
- Top pages
- Traffic by device type

### Feedback Widget

Add a page rating widget to any page:

```html
<div data-minimaldoc-feedback data-path="/docs/getting-started"></div>
```

The widget renders:
- 5-star rating buttons
- Optional feedback text area
- Thank you message on submit

Include the optional CSS for styling:

```html
<link rel="stylesheet" href="http://your-server/minimaldoc.css">
```

### Newsletter

Add a signup form to any page:

```html
<form data-minimaldoc-newsletter>
  <input type="email" placeholder="Enter your email" required>
  <button type="submit">Subscribe</button>
</form>
```

Features:
- Email verification (double opt-in)
- Subscriber management in admin portal
- Unsubscribe support

### Email System

The backend includes a pluggable email system:

**Providers:**
- `mock` - Logs emails to console (for testing)
- `smtp` - Standard SMTP delivery

**Email Templates:**
- Verification email - Sent when user subscribes
- Welcome email - Sent after email is verified
- Unsubscribe confirmation - Sent when user unsubscribes

All emails include both HTML and plain text versions.

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_DRIVER` | Yes | - | `sqlite` or `postgres` |
| `DATABASE_URL` | Yes | - | Database connection string |
| `AUTH_JWT_SECRET` | Yes | - | Secret for JWT tokens (min 32 chars) |
| `SERVER_PORT` | No | `8080` | HTTP port |
| `SERVER_HOST` | No | `0.0.0.0` | Bind address |
| `CORS_ORIGINS` | No | `*` | Allowed CORS origins |

### Email Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `EMAIL_PROVIDER` | No | `mock` | `smtp` or `mock` |
| `SMTP_HOST` | If smtp | - | SMTP server hostname |
| `SMTP_PORT` | If smtp | `587` | SMTP server port |
| `SMTP_USER` | If smtp | - | SMTP username |
| `SMTP_PASS` | If smtp | - | SMTP password |
| `EMAIL_FROM_ADDRESS` | No | `noreply@example.com` | Sender email |
| `EMAIL_FROM_NAME` | No | `MinimalDoc` | Sender name |
| `EMAIL_BASE_URL` | No | `http://localhost:8080` | Base URL for verification links |

### PostgreSQL

For production, use PostgreSQL:

```bash
export DB_DRIVER=postgres
export DB_URL="postgres://user:pass@localhost/minimaldoc?sslmode=disable"
```

### Docker

```yaml
version: '3.8'
services:
  server:
    build:
      context: .
      dockerfile: docker/Dockerfile.server
    environment:
      DB_DRIVER: postgres
      DB_URL: postgres://postgres:postgres@db/minimaldoc
      AUTH_JWT_SECRET: your-secret-key-at-least-32-chars
    ports:
      - "8090:8080"
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: minimaldoc
      POSTGRES_PASSWORD: postgres
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

## API Reference

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/bootstrap` | Initial setup (first run only) |
| POST | `/api/analytics/track` | Track page view |
| POST | `/api/feedback` | Submit page feedback |
| POST | `/api/newsletter/subscribe` | Subscribe to newsletter |
| GET | `/api/newsletter/verify` | Verify email subscription |
| GET | `/minimaldoc.js` | Client JavaScript |
| GET | `/minimaldoc.css` | Widget CSS |

### Authenticated Endpoints

Require `Authorization: Bearer <token>` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/logout` | Logout |
| POST | `/api/auth/refresh` | Refresh token |
| GET | `/api/auth/me` | Current user info |
| GET | `/api/analytics/summary` | Analytics dashboard data |
| GET | `/api/analytics/pages` | Per-page statistics |
| GET | `/api/feedback/stats` | Feedback statistics |
| GET | `/api/feedback/list` | All feedback entries |
| GET | `/api/newsletter/subscribers` | Subscriber list |

## Client JavaScript

The client library (`minimaldoc.js`) can be used in two ways:

### Auto-initialization

Via script tag attributes:

```html
<script
  src="http://your-server/minimaldoc.js"
  data-endpoint="http://your-server"
  data-site-id="your-site-id"
  data-features="analytics,feedback,newsletter"
  defer>
</script>
```

### Manual Initialization

```html
<script src="http://your-server/minimaldoc.js"></script>
<script>
  MinimalDoc.init({
    endpoint: 'http://your-server',
    siteId: 'your-site-id',
    features: ['analytics', 'feedback', 'newsletter'],
    debug: true
  });
</script>
```

### SPA Support

The client automatically tracks page views for single-page applications by intercepting `history.pushState` and `popstate` events.

## Admin Portal

Access the admin UI at `http://your-server/admin`.

Features:
- Dashboard with analytics overview
- Feedback management
- Subscriber management
- Site settings

Login with the credentials from bootstrap or any user created via API.

## Privacy

The backend is designed for privacy:

- **No cookies** - Session tracking uses hashed identifiers
- **No PII** - Email addresses stored only for newsletter (with consent)
- **IP anonymization** - Country extracted then IP discarded
- **Self-hosted** - Full control over your data
- **GDPR-friendly** - No third-party data sharing

## Deployment

### Binary

Download from releases or build from source:

```bash
make build-server
./minimaldoc-server
```

### Systemd

```ini
[Unit]
Description=MinimalDoc Server
After=network.target

[Service]
Type=simple
User=minimaldoc
Environment=DB_DRIVER=sqlite
Environment=DB_URL=/var/lib/minimaldoc/data.db
Environment=AUTH_JWT_SECRET=your-secret
ExecStart=/usr/local/bin/minimaldoc-server
Restart=always

[Install]
WantedBy=multi-user.target
```

### Behind Nginx

```nginx
server {
    listen 443 ssl;
    server_name api.docs.example.com;

    location / {
        proxy_pass http://127.0.0.1:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
