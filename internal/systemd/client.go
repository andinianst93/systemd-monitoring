package systemd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/andinianst93/systemd-monitoring/internal/models"
)

type Client struct {
	useSudo bool
}

func NewClient(useSudo bool) *Client {
	return &Client{
		useSudo: useSudo,
	}
}

func (c *Client) ListServices() (*models.ServiceList, error) {
	// 1. Build command: systemctl list-units --type=service --all --no-pager
	cmd := c.buildCommand("systemctl", "list-units", "--type=service", "--all", "--no-pager")

	// 2. Execute command using exec.Command()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute systemctl: %w", err)
	}

	// 3. Parse output line by line
	serviceList := models.NewServiceList()
	lines := strings.Split(string(output), "\n")

	// 4. For each service line:
	//    - Extract: name, load-state, active-state, sub-state
	//    - Create ServiceInfo
	//    - Add to ServiceList

	for _, line := range lines {
		// Skip empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip header lines
		if strings.Contains(line, "UNIT") || strings.Contains(line, "●") {
			continue
		}

		// Skip summary line
		if strings.Contains(line, "loaded units") {
			break
		}

		// Parse setiap field dari line
		// Format output systemctl:
		// UNIT                    LOAD   ACTIVE SUB     DESCRIPTION
		// ssh.service             loaded active running OpenSSH server
		fields := strings.Fields(line)

		// There must 4 fields: name, load, active, sub
		if len(fields) < 4 {
			continue
		}

		// Extract fields
		name := fields[0]
		// loadState := fields[1]
		activeState := fields[2]
		subState := fields[3]

		// Create serviceinfo for each services
		serviceInfo := models.NewServiceInfo(name)
		serviceInfo.ActiveState = activeState
		serviceInfo.SubState = subState
		serviceInfo.Status = c.parseStatus(activeState, subState)

		serviceList.AddService(serviceInfo)

	}
	// 5. Return ServiceList

	return serviceList, nil
}

func (c *Client) GetServiceStatus(serviceName string) (*models.ServiceInfo, error) {
	// PSEUDOCODE:
	// 1. Add .service suffix if not present
	if !strings.HasSuffix(serviceName, ".service") {
		serviceName = serviceName + ".service"
	}

	// 2.  Build and execute command
	cmd := c.buildCommand("systemctl", "show", serviceName, "--no-pager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get service status for %s: %w", serviceName, err)
	}

	// Create ServiceInfo
	serviceInfo := models.NewServiceInfo(serviceName)

	// Parse key=value output line by line
	lines := strings.Split(string(output), "\n")

	// Variables to store parsed data
	var (
		activeState     string
		subState        string
		activeEnterTime string
	)

	// STEP 5: Loop dan extract values
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split by first "=" to get key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// STEP 6: Extract fields yang kita butuhkan
		switch key {
		case "ActiveState":
			activeState = value
			serviceInfo.ActiveState = value

		case "SubState":
			subState = value
			serviceInfo.SubState = value

		case "MainPID":
			// Convert string to int
			if pid, err := strconv.Atoi(value); err == nil {
				serviceInfo.PID = pid
			}

		case "MemoryCurrent":
			// Convert bytes to human readable (MB)
			if memBytes, err := strconv.ParseInt(value, 10, 64); err == nil {
				serviceInfo.MemoryUsage = formatMemory(memBytes)
			} else {
				serviceInfo.MemoryUsage = value
			}

		case "ActiveEnterTimestamp":
			activeEnterTime = value
			// We'll calculate uptime later
		}
	}

	serviceInfo.Status = c.parseStatus(activeState, subState)

	if activeEnterTime != "" && activeEnterTime != "0" {
		uptime, err := calculateUptime(activeEnterTime)
		if err == nil {
			serviceInfo.Uptime = uptime
		}
	}

	serviceInfo.CheckedAt = time.Now()

	return serviceInfo, nil
}

func (c *Client) buildCommand(name string, args ...string) *exec.Cmd {
	// If useSudo, prepend "sudo"
	// Return exec.Command()
	if c.useSudo {
		allArgs := append([]string{name}, args...)
		return exec.Command("sudo", allArgs...)
	}
	return exec.Command(name, args...)
}

func (c *Client) parseStatus(activeState, subState string) models.ServiceStatus {
	// Convert systemd states to our ServiceStatus
	if activeState == "active" && subState == "running" {
		return models.StatusRunning
	}

	if activeState == "failed" {
		return models.StatusFailed
	}

	if activeState == "inactive" || subState == "dead" {
		return models.StatusStopped
	}

	return models.StatusUnknown
}

// formatMemory converts bytes to human readable format
func formatMemory(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	// KB, MB, GB, TB
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGT"[exp])
}

// calculateUptime calculates uptime from timestamp string
func calculateUptime(timestamp string) (time.Duration, error) {
	// ActiveEnterTimestamp format: "Mon 2024-01-15 10:30:45 WIB"
	// Try multiple time formats systemd might use
	formats := []string{
		"Mon 2006-01-02 15:04:05 MST",
		"Mon 2006-01-02 15:04:05 -0700",
		time.RFC3339,
		time.RFC1123,
	}

	var startTime time.Time
	var parseErr error

	for _, format := range formats {
		startTime, parseErr = time.Parse(format, timestamp)
		if parseErr == nil {
			break
		}
	}

	// If all formats failed, try Unix timestamp (microseconds)
	if parseErr != nil {
		if usec, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
			startTime = time.Unix(usec/1000000, 0)
			parseErr = nil
		}
	}

	if parseErr != nil {
		return 0, fmt.Errorf("failed to parse timestamp %s: %w", timestamp, parseErr)
	}

	// Calculate duration from start time until now
	uptime := time.Since(startTime)

	return uptime, nil
}

