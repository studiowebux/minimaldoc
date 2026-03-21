---
title: "Deployment Issues"
description: "Fixing common deployment and hosting problems"
tags: ["deployment", "hosting", "github-pages"]
---

## GitHub Pages 404 Errors

**Problem**: Pages work locally but show 404 on GitHub Pages.

**Solution**: Ensure your `base_url` includes the repository name:

```yaml
base_url: https://username.github.io/repo-name
```

Also verify the `gh-pages` branch contains the built files in the root, not in a subdirectory.

## Broken Links After Deploy

**Problem**: Internal links work locally but break in production.

**Solution**: Check that all links use the correct format:

```markdown
# Correct
[Configuration](/getting-started/configuration.html)

# Incorrect (missing leading slash)
[Configuration](getting-started/configuration.html)
```

## Assets Not Loading (HTTPS)

**Problem**: CSS/JS files blocked due to mixed content.

**Solution**: Ensure `base_url` uses `https://`:

```yaml
base_url: https://docs.example.com
```

## Clean URLs Not Working

**Problem**: Using `clean_urls: true` but getting 404s.

**Solution**: Your web server must be configured to serve `index.html` for directory requests. For static hosts like Netlify or Vercel, this works automatically. For nginx:

```nginx
location / {
    try_files $uri $uri/ $uri/index.html =404;
}
```

## Large Build Output

**Problem**: The `public/` directory is unexpectedly large.

**Solution**: Check for accidentally included files:

1. Large images should be optimized
2. Remove unused OpenAPI spec files
3. Check for duplicate assets
