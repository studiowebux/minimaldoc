---
title: Markdown Syntax Guide
description: Complete guide to writing Markdown in MinimalDoc
tags:
  - markdown
  - syntax
  - writing
---

# Markdown Syntax Guide

MinimalDoc supports GitHub Flavored Markdown with additional features.

## Headings

Create headings using `#` symbols:

```markdown
# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6
```

## Text Formatting

**Bold text** using `**bold**` or `__bold__`

*Italic text* using `*italic*` or `_italic_`

***Bold and italic*** using `***text***`

~~Strikethrough~~ using `~~text~~`

## Lists

### Unordered Lists

```markdown
- Item 1
- Item 2
  - Nested item
  - Another nested item
- Item 3
```

Result:
- Item 1
- Item 2
  - Nested item
  - Another nested item
- Item 3

### Ordered Lists

```markdown
1. First item
2. Second item
3. Third item
   1. Nested item
   2. Another nested
```

Result:
1. First item
2. Second item
3. Third item
   1. Nested item
   2. Another nested

### Task Lists

```markdown
- [x] Completed task
- [ ] Incomplete task
- [ ] Another task
```

Result:
- [x] Completed task
- [ ] Incomplete task
- [ ] Another task

## Links

```markdown
[Link text](https://example.com)
[Link with title](https://example.com "Title")
[Internal link](../getting-started/installation.html)
```

[External link](https://github.com)
[Internal link](../getting-started/installation.html)

:::info
External links automatically open in a new tab!
:::

## Code

### Inline Code

Use backticks for `inline code`.

### Code Blocks

Use triple backticks with language specification:

````markdown
```javascript
function hello(name) {
  console.log(`Hello, ${name}!`);
}
```
````

Result:

```javascript
function hello(name) {
  console.log(`Hello, ${name}!`);
}
```

### Supported Languages

Common languages with syntax highlighting:

```python
# Python
def factorial(n):
    return 1 if n <= 1 else n * factorial(n - 1)

print(factorial(5))
```

```go
// Go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

```rust
// Rust
fn main() {
    let greeting = "Hello, Rust!";
    println!("{}", greeting);
}
```

```typescript
// TypeScript
interface User {
  name: string;
  age: number;
}

const user: User = {
  name: "Alice",
  age: 30
};
```

## Tables

Create tables using pipes and hyphens:

```markdown
| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Row 1    | Data     | More     |
| Row 2    | Data     | More     |
```

Result:

| Feature | Supported | Notes |
|---------|-----------|-------|
| Tables | Yes | GFM tables |
| Alignment | Yes | Left, center, right |
| Complex cells | Limited | Keep it simple |

### Table Alignment

```markdown
| Left | Center | Right |
|:-----|:------:|------:|
| L    | C      | R     |
```

| Left | Center | Right |
|:-----|:------:|------:|
| L    | C      | R     |

## Blockquotes

```markdown
> This is a blockquote.
> It can span multiple lines.
>
> > Nested blockquotes are also supported.
```

Result:

> This is a blockquote.
> It can span multiple lines.
>
> > Nested blockquotes are also supported.

## Horizontal Rules

Create horizontal rules with three or more hyphens, asterisks, or underscores:

```markdown
---
***
___
```

---

## Images

```markdown
![Alt text](https://via.placeholder.com/150)
![Alt with title](https://via.placeholder.com/150 "Image title")
```

## Admonitions {#admonitions}

MinimalDoc supports custom admonition blocks for callouts:

### Info

```markdown
:::info
This is an informational message.
:::
```

:::info
This is an informational message.
:::

### Warning

```markdown
:::warning
This is a warning message.
:::
```

:::warning
This is a warning message.
:::

### Danger

```markdown
:::danger
This is a danger/error message.
:::
```

:::danger
This is a danger/error message.
:::

### Success

```markdown
:::success
This is a success message.
:::
```

:::success
This is a success message.
:::

### Question

```markdown
:::question
This is for FAQ or question blocks.
:::
```

:::question
This is for FAQ or question blocks.
:::

### Note

```markdown
:::note
This is a side note.
:::
```

:::note
This is a side note.
:::

### Custom Titles

```markdown
:::warning Important Security Notice
Always validate user input before processing.
:::
```

:::warning Important Security Notice
Always validate user input before processing.
:::

## HTML Support

You can embed raw HTML when needed:

```html
<div style="background: #f0f0f0; padding: 1rem; border-radius: 4px;">
  Custom HTML content
</div>
```

<div style="background: #f0f0f0; padding: 1rem; border-radius: 4px;">
  Custom HTML content
</div>

:::warning
Use HTML sparingly. Stick to Markdown when possible for better maintainability.
:::

## Special Characters

Use backslash to escape special characters:

```markdown
\* Not italic \*
\_ Not italic \_
\# Not a heading
```

## Best Practices

:::success Tips for Great Documentation

1. **Use descriptive headings** - Make them scannable
2. **Keep paragraphs short** - 3-5 sentences maximum
3. **Use code blocks** - Always specify the language
4. **Add examples** - Show, don't just tell
5. **Use admonitions** - Highlight important information
6. **Test your links** - Ensure all links work
:::

## Next Steps

- [Learn about navigation](navigation.html)
- [Explore search features](search.html)
- [See the API reference](../api/reference.html)
