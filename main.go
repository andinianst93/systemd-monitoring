package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andinianst93/systemd-monitoring/internal/logger"
	"github.com/andinianst93/systemd-monitoring/internal/models"
	"github.com/andinianst93/systemd-monitoring/internal/output"
	"github.com/andinianst93/systemd-monitoring/internal/systemd"
)

func main() {
	// Parse command dan route ke handler

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		handleList()
	case "check":
		handleCheck()
	case "monitor":
		handleMonitor()
	case "logs":
		handleLogs()
	case "write-log":
		handleWriteLog()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleList() {
	// TUGAS KAMU: Implement list command
	// 1. Parse flags
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	statusFilter := listCmd.String("status", "all", "Filter by status (running/failed/stopped/all)")
	outputFormat := listCmd.String("output", "table", "Output format (table/json)")
	useSudo := listCmd.Bool("sudo", false, "Use sudo for systemctl")

	listCmd.Parse(os.Args[2:])

	// 2. Create client
	client := systemd.NewClient(*useSudo)

	// 3. Get services
	serviceList, err := client.ListServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// 4. Filter if needed
	if *statusFilter != "all" {
		// Convert string to ServiceStatus
		var filterStatus models.ServiceStatus
		switch *statusFilter {
		case "running":
			filterStatus = models.StatusRunning
		case "failed":
			filterStatus = models.StatusFailed
		case "stopped":
			filterStatus = models.StatusStopped
		default:
			fmt.Fprintf(os.Stderr, "Invalid status filter: %s\n", *statusFilter)
			os.Exit(2)
		}

		// Filter serviceList.Services
		filteredServices := make([]*models.ServiceInfo, 0)
		for _, service := range serviceList.Services {
			if service.Status == filterStatus {
				filteredServices = append(filteredServices, service)
			}
		}
		serviceList.Services = filteredServices
	}

	// 5. Print output
	// if *outputFormat == "json":
	//    output.PrintJSON(serviceList)
	// else:
	//    output.PrintTable(serviceList)

	if *outputFormat == "json" {
		output.PrintJSON(serviceList)
	} else {
		output.PrintTable(serviceList)
	}

	// 6. Exit with code 1 if any failures
	// if serviceList.HasFailures():
	//    os.Exit(1)

	if serviceList.HasFailures() {
		os.Exit(1)
	}
}

func handleCheck() {
	// 1. Parse flags
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	useSudo := checkCmd.Bool("sudo", false, "Use sudo")

	checkCmd.Parse(os.Args[2:])

	// 2. Get service names from remaining args
	serviceNames := checkCmd.Args()
	if len(serviceNames) == 0 {
		fmt.Println("Error: No services specified")
		os.Exit(1)
	}

	// 3. Check each service
	client := systemd.NewClient(*useSudo)
	hasFailures := false
	for _, name := range serviceNames {
		service, err := client.GetServiceStatus(name)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		if service.IsFailed() {
			hasFailures = true
		}
		output.PrintService(service)
	}

	// 4. Exit with code 1 if any failed
	// if hasFailures: os.Exit(1)
	if hasFailures {
		os.Exit(1)
	}
}

func handleMonitor() {
	// 1. Parse flags
	monitorCmd := flag.NewFlagSet("monitor", flag.ExitOnError)
	services := monitorCmd.String("services", "", "Comma-separated service names")
	interval := monitorCmd.Duration("interval", 30*time.Second, "Check interval")
	logFile := monitorCmd.String("log-file", "logs/monitor.log", "Log file path")
	useSudo := monitorCmd.Bool("sudo", false, "Use sudo")

	monitorCmd.Parse(os.Args[2:])

	// 2. Validate services parameter
	if *services == "" {
		fmt.Println("Error: --services parameter is required")
		os.Exit(1)
	}

	// 3. Parse service list
	serviceList := strings.Split(*services, ",")

	// 4. Create logger
	fileLogger, err := logger.NewFileLogger(*logFile)
	if err != nil {
		fmt.Println("Error creating logger:", err)
		os.Exit(1)
	}
	defer fileLogger.Close()

	// 5. Create client
	client := systemd.NewClient(*useSudo)

	// 6. Create ticker
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	fmt.Println("Monitoring services. Press Ctrl+C to stop...")

	// 7. Loop
	for range ticker.C {
		fmt.Println("\n--- Checking services ---")

		// Check services
		for _, serviceName := range serviceList {
			serviceName = strings.TrimSpace(serviceName)
			if serviceName == "" {
				continue
			}

			service, err := client.GetServiceStatus(serviceName)
			if err != nil {
				fileLogger.Error(err)
				fmt.Printf("Error checking %s: %v\n", serviceName, err)
				continue
			}

			// Log to file
			fileLogger.WriteServiceStatus(service.Name, string(service.Status))

			// Print to console
			output.PrintService(service)
		}

		// Log summary
		fileLogger.Info(fmt.Sprintf("Checked %d services", len(serviceList)))
	}
}

func handleLogs() {
	// Parse flags
	logsCmd := flag.NewFlagSet("logs", flag.ExitOnError)
	lines := logsCmd.Int("lines", 50, "Number of lines to show")
	follow := logsCmd.Bool("follow", false, "Follow log output (like -f)")
	since := logsCmd.String("since", "", "Show logs since (e.g., '1 hour ago', 'today')")
	until := logsCmd.String("until", "", "Show logs until")
	priority := logsCmd.String("priority", "", "Filter by priority (emerg, alert, crit, err, warning, notice, info, debug)")
	grep := logsCmd.String("grep", "", "Filter logs by pattern")
	useSudo := logsCmd.Bool("sudo", false, "Use sudo")

	logsCmd.Parse(os.Args[2:])

	// Get service name
	args := logsCmd.Args()
	if len(args) == 0 {
		fmt.Println("Error: No service specified")
		fmt.Println("\nUsage: systemd-monitor logs <service> [options]")
		os.Exit(1)
	}

	serviceName := args[0]

	// Create client
	client := systemd.NewClient(*useSudo)

	// Create log options
	opts := &models.LogOptions{
		Lines:    *lines,
		Follow:   *follow,
		Since:    *since,
		Until:    *until,
		Priority: *priority,
		Grep:     *grep,
	}

	// Follow mode (real-time)
	if *follow {
		fmt.Printf("Following logs for %s (Ctrl+C to stop)...\n\n", serviceName)

		logChan, errChan, err := client.GetServiceLogsStream(serviceName, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		// Print logs as they come
		for {
			select {
			case entry, ok := <-logChan:
				if !ok {
					return
				}
				printLogEntry(entry)
			case err, ok := <-errChan:
				if ok && err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
				return
			}
		}
	} else {
		// One-time fetch
		entries, err := client.GetServiceLogs(serviceName, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		if len(entries) == 0 {
			fmt.Println("No logs found")
			return
		}

		fmt.Printf("Showing %d log entries for %s:\n\n", len(entries), serviceName)

		// Print all entries
		for _, entry := range entries {
			printLogEntry(entry)
		}
	}
}

func printLogEntry(entry *models.LogEntry) {
	color := entry.GetColorForLevel()
	icon := entry.GetLevelIcon()
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

	fmt.Printf("%s[%s]%s %s %s %s\n",
		color,
		timestamp,
		"\033[0m",
		icon,
		entry.Level,
		entry.Message)
}

func printUsage() {
	fmt.Println("Usage: systemd-monitor <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  list              List all systemd services")
	fmt.Println("  check <services>  Check specific services")
	fmt.Println("  monitor           Monitor services continuously")
	fmt.Println("  logs <service>    View service logs")
	fmt.Println("  write-log         Write message to systemd journal")
	fmt.Println("\nList Options:")
	fmt.Println("  --status string   Filter by status (running/failed/stopped/all)")
	fmt.Println("  --output string   Output format (table/json)")
	fmt.Println("  --sudo            Use sudo for systemctl")
	fmt.Println("\nCheck Options:")
	fmt.Println("  --output string   Output format (table/json)")
	fmt.Println("  --sudo            Use sudo")
	fmt.Println("\nMonitor Options:")
	fmt.Println("  --services string Comma-separated service names")
	fmt.Println("  --interval duration Check interval (default 30s)")
	fmt.Println("  --log-file string   Log file path")
	fmt.Println("  --sudo            Use sudo")
	fmt.Println("\nLogs Options:")
	fmt.Println("  --lines int       Number of lines to show (default 50)")
	fmt.Println("  --follow          Follow log output in real-time")
	fmt.Println("  --since string    Show logs since (e.g., '1 hour ago', 'today')")
	fmt.Println("  --until string    Show logs until")
	fmt.Println("  --priority string Filter by priority (info, warning, error, etc)")
	fmt.Println("  --grep string     Filter logs by pattern")
	fmt.Println("  --sudo            Use sudo")
	fmt.Println("\nWrite-Log Options:")
	fmt.Println("  --message string  Message to write (required)")
	fmt.Println("  --priority string Priority level (info, warning, err, crit, debug)")
	fmt.Println("  --identifier string Application identifier (default: systemd-monitor)")
	fmt.Println("\nExamples:")
	fmt.Println("  systemd-monitor list")
	fmt.Println("  systemd-monitor list --status running --output json")
	fmt.Println("  systemd-monitor check nginx mysql redis")
	fmt.Println("  systemd-monitor monitor --services nginx,mysql --interval 1m")
	fmt.Println("  systemd-monitor logs clash")
	fmt.Println("  systemd-monitor logs clash --follow")
	fmt.Println("  systemd-monitor logs clash --lines 100 --since '1 hour ago'")
	fmt.Println("  systemd-monitor logs nginx --grep error --priority err")
	fmt.Println("  systemd-monitor write-log --message 'Service started' --priority info")
	fmt.Println("  systemd-monitor write-log --message 'Critical error' --priority crit")
}

func handleWriteLog() {
	// Parse flags
	writeCmd := flag.NewFlagSet("write-log", flag.ExitOnError)
	message := writeCmd.String("message", "", "Message to write to journal (required)")
	priority := writeCmd.String("priority", "info", "Priority level (info, warning, err, crit, debug)")
	identifier := writeCmd.String("identifier", "systemd-monitor", "Application identifier")

	writeCmd.Parse(os.Args[2:])

	// Validate message
	if *message == "" {
		fmt.Println("Error: --message is required")
		fmt.Println("\nUsage: systemd-monitor write-log --message <text> [options]")
		os.Exit(1)
	}

	// Check if systemd-cat is available
	if !logger.IsJournalAvailable() {
		fmt.Println("Error: systemd-cat not found. Make sure systemd is installed.")
		os.Exit(2)
	}

	// Create journal logger
	journalLogger := logger.NewJournalLogger(*identifier)

	// Write to journal
	err := journalLogger.WriteToJournal(*message, *priority)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to journal: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("✅ Message written to systemd journal\n")
	fmt.Printf("   Priority: %s\n", *priority)
	fmt.Printf("   Identifier: %s\n", *identifier)
	fmt.Printf("\nView with: journalctl -t %s -n 10\n", *identifier)
}
