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
| Blog CMS | Markdown editor with live preview, RBAC, comment moderation |
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

### Blog CMS

Full-featured blog with markdown editor and comment system.

**Admin UI Features:**
- Post list with status filters (draft/published/archived)
- Markdown editor with toolbar (bold, italic, heading, link, code)
- Live preview via HTMX
- Auto-slug generation from title
- Publish/unpublish controls

**Public API:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/blog/posts` | List published posts |
| GET | `/api/blog/posts/:slug` | Get post with rendered HTML |
| GET | `/api/blog/posts/:slug/comments` | List approved comments |
| POST | `/api/blog/posts/:slug/comments` | Submit comment (pending moderation) |

**RBAC Roles:**

| Role | Create | Edit Own | Edit Any | Delete | Moderate |
|------|--------|----------|----------|--------|----------|
| admin | Yes | Yes | Yes | Yes | Yes |
| editor | Yes | Yes | Yes | No | Yes |
| author | Yes | Yes | No | No | No |
| viewer | No | No | No | No | No |

**Comment Moderation:**
- Comments require approval before displaying
- Moderation queue in admin panel
- Approve, reject, or mark as spam
- Moderator tracking for audit

**Scheduled Publishing:**

Schedule posts for future publication:

1. In the blog editor, use the "Schedule" section
2. Select date and time
3. Click "Schedule"

The server automatically publishes posts when the scheduled time arrives. No external cron needed - the scheduler runs inside the server process.

**Schedule via API:**

```bash
curl -X POST http://your-server/api/blog/posts/:id/schedule \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"scheduled_at": "2024-02-15T10:00:00Z"}'
```

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

### Storage Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STORAGE_PROVIDER` | No | `local` | `local` or `s3` |
| `STORAGE_LOCAL_PATH` | No | `./uploads` | Local upload directory |
| `STORAGE_S3_BUCKET` | If s3 | - | S3 bucket name |
| `STORAGE_S3_REGION` | If s3 | `us-east-1` | S3 region |
| `STORAGE_S3_ACCESS_KEY` | If s3 | - | AWS access key |
| `STORAGE_S3_SECRET_KEY` | If s3 | - | AWS secret key |
| `STORAGE_S3_ENDPOINT` | No | - | Custom S3 endpoint (MinIO/R2) |
| `STORAGE_S3_PUBLIC_URL` | No | - | CDN URL prefix |
| `STORAGE_MAX_FILE_SIZE` | No | `5242880` | Max upload size (bytes) |
| `STORAGE_ALLOWED_TYPES` | No | `image/jpeg,...` | Allowed MIME types |

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
| POST | `/api/analytics/event` | Track custom event |
| POST | `/api/analytics/duration` | Update page view duration |
| POST | `/api/feedback` | Submit page feedback |
| POST | `/api/newsletter/subscribe` | Subscribe to newsletter |
| GET | `/api/newsletter/verify` | Verify email subscription |
| GET | `/api/blog/posts` | List published blog posts |
| GET | `/api/blog/posts/:slug` | Get published post by slug |
| GET | `/api/blog/posts/:slug/comments` | List approved comments |
| POST | `/api/blog/posts/:slug/comments` | Submit comment |
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

### Blog Admin Endpoints

Require authentication and appropriate role.

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/blog/posts` | Any | List all posts |
| POST | `/api/blog/posts` | Author+ | Create post |
| GET | `/api/blog/posts/:id` | Any | Get post by ID |
| PUT | `/api/blog/posts/:id` | Author+ | Update post |
| DELETE | `/api/blog/posts/:id` | Admin | Delete post |
| POST | `/api/blog/posts/:id/publish` | Author+ | Publish post |
| POST | `/api/blog/posts/:id/unpublish` | Author+ | Unpublish post |
| GET | `/api/blog/comments` | Editor+ | List all comments |
| PUT | `/api/blog/comments/:id/approve` | Editor+ | Approve comment |
| PUT | `/api/blog/comments/:id/reject` | Editor+ | Reject comment |
| PUT | `/api/blog/comments/:id/spam` | Editor+ | Mark as spam |
| DELETE | `/api/blog/comments/:id` | Admin | Delete comment |

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

## Blog Client (Headless CMS)

The backend serves as a headless CMS. Your static site consumes blog content via the client API.

### Auto-Render Blog Posts

Add containers to your HTML and the client renders posts automatically:

```html
<!-- List recent posts -->
<div data-minimaldoc-blog="list" data-limit="5"></div>

