[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/treepeck/chego.svg)](https://pkg.go.dev/github.com/treepeck/chego)

Chego is a Go module that implements the rules of chess (board representation,<br/>
legal move generator, game state management, etc.).

## Usage

To use chego in your chess server, run `go get`:

```
go get github.com/treepeck/chego
```

Here is a simple example:

```go
package main

import (
	"fmt"

	"github.com/treepeck/chego"
)

func main() {
	g := chego.NewGame()
	// Scholar's mate.
	g.PushMove(chego.NewMove(chego.SF3, chego.SF2, chego.MoveNormal))
	g.PushMove(chego.NewMove(chego.SE5, chego.SE7, chego.MoveNormal))
	g.PushMove(chego.NewMove(chego.SG4, chego.SG2, chego.MoveNormal))
	g.PushMove(chego.NewMove(chego.SH4, chego.SD8, chego.MoveNormal))
 	// Prints "Is checkmate: true"
	fmt.Printf("Is checkmate: %t\n", g.IsCheckmate())
}
```

## License

Copyright (c) 2025-2026 Artem Bielikov

This project is available under the Mozilla Public License, v. 2.0.<br/>
See the [LICENSE](LICENSE) file for details.

## More

- [Tests and benchmarks](internal/perft/README.md)
- [Move compression](internal/codegen/README.md)