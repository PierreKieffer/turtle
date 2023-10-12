# turtle

```
                 _____     ____
                /      \  |  o |
               |        |/ ___\|  
               |_________/
               |_|_| |_|_|       
```
--- 

The goal was to write a basic, fast, structured logger with optimized memory allocation to avoid pressure on the garbage collector.

## Usage 
```go 
import (
    "github.com/PierreKieffer/turtle"
)

logger, _ := turtle.New("out.log") // default is stdout
logger.Info("test message", turtle.Label{Key: "foo", Value: "bar"})
```

## Benchmark
```
goos: darwin
goarch: arm64
pkg: github.com/PierreKieffer/turtle
BenchmarkTurtle-8   	  813145	      1730 ns/op	       0 B/op	       0 allocs/op
```