<!-- List posts by category -->
<div data-minimaldoc-blog="list" data-category="tutorials" data-limit="10"></div>

<!-- Single post by slug -->
<div data-minimaldoc-blog="single" data-slug="getting-started"></div>

<!-- Related posts -->
<div data-minimaldoc-blog="related" data-slug="getting-started" data-limit="3"></div>
```

Enable the blog feature:

```html
<script
  src="http://your-server/minimaldoc.js"
  data-endpoint="http://your-server"
  data-site-id="your-site-id"
  data-features="analytics,blog"
  defer>
</script>
```

### Programmatic API

For custom rendering:

```javascript
// List posts with filters
MinimalDoc.blog.list({
  limit: 10,
  offset: 0,
  category: 'news',
  tag: 'release',
  search: 'query'
}).then(response => {
  console.log(response.posts);
  console.log(response.total);
});

// Get single post
MinimalDoc.blog.get('my-post-slug').then(post => {
  console.log(post.title);
  console.log(post.content_html);
});

// Get related posts
MinimalDoc.blog.related('my-post-slug', 5).then(response => {
  console.log(response.posts);
});
```

### Response Format

**Post object:**

```json
{
  "id": "uuid",
  "slug": "my-post",
  "title": "My Post Title",
  "description": "Brief description",
  "content_html": "<p>Rendered HTML...</p>",
  "featured_image": "https://...",
  "category": "tutorials",
  "tags": ["go", "api"],
  "reading_time": 5,
  "published_at": "2024-01-15T10:30:00Z",
  "author": {
    "name": "John Doe",
    "avatar_url": "https://..."
  }
}
```

### Styling

The auto-rendered blog uses these CSS classes:

| Class | Element |
|-------|---------|
| `.minimaldoc-blog-list` | Container for post list |
| `.minimaldoc-blog-post` | Individual post article |
| `.minimaldoc-blog-title` | Post title (h2) |
| `.minimaldoc-blog-meta` | Date, reading time, category |
| `.minimaldoc-blog-excerpt` | Description/excerpt |
| `.minimaldoc-blog-image` | Featured image |
| `.minimaldoc-blog-tags` | Tags container |
| `.minimaldoc-blog-tag` | Individual tag |
| `.minimaldoc-blog-body` | Full post content (single view) |
| `.minimaldoc-blog-pagination` | Pagination controls |

## Image Uploads

The backend supports image uploads for blog posts.

### Upload via Admin UI

The blog editor includes an "Upload" button in the toolbar. Click it to select an image file, which uploads automatically and inserts the markdown image syntax at your cursor position.

### Upload via API

```bash
curl -X POST http://your-server/api/uploads \
  -H "Authorization: Bearer <token>" \
  -F "file=@image.jpg"
```

Response:

```json
{
  "upload": {
    "id": "uuid",
    "filename": "image.jpg",
    "mime_type": "image/jpeg",
    "size_bytes": 102400,
    "url": "/uploads/2024/01/uuid.jpg"
  },
  "url": "/uploads/2024/01/uuid.jpg"
}
```

### Storage Configuration

Configure storage via environment variables:

**Local Storage (default):**

```bash
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=./uploads
```

**S3 / S3-Compatible (AWS, MinIO, Cloudflare R2):**

```bash
STORAGE_PROVIDER=s3
STORAGE_S3_BUCKET=my-bucket
STORAGE_S3_REGION=us-east-1
STORAGE_S3_ACCESS_KEY=AKIA...
STORAGE_S3_SECRET_KEY=secret
STORAGE_S3_ENDPOINT=          # Optional: for MinIO/R2
STORAGE_S3_PUBLIC_URL=        # Optional: CDN URL
```

**Upload Limits:**

```bash
STORAGE_MAX_FILE_SIZE=5242880                           # 5MB default
STORAGE_ALLOWED_TYPES=image/jpeg,image/png,image/gif,image/webp
```

### Storage Providers

| Provider | Config | Use Case |
|----------|--------|----------|
| `local` | `STORAGE_LOCAL_PATH` | Development, simple deployments |
| `s3` | `STORAGE_S3_*` | Production, scalable storage |

S3-compatible services supported:
- AWS S3
- MinIO (self-hosted)
- Cloudflare R2
- DigitalOcean Spaces
- Backblaze B2

## Admin Portal

Access the admin UI at `http://your-server/admin`.

