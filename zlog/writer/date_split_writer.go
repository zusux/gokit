// zlog/writer/date_split_writer.go
package writer

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DateSplitWriter struct {
	mu          sync.Mutex // 控制写文件的锁
	cleanMutex  sync.Mutex // 控制清理的锁
	logDir      string
	baseName    string // 如 out.log / err.log
	currentDate string
	file        *os.File
	rotation    int64
}

func NewDateSplitWriter(logDir, baseName string, rotation int64) *DateSplitWriter {
	return &DateSplitWriter{
		logDir:   logDir,
		baseName: baseName,
		rotation: rotation,
	}
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

func (w *DateSplitWriter) currentTimeKey(rotation int64) string {
	now := time.Now()
	if rotation == 24 {
		return now.Format("2006-01-02") // 按天切割
	} else if rotation == 1 {
		return now.Format("2006-01-02-15") // 按小时切割
	}
	// 你可以根据需要增加更多切割规则
	return now.Format("2006-01-02")
}

func (w *DateSplitWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	nowKey := w.currentTimeKey(w.rotation)
	if w.file == nil || w.currentDate != nowKey {
		_ = w.rotateFile(nowKey)
		w.currentDate = nowKey
	}
	return w.file.Write(p)
}

func (w *DateSplitWriter) cleanExpiredLogs(ageDays int64) {
	if !w.cleanMutex.TryLock() {
		// 已有清理在执行，跳过本次
		return
	}
	defer w.cleanMutex.Unlock()

	cutoff := time.Now().AddDate(0, 0, int(-ageDays))
	files, _ := os.ReadDir(w.logDir)
	for _, file := range files {
		if file.IsDir() {
			dirTime, err := time.Parse("2006-01-02", file.Name())
			if err == nil && dirTime.Before(cutoff) {
				os.RemoveAll(filepath.Join(w.logDir, file.Name()))
			}
		}
	}
}

func (w *DateSplitWriter) StartCleaner(ageDays int64, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.cleanExpiredLogs(ageDays)
			}
		}
	}()
}
