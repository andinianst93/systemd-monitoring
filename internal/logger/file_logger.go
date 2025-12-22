package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileLogger handles writing logs to a file
type FileLogger struct {
	filepath string
	file     *os.File
}

// NewFileLogger creates a new file logger
func NewFileLogger(filePath string) (*FileLogger, error) {
	// 1. Create directory if not exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 2. Open file using os.OpenFile()
	//    flags: os.O_APPEND | os.O_CREATE | os.O_WRONLY
	//    mode: 0644
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// 3. Return &FileLogger{filepath: filePath, file: file}
	return &FileLogger{filepath: filePath, file: file}, nil
}

// WriteLog writes a log message with timestamp
func (fl *FileLogger) WriteLog(message string) error {
	// 1. Format: timestamp := time.Now().Format("2006-01-02 15:04:05")
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// 2. Create log line: line := fmt.Sprintf("[%s] %s\n", timestamp, message)
	line := fmt.Sprintf("[%s] %s\n", timestamp, message)
	_, err := fl.file.WriteString(line)
	if err != nil {
		return err
	}
	return nil
}

// WriteServiceStatus writes service status to log
func (fl *FileLogger) WriteServiceStatus(serviceName string, status string) error {
	// Create message dan call WriteLog()
	// Format: "Service nginx.service is running"
	message := fmt.Sprintf("Service %s is %s", serviceName, status)
	err := fl.WriteLog(message)
	if err != nil {
		return err
	}
	return nil
}

// Error logs an error message
func (fl *FileLogger) Error(err error) error {
	message := fmt.Sprintf("ERROR: %v", err)
	return fl.WriteLog(message)
}

// Info logs an info message
func (fl *FileLogger) Info(message string) error {
	logMessage := fmt.Sprintf("INFO: %s", message)
	return fl.WriteLog(logMessage)
}

// Close closes the log file
func (fl *FileLogger) Close() error {
	if fl.file != nil {
		return fl.file.Close()
	}
	return nil
}
