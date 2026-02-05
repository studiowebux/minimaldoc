---
title: Code Highlighting
description: Syntax highlighting for 100+ programming languages
tags:
  - features
  - code
  - syntax-highlighting
---

# Code Highlighting

MinimalDoc uses Chroma for syntax highlighting, supporting 100+ languages.

## Basic Usage

Wrap code in triple backticks with language identifier:

````markdown
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```
````

Renders as:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## Supported Languages

### Web

| Language | Identifier |
|----------|------------|
| HTML | `html` |
| CSS | `css` |
| JavaScript | `javascript`, `js` |
| TypeScript | `typescript`, `ts` |
| JSON | `json` |
| XML | `xml` |
| YAML | `yaml` |
| TOML | `toml` |

### Systems

| Language | Identifier |
|----------|------------|
| Go | `go` |
| Rust | `rust` |
| C | `c` |
| C++ | `cpp`, `c++` |
| Swift | `swift` |
| Objective-C | `objectivec` |

### Scripting

| Language | Identifier |
|----------|------------|
| Python | `python`, `py` |
| Ruby | `ruby`, `rb` |
| PHP | `php` |
| Perl | `perl` |
| Lua | `lua` |

### Shell

| Language | Identifier |
|----------|------------|
| Bash | `bash`, `sh` |
| PowerShell | `powershell` |
| Fish | `fish` |
| Zsh | `zsh` |

### Data

| Language | Identifier |
|----------|------------|
| SQL | `sql` |
| GraphQL | `graphql` |
| Protobuf | `protobuf` |

### Other

| Language | Identifier |
|----------|------------|
| Markdown | `markdown`, `md` |
| Diff | `diff` |
| Makefile | `makefile` |
| Dockerfile | `dockerfile` |
| Terraform | `hcl`, `terraform` |

## Examples

### JavaScript

```javascript
const greet = (name) => {
  return `Hello, ${name}!`;
};

export default greet;
```

### Python

```python
def fibonacci(n: int) -> list[int]:
    if n <= 0:
        return []
    sequence = [0, 1]
    while len(sequence) < n:
        sequence.append(sequence[-1] + sequence[-2])
    return sequence[:n]
```

### Rust

```rust
fn main() {
    let numbers: Vec<i32> = (1..=10).collect();
    let sum: i32 = numbers.iter().sum();
    println!("Sum: {}", sum);
}
```

### SQL

```sql
SELECT u.name, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.created_at > '2025-01-01'
GROUP BY u.id
HAVING order_count > 5
ORDER BY order_count DESC;
```

### YAML

```yaml
server:
  host: localhost
  port: 8080
  tls:
    enabled: true
    cert: /path/to/cert.pem
    key: /path/to/key.pem
```

### Bash

```bash
#!/bin/bash
set -euo pipefail

for file in *.md; do
    echo "Processing: $file"
    minimaldoc build "$file"
done
```

### Diff

```diff
- const oldValue = "deprecated";
+ const newValue = "updated";

  function process() {
-     return oldValue;
+     return newValue;
  }
```

## Inline Code

Use single backticks for inline code:

```markdown
Use the `build` command to generate HTML.
```

Renders as: Use the `build` command to generate HTML.

## No Language

Omit language for plain text:

````markdown
```
Plain text without highlighting
```
````

```
Plain text without highlighting
```

## Escaping

To show code block syntax, use four backticks:

`````markdown
````markdown
```go
// This shows the code block syntax
```
````
`````

## Theme Integration

Code blocks adapt to light/dark mode:

| Mode | Background | Text |
|------|------------|------|
| Light | `#f5f5f5` | Dark syntax colors |
| Dark | `#0d1117` | Light syntax colors |

## Performance

Syntax highlighting happens at build time:
- No JavaScript required for highlighting
- Fast page loads
- Works without JS enabled
