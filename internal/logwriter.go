package internal

import (
	"io"
	"strings"

	"github.com/charmbracelet/log"
)

type LogWriter struct {
	Logger *log.Logger
	Level  log.Level
}

func (w *LogWriter) Write(p []byte) (n int, err error) {

	go func() {
		if w.Logger.GetLevel() > w.Level {
			return
		}
		w.Logger.Print(strings.TrimSuffix(string(p), "\n"))
	}()

	return len(p), nil
}
func (w *LogWriter) Close() error {
	return nil
}

func NewLogWriter(logger *log.Logger, level log.Level) io.Writer {

	outWriter := &LogWriter{Logger: logger, Level: level}

	return outWriter
}
