<img src="https://img.shields.io/badge/lifecycle-experimental-orange.svg" alt="Life cycle: experimental"> [![GoDoc](https://godoc.org/github.com/openbiox/ligo?status.svg)](https://godoc.org/github.com/openbiox/ligo)

# ligo

A set of Golang utils function (logging...).

## Installation

```bash
go get -u github.com/openbiox/ligo/...
```

## Usage

```golang
package main

import (
	"fmt"

	"github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/stringo"
)

func main() {
	prog := "ligo"
	log.Infof("Starting %s...", ligo)

	v := stringo.StrDetect("AACCDD#AC", "#.*")
	fmt.Println(v)
}
```

## License

Apache 2.0
