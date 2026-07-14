[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/treepeck/chego.svg)](https://pkg.go.dev/github.com/treepeck/chego)

Chego is a Go module that implements the rules of chess (board representation,
legal move generator, position state management, etc.).

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

func play(p *chego.Position, m chego.Move) {
	moved := p.GetPieceFromSquare(1 << m.From())
	captured := p.GetPieceFromSquare(1 << m.To())
	p.MakeMove(m, moved, captured)
}

func main() {
	l := &chego.MoveList{} // To store generated legal moves.
	p := chego.ParseFen(chego.InitialPos)
	chego.GenLegalMoves(*p, l)
	// Prints "Number of legal moves: 20"
	fmt.Printf("Number of legal moves: %d\n", l.Len)
	// Scholar's mate.
	play(p, chego.NewMove(chego.SF3, chego.SF2, chego.MoveNormal))
	play(p, chego.NewMove(chego.SE5, chego.SE7, chego.MoveNormal))
	play(p, chego.NewMove(chego.SG4, chego.SG2, chego.MoveNormal))
	play(p, chego.NewMove(chego.SH4, chego.SD8, chego.MoveNormal))
	// Re-generate legal moves.
	chego.GenLegalMoves(*p, l)
	// Prints "Is checkmate: true" since the king is under attack and there are
	// no legal moves to save it.
	fmt.Printf(
		"Is checkmate: %t\n",
		chego.GenChecksCounter(p.Bitboards, 1^p.ActiveColor) > 0 && l.Len == 0,
	)
}
```

## License

Copyright (c) 2025-2026 Artem Bielikov

This project is available under the Mozilla Public License, v. 2.0.<br/>
See the [LICENSE](LICENSE) file for details.

## More

- [Tests and benchmarks](internal/perft/README.md)
- [Move compression](internal/codegen/README.md)