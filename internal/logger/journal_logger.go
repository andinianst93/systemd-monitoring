package logger

import (
	"fmt"
	"os/exec"
	"strings"
)

// JournalLogger handles writing logs to systemd journal
type JournalLogger struct {
	identifier string // Application identifier for journal
}

// NewJournalLogger creates a new journal logger
func NewJournalLogger(identifier string) *JournalLogger {
	return &JournalLogger{
		identifier: identifier,
	}
}

// WriteToJournal writes a message to systemd journal using systemd-cat
func (jl *JournalLogger) WriteToJournal(message string, priority string) error {
	// Priority levels: emerg, alert, crit, err, warning, notice, info, debug
	// Default to info if not specified
	if priority == "" {
		priority = "info"
	}

	// Validate priority
	validPriorities := []string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"}
	isValid := false
	for _, p := range validPriorities {
		if priority == p {
			isValid = true
			break
		}
	}
	if !isValid {
		priority = "info"
	}

	// Use systemd-cat to write to journal
	// Format: echo "message" | systemd-cat -t identifier -p priority
	cmd := exec.Command("systemd-cat", "-t", jl.identifier, "-p", priority)
	cmd.Stdin = strings.NewReader(message)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to write to journal: %w (output: %s)", err, string(output))
	}

	return nil
}

// Info writes an info-level message to journal
func (jl *JournalLogger) Info(message string) error {
	return jl.WriteToJournal(message, "info")
}

// Warning writes a warning-level message to journal
func (jl *JournalLogger) Warning(message string) error {
	return jl.WriteToJournal(message, "warning")
}

// Error writes an error-level message to journal
func (jl *JournalLogger) Error(message string) error {
	return jl.WriteToJournal(message, "err")
}

// Critical writes a critical-level message to journal
func (jl *JournalLogger) Critical(message string) error {
	return jl.WriteToJournal(message, "crit")
}

// Debug writes a debug-level message to journal
func (jl *JournalLogger) Debug(message string) error {
	return jl.WriteToJournal(message, "debug")
}

// Notice writes a notice-level message to journal
func (jl *JournalLogger) Notice(message string) error {
	return jl.WriteToJournal(message, "notice")
}

// WriteServiceStatus writes service status change to journal
func (jl *JournalLogger) WriteServiceStatus(serviceName, status string) error {
	message := fmt.Sprintf("Service %s status: %s", serviceName, status)

	// Use appropriate priority based on status
	var priority string
	switch strings.ToLower(status) {
	case "failed":
		priority = "err"
	case "running":
		priority = "info"
	case "stopped":
		priority = "warning"
	default:
		priority = "notice"
	}

	return jl.WriteToJournal(message, priority)
}

// WriteMonitoringEvent writes a monitoring event to journal
func (jl *JournalLogger) WriteMonitoringEvent(event string, details string) error {
	message := fmt.Sprintf("[MONITORING] %s: %s", event, details)
	return jl.WriteToJournal(message, "info")
}

// WriteBulk writes multiple messages to journal
func (jl *JournalLogger) WriteBulk(messages []string, priority string) error {
	for _, msg := range messages {
		if err := jl.WriteToJournal(msg, priority); err != nil {
			return fmt.Errorf("failed to write bulk message '%s': %w", msg, err)
		}
	}
	return nil
}

// IsJournalAvailable checks if systemd-cat is available
func IsJournalAvailable() bool {
	_, err := exec.LookPath("systemd-cat")
	return err == nil
}
