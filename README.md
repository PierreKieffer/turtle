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

logger, _ := turtle.New("out.log") // default is stdout : logger, _ := turtle.New()
logger.Info("test message", turtle.Label{Key: "foo", Value: "bar"})
logger.Warn("test message", turtle.Label{Key: "foo", Value: "bar"})
logger.Debug("test message", turtle.Label{Key: "foo", Value: "bar"})
logger.Error("test message", turtle.Label{Key: "foo", Value: "bar"})
```
```
2023-10-12 14:08:34 [INFO] test message foo:bar
2023-10-12 14:08:34 [WARN] test message foo:bar
2023-10-12 14:08:34 [DEBUG] test message foo:bar
2023-10-12 14:08:34 [ERROR] test message foo:bar
```

## Benchmark

```go
func BenchmarkTurtle(b *testing.B) {

    b.ReportAllocs()
    b.ResetTimer()

    logger, _ := New("turtle.log")
    for i := 0; i < b.N; i++ {
        logger.Info("This is a test log message", Label{Key: "foo", Value: "bar"})
    }
}

func BenchmarkSlog(b *testing.B) {

    b.ReportAllocs()
    b.ResetTimer()

    file, _ := os.OpenFile("slog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    logger := slog.New(slog.NewTextHandler(file, nil))

    for i := 0; i < b.N; i++ {
        logger.Info("This is a test log message", "foo", "bar")
    }
}

func BenchmarkZap(b *testing.B) {
    b.ReportAllocs()
    b.ResetTimer()
    config := zap.NewProductionConfig()
    config.Sampling = nil
    config.OutputPaths = []string{"zap.log"}
    logger, _ := config.Build()
    defer logger.Sync()
    for i := 0; i < b.N; i++ {
        logger.Info("This is a test log message",
            zap.String("foo", "bar"),
        )
    }
}
```
```
goos: darwin
goarch: arm64
pkg: github.com/PierreKieffer/turtle
BenchmarkTurtle-8   	  876679	      1365 ns/op	       0 B/op	       0 allocs/op
BenchmarkSlog-8     	  583070	      1989 ns/op	       0 B/op	       0 allocs/op
BenchmarkZap-8      	  509536	      2447 ns/op	     320 B/op	       3 allocs/op
```
