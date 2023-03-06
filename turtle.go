package turtle

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	BUFFER_CAP int = 256 // Buffer capacity for memory allocation
)

var (
	INFO  []byte = []byte{32, 91, 73, 78, 70, 79, 93, 32}     //" [INFO] "
	DEBUG []byte = []byte{32, 91, 68, 69, 66, 85, 71, 93, 32} //" [DEBUG] "
	WARN  []byte = []byte{32, 91, 87, 65, 82, 78, 93, 32}     //" [WARN] "
	ERROR []byte = []byte{32, 91, 69, 82, 82, 79, 82, 93, 32} //" [ERROR] "
)

// LoggerInterface
// Holds Logger public methods
type LoggerInterface interface {
	Info()
	Debug()
	Error()
	InfoEndpoint()
}

// Logger
// Logger type holds a writer to write logs to output
type Logger struct {
	writer io.Writer
}

// Label
// Label is used to passed typed structured data to logger
type Label struct {
	Key   string
	Value string
}

// bufferPool
// A Pool is a set of temporary objects that may be individually saved and retrieved.
// Pool's purpose is to cache allocated but unused items for later reuse, relieving pressure on the garbage collector
// Here sync.Pool is used to cache allocated buffers
var bufferPool = &sync.Pool{
	New: func() interface{} {
		return newBuffer()
	},
}

// New
// Initialize logger
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
			return &l, fmt.Errorf("invalid outputPath")
		}
	}

	l.writer = os.Stdout
	return &l, nil
}

// serialize
// serialize format the timestamp, level, log message and associated labels
// and write content to the buffer
func serialize(buffer *[]byte, level string, msg string, labels ...Label) {

	// time
	*buffer = append(*buffer, []byte(time.Now().Format("2006-01-02 15:04:05"))...)

	// level
	switch level {
	case "info":
		*buffer = append(*buffer, INFO...)
	case "debug":
		*buffer = append(*buffer, DEBUG...)
	case "warn":
		*buffer = append(*buffer, WARN...)
	case "error":
		*buffer = append(*buffer, ERROR...)
	}

	// message
	*buffer = append(*buffer, []byte(msg)...)

	// labels
	for i := range labels {
		*buffer = append(*buffer, 32) // byte 32 = " "
		*buffer = append(*buffer, []byte(labels[i].Key)...)
		*buffer = append(*buffer, 58) // byte 58  = ":"
		*buffer = append(*buffer, []byte(labels[i].Value)...)
	}
	// separator
	*buffer = append(*buffer, 10) // byte 10 = \n
}

// newBuffer
// Init a new byte buffer of capacity BUFFER_CAP
func newBuffer() *[]byte {
	var b = make([]byte, 0, BUFFER_CAP)
	return &b
}

// resetBuffer
// resetBuffer resets the buffer to be empty,
// but it retains the underlying storage for use by future writes
func resetBuffer(b *[]byte) {
	*b = (*b)[:0]
}

// Info
func (l *Logger) Info(msg string, labels ...Label) {

	b := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(b)
	defer resetBuffer(b)

	serialize(b, "info", msg, labels...)
	l.writer.Write(*b)
}

// Debug
func (l *Logger) Debug(msg string, labels ...Label) {

	b := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(b)
	defer resetBuffer(b)

	serialize(b, "debug", msg, labels...)
	l.writer.Write(*b)
}

// Error
func (l *Logger) Error(msg string, labels ...Label) {

	b := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(b)
	defer resetBuffer(b)

	serialize(b, "error", msg, labels...)
	l.writer.Write(*b)
}

// Warn
func (l *Logger) Warn(msg string, labels ...Label) {

	b := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(b)
	defer resetBuffer(b)

	serialize(b, "warn", msg, labels...)
	l.writer.Write(*b)
}
