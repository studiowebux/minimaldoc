---
title: Deployment
description: Deploy MinimalDoc sites to any static hosting
tags:
  - guides
  - deployment
---

# Deployment

MinimalDoc generates static HTML that deploys anywhere.

## Build for Production

```bash
minimaldoc build ./docs \
  --base-url "https://docs.example.com" \
  --output dist
```

Output in `dist/` is ready to deploy.

## Base URL

Set `base_url` for correct absolute URLs:

```yaml
# config.yaml
base_url: https://docs.example.com
```

Or via CLI:

```bash
minimaldoc build --base-url "https://docs.example.com"
```

Affects:
- Sitemap URLs
- Canonical links
- Open Graph URLs
- RSS feed links

### Subdirectory Deployment

For `/docs/` subdirectory:

```yaml
base_url: https://example.com/docs
```

## Clean URLs

Enable for `/page/` instead of `/page.html`:

```yaml
clean_urls: true
```

Requires server configuration for fallback.

### Nginx

```nginx
location / {
    try_files $uri $uri/ $uri.html =404;
}
```

### Apache

```apache
RewriteEngine On
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule ^(.*)$ $1.html [L]
```

### Netlify

```toml
# netlify.toml
[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

## Hosting Options

### Static File Servers

| Platform | Free Tier | Custom Domain |
|----------|-----------|---------------|
| GitHub Pages | Yes | Yes |
| Netlify | Yes | Yes |
| Vercel | Yes | Yes |
| Cloudflare Pages | Yes | Yes |
| AWS S3 + CloudFront | Pay-as-you-go | Yes |
| Firebase Hosting | Yes | Yes |

### Self-Hosted

Any web server serving static files:

**Nginx:**

```nginx
server {
    listen 80;
    server_name docs.example.com;
    root /var/www/docs;
    index index.html;

    location / {
        try_files $uri $uri/ =404;
    }
}
```

**Apache:**

```apache
<VirtualHost *:80>
    ServerName docs.example.com
    DocumentRoot /var/www/docs

    <Directory /var/www/docs>
        Options Indexes FollowSymLinks
        AllowOverride None
        Require all granted
    </Directory>
</VirtualHost>
```

**Caddy:**

```
docs.example.com {
    root * /var/www/docs
    file_server
}
```

## Netlify

### Deploy via CLI

```bash
# Install
npm install -g netlify-cli

# Build
minimaldoc build ./docs --output dist

# Deploy
netlify deploy --prod --dir=dist
```

### Deploy via Git

1. Connect repository to Netlify
2. Set build command: `minimaldoc build ./docs --output dist`
3. Set publish directory: `dist`

### netlify.toml

```toml
[build]
  command = "minimaldoc build ./docs --output dist"
  publish = "dist"

[build.environment]
  GO_VERSION = "1.24"

[[headers]]
  for = "/*"
  [headers.values]
    X-Frame-Options = "DENY"
    X-Content-Type-Options = "nosniff"
```

## Vercel

### vercel.json

```json
{
  "buildCommand": "minimaldoc build ./docs --output dist",
  "outputDirectory": "dist",
  "installCommand": "go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest"
}
```

## AWS S3

### Upload

```bash
# Build
minimaldoc build ./docs --output dist

# Sync to S3
aws s3 sync dist/ s3://my-docs-bucket/ --delete
```

### CloudFront Invalidation

```bash
aws cloudfront create-invalidation \
  --distribution-id DISTRIBUTION_ID \
  --paths "/*"
```

### Terraform

```hcl
resource "aws_s3_bucket" "docs" {
  bucket = "my-docs-bucket"
}

resource "aws_s3_bucket_website_configuration" "docs" {
  bucket = aws_s3_bucket.docs.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "404.html"
  }
}
```

## Firebase

### firebase.json

```json
{
  "hosting": {
    "public": "dist",
    "ignore": ["firebase.json", "**/.*"],
    "rewrites": [
      {
        "source": "**",
        "destination": "/index.html"
      }
    ]
  }
}
```

### Deploy

```bash
minimaldoc build ./docs --output dist
firebase deploy --only hosting
```

## Docker

### Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder
RUN go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

FROM nginx:alpine
COPY --from=builder /go/bin/minimaldoc /usr/local/bin/
COPY docs/ /docs/
RUN minimaldoc build /docs --output /usr/share/nginx/html
```

### docker-compose.yml

```yaml
version: '3'
services:
  docs:
    build: .
    ports:
      - "8080:80"
```

## Optimization

### Compression

Enable gzip/brotli on your server:

**Nginx:**

```nginx
gzip on;
gzip_types text/html text/css application/javascript application/json;
```

### Caching

Set cache headers for static assets:

**Nginx:**

```nginx
location ~* \.(css|js|png|jpg|gif|ico|woff2)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### CDN

Use CDN for global distribution:
- Cloudflare
- CloudFront
- Fastly
- Bunny CDN

## Security Headers

```nginx
add_header X-Frame-Options "DENY";
add_header X-Content-Type-Options "nosniff";
add_header X-XSS-Protection "1; mode=block";
add_header Referrer-Policy "strict-origin-when-cross-origin";
add_header Content-Security-Policy "default-src 'self'";
```
