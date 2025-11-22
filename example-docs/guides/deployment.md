---
title: Deployment Guide
description: Deploy your MinimalDoc site to various platforms
tags:
  - deployment
  - hosting
  - production
---

# Deployment Guide

Deploy your MinimalDoc site to popular static hosting platforms.

## Build for Production

First, build your site:

```bash
minimaldoc build docs/ \
  --title="Your Site Title" \
  --base-url="https://your-domain.com" \
  --output="public"
```

This generates static files in the `public/` directory.

## GitHub Pages

### Deploy with GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'

      - name: Install MinimalDoc
        run: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

      - name: Build site
        run: minimaldoc build docs/ --base-url="https://username.github.io/repo"

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

### Manual Deployment

```bash
# Build
minimaldoc build docs/ --base-url="https://username.github.io/repo"

# Deploy
cd public
git init
git add -A
git commit -m "Deploy"
git push -f git@github.com:username/repo.git main:gh-pages
```

:::info
Enable GitHub Pages in your repository settings, selecting the `gh-pages` branch.
:::

## Netlify

### Automatic Deployment

1. Push your repo to GitHub/GitLab
2. Connect to Netlify
3. Configure build settings:

```
Build command: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest && minimaldoc build docs/
Publish directory: public
```

### Manual Deployment

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Build
minimaldoc build docs/ --base-url="https://your-site.netlify.app"

# Deploy
netlify deploy --prod --dir=public
```

### netlify.toml

Create `netlify.toml` in your repo root:

```toml
[build]
  command = "go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest && minimaldoc build docs/ --base-url='https://your-site.netlify.app'"
  publish = "public"

[[headers]]
  for = "/*"
  [headers.values]
    X-Frame-Options = "DENY"
    X-XSS-Protection = "1; mode=block"
    X-Content-Type-Options = "nosniff"
```

## Vercel

### Automatic Deployment

1. Push to GitHub/GitLab
2. Import project to Vercel
3. Configure:

```
Framework Preset: Other
Build Command: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest && minimaldoc build docs/
Output Directory: public
```

### vercel.json

Create `vercel.json`:

```json
{
  "buildCommand": "go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest && minimaldoc build docs/",
  "outputDirectory": "public",
  "headers": [
    {
      "source": "/(.*)",
      "headers": [
        {
          "key": "X-Frame-Options",
          "value": "DENY"
        }
      ]
    }
  ]
}
```

## Cloudflare Pages

### Automatic Deployment

1. Connect your repository
2. Configure build:

```
Build command: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest && minimaldoc build docs/
Build output directory: public
```

### Manual Deployment

```bash
# Build
minimaldoc build docs/

# Deploy with Wrangler
npx wrangler pages publish public
```

## Self-Hosted

### Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/your-site/public;
    index index.html;

    location / {
        try_files $uri $uri.html $uri/ =404;
    }

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

### Apache

Create `.htaccess` in your `public/` directory:

```apache
# Enable clean URLs
RewriteEngine On
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule ^(.*)$ $1.html [L]

# Gzip compression
<IfModule mod_deflate.c>
    AddOutputFilterByType DEFLATE text/html text/plain text/xml text/css text/javascript application/javascript
</IfModule>
```

### Docker

Create `Dockerfile`:

```dockerfile
FROM nginx:alpine

COPY public /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

Build and run:

```bash
# Build site
minimaldoc build docs/

# Build Docker image
docker build -t my-docs .

# Run
docker run -p 8080:80 my-docs
```

## Custom Domain

### DNS Configuration

Add these DNS records:

**For root domain (example.com):**
```
A     @     IP_ADDRESS
```

**For subdomain (docs.example.com):**
```
CNAME docs  your-host.netlify.app
```

### HTTPS/SSL

Most platforms provide free SSL:

- **GitHub Pages**: Automatic with custom domains
- **Netlify**: Automatic Let's Encrypt
- **Vercel**: Automatic
- **Cloudflare Pages**: Automatic

For self-hosted, use Let's Encrypt:

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com
```

## Optimization

### Pre-deployment Checklist

- [ ] Set correct `--base-url`
- [ ] Test all internal links
- [ ] Verify search works
- [ ] Check mobile responsiveness
- [ ] Test in multiple browsers
- [ ] Validate HTML
- [ ] Check sitemap.xml
- [ ] Test 404 page

### Performance Tips

```bash
# Enable clean URLs for better performance
minimaldoc build docs/ --clean-urls

# Verify output size
du -sh public/

# Test locally
python3 -m http.server --directory public 8000
```

## Continuous Deployment

### Best Practices

:::success Deployment Tips

1. **Automate builds** - Use CI/CD pipelines
2. **Test before deploy** - Run local builds first
3. **Use preview deploys** - Test changes safely
4. **Monitor builds** - Watch for failures
5. **Cache dependencies** - Speed up builds
:::

### Build Script

Create `build.sh`:

```bash
#!/bin/bash

echo "Building MinimalDoc site..."

# Clean previous build
rm -rf public/

# Build site
minimaldoc build docs/ \
  --title="My Documentation" \
  --base-url="https://docs.example.com" \
  --output="public"

# Verify build
if [ -f "public/index.html" ]; then
    echo "✓ Build successful!"
    exit 0
else
    echo "✗ Build failed!"
    exit 1
fi
```

Make it executable:

```bash
chmod +x build.sh
./build.sh
```

## Troubleshooting

### Build Failures

If builds fail:

1. Check Go version (1.20+)
2. Verify markdown syntax
3. Check file paths
4. Review build logs
5. Test locally first

### 404 Errors

If pages show 404:

1. Verify file exists in public/
2. Check `--base-url` matches deployment
3. Configure server for .html extension
4. Clear browser cache

### Broken Links

If links don't work:

1. Use relative paths
2. Include .html extension
3. Check case sensitivity
4. Verify paths in TOC.md

## Next Steps

- [Learn about API reference](../api/reference.html)
- [Read the FAQ](faq.html)
- [Join the community](https://github.com/studiowebux/minimaldoc)
