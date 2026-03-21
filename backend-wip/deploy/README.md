# MinimalDoc Server Deployment

Production deployment guide for minimaldoc-server.

## Prerequisites

- Linux server (Ubuntu 22.04+ recommended)
- nginx installed
- certbot for TLS certificates
- PostgreSQL (recommended) or SQLite

## Quick Start

### 1. Build the Server

```bash
git clone https://github.com/studiowebux/minimaldoc
cd minimaldoc
make build-server
sudo cp minimaldoc-server /usr/local/bin/
```

### 2. Create System User

```bash
sudo useradd -r -s /bin/false minimaldoc
sudo mkdir -p /var/lib/minimaldoc /var/log/minimaldoc
sudo chown minimaldoc:minimaldoc /var/lib/minimaldoc /var/log/minimaldoc
```

### 3. Configure Environment

Create `/etc/minimaldoc/env`:

```bash
sudo mkdir -p /etc/minimaldoc
sudo tee /etc/minimaldoc/env << 'EOF'
# Database
DB_DRIVER=postgres
DATABASE_URL=postgres://minimaldoc:password@localhost/minimaldoc?sslmode=require

# Authentication
AUTH_JWT_SECRET=your-random-secret-at-least-32-characters-long
AUTH_SECURE_COOKIES=true
BOOTSTRAP_TOKEN=your-bootstrap-secret

# Server
SERVER_PORT=8080
SERVER_ADMIN_PORT=8090
SERVER_CORS_ORIGINS=https://api.example.com,https://docs.example.com

# Email (optional)
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-user
SMTP_PASS=your-password
EMAIL_FROM_ADDRESS=noreply@example.com
EMAIL_BASE_URL=https://api.example.com

# Storage
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=/var/lib/minimaldoc/uploads
EOF

sudo chmod 600 /etc/minimaldoc/env
sudo chown minimaldoc:minimaldoc /etc/minimaldoc/env
```

### 4. Create Systemd Service

Create `/etc/systemd/system/minimaldoc.service`:

```ini
[Unit]
Description=MinimalDoc Server
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=minimaldoc
Group=minimaldoc
WorkingDirectory=/var/lib/minimaldoc
EnvironmentFile=/etc/minimaldoc/env
ExecStart=/usr/local/bin/minimaldoc-server
Restart=always
RestartSec=5

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/minimaldoc /var/log/minimaldoc

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable minimaldoc
sudo systemctl start minimaldoc
```

### 5. Configure Nginx

```bash
# Copy config
sudo cp deploy/nginx/minimaldoc.conf /etc/nginx/sites-available/

# Edit server_name and paths
sudo nano /etc/nginx/sites-available/minimaldoc.conf

# Enable
sudo ln -s /etc/nginx/sites-available/minimaldoc.conf /etc/nginx/sites-enabled/

# Get TLS certificates
sudo certbot certonly --webroot -w /var/www/certbot -d api.example.com -d admin.example.com

# Test and reload
sudo nginx -t
sudo systemctl reload nginx
```

### 6. Bootstrap

Create initial site and admin user:

```bash
curl -X POST https://api.example.com/api/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "secure-password",
    "site_name": "My Docs",
    "domain": "docs.example.com",
    "bootstrap_token": "your-bootstrap-secret"
  }'
```

Save the returned `site_id` for your static site configuration.

## Configuration Options

### Two-Domain Setup (Recommended)

Separate domains for API and admin:

| Domain | Purpose | Access |
|--------|---------|--------|
| `api.example.com` | Public API | Open |
| `admin.example.com` | Admin portal | IP restricted |

Use `nginx/minimaldoc.conf`.

### Single-Domain Setup

Both on same domain with path-based routing:

| Path | Purpose | Access |
|------|---------|--------|
| `/api/*` | Public API | Open |
| `/admin/*` | Admin portal | IP restricted |

Use `nginx/minimaldoc-single-domain.conf`.

### Hybrid Setup (Static + Backend)

Static site generator as core, with optional backend for blog/forum:

| Path | Source | Purpose |
|------|--------|---------|
| `/` | nginx (static) | Landing, docs, changelog, kb, etc. |
| `/blog/*` | backend:8080 | Blog pages |
| `/forum/*` | backend:8080 | Forum pages |
| `/api/*` | backend:8080 | Public API |
| `/admin/*` | backend:8090 | Admin portal (IP restricted) |

Use `nginx/minimaldoc-hybrid.conf`.

## Security Checklist

| Item | Status |
|------|--------|
| HTTPS enabled | Required |
| Strong JWT secret (32+ chars) | Required |
| Secure cookies enabled | Required for HTTPS |
| Bootstrap token set | Recommended |
| CORS restricted to your domains | Recommended |
| Admin port firewalled or IP restricted | Recommended |
| PostgreSQL with SSL | Recommended |
| Systemd security hardening | Included |

## Monitoring

### Logs

```bash
# Application logs
sudo journalctl -u minimaldoc -f

# Nginx access logs
sudo tail -f /var/log/nginx/minimaldoc_api_access.log
sudo tail -f /var/log/nginx/minimaldoc_admin_access.log
```

### Health Check

```bash
curl https://api.example.com/health
```

### Database

```bash
# PostgreSQL
sudo -u postgres psql minimaldoc -c "SELECT COUNT(*) FROM page_views;"
```

## Backup

### Database

```bash
# PostgreSQL
pg_dump minimaldoc > backup.sql

# SQLite
cp /var/lib/minimaldoc/minimaldoc.db backup.db
```

### Uploads

```bash
tar -czf uploads-backup.tar.gz /var/lib/minimaldoc/uploads
```

## Troubleshooting

### Server won't start

```bash
# Check logs
sudo journalctl -u minimaldoc -n 50

# Verify config
sudo -u minimaldoc env $(cat /etc/minimaldoc/env | xargs) /usr/local/bin/minimaldoc-server
```

### 502 Bad Gateway

```bash
# Check if server is running
sudo systemctl status minimaldoc

# Check ports
sudo ss -tlnp | grep -E '8080|8090'
```

### TLS Certificate Issues

```bash
# Renew certificates
sudo certbot renew

# Check certificate
openssl s_client -connect api.example.com:443 -servername api.example.com
```
