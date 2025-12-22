package models

import (
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	ServiceName string    `json:"service_name"`
	Message     string    `json:"message"`
	Level       string    `json:"level"` // e.g., "info", "warning", "error", "critical"
}

func NewLogEntry(serviceName, message string) *LogEntry {
	return &LogEntry{
		Timestamp:   time.Now(),
		ServiceName: serviceName,
		Message:     message,
		Level:       "info",
	}
}

// GetColorForLevel returns ANSI color code for log level
func (l *LogEntry) GetColorForLevel() string {
	upper := strings.ToUpper(l.Level)
	switch upper {
	case "EMERG", "ALERT", "CRIT", "ERROR", "ERR", "CRITICAL":
		return "\033[31m" // Red
	case "WARN", "WARNING":
		return "\033[33m" // Yellow
	case "NOTICE":
		return "\033[36m" // Cyan
	case "INFO":
		return "\033[32m" // Green
	case "DEBUG":
		return "\033[37m" // White
	default:
		return "\033[0m" // Reset
	}
}

// GetLevelIcon returns emoji icon for log level
func (l *LogEntry) GetLevelIcon() string {
	upper := strings.ToUpper(l.Level)
	switch upper {
	case "EMERG", "ALERT", "CRIT", "ERROR", "ERR", "CRITICAL":
		return "❌"
	case "WARN", "WARNING":
		return "⚠️"
	case "NOTICE":
		return "ℹ️"
	case "INFO":
		return "✅"
	case "DEBUG":
		return "🔍"
	default:
		return "📝"
	}
}

type LogBuffer struct {
	Entries []*LogEntry `json:"entries"`
	MaxSize int         `json:"max_size"`
}

// Add adds a log entry to the buffer
func (lb *LogBuffer) Add(entry *LogEntry) {
	lb.Entries = append(lb.Entries, entry)

	// Keep buffer size under MaxSize
	if lb.MaxSize > 0 && len(lb.Entries) > lb.MaxSize {
		lb.Entries = lb.Entries[len(lb.Entries)-lb.MaxSize:]
	}
}

// Clear clears all entries
func (lb *LogBuffer) Clear() {
	lb.Entries = []*LogEntry{}
}

// GetLatest returns the latest N entries
func (lb *LogBuffer) GetLatest(n int) []*LogEntry {
	if n <= 0 || n >= len(lb.Entries) {
		return lb.Entries
	}
	return lb.Entries[len(lb.Entries)-n:]
}

// LogOptions configures log retrieval
type LogOptions struct {
	Lines    int    `json:"lines"`    // Number of lines to retrieve
	Follow   bool   `json:"follow"`   // Follow logs (like -f)
	Since    string `json:"since"`    // Time since (e.g., "1 hour ago", "today")
	Until    string `json:"until"`    // Time until
	Priority string `json:"priority"` // Log priority (emerg, alert, crit, err, warning, notice, info, debug)
	Grep     string `json:"grep"`     // Filter logs by pattern
}

// NewLogOptions creates default log options
func NewLogOptions() *LogOptions {
	return &LogOptions{
		Lines:    50,
		Follow:   false,
		Since:    "",
		Until:    "",
		Priority: "",
		Grep:     "",
	}
}
