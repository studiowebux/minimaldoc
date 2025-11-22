---
title: Installation
description: How to install MinimalDoc
tags:
  - getting-started
  - installation
---

# Installation

MinimalDoc is easy to install and get running. Follow the steps below for your platform.

## Prerequisites

Before installing MinimalDoc, ensure you have:

- **Go 1.20 or later** - [Download Go](https://golang.org/dl/)
- **Git** (optional) - For cloning examples

## Installation Methods

### Method 1: Install from Source

Clone the repository and build from source:

```bash
git clone https://github.com/studiowebux/minimaldoc.git
cd minimaldoc
go build -o minimaldoc ./cmd/minimaldoc
```

### Method 2: Go Install

Install directly using Go:

```bash
go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
```

### Method 3: Download Binary

Download pre-built binaries from the [releases page](https://github.com/studiowebux/minimaldoc/releases).

## Verify Installation

Verify the installation by running:

```bash
minimaldoc --version
```

You should see the version number displayed.

## Next Steps

Now that you have MinimalDoc installed, you can:

1. [Create your first site](quick-start.html)
2. [Learn about configuration](configuration.html)
3. [Explore features](../features/overview.html)

:::info
**Tip**: Add MinimalDoc to your PATH to use it from anywhere in your terminal.
:::
