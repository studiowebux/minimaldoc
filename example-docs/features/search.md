---
title: Search Functionality
description: How to use the built-in search feature
tags:
  - search
  - features
---

# Search Functionality

MinimalDoc includes powerful client-side search with keyboard shortcuts.

## How It Works

Search is completely client-side - no server required! When you build your site, MinimalDoc generates a `search-index.json` file containing:

- Page titles
- Descriptions
- Tags
- Content excerpts

The search runs entirely in your browser using fuzzy matching.

## Opening Search

### Keyboard Shortcut

Press `Cmd+K` (Mac) or `Ctrl+K` (Windows/Linux) to open search from anywhere.

:::tip Pro Tip
The keyboard shortcut works on any page, making it super fast to find what you need!
:::

### Future: Search Button

(Coming soon) Click the search icon in the navigation bar.

## Using Search

### Basic Search

1. Open search with `Cmd+K` / `Ctrl+K`
2. Type your query
3. Results appear instantly as you type
4. Click a result or use keyboard navigation

### Keyboard Navigation

Once search is open:

- `↑` - Move up in results
- `↓` - Move down in results
- `Enter` - Go to selected page
- `Esc` - Close search

:::success
Navigate search entirely with your keyboard for maximum efficiency!
:::

## Search Features

### Fuzzy Matching

Search finds matches even with typos or partial words:

- "instl" finds "installation"
- "qick" finds "quick"
- "markdwn" finds "markdown"

### Multi-Field Search

Search looks in multiple places:

1. **Page titles** (highest priority)
2. **Descriptions**
3. **Tags**
4. **Content**

### Smart Ranking

Results are ranked by relevance:

- Exact title matches appear first
- Partial title matches next
- Description matches
- Content matches last

## Search Index

The search index is generated automatically during build:

```bash
minimaldoc build docs/
```

This creates `public/search-index.json` containing:

```json
{
  "entries": [
    {
      "title": "Installation",
      "description": "How to install MinimalDoc",
      "url": "/getting-started/installation.html",
      "content": "MinimalDoc is easy to install...",
      "tags": ["getting-started", "installation"]
    }
  ]
}
```

## Optimizing for Search

### Use Descriptive Titles

Good titles improve search results:

```markdown
---
title: Installing MinimalDoc on macOS
---
```

Better than:

```markdown
---
title: Installation
---
```

### Add Descriptions

Descriptions appear in search results:

```yaml
---
title: Quick Start Guide
description: Get your documentation site running in 5 minutes
---
```

### Use Tags

Tags boost search relevance:

```yaml
---
title: API Reference
tags:
  - api
  - reference
  - rest
  - graphql
---
```

### Write Clear Content

The first 500 characters of each page are indexed:

- Start with key information
- Use important terms early
- Be concise and clear

:::note
Search only indexes the first 500 characters to keep the index size manageable.
:::

## Search Best Practices

### Do's

✅ Use descriptive, unique titles
✅ Add meaningful descriptions
✅ Tag pages appropriately
✅ Use keywords naturally in content
✅ Keep titles scannable

### Don'ts

❌ Use generic titles like "Page 1"
❌ Leave descriptions empty
❌ Overuse tags
❌ Keyword stuff
❌ Make titles too long

## Examples

### Good Search-Friendly Page

```markdown
---
title: Deploying to Netlify
description: Step-by-step guide to deploy your MinimalDoc site to Netlify
tags:
  - deployment
  - netlify
  - hosting
---

# Deploying to Netlify

Learn how to deploy your MinimalDoc site to Netlify in minutes.
Netlify offers free hosting for static sites with automatic builds...
```

### Poor Search-Friendly Page

```markdown
---
title: Deployment
---

# Deployment

This page is about deployment.
```

## Performance

### Index Size

The search index is optimized for size:

- Compressed JSON format
- Limited content per page (500 chars)
- Efficient data structure

Typical sizes:
- 10 pages: ~5 KB
- 50 pages: ~25 KB
- 100 pages: ~50 KB

### Search Speed

Search is fast:
- Instant results as you type
- No network requests
- Runs entirely in browser
- Works offline

## Troubleshooting

### Search Not Working

If search doesn't work:

1. **Rebuild your site** - Ensure search-index.json exists
2. **Check browser console** - Look for JavaScript errors
3. **Verify file exists** - Check `public/search-index.json`
4. **Try another browser** - Rule out browser issues

### No Results Found

If search returns no results:

1. **Check spelling** - Try different terms
2. **Use simpler queries** - Single words work best
3. **Check content exists** - Ensure pages have content
4. **Rebuild** - Regenerate the search index

### Keyboard Shortcut Not Working

If `Cmd+K` doesn't work:

1. **Check browser conflicts** - Some extensions use Cmd+K
2. **Try Ctrl+K** - Alternative shortcut
3. **Click search button** - When available
4. **Reload page** - JavaScript may not have loaded

## Future Enhancements

Planned improvements:

- [ ] Search history
- [ ] Advanced filters (by tag, section)
- [ ] Search suggestions
- [ ] Highlighted matches in results
- [ ] Search analytics

## Technical Details

### Implementation

MinimalDoc uses a custom JavaScript search engine with:

- Fuzzy string matching
- Multi-field scoring
- Result ranking algorithm
- Keyboard navigation
- Modal UI

No external dependencies - everything is built-in!

### Customization

(Coming soon) Customize search behavior:

```javascript
// Future configuration
{
  maxResults: 10,
  fuzzyThreshold: 0.6,
  searchFields: ['title', 'description', 'content']
}
```

## Next Steps

- [Learn about theming](theming.html)
- [Explore the API](../api/reference.html)
- [Read deployment guides](../guides/deployment.html)

:::info
Search makes your documentation much more accessible. Make sure to rebuild after adding new content!
:::
