## Compression

Chego allows to generate Huffman codes for legal moves, which helps to compress<br/>
the completed move storage and drastically reduce the database size.

1. Download the PGN file containing one or more games (https://database.lichess.org/).

2. Prepare the file for code generation by executing this command in the chego folder:

```
go run ./internal/codegen/codegen.go -input {Filename.pgn} -output {Clean.txt} -task clean
```

3. Generate Huffman codes after cleaning:

```
go run ./internal/codegen/codegen.go -input {Clean.txt} -output {Codes.txt} -task generate -workers {IntValue}
```

The workers flag defines how many concurrent goroutines will perform Huffman<br/>
code generation.  Higher values reduce execution time, but increase CPU usage.