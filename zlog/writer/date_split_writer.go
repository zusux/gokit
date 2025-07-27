// zlog/writer/date_split_writer.go
package writer

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DateSplitWriter struct {
	mu          sync.Mutex
	logDir      string
	baseName    string // 如 out.log / err.log
	currentDate string
	file        *os.File
}

func NewDateSplitWriter(logDir, baseName string) *DateSplitWriter {
	return &DateSplitWriter{
		logDir:   logDir,
		baseName: baseName,
	}
}

func (w *DateSplitWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().Format("2006-01-02")
	if w.file == nil || w.currentDate != now {
		_ = w.rotateFile(now)
	}
	return w.file.Write(p)
}

func (w *DateSplitWriter) Sync() error {
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

func (w *DateSplitWriter) rotateFile(date string) error {
	if w.file != nil {
		_ = w.file.Close()
	}

	dir := filepath.Join(w.logDir, date)
	_ = os.MkdirAll(dir, 0755)

	path := filepath.Join(dir, w.baseName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.currentDate = date
	w.file = file
	return nil
}
