## Perfromance

[Perft](https://www.chessprogramming.org/Perft) is used to validate the correctness
and productivity of the move generator.

To execute the performance test, run this command in the chego folder:

```
go run ./internal/perft/perft.go -depth {IntValue}
```

## Profiling

Using the `go tool pprof` the performance test can be profiled to detect the
performance bottlenecks.

`Makefile` contains the build commands.