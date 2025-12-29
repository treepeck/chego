[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/treepeck/chego.svg)](https://pkg.go.dev/github.com/treepeck/chego)

Chego implements chessboard state management, legal move generation, and move<br/>
compression.

Piece positions are stored as bitboards.

Move generation is implemented using the Magic Bitboards method.

Compression is implemented using the Huffman coding method.

It is assigned to use in the web-servers (for example, [justchess.org](https://justchess.org/)), hence it doesn't<br/> provide any GUI or CLI.

## Usage

To install chego, run `go get`:

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

## Local installation

First install the Go compiler version 1.24.1 or newer (see https://go.dev/dl).

Once the compiler is installed, clone this repository:

```
git clone https://github.com/treepeck/chego
cd chego
```

## Tests and benchmarks

To execute the performance test, run this command in the chego folder:

```
go run ./internal/perft/perft.go -depth {IntValue}
```

Chego generates 119060324 moves at depth 6 in approximately 6 seconds on an<br/>
Intel i7-10750H CPU.

## Compression

Chego allows to generate Huffman codes for legal moves, which helps to compress<br/>
the completed move storage and drastically reduce the size of database.

1. Download the PGN file containing one or more games (https://database.lichess.org/).

2. Prepare the file to code generation by executing this command in the chego folder:

```
go run ./internal/codegen/codegen.go -input {Filename.pgn} -output {Clean.txt} -task clean
```

3. Generate Huffman codes after cleaning:

```
go run ./internal/codegen/codegen.go -input {Clean.txt} -output {Codes.txt} -task generate -workers {IntValue}
```

The workers flag defines how many concurrent goroutines will perform Huffman<br/>
code generation.  Higher values reduce execution time, but increase CPU usage.

## License

Copyright (c) 2025 Artem Bielikov

This project is available under the Mozilla Public License, v. 2.0.<br/>
See the [LICENSE](LICENSE) file for details.
