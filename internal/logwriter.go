package internal

import (
	"bytes"
	"io"
	"sync"

	"github.com/charmbracelet/log"
)

type LogWriter struct {
	Logger *log.Logger
	Level  log.Level
	mu     sync.Mutex
	buf    bytes.Buffer
}

// Write buffers p and emits any complete lines to the logger.
// Partial lines (no trailing newline) are held until the next write or Close.
func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.buf.Write(p)

	for {
		b := w.buf.Bytes()
		idx := bytes.IndexByte(b, '\n')
		if idx < 0 {
			break
		}
		line := string(b[:idx])
		w.buf.Next(idx + 1) // consume line + newline
		if w.Logger.GetLevel() <= w.Level {
			w.Logger.Print(line)
		}
	}

	return len(p), nil
}

// Close flushes any buffered partial line to the logger.
func (w *LogWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.buf.Len() > 0 && w.Logger.GetLevel() <= w.Level {
		w.Logger.Print(w.buf.String())
		w.buf.Reset()
	}
	return nil
}

func NewLogWriter(logger *log.Logger, level log.Level) io.Writer {
	return &LogWriter{Logger: logger, Level: level}
}
