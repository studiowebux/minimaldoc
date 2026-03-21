---
title: Markdown Syntax
description: Supported Markdown syntax and GitHub Flavored Markdown extensions
tags:
  - core-concepts
  - markdown
---

# Markdown Syntax

MinimalDoc supports standard Markdown plus GitHub Flavored Markdown (GFM) extensions.

## Headings

```markdown
# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6
```

Headings generate anchor IDs automatically:
- `## Getting Started` → `#getting-started`
- `### API Reference` → `#api-reference`

## Text Formatting

```markdown
**Bold text**
*Italic text*
***Bold and italic***
~~Strikethrough~~
`Inline code`
```

**Bold text**, *Italic text*, ***Bold and italic***, ~~Strikethrough~~, `Inline code`

## Links

```markdown
[Internal link](../guides/deployment.md)
[Anchor link](#headings)
[External link](https://example.com)
```

External links automatically open in new tabs with `target="_blank"`.

## Images

```markdown
![Alt text](/images/screenshot.png)
![Alt text](/images/photo.jpg "Optional title")
```

## Lists

### Unordered

```markdown
- Item one
- Item two
  - Nested item
  - Another nested
- Item three
```

- Item one
- Item two
  - Nested item
  - Another nested
- Item three

### Ordered

```markdown
1. First item
2. Second item
   1. Nested first
   2. Nested second
3. Third item
```

1. First item
2. Second item
   1. Nested first
   2. Nested second
3. Third item

### Task Lists

```markdown
- [x] Completed task
- [ ] Incomplete task
- [ ] Another task
```

- [x] Completed task
- [ ] Incomplete task
- [ ] Another task

## Blockquotes

```markdown
> Single line quote

> Multi-line quote
> continues here
>
> > Nested quote
```

> Single line quote

> Multi-line quote
> continues here
>
> > Nested quote

## Code

### Inline

```markdown
Use `config.yaml` for configuration.
```

Use `config.yaml` for configuration.

### Blocks

````markdown
```go
package main

func main() {
    fmt.Println("Hello")
}
```
````

```go
package main

func main() {
    fmt.Println("Hello")
}
```

Supported languages include: go, javascript, typescript, python, bash, yaml, json, html, css, sql, rust, swift, and 100+ more via Chroma.

## Tables

```markdown
| Header 1 | Header 2 | Header 3 |
|----------|:--------:|---------:|
| Left     | Center   | Right    |
| Data     | Data     | Data     |
```

| Header 1 | Header 2 | Header 3 |
|----------|:--------:|---------:|
| Left     | Center   | Right    |
| Data     | Data     | Data     |

Column alignment:
- `:---` Left (default)
- `:---:` Center
- `---:` Right

## Horizontal Rule

```markdown
---
```

---

## HTML

Raw HTML is supported:

```html
<details>
<summary>Click to expand</summary>

Hidden content here.

</details>
```

<details>
<summary>Click to expand</summary>

Hidden content here.

</details>

## Escaping

Use backslash to escape special characters:

```markdown
\*not italic\*
\# not a heading
\[not a link\]
```

\*not italic\*, \# not a heading, \[not a link\]

## Line Breaks

Two trailing spaces create a line break:

```markdown
Line one
Line two
```

Or use `<br>`:

```markdown
Line one<br>Line two
```