// GetServiceLogs retrieves logs from systemd journal
func (c *Client) GetServiceLogs(serviceName string, opts *models.LogOptions) ([]*models.LogEntry, error) {
	// Add .service suffix if not present
	if !strings.HasSuffix(serviceName, ".service") {
		serviceName = serviceName + ".service"
	}

	// Build journalctl command
	args := []string{"journalctl", "-u", serviceName, "--no-pager"}

	// Add options
	if opts != nil {
		if opts.Lines > 0 {
			args = append(args, "-n", fmt.Sprintf("%d", opts.Lines))
		}
		if opts.Since != "" {
			args = append(args, "--since", opts.Since)
		}
		if opts.Until != "" {
			args = append(args, "--until", opts.Until)
		}
		if opts.Priority != "" {
			args = append(args, "-p", opts.Priority)
		}
	}

	// Build command
	cmd := c.buildCommand(args[0], args[1:]...)

	// Execute
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for %s: %w", serviceName, err)
	}

	// Parse output
	entries := parseJournalOutput(string(output), serviceName)

	// Apply grep filter if specified
	if opts != nil && opts.Grep != "" {
		entries = filterLogs(entries, opts.Grep)
	}

	return entries, nil
}

// GetServiceLogsStream returns a channel for following logs in real-time
func (c *Client) GetServiceLogsStream(serviceName string, opts *models.LogOptions) (<-chan *models.LogEntry, <-chan error, error) {
	// Add .service suffix if not present
	if !strings.HasSuffix(serviceName, ".service") {
		serviceName = serviceName + ".service"
	}

	// Build journalctl command with -f (follow)
	args := []string{"journalctl", "-u", serviceName, "-f", "--no-pager"}

	// Add options
	if opts != nil {
		if opts.Lines > 0 {
			args = append(args, "-n", fmt.Sprintf("%d", opts.Lines))
		}
		if opts.Since != "" {
			args = append(args, "--since", opts.Since)
		}
		if opts.Priority != "" {
			args = append(args, "-p", opts.Priority)
		}
	}

	// Build command
	cmd := c.buildCommand(args[0], args[1:]...)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start journalctl: %w", err)
	}

	// Create channels
	logChan := make(chan *models.LogEntry, 100)
	errChan := make(chan error, 1)

	// Start goroutine to read logs
	go func() {
		defer close(logChan)
		defer close(errChan)
		defer cmd.Wait()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			entry := parseJournalLine(line, serviceName)

			// Apply grep filter if specified
			if opts != nil && opts.Grep != "" {
				if !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(opts.Grep)) {
					continue
				}
			}

			logChan <- entry
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
		}
	}()

	return logChan, errChan, nil
}

// parseJournalOutput parses journalctl output into LogEntry structs
func parseJournalOutput(output string, serviceName string) []*models.LogEntry {
	var entries []*models.LogEntry
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		entry := parseJournalLine(line, serviceName)
		entries = append(entries, entry)
	}

	return entries
}

// parseJournalLine parses a single journalctl line
func parseJournalLine(line string, serviceName string) *models.LogEntry {
	// Journal format examples:
	// Dec 22 19:32:42 msi clash[1397]: time="2025-12-22T19:32:42.832776698+08:00" level=info msg="Start initial configuration"
	// Dec 22 19:32:42 msi systemd[1]: Started clash.service - Clash daemon.

	entry := models.NewLogEntry(serviceName, line)

	// Try to extract timestamp from beginning (Dec 22 19:32:42)
	parts := strings.Fields(line)
	if len(parts) >= 3 {
		// Try to parse first 3 fields as date
		dateStr := strings.Join(parts[0:3], " ")
		// Add current year
		year := time.Now().Year()
		fullDateStr := fmt.Sprintf("%s %d", dateStr, year)

		if t, err := time.Parse("Jan 2 15:04:05 2006", fullDateStr); err == nil {
			entry.Timestamp = t
		}
	}

	// Try to find the actual message after the host/process info
	// Format: Dec 22 19:32:42 msi clash[1397]: MESSAGE_HERE
	colonIndex := strings.Index(line, "]: ")
	if colonIndex > 0 && colonIndex+3 < len(line) {
		entry.Message = line[colonIndex+3:]
	} else {
		// Try finding just ": " separator
		colonIndex = strings.Index(line, ": ")
		if colonIndex > 0 && colonIndex+2 < len(line) {
			entry.Message = line[colonIndex+2:]
		}
	}

	// Extract log level from message if present
	// Format 1: level=info msg="..."
	if strings.Contains(entry.Message, "level=") {
		levelStart := strings.Index(entry.Message, "level=")
		if levelStart != -1 {
			remaining := entry.Message[levelStart+6:]
			spaceIdx := strings.IndexAny(remaining, " \"")
			if spaceIdx > 0 {
				level := remaining[:spaceIdx]
				entry.Level = strings.ToUpper(level)
			}
		}
	}

	// Format 2: [INFO] message or INFO: message
	msgUpper := strings.ToUpper(entry.Message)
	levels := []string{"ERROR", "WARN", "WARNING", "INFO", "DEBUG", "CRITICAL", "FATAL"}
	for _, level := range levels {
		if strings.HasPrefix(msgUpper, "["+level+"]") || strings.HasPrefix(msgUpper, level+":") {
			entry.Level = level
			break
		}
	}

	return entry
}

// filterLogs filters log entries by pattern
func filterLogs(entries []*models.LogEntry, pattern string) []*models.LogEntry {
	var filtered []*models.LogEntry
	pattern = strings.ToLower(pattern)

	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Message), pattern) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}
