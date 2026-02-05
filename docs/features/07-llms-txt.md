---
title: LLMs.txt
description: Generate LLM-friendly documentation for AI assistants
tags:
  - features
  - llm
  - ai
---

# LLMs.txt

MinimalDoc generates `llms.txt`, a single file containing all documentation for AI assistants.

## Purpose

LLM-friendly format enables:
- AI chatbots to answer questions about your project
- IDE assistants (Copilot, Cursor) to understand your docs
- Custom AI integrations
- Documentation search tools

## Output Format

```
# MinimalDoc Documentation

## Table of Contents
- Installation
- Quick Start
- Configuration
...

---

# Installation

MinimalDoc is easy to install...

---

# Quick Start

Initialize a new site...

---
```

All pages concatenated with navigation structure.

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
├── llms.txt          # LLM-friendly content
├── sitemap.xml
└── ...
```

## Content Structure

### Header

```
# {Site Title}

{Site Description}

Generated: {Build Date}
Source: {Base URL}
```

### Navigation

```
## Table of Contents

### Getting Started
- Installation
- Quick Start
- Configuration

### Features
- Search
- Theming
...
```

### Pages

Each page includes:

```
---

# {Page Title}

{Page Description}

Tags: {comma-separated tags}

{Full page content in Markdown}
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

Typical sizes:

| Pages | llms.txt Size |
|-------|---------------|
| 10 | ~20 KB |
| 50 | ~100 KB |
| 100 | ~200 KB |

Most LLMs handle 100KB+ easily.

## Comparison with Other Formats

| Format | Use Case |
|--------|----------|
| `llms.txt` | AI assistants, full context |
| `sitemap.xml` | Search engines, URL discovery |
| `search-index.json` | Client-side search |
| RSS/Atom | Feed readers, updates |

## Validation

Check output:

```bash
head -100 public/llms.txt
wc -l public/llms.txt
```

Ensure:
- All pages included
- Navigation structure correct
- Content readable
