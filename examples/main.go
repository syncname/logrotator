package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotationInterval int

const (
	Daily RotationInterval = iota
	Weekly
	Monthly
	Yearly
)

type Rotator struct {
	baseDir      string
	fileName     string
	currentFile  *os.File
	maxSize      int64
	interval     RotationInterval
	dirLayout    string
	fileLayout   string
	nextRotation time.Time
	currentSize  int64
	currentPath  string
	mu           sync.Mutex
}

func NewRotator(baseDir, fileName string, maxSizeMB int, interval RotationInterval) (*Rotator, error) {
	r := &Rotator{
		baseDir:    baseDir,
		fileName:   fileName,
		maxSize:    int64(maxSizeMB) * 1024 * 1024,
		interval:   interval,
		dirLayout:  "2006/01/02",      // Год/Месяц/День
		fileLayout: "20060102-150405", // Таймстамп для файлов
	}

	if err := r.init(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Rotator) init() error {
	// Формируем путь с текущей датой
	dateDir := time.Now().Format(r.dirLayout)
	fullDir := filepath.Join(r.baseDir, dateDir)
	r.currentPath = filepath.Join(fullDir, r.fileName)

	// Создаем директорию
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return err
	}

	// Открываем файл
	f, err := os.OpenFile(r.currentPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Получаем информацию о файле
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	// Инициализируем состояние
	r.currentFile = f
	r.currentSize = fi.Size()
	r.scheduleNextRotation()
	return nil
}

func (r *Rotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.shouldRotate(len(p)) {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = r.currentFile.Write(p)
	r.currentSize += int64(n)
	return
}

func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.currentFile.Close()
}

func (r *Rotator) shouldRotate(bytesToWrite int) bool {
	now := time.Now()

	// Проверка по размеру
	if r.currentSize+int64(bytesToWrite) > r.maxSize {
		return true
	}

	// Проверка по времени
	return now.After(r.nextRotation)
}

func (r *Rotator) scheduleNextRotation() {
	now := time.Now()
	switch r.interval {
	case Daily:
		r.nextRotation = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	case Weekly:
		daysUntilMonday := (7 - int(now.Weekday()) + 1) % 7
		r.nextRotation = time.Date(now.Year(), now.Month(), now.Day()+daysUntilMonday, 0, 0, 0, 0, now.Location())
	case Monthly:
		r.nextRotation = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	case Yearly:
		r.nextRotation = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
	}
}

func (r *Rotator) rotate() error {
	// Закрываем текущий файл
	if err := r.currentFile.Close(); err != nil {
		return err
	}

	// Создаем новое имя файла с таймстампом
	timestamp := time.Now().Format(r.fileLayout)
	newFileName := r.fileName + "." + timestamp
	oldPath := r.currentPath
	newPath := filepath.Join(filepath.Dir(oldPath), newFileName)

	// Переименовываем файл
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	// Инициализируем новый файл
	return r.init()
}

func main() {
	rotator, err := NewRotator("logs", "app.log", 10, Daily)
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		_, err = rotator.Write([]byte(fmt.Sprintf("Log entry %d\n", i)))
		if err != nil {
			panic(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
