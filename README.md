# LogRotator

LogRotator is a custom package for log rotation that organizes logs into hierarchical folders (e.g., year/month/day) based on the selected strategy. This package is compatible with various popular logging libraries such as `zerolog`, `logrus`, `zap`, `slog`, and Go's standard logging.

## Table of Contents

1. [Folder Structure](#folder-structure)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Usage with Different Loggers](#usage-with-different-loggers)
    - [Zerolog](#zerolog)
    - [Logrus](#logrus)
    - [Zap](#zap)
    - [Slog](#slog)
    - [Go Standard Logs](#go-standard-logs)
5. [Rotation Strategies](#rotation-strategies)
6. [License](#license)

---

## Rotation Strategies

1. DailyStrategy: Creates folders in the format year/month/day.
2. WeeklyStrategy: Creates folders in the format year-Wweek.
3. MonthlyStrategy: Creates folders in the format year/month.
4. YearlyStrategy: Creates folders in the format year.


## Folder Structure

Depending on the chosen rotation strategy, the following folder structures will be created:

- **DailyStrategy**:  
- ./logs/2023/10/05/log_1234567890.log

- **WeeklyStrategy**:  
- ./logs/2023-W40/log_1234567890.log


- **MonthlyStrategy**:  
  ./logs/2023/10/log_1234567890.log


- **YearlyStrategy**:  
  ./logs/2023/log_1234567890.log


---

## Installation

To use this package, add it to your project:

```bash
go get github.com/syncname/logrotator
```

# Configuration

To get started with LogRotator, create an instance of LogRotator by specifying the base path, rotation interval, maximum file size, and rotation strategy.

```go
package main

import (
	"github.com/syncname/logrotator"
	"log"
	"time"
)

func main() {
	basePath := "./logs" // Path to the logs folder
	interval := 24 * time.Hour
	maxSize := int64(10 << 20) // 10 MB
	strategy := logrotator.DailyStrategy

	rotator, err := logrotator.NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}
	defer rotator.CurrentFile().Close()
}
```

## Usage with Different Loggers
### Zerolog

To integrate with zerolog, create a custom io.Writer that uses your LogRotator.

```go
package main

import (
    "log"
    "time"
    "github.com/rs/zerolog"
    "github.com/syncname/logrotator"
)



func main() {
	basePath := "./logs" // Путь к папке для логов
	interval := 24 * time.Hour
	maxSize := int64(10 << 20) // 10 MB
	strategy := logrotator.DailyStrategy

	rotator, err := logrotator.NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	// Создаем кастомный Writer
	writer := &logrotator.ZeroLogWriter{}
	writer.SetRotator(rotator)

	// Настройка zerolog
	logger := zerolog.New(writer).With().Any("loggerName", "zerolog").Timestamp().Logger()

	// Пример использования zerolog
	logger.Info().Msg("This is an info message")
	logger.Warn().Msg("This is a warning message")
	logger.Error().Msg("This is an error message")

	// Закрытие текущего файла лога при завершении программы
	defer func() {
		if err := rotator.CurrentFile().Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()
}
```
### Logrus

To integrate with logrus, use the same approach with a custom io.Writer.
```go
package main

import (
	"log"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/syncname/logrotator"
)

func main() {
	// Настройка LogRotator
	basePath := "./logs_logrus" // Путь к папке для логов
	interval := 24 * time.Hour
	maxSize := int64(10 << 20) // 10 MB
	strategy := logrotator.DailyStrategy

	rotator, err := logrotator.NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	// Создаем кастомный Writer
	writer := &logrotator.LogrusWriter{}
	writer.SetRotator(rotator)
	logrus.SetOutput(writer)                     // Устанавливаем наш кастомный Writer
	logrus.SetFormatter(&logrus.JSONFormatter{}) // Используем JSON формат для логов
	logrus.SetLevel(logrus.DebugLevel)           // Устанавливаем уровень логирования

	// Пример использования logrus
	logrus.Info("This is an info message")
	logrus.Warn("This is a warning message")
	logrus.Error("This is an error message")

	// Закрытие текущего файла лога при завершении программы
	defer func() {
		if err := rotator.CurrentFile().Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()
}

```

### Zap
To integrate with zap, use the ZapCoreAdapter that you already implemented.

```go
package main

import (
    "log"
    "time"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/syncname/logrotator"
)

func main() {
    basePath := "./logs"
    interval := 24 * time.Hour
    maxSize := int64(10 << 20)
    strategy := logrotator.DailyStrategy

    rotator, err := logrotator.NewLogRotator(basePath, interval, maxSize, strategy)
    if err != nil {
        log.Fatalf("Failed to create log rotator: %v", err)
    }

    zapCore := logrotator.NewZapCore(rotator, zapcore.InfoLevel)
    logger := zap.New(zapCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

    logger.Info("This is an info message")
    logger.Warn("This is a warning message")
    logger.Error("This is an error message")

    defer rotator.CurrentFile().Close()
}
```

### Slog
To integrate with slog, use a custom io.Writer.

```go
package main

import (
    "log"
    "time"

	"log/slog"
    "github.com/syncname/logrotator"
)

func main() {
	basePath := "./logs"
	interval := 24 * time.Hour
	maxSize := int64(10 << 20)
	strategy := logrotator.DailyStrategy

	rotator, err := logrotator.NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	writer := &logrotator.SlogWriter{}
	writer.SetRotator(rotator)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	defer rotator.CurrentFile().Close()
}
```

### Go Standard Logs

To use with Go's standard logs, simply pass LogRotator as an io.Writer.

```go
package main

import (
	"github.com/syncname/logrotator"
	"log"
	"time"
)

func main() {
	basePath := "./logs"
	interval := 24 * time.Hour
	maxSize := int64(10 << 20)
	strategy := logrotator.DailyStrategy

	rotator, err := logrotator.NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	writer := &logrotator.StdWriter{}
	writer.SetRotator(rotator)
	log.SetOutput(writer)

	log.Println("This is a standard log message")

	defer rotator.CurrentFile().Close()
}

```


