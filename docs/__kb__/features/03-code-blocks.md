---
title: "Code Blocks"
description: "Syntax highlighting and code block features"
tags: ["code", "syntax", "highlighting"]
---

## Basic Usage

Wrap code in triple backticks with the language identifier:

````markdown
```javascript
function hello() {
  console.log("Hello, world!");
}
```
````

## Supported Languages

MinimalDoc uses Chroma for syntax highlighting, supporting 100+ languages including:

- JavaScript, TypeScript
- Python, Ruby, Go, Rust
- HTML, CSS, SCSS
- Bash, PowerShell
- YAML, JSON, TOML
- SQL, GraphQL
- And many more

## Copy Button

All code blocks include a copy button in the top-right corner. Click to copy the code to your clipboard.

## Line Numbers

Code blocks automatically display line numbers for blocks longer than 3 lines.

## Inline Code

Use single backticks for inline code: `const x = 42`

## Code in Admonitions

Code blocks work inside admonitions:

:::tip
Here's a useful snippet:
```bash
minimaldoc build --clean-urls
```
:::
