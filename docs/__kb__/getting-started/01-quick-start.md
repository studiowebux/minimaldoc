---
title: "Quick Start Guide"
description: "Get up and running with MinimalDoc in minutes"
tags: ["beginner", "setup", "installation"]
---

## Prerequisites

Before installing MinimalDoc, ensure you have:

- A terminal/command line interface
- Basic familiarity with Markdown

## Installation

Download the latest binary for your platform:

```bash
curl -sSL \
  https://github.com/studiowebux/minimaldoc/releases/latest/download/minimaldoc-$(uname -s)-$(uname -m) \
  -o minimaldoc
chmod +x minimaldoc
```

## Create Your First Site

Initialize a new documentation project:

```bash
./minimaldoc init my-docs
cd my-docs
```

## Build and Preview

Generate your static site:

```bash
../minimaldoc build
```

Your documentation is now available in the `public/` directory.

## Next Steps

- Learn about [configuration options](/getting-started/configuration.html)
- Explore [available features](/features/overview.html)
- Check out [deployment guides](/guides/deployment.html)
