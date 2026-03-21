# MinimalDoc

A minimal static site generator for documentation. Fast, clean, and easy to use.

> **Work in Progress** — Most features are stable. Some newer additions
> (waitlist, claude-assist, PDF export) are still being refined.

> **Note:** The optional backend server (`minimaldoc-server`) has been removed
> from this repository. The scope was expanding too fast and the code was not
> production-ready. This project now focuses exclusively on the static site
> generator. The server effort will be revisited later.

## Quick Start

```bash
# Install
git clone https://github.com/studiowebux/minimaldoc
cd minimaldoc && make build

# Create a new docs site
./minimaldoc init my-docs
cd my-docs

# Build to static HTML
../minimaldoc build
```

## Documentation

Full documentation: [minimaldoc.studiowebux.com](https://minimaldoc.studiowebux.com)

The `docs/` directory is a demo site that showcases what MinimalDoc can generate.
Build it locally with `make docs` then open `public/index.html`.

## Development

```bash
make build    # Build the CLI binary
make test     # Run tests
make lint     # gofmt + go vet
make ci       # lint + test
make docs     # Build the demo docs site
make check    # staticcheck
make clean    # Remove build artifacts
```

## License

See [LICENSE](LICENSE) file.

## Credits

Created by [Studio Webux](https://studiowebux.com)
