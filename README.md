# go-ruby-prism

[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-ruby-prism?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-ruby-prism)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-ruby-prism/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-ruby-prism)

The go-ruby-prism is package that leverages the Ruby Prism parser compiled to WebAssembly for parsing Ruby code without the need for CGO.

## Features

- **CGO-Free**: Go-Ruby-Prism utilizes the [Ruby Prism parser](https://github.com/ruby/prism) compiled to WebAssembly, eliminating the need for CGO bindings.
- **Simplified Integration**: Seamlessly integrate Ruby code parsing into your Go applications with minimal setup.
- **High Performance**: Harnesses the efficiency of WebAssembly for speedy and efficient parsing of Ruby code.
- **Cross-Platform**: Works across various platforms supported by Go, ensuring compatibility in diverse environments.

## Usage

Here's a basic example demonstrating how to use this package:

```go
package main

import (
	"context"
	"fmt"

	parser "github.com/danielgatis/go-ruby-prism/parser"
)

func main() {
	ctx := context.Background()

	p, _ := parser.NewParser(ctx)
	source := "puts 'Hello, World!'"
	result, _ := p.Parse(ctx, source)
	fmt.Println(result)
}
```

You can find more examples in the examples folder.


## License

Copyright (c) 2024-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE.txt)

