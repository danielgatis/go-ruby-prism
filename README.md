# go-ruby-prism

[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-ruby-prism?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-ruby-prism)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-ruby-prism/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-ruby-prism)

**go-ruby-prism** is a Go library that leverages the [Ruby Prism parser](https://github.com/ruby/prism) compiled to WebAssembly, enabling Ruby code parsing and analysis without the need for CGO bindings. This solution provides native and efficient integration for Ruby code analysis in Go applications.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Detailed Examples](#detailed-examples)
- [API Reference](#api-reference)
- [Advanced Configuration](#advanced-configuration)
- [Project Structure](#project-structure)
- [Build and Development](#build-and-development)
- [Performance](#performance)
- [Contributing](#contributing)
- [License](#license)

## Features

- **🚀 CGO-Free**: Uses the Ruby Prism parser compiled to WebAssembly, completely eliminating the need for CGO bindings
- **⚡ High Performance**: Leverages WebAssembly efficiency for fast and efficient Ruby code parsing
- **🔄 Simplified Integration**: Seamless integration into Go applications with minimal configuration
- **🌐 Cross-Platform**: Works on any platform supported by Go, ensuring universal compatibility
- **🎯 Flexible API**: Offers multiple ways to work with the generated AST (Abstract Syntax Tree)
- **📊 Complete Ruby Support**: Compatible with Ruby versions 3.3 and 3.4
- **🛠️ Advanced Configuration**: Support for encoding, scopes, frozen string literals, and other Ruby options

## Installation

### Prerequisites

- Go 1.23 or higher

### Installation via go get

```bash
go get github.com/danielgatis/go-ruby-prism
```

### Manual installation

```bash
git clone https://github.com/danielgatis/go-ruby-prism.git
cd go-ruby-prism
make all
```

## Quick Start

### Basic Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	parser "github.com/danielgatis/go-ruby-prism/parser"
)

func main() {
	ctx := context.Background()

	// Create a new parser instance
	p, err := parser.NewParser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close(ctx)

	// Ruby code to analyze
	source := "puts 'Hello, World!'"

	// Parse the code
	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AST Root: %T\n", result.Value)
	fmt.Printf("Errors: %d\n", len(result.Errors))
	fmt.Printf("Warnings: %d\n", len(result.Warnings))
}
```

## Detailed Examples

### 1. Converting AST to JSON

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	parser "github.com/danielgatis/go-ruby-prism/parser"
)

func main() {
	ctx := context.Background()

	p, err := parser.NewParser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close(ctx)

	source := `
class User
  attr_reader :name, :email

  def initialize(name, email)
    @name = name
    @email = email
  end

  def greet
    puts "Hello, #{@name}!"
  end
end
`

	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		log.Fatal(err)
	}

	// Convert to JSON
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonResult))
}
```

### 2. Using the Visitor Pattern

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/danielgatis/go-ruby-prism/parser"
)

type CodeAnalyzer struct {
	parser.DefaultVisitor
	methodCount int
	classCount  int
}

func (v *CodeAnalyzer) Visit(node parser.Node) {
	switch node.(type) {
	case *parser.DefNode:
		v.methodCount++
		fmt.Printf("📍 Found method at line %d\n", node.StartLine())
	case *parser.ClassNode:
		v.classCount++
		fmt.Printf("📍 Found class at line %d\n", node.StartLine())
	}
	v.DefaultVisitor.Visit(node)
}

func (v *CodeAnalyzer) Analyze(node parser.Node) {
	for _, child := range node.ChildNodes() {
		v.Visit(child)
		if len(child.ChildNodes()) > 0 {
			v.Analyze(child)
		}
	}
}

func main() {
	ctx := context.Background()

	p, err := parser.NewParser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close(ctx)

	source := `
class Calculator
  def add(a, b)
    a + b
  end

  def multiply(a, b)
    a * b
  end
end

class MathUtils
  def self.pi
    3.14159
  end
end
`

	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		log.Fatal(err)
	}

	analyzer := &CodeAnalyzer{}
	analyzer.Analyze(result.Value)

	fmt.Printf("\n📊 Analysis complete:\n")
	fmt.Printf("   Classes found: %d\n", analyzer.classCount)
	fmt.Printf("   Methods found: %d\n", analyzer.methodCount)
}
```

### 3. Rails Application Analysis

The project includes a complete example that downloads and analyzes the entire Rails codebase:

```bash
go run example/parse_rails/main.go
```

This example demonstrates:
- Automatic Rails source code download
- Parsing thousands of Ruby files

## API Reference

### Parser

#### `NewParser(ctx context.Context, options ...ParserOption) (*Parser, error)`

Creates a new parser instance with the specified options.

#### `Parse(ctx context.Context, source []byte) (*ParseResult, error)`

Parses the provided Ruby code and returns the resulting AST.

#### `Close(ctx context.Context) error`

Releases WebAssembly runtime resources. Should always be called when the parser is no longer needed.

### Configuration Options

```go
// Configure source file
parser.WithFilePath("app/models/user.rb")

// Configure starting line
parser.WithLine(10)

// Configure encoding
parser.WithEncoding("UTF-8")

// Enable frozen string literals
parser.WithFrozenStringLiteral(true)

// Configure Ruby version
parser.WithVersion(parser.SyntaxVersionV3_4)

// Configure as main script
parser.WithMainScript(true)

// Configure local variable scopes
parser.WithScopes([][][]byte{{[]byte("local_var")}})

// Configure custom logger
parser.WithLogger(customLogger)
```

### ParseResult

The `ParseResult` structure contains:

```go
type ParseResult struct {
    Value    Node      // Root AST node
    Errors   []Error   // Parsing errors
    Warnings []Warning // Parser warnings
    Source   []byte    // Original source code
}
```

### Supported Syntax Versions

```go
const (
    SyntaxVersionLatest SyntaxVersion = iota  // Latest version
    SyntaxVersionV3_3                         // Ruby 3.3
    SyntaxVersionV3_4                         // Ruby 3.4
)
```

## Advanced Configuration

### Custom Logging

```go
type CustomLogger struct{}

func (l *CustomLogger) Debug(format string, args ...interface{}) {
    log.Printf("[DEBUG] "+format, args...)
}

// Use custom logger
p, err := parser.NewParser(ctx, parser.WithLogger(&CustomLogger{}))
```

### Rails Code Configuration

```go
p, err := parser.NewParser(ctx,
    parser.WithMainScript(true),
    parser.WithVersion(parser.SyntaxVersionV3_4),
    parser.WithFrozenStringLiteral(true),
    parser.WithEncoding("UTF-8"),
)
```

### Parsing with File Context

```go
p, err := parser.NewParser(ctx,
    parser.WithFilePath("app/controllers/users_controller.rb"),
    parser.WithLine(1),
    parser.WithMainScript(true),
)
```

## Project Structure

```
go-ruby-prism/
├── example/                 # Usage examples
│   ├── json/                # JSON conversion
│   ├── parse_rails/         # Rails application analysis
│   └── visitor/             # Visitor pattern
├── parser/                  # Main parser API
│   ├── parser.go            # Main interface
│   ├── gen_nodes.go         # Generated AST nodes
│   ├── gen_visitor.go       # Generated visitor pattern
│   └── parsing_options.go   # Configuration options
├── prism/                   # Ruby Prism submodule
├── wasm/                    # WebAssembly runtime
├── templates/               # Code generation templates
└── test/                    # Tests and fixtures
```

## Build and Development

### Available Make Commands

```bash
# Complete build (recommended for first time)
make all

# Generate Go code only
make generate

# Build WASM only
make wasm_build

# Format code
make format

# Clean temporary files
make clean
```

### Development Dependencies

- Ruby with Bundler (for Prism build)
- WASI SDK (downloaded automatically)
- Go 1.23+

### Build Process

1. **Initialization**: `git submodule update --init --recursive`
2. **Prism Compilation**: Compiles Ruby parser to WebAssembly
3. **Code Generation**: Generates Go structs from Prism schema
4. **Formatting**: Applies `go fmt` to all files

### Optimizations

- Parser pools for concurrent usage
- Instance reuse when possible
- Automatic WebAssembly memory management

### Parser Pool Example

```go
type ParserPool struct {
    parsers chan *parser.Parser
    ctx     context.Context
}

func NewParserPool(ctx context.Context, size int) *ParserPool {
    pool := &ParserPool{
        parsers: make(chan *parser.Parser, size),
        ctx:     ctx,
    }

    for i := 0; i < size; i++ {
        p, _ := parser.NewParser(ctx)
        pool.parsers <- p
    }

    return pool
}

func (p *ParserPool) Parse(source []byte) (*parser.ParseResult, error) {
    parser := <-p.parsers
    defer func() { p.parsers <- parser }()

    return parser.Parse(p.ctx, source)
}
```

## Practical Examples

All examples are available in the `/example` directory:

1. **JSON Export** (`example/json/`): Converts AST to JSON format
2. **Visitor Pattern** (`example/visitor/`): Demonstrates custom AST traversal
3. **Rails Analysis** (`example/parse_rails/`): Complete Rails application analysis

To run any example:

```bash
go run example/[example-name]/main.go
```

### Development

```bash
# Set up development environment
git clone https://github.com/danielgatis/go-ruby-prism.git
cd go-ruby-prism
make all

# Run tests
go test ./...

# Check linting
golangci-lint run
```

## License

Copyright (c) 2024-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE.txt)

