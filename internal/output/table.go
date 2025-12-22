package output

import (
	"fmt"

	"github.com/andinianst93/systemd-monitoring/internal/models"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorWhite  = "\033[37m"
)

// PrintTable prints services in a formatted table
func PrintTable(serviceList *models.ServiceList) {
	// 1. Print header
	printHeader()

	// 2. Print services
	for _, service := range serviceList.Services {
		// Get color based on status
		color := colorizeStatus(service.Status)

		// Format: Name (20 chars), Status with icon (12 chars), ActiveState (8 chars), Uptime (17 chars)
		fmt.Printf("║ %-20s │ %s%-12s%s │ %-8s │ %-17s ║\n",
			truncateString(service.Name, 20),
			color,
			service.GetStatusIcon()+" "+string(service.Status),
			ColorReset,
			service.ActiveState,
			service.GetUptimeString())
	}

	// 3. Print footer with summary
	printFooter(serviceList)
}

func printHeader() {
	// Print box drawing characters untuk header
	// Example:
	// ╔══════════════════════════════════════════════════════════╗
	// ║           SYSTEMD SERVICE MONITOR                        ║
	// ╠══════════════════════════════════════════════════════════╣
	// ║ Service          │ Status    │ Active  │ Uptime          ║
	// ╠══════════════════════════════════════════════════════════╣

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              SYSTEMD SERVICE MONITOR                         ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ Service          │ Status      │ Active  │ Uptime           ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
}

func printFooter(sl *models.ServiceList) {
	// Print separator dan summary
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Total: %d  │ Running: %d  │ Failed: %d  │ Stopped: %d        ║\n",
		sl.Total, sl.Running, sl.Failed, sl.Stopped)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

// colorizeStatus returns ANSI color code based on service status
func colorizeStatus(status models.ServiceStatus) string {
	switch status {
	case models.StatusRunning:
		return ColorGreen
	case models.StatusFailed:
		return ColorRed
	case models.StatusStopped:
		return ColorYellow
	default:
		return ColorWhite
	}
}

// truncateString truncates a string to maxLen and adds "..." if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// PrintService prints a single service info
func PrintService(service *models.ServiceInfo) {
	color := colorizeStatus(service.Status)

	fmt.Printf("%s[%s]%s %s - %s (%s)\n",
		color,
		service.GetStatusIcon(),
		ColorReset,
		service.Name,
		service.Status,
		service.ActiveState)

	if service.PID > 0 {
		fmt.Printf("  PID: %d\n", service.PID)
	}
	if service.MemoryUsage != "" {
		fmt.Printf("  Memory: %s\n", service.MemoryUsage)
	}
	if service.Uptime > 0 {
		fmt.Printf("  Uptime: %s\n", service.GetUptimeString())
	}
	fmt.Println()
}
