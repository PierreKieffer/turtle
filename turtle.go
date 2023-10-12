package turtle

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Buffer []byte

// LoggerInterface holds Logger public methods
type LoggerInterface interface {
	Info()
	Debug()
	Warn()
	Error()
	InfoEndpoint()
}

// Logger type holds a writer to write logs to output
type Logger struct {
	writer io.Writer
}

// Label is used to passed typed structured data to logger
type Label struct {
	Key   string
	Value string
}

// bufferPool
// A Pool is a set of temporary objects that may be individually saved and retrieved.
// Pool's purpose is to cache allocated but unused items for later reuse, relieving pressure on the garbage collector
// Here sync.Pool is used to cache allocated buffers
var bufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*Buffer)(&b)
	},
}

// resetBuffer resets the buffer to be empty,
// but it retains the underlying storage for use by future writes
func (b *Buffer) resetBuffer() {
	*b = (*b)[:0]
}

// New initializes logger
// Default output is Stdout
func New(outputPath ...string) (*Logger, error) {
	var l Logger

	if len(outputPath) > 0 {
		if outputPath[0] != "" {
			logfile := filepath.Join(outputPath[0])
			file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return &l, err
			}
			l.writer = file
			return &l, nil
		} else {
			return &l, fmt.Errorf("Invalid output file path")
		}
	}

	l.writer = os.Stdout
	return &l, nil
}

// Info
func (l *Logger) Info(msg string, labels ...Label) {

	b := bufferPool.Get().(*Buffer)
	defer bufferPool.Put(b)
	defer b.resetBuffer()

	b.serialize("info", msg, labels...)
	l.writer.Write(*b)
}

// Debug
func (l *Logger) Debug(msg string, labels ...Label) {

	b := bufferPool.Get().(*Buffer)
	defer bufferPool.Put(b)
	defer b.resetBuffer()

	b.serialize("debug", msg, labels...)
	l.writer.Write(*b)
}

// Error
func (l *Logger) Error(msg string, labels ...Label) {

	b := bufferPool.Get().(*Buffer)
	defer bufferPool.Put(b)
	defer b.resetBuffer()

	b.serialize("error", msg, labels...)
	l.writer.Write(*b)
}

// Warn
func (l *Logger) Warn(msg string, labels ...Label) {

	b := bufferPool.Get().(*Buffer)
	defer bufferPool.Put(b)
	defer b.resetBuffer()

	b.serialize("warn", msg, labels...)
	l.writer.Write(*b)
}

// serialize formats the timestamp, level, log message and associated labels
// and write content to the buffer
func (b *Buffer) serialize(level string, msg string, labels ...Label) {

	// time
	b.writeTime()

	// level
	switch level {
	case "info":
		*b = append(*b, " [INFO] "...)
	case "debug":
		*b = append(*b, " [DEBUG] "...)
	case "warn":
		*b = append(*b, " [WARN] "...)
	case "error":
		*b = append(*b, " [ERROR] "...)
	}

	// message
	*b = append(*b, msg...)

	// labels
	for i := range labels {
		*b = append(*b, 32) // byte 32 = " "
		*b = append(*b, labels[i].Key...)
		*b = append(*b, 58) // byte 58  = ":"
		*b = append(*b, labels[i].Value...)
	}
	// separator
	*b = append(*b, 10) // byte 10 = \n
}

// writeTime computes the time in format "2006-01-02 15:04:05" and append to b
func (b *Buffer) writeTime() {
	t := time.Now()
	year, month, day := t.Date()
	h, m, s := t.Clock()

	b.writeInt(year)
	*b = append(*b, 45)
	b.writeInt(int(month))
	*b = append(*b, 45)
	b.writeInt(day)
	*b = append(*b, 32)
	b.writeInt(h)
	*b = append(*b, 58)
	b.writeInt(m)
	*b = append(*b, 58)
	b.writeInt(s)
}

// writeInt appends the decimal form of x to b
func (b *Buffer) writeInt(x int) {
	// Compute number of decimals n
	var n int
	if x == 0 {
		n = 1
	}
	for x2 := x; x2 > 0; x2 /= 10 {
		n++
	}
	if n == 1 {
		n++
	}

	// Increase buffer len
	*b = (*b)[:len(*b)+n]
	i := len(*b) - 1
	if x < 10 {
		(*b)[i] = byte(48 + x)
		(*b)[i-1] = byte(48)
	} else {
		// Add decimals bytes in reverse
		for x >= 10 {
			q := x / 10
			(*b)[i] = byte(48 + x - q*10)
			x = q
			i--
		}
		(*b)[i] = byte(48 + x)
	}
}
