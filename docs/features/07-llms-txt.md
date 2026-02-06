---
title: LLMs.txt
description: Generate LLM-friendly documentation for AI assistants
tags:
  - features
  - llm
  - ai
---

# LLMs.txt

MinimalDoc generates LLM-friendly output for AI assistants.

## Generated Files

| File | Purpose |
|------|---------|
| `llms.txt` | Index with links to individual pages |
| `llms-full.txt` | All documentation in a single file |
| `*.md` | Companion markdown files alongside HTML |

## Purpose

LLM-friendly format enables:
- AI chatbots to answer questions about your project
- IDE assistants (Copilot, Cursor) to understand your docs
- Custom AI integrations
- Documentation search tools

## Configuration

Enabled by default:

```yaml
enable_llms: true
```

Disable:

```yaml
enable_llms: false
```

Or via CLI:

```bash
minimaldoc build --llms=false
```

## File Location

Generated at output root:

```
public/
├── index.html
├── llms.txt              # Index with page links
├── llms-full.txt         # All content in one file
├── getting-started/
│   ├── installation.html
│   └── installation.md   # Companion markdown
├── sitemap.xml
└── ...
```

## llms.txt Format

Index file with links to documentation:

```markdown
# Site Title

> Site description

## Documentation

- [Installation](/getting-started/installation.md): Install MinimalDoc
- [Quick Start](/getting-started/quick-start.md): Create your first site

## Optional

- [Complete Documentation](/llms-full.txt): All documentation in a single file
```

## llms-full.txt Format

All content concatenated:

```markdown
# Site Title

> Site description

---

# Installation

Full page content here...

---

# Quick Start

Full page content here...
```

## Companion .md Files

Each HTML page gets a markdown companion for direct LLM access:

```
/getting-started/installation.html  →  /getting-started/installation.md
```

## Use Cases

### ChatGPT / Claude

Upload `llms.txt` as context:

```
Based on the attached documentation, how do I configure dark mode?
```

### GitHub Copilot

Reference in your codebase:

```typescript
// See llms.txt for MinimalDoc configuration options
const config = {
  // ...
};
```

### Custom AI Integration

```python
import openai

with open('llms.txt', 'r') as f:
    docs = f.read()

response = openai.ChatCompletion.create(
    model="gpt-4",
    messages=[
        {"role": "system", "content": f"Documentation:\n{docs}"},
        {"role": "user", "content": "How do I add a new page?"}
    ]
)
```

### RAG (Retrieval Augmented Generation)

Split and embed for vector search:

```python
from langchain.text_splitter import MarkdownTextSplitter
from langchain.vectorstores import Chroma

splitter = MarkdownTextSplitter(chunk_size=1000)
chunks = splitter.split_text(docs)

vectorstore = Chroma.from_texts(chunks, embeddings)
```

## Optimization

### Exclude Content

Hidden pages are excluded:

```yaml
---
title: Internal Notes
hidden: true
---
```

### Prioritize Content

Structure pages so important content appears early:
- Clear, concise introductions
- Key information first
- Examples after explanations

### Token Efficiency

LLM context windows are limited. Keep documentation:
- Focused and relevant
- Free of redundancy
- Well-organized

## File Size

`llms-full.txt` size depends on documentation volume. Most LLMs handle 100KB+ easily.

## Comparison

| Format | Use Case |
|--------|----------|
| `llms.txt` | Index with links to pages |
| `llms-full.txt` | All content for full context |
| `*.md` companions | Individual page access |
| `sitemap.xml` | Search engines |
| `search-index.json` | Client-side search |
