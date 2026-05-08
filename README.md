# gograph

[![Go Report Card](https://goreportcard.com/badge/github.com/ozgurcd/gograph)](https://goreportcard.com/report/github.com/ozgurcd/gograph)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gograph` is a local-only CLI tool designed to generate repository structures and improve IDE context awareness for codebases.

> **Note on Language Support:** I originally built `gograph` specifically for **Golang** because that is what I needed for my own workflows. It currently only parses and maps Go codebases. However, the architecture is extensible! If you want to add support for other languages (Python, TypeScript, Rust, etc.), **contributions are more than welcome.** Please see the [Contributing Guide](CONTRIBUTING.md) to get started.

## Features
- **Local Only:** No network calls or external API dependencies. All analysis is done securely on your machine.
- **Go Focused:** Deeply understands Go project structures, packages, and dependencies.
- **Fast:** Written in Go for high performance.

## Installation

```bash
go install github.com/ozgurcd/gograph@latest
```

## Usage

Navigate to your Go project and run:

```bash
gograph path/to/repo
```

## Contributing

We love pull requests! See the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to build, test, and contribute to the project. If you are adding support for a new language, please open an issue first to discuss the design.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