Features:
- Dashboard with analytics overview
- Blog management with markdown editor
- Comment moderation queue
- Feedback management
- Subscriber management
- Site settings

Login with the credentials from bootstrap or any user created via API.

## Post Visibility

Blog posts can have different visibility levels:

| Visibility | Description |
|------------|-------------|
| `public` | Visible to everyone (default) |
| `authenticated` | Requires login |
| `role_viewer` | Requires viewer role or higher |
| `role_author` | Requires author role or higher |
| `role_editor` | Requires editor role or higher |
| `role_admin` | Admins only |

Set visibility in the blog editor sidebar or via API:

```bash
curl -X POST http://your-server/api/blog/posts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Internal Update",
    "slug": "internal-update",
    "content": "...",
    "visibility": "authenticated"
  }'
```

Non-public posts are hidden from the public blog API. Users must authenticate via the admin portal to view them.

## Custom Events

Track custom user interactions (button clicks, downloads, form submissions):

### Client API

```javascript
// Track a simple event
MinimalDoc.trackEvent('signup_click');

// Track with category
MinimalDoc.trackEvent('download', 'engagement');

// Track with value
MinimalDoc.trackEvent('purchase', 'conversion', '49.99');
```

### Event Schema

| Field | Max Length | Required | Description |
|-------|------------|----------|-------------|
| `name` | 128 | Yes | Event name |
| `category` | 64 | No | Event category |
| `path` | 512 | No | Page where event occurred (auto-filled) |
| `value` | 1024 | No | Additional data |

### Admin Dashboard

Events appear in the Analytics section:
- Total event count
- Unique event names
- Events grouped by name
- Recent events table

## Security

### Security Headers

The server includes security headers on all responses:

| Header | Value | Purpose |
|--------|-------|---------|
| `X-Content-Type-Options` | `nosniff` | Prevent MIME sniffing |
| `X-Frame-Options` | `DENY` | Prevent clickjacking |
| `X-XSS-Protection` | `1; mode=block` | XSS filter |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Limit referrer leakage |
| `Permissions-Policy` | `geolocation=(), ...` | Disable unused APIs |
| `Content-Security-Policy` | See below | Restrict resource loading |

### Security Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `BOOTSTRAP_TOKEN` | - | If set, required for `/api/bootstrap` |
| `AUTH_SECURE_COOKIES` | `false` | Set `Secure` flag on cookies (requires HTTPS) |
| `SERVER_CORS_ORIGINS` | `*` | Allowed CORS origins (comma-separated) |

### Rate Limiting

Built-in rate limiting protects against abuse:

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_ENABLED` | `true` | Enable rate limiting |
| `RATE_LIMIT_LOGIN_LIMIT` | `5` | Max login attempts per window |
| `RATE_LIMIT_LOGIN_WINDOW` | `15m` | Login rate limit window |
| `RATE_LIMIT_API_LIMIT` | `100` | Max API requests per window |
| `RATE_LIMIT_API_WINDOW` | `1m` | API rate limit window |
| `RATE_LIMIT_SUBMIT_LIMIT` | `10` | Max submissions per window |
| `RATE_LIMIT_SUBMIT_WINDOW` | `1m` | Submit rate limit window |

Rate limit headers returned:
- `X-RateLimit-Limit`: Request limit
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Reset timestamp
- `Retry-After`: Seconds until reset (on 429)

### Production Checklist

Before deploying to production:

```bash
# Required: Strong JWT secret (32+ chars)
AUTH_JWT_SECRET=<random-secret-at-least-32-characters>

# Required for HTTPS: Enable secure cookies
AUTH_SECURE_COOKIES=true

# Recommended: Protect bootstrap endpoint
BOOTSTRAP_TOKEN=<random-secret>

# Recommended: Restrict CORS to your domain
SERVER_CORS_ORIGINS=https://yourdomain.com,https://api.yourdomain.com

# Recommended: Use PostgreSQL for production
DB_DRIVER=postgres
DATABASE_URL=postgres://user:pass@host/db?sslmode=require
```

### Deployment Security

1. **Use HTTPS** - Deploy behind a reverse proxy (nginx, Caddy) with TLS
2. **Firewall admin port** - Port 8090 should not be publicly accessible
3. **Separate ports** - Public API (8080) and admin (8090) run on different ports
4. **Database security** - Use strong passwords, enable SSL for PostgreSQL

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
