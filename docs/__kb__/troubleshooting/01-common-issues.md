---
title: "Common Issues"
description: "Solutions to frequently encountered problems"
tags: ["troubleshooting", "errors", "help"]
---

## Build Fails with "config.yaml not found"

**Problem**: Running `minimaldoc build` shows an error about missing config.

**Solution**: Ensure you're running the command from the directory containing your `docs/` folder, or specify the path explicitly:

```bash
minimaldoc build ./path/to/docs
```

## Styles Not Loading

**Problem**: The generated site shows unstyled HTML.

**Solution**: Check your `base_url` in `config.yaml`. It should match your deployment URL:

```yaml
# For GitHub Pages
base_url: https://username.github.io/repo-name

# For local development
base_url: http://localhost:8889
```

## Search Not Working

**Problem**: The search feature returns no results.

**Solution**:

1. Ensure `enable_search: true` is set in your config
2. Rebuild the site to regenerate the search index
3. Check browser console for JavaScript errors

## Images Not Displaying

**Problem**: Images in markdown don't appear.

**Solution**: Use relative paths from your markdown file's location:

```markdown
![Alt text](./images/screenshot.png)
```

Or use absolute paths from the docs root:

```markdown
![Alt text](/images/screenshot.png)
```
