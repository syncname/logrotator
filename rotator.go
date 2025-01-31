package logrotator

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DailyStrategy   = "daily"
	WeeklyStrategy  = "weekly"
	MonthlyStrategy = "monthly"
	YearlyStrategy  = "yearly"
)

type LogRotator struct {
	mu               sync.Mutex
	basePath         string
	interval         time.Duration
	maxSize          int64
	currentFile      *os.File
	currentFileSize  int64
	rotationStrategy string
	lastRotationTime time.Time
}

// NewLogRotator создает новый экземпляр LogRotator
func NewLogRotator(basePath string, interval time.Duration, maxSize int64, strategy string) (*LogRotator, error) {
	rotator := &LogRotator{
		basePath:         basePath,
		interval:         interval,
		maxSize:          maxSize,
		rotationStrategy: strategy,
		lastRotationTime: time.Now(),
	}

	//create new log file
	if err := rotator.rotate(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (r *LogRotator) CurrentFile() *os.File {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.currentFile
}

// Write data to log file
func (r *LogRotator) Write(data []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//time log rotation
	if time.Since(r.lastRotationTime) >= r.interval {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	//size log rotation
	if r.currentFileSize+int64(len(data)) > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := r.currentFile.Write(data)
	if err != nil {
		return 0, err
	}
	r.currentFileSize += int64(n)

	return n, nil
}

func (r *LogRotator) rotate() error {

	if r.currentFile != nil {
		if err := r.currentFile.Close(); err != nil {
			return err
		}
	}

	folder := r.getRotationFolder()
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}
	filename := filepath.Join(folder, fmt.Sprintf("log_%d.log", time.Now().Unix()))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("create log file error: %w", err)
	}

	r.currentFile = file
	r.currentFileSize = 0
	r.lastRotationTime = time.Now()

	return nil
}

// getRotationFolder возвращает путь для текущей ротации в зависимости от стратегии
func (r *LogRotator) getRotationFolder() string {
	now := time.Now()
	switch r.rotationStrategy {
	case DailyStrategy:
		return filepath.Join(r.basePath, now.Format("2006"), now.Format("01"), now.Format("02"))
	case WeeklyStrategy:
		year, week := now.ISOWeek()
		return filepath.Join(r.basePath, fmt.Sprintf("%d-W%02d", year, week))
	case MonthlyStrategy:
		return filepath.Join(r.basePath, now.Format("2006"), now.Format("01"))
	case YearlyStrategy:
		return filepath.Join(r.basePath, now.Format("2006"))
	default:
		return r.basePath
	}
}
