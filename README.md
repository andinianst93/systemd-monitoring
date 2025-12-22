# 🔍 Systemd Service Monitor

A user-friendly CLI tool for monitoring and managing systemd services on Linux systems. Built with Go, featuring colorized output, real-time monitoring, and comprehensive logging capabilities.

---

## 📋 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Usage](#-usage)
  - [List Services](#1-list-services)
  - [Check Services](#2-check-specific-services)
  - [Monitor Services](#3-monitor-continuously)
  - [View Logs](#4-view-service-logs)
  - [Write Logs](#5-write-logs-to-journal)
- [Command Reference](#-command-reference)
- [Examples](#-examples)
- [Configuration](#-configuration)
- [Troubleshooting](#-troubleshooting)
- [Development](#-development)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

### Core Features
- 📊 **List All Services** - View all systemd services with status, uptime, and more
- 🔍 **Check Service Status** - Get detailed information about specific services
- 📡 **Real-time Monitoring** - Continuously monitor services with configurable intervals
- 📝 **Log Viewing** - Read and follow systemd journal logs with filtering
- ✍️ **Log Writing** - Write custom messages directly to systemd journal
- 🎨 **Colorized Output** - Beautiful, color-coded terminal output
- 📈 **Multiple Output Formats** - Table view or JSON export

### Advanced Features
- 🔐 **Sudo Support** - Optional sudo for privileged operations
- 📊 **Service Statistics** - Total, running, failed, and stopped counts
- 🎯 **Filtering** - Filter services by status (running/failed/stopped)
- ⏱️ **Uptime Tracking** - See how long services have been running
- 💾 **Memory Usage** - View memory consumption per service
- 🔴 **PID Display** - Process ID information
- 📋 **Log Export** - Save monitoring results to files

### Output Examples

**Table View:**
```
╔══════════════════════════════════════════════════════════════╗
║              SYSTEMD SERVICE MONITOR                         ║
╠══════════════════════════════════════════════════════════════╣
║ Service          │ Status      │ Active  │ Uptime           ║
╠══════════════════════════════════════════════════════════════╣
║ nginx.service    │ ✅ running  │ active  │ 2h 15m           ║
║ ssh.service      │ ✅ running  │ active  │ 5d 3h            ║
║ mysql.service    │ ❌ failed   │ failed  │ 0s               ║
╠══════════════════════════════════════════════════════════════╣
║ Total: 3  │ Running: 2  │ Failed: 1  │ Stopped: 0        ║
╚══════════════════════════════════════════════════════════════╝
```

**Service Details:**
```
✅ [✅] nginx.service - running (active)
  PID: 1234
  Memory: 45.6 MB
  Uptime: 2h 15m
```

**Log Output:**
```
[2025-12-22 19:32:42] ✅ INFO Started nginx.service
[2025-12-22 19:32:42] ✅ INFO nginx: configuration file loaded
[2025-12-22 19:31:20] ⚠️ WARN nginx: [warn] conflicting server name
[2025-12-22 19:30:15] ❌ ERROR nginx: connection timeout
```

---

## 🔧 Prerequisites

### System Requirements
- **Operating System:** Linux (any distribution with systemd)
- **Go:** Version 1.25
- **systemd:** Systemd init system
- **Permissions:** Root/sudo access for full functionality

## 📦 Installation

### Method 1: Clone and Build (Recommended)

```bash
# Clone the repository
git clone https://github.com/andinianst93/systemd-monitoring.git
cd systemd-monitoring

# Build the binary
go build -o bin/monitor main.go

# Make it executable
chmod +x bin/monitor

# Optional: Add to PATH
sudo cp bin/monitor /usr/local/bin/monitor
```

### Method 2: Direct Install

```bash
# Clone and install in one go
git clone https://github.com/andinianst93/systemd-monitoring.git
cd systemd-monitoring
go build -o bin/monitor main.go
sudo cp bin/monitor /usr/local/bin/monitor

# Verify installation
monitor --help
```

### Method 3: Development Mode

```bash
# Clone repository
git clone https://github.com/yourusername/systemd-monitoring.git
cd systemd-monitoring

# Run directly with go
go run main.go list
```

### Verify Installation

```bash
# Check if binary exists
./bin/monitor --help

# Or if installed globally
monitor --help

# Test basic functionality
./bin/monitor list
```

---

## 🚀 Quick Start

### 1. List All Services

```bash
# Basic list
sudo ./bin/monitor list

# Filter running services only
sudo ./bin/monitor list --status running

# Export to JSON
sudo ./bin/monitor list --output json > services.json
```

### 2. Check Service Status

```bash
# Check single service
sudo ./bin/monitor check nginx

# Check multiple services
sudo ./bin/monitor check nginx mysql redis

# Without sudo (if you're in systemd-journal group)
./bin/monitor check nginx
```

### 3. Monitor Continuously

```bash
# Monitor with 30s interval
sudo ./bin/monitor monitor --services nginx --interval 30s

# Monitor multiple services
sudo ./bin/monitor monitor --services nginx,mysql,redis --interval 10s

# Save to log file
sudo ./bin/monitor monitor --services nginx --interval 30s --log-file logs/monitor.log
```

### 4. View Logs

```bash
# View last 50 logs
sudo ./bin/monitor logs nginx

# Follow logs in real-time
sudo ./bin/monitor logs nginx --follow

# Filter by time
sudo ./bin/monitor logs nginx --since "1 hour ago"

# Search in logs
sudo ./bin/monitor logs nginx --grep "error"
```

### 5. Write Logs to Journal

```bash
# Write a simple message
./bin/monitor write-log --message "Application started"

# Write with priority level
./bin/monitor write-log --message "Critical error" --priority crit

# Write with custom identifier
./bin/monitor write-log --message "Deployment complete" --identifier my-app

# View your logs
journalctl -t systemd-monitor -n 20
```

---

## 📖 Usage

### 1. List Services

Display all systemd services with their current status.

**Syntax:**
```bash
./bin/monitor list [options]
```

**Options:**
- `--status <filter>` - Filter by status: `running`, `failed`, `stopped`, or `all` (default: `all`)
- `--output <format>` - Output format: `table` or `json` (default: `table`)
- `--sudo` - Use sudo for systemctl commands

**Examples:**

```bash
# List all services (table format)
sudo ./bin/monitor list

# List only running services
sudo ./bin/monitor list --status running

# List failed services
sudo ./bin/monitor list --status failed

# Export to JSON
sudo ./bin/monitor list --output json

# Combine filters
sudo ./bin/monitor list --status running --output json > running-services.json
```

**Output Fields:**
- **Service Name** - Name of the systemd service
- **Status** - Current status (✅ running, ❌ failed, ⏸️ stopped)
- **Active State** - Systemd active state
- **Uptime** - How long the service has been running

---

### 2. Check Specific Services

Get detailed information about one or more services.

**Syntax:**
```bash
./bin/monitor check <service1> [service2] [...] [options]
```

**Options:**
- `--sudo` - Use sudo for systemctl commands

**Examples:**

```bash
# Check single service
sudo ./bin/monitor check nginx

# Check multiple services
sudo ./bin/monitor check nginx mysql redis sshd

# Without .service suffix (auto-added)
sudo ./bin/monitor check nginx
```

**Output Information:**
- Service name and status
- PID (Process ID)
- Memory usage
- Uptime
- Active and sub states

**Exit Codes:**
- `0` - All services are healthy
- `1` - One or more services have failed

---

### 3. Monitor Continuously

Continuously monitor services at specified intervals.

**Syntax:**
```bash
./bin/monitor monitor --services <service1,service2,...> [options]
```

**Options:**
- `--services <list>` - Comma-separated list of services (required)
- `--interval <duration>` - Check interval (default: `30s`)
  - Examples: `10s`, `1m`, `5m`, `1h`
- `--log-file <path>` - Log file path (default: `logs/monitor.log`)
- `--sudo` - Use sudo for systemctl commands

**Examples:**

```bash
# Monitor single service, 30s interval
sudo ./bin/monitor monitor --services nginx --interval 30s

# Monitor multiple services
sudo ./bin/monitor monitor --services nginx,mysql,redis --interval 1m

# With custom log file
sudo ./bin/monitor monitor --services nginx --interval 10s --log-file /var/log/nginx-monitor.log

# Monitor every 5 seconds
sudo ./bin/monitor monitor --services redis --interval 5s
```

**Features:**
- Real-time status updates
- Automatic log file creation
- Service status changes logged
- Continuous until Ctrl+C
- Service statistics

**Log File Format:**
```
[2024-12-22 15:30:45] Service nginx.service is running
[2024-12-22 15:30:45] INFO: Checked 1 services
[2024-12-22 15:31:15] Service nginx.service is running
[2024-12-22 15:31:15] INFO: Checked 1 services
```

---

### 4. View Service Logs

Read and follow systemd journal logs for services.

**Syntax:**
```bash
./bin/monitor logs <service> [options]
```

**Options:**
- `--lines <n>` - Number of lines to show (default: `50`)
- `--follow` - Follow logs in real-time (like `tail -f`)
- `--since <time>` - Show logs since specific time
  - Examples: `"1 hour ago"`, `"today"`, `"2024-12-22 15:00"`
- `--until <time>` - Show logs until specific time
- `--priority <level>` - Filter by priority level
  - Levels: `emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`
- `--grep <pattern>` - Filter logs by search pattern (case-insensitive)
- `--sudo` - Use sudo for journalctl

**Examples:**

```bash
# View last 50 logs
sudo ./bin/monitor logs nginx

# View last 100 logs
sudo ./bin/monitor logs nginx --lines 100

# Follow logs in real-time
sudo ./bin/monitor logs nginx --follow

# Logs from last hour
sudo ./bin/monitor logs nginx --since "1 hour ago"

# Logs from today
sudo ./bin/monitor logs nginx --since "today"

# Only errors
sudo ./bin/monitor logs nginx --priority err

# Search for specific text
sudo ./bin/monitor logs nginx --grep "timeout"

# Combine filters
sudo ./bin/monitor logs nginx --since "1 hour ago" --grep "error" --priority err

# Follow with initial context
sudo ./bin/monitor logs nginx --follow --lines 20
```

**Log Levels & Colors:**
- 🔍 **DEBUG** (White) - Debug messages
- ✅ **INFO** (Green) - Informational messages
- ℹ️ **NOTICE** (Cyan) - Normal but significant
- ⚠️ **WARNING** (Yellow) - Warning conditions
- ❌ **ERROR** (Red) - Error conditions
- ❌ **CRITICAL** (Red) - Critical conditions

---

### 5. Write Logs to Journal

Write custom messages directly to the systemd journal for centralized logging.

**Syntax:**
```bash
./bin/monitor write-log --message <text> [options]
```

**Options:**
- `--message <text>` - Message to write (required)
- `--priority <level>` - Priority level: `info`, `warning`, `err`, `crit`, `debug` (default: `info`)
- `--identifier <name>` - Application identifier tag (default: `systemd-monitor`)

**Priority Levels:**

| Level | Description | Use Case |
|-------|-------------|----------|
| `emerg` | Emergency | System is unusable |
| `alert` | Alert | Action must be taken immediately |
| `crit` | Critical | Critical conditions |
| `err` | Error | Error conditions |
| `warning` | Warning | Warning conditions |
| `notice` | Notice | Normal but significant |
| `info` | Information | Informational messages (default) |
| `debug` | Debug | Debug-level messages |

**Examples:**

```bash
# Basic info message
./bin/monitor write-log --message "Service started successfully"

# Error message
./bin/monitor write-log --message "Database connection failed" --priority err

# Critical alert
./bin/monitor write-log --message "Disk space critically low" --priority crit

# Custom identifier
./bin/monitor write-log --message "User login" --identifier my-web-app

# Warning message
./bin/monitor write-log --message "High memory usage: 85%" --priority warning
```

**View Your Logs:**

```bash
# View logs by identifier
journalctl -t systemd-monitor -n 20

# View custom identifier
journalctl -t my-app -n 50

# Follow logs in real-time
journalctl -t systemd-monitor -f

# View by priority
journalctl -t systemd-monitor -p err

# View from specific time
journalctl -t systemd-monitor --since "1 hour ago"
```

**Integration Example (Shell Script):**

```bash
#!/bin/bash
# backup.sh - Backup script with logging

APP_NAME="backup-service"

# Log start
./bin/monitor write-log \
  --message "Backup started" \
  --priority info \
  --identifier $APP_NAME

# Perform backup
if rsync -av /data /backup; then
    ./bin/monitor write-log \
      --message "Backup completed successfully" \
      --priority info \
      --identifier $APP_NAME
else
    ./bin/monitor write-log \
      --message "Backup failed with exit code $?" \
      --priority err \
      --identifier $APP_NAME
    exit 1
fi
```

**Use Cases:**
- Application logging from scripts
- Custom monitoring alerts
- Integration with external systems
- Deployment tracking
- Security event logging
- Resource monitoring
- Automated task logging

**Technical Details:**
- Uses `systemd-cat` to write to journal
- Messages stored in systemd journal (`/var/log/journal/`)
- Automatically indexed by identifier and priority
- Logs persist across reboots (if journal persistence enabled)

---

## 📚 Command Reference

### Complete Command List

```bash
# LISTING COMMANDS
./bin/monitor list                                    # List all services
./bin/monitor list --status running                   # List running services only
./bin/monitor list --status failed                    # List failed services only
./bin/monitor list --output json                      # Export as JSON
./bin/monitor list --sudo                             # Use sudo

# CHECKING COMMANDS
./bin/monitor check <service>                         # Check single service
./bin/monitor check nginx mysql redis                 # Check multiple services
./bin/monitor check nginx --sudo                      # Check with sudo

# MONITORING COMMANDS
./bin/monitor monitor --services nginx                # Monitor with defaults (30s)
./bin/monitor monitor --services nginx,mysql          # Monitor multiple services
./bin/monitor monitor --services nginx --interval 10s # Custom interval
./bin/monitor monitor --services nginx --log-file ./monitor.log  # Custom log file
./bin/monitor monitor --services nginx --sudo         # Monitor with sudo

# LOGS COMMANDS
./bin/monitor logs nginx                              # View last 50 logs
./bin/monitor logs nginx --lines 100                  # View last 100 logs
./bin/monitor logs nginx --follow                     # Follow logs
./bin/monitor logs nginx --since "1 hour ago"        # Time filter
./bin/monitor logs nginx --priority err               # Priority filter
./bin/monitor logs nginx --grep "error"               # Text search
./bin/monitor logs nginx --follow --lines 20          # Follow with context
./bin/monitor logs nginx --sudo                       # Use sudo

# WRITE-LOG COMMANDS
./bin/monitor write-log --message "Service started"  # Write info message
./bin/monitor write-log --message "Error" --priority err  # Write error
./bin/monitor write-log --message "Alert" --priority crit # Write critical
./bin/monitor write-log --message "Event" --identifier my-app # Custom identifier
```

### Global Flags

- `--sudo` - Use sudo for privileged operations (available for all commands)
- `--help` - Show help message
- `-h` - Short form of --help

---

## 💡 Examples

### Example 1: Daily Health Check

```bash
#!/bin/bash
# daily-check.sh - Check critical services

echo "=== Daily Service Health Check ==="
echo "Date: $(date)"
echo ""

# Check critical services
./bin/monitor check nginx mysql redis postgresql

# Check for failed services
failed=$(./bin/monitor list --status failed --output json | jq '.Failed')

if [ "$failed" -gt 0 ]; then
    echo "⚠️  WARNING: $failed services have failed!"
    ./bin/monitor list --status failed
    exit 1
else
    echo "✅ All services are healthy"
    exit 0
fi
```

### Example 2: Monitor with Alerts

```bash
#!/bin/bash
# monitor-with-alerts.sh - Monitor and send alerts

SERVICE="nginx"
INTERVAL="30s"

while true; do
    status=$(./bin/monitor check $SERVICE 2>/dev/null | grep -o "running\|failed\|stopped")
    
    if [ "$status" = "failed" ]; then
        echo "⚠️  ALERT: $SERVICE has failed!"
        # Send email, Slack notification, etc.
        # mail -s "Service Alert" admin@example.com < /dev/null
    fi
    
    sleep 30
done
```

### Example 3: Export Service Inventory

```bash
#!/bin/bash
# export-inventory.sh - Create service inventory

DATE=$(date +%Y%m%d)
OUTFILE="service-inventory-$DATE.json"

echo "Exporting service inventory..."

./bin/monitor list --output json > "$OUTFILE"

echo "✅ Inventory saved to $OUTFILE"

# Optional: Upload to cloud storage
# aws s3 cp "$OUTFILE" s3://my-bucket/inventories/
```

### Example 4: Log Analysis

```bash
#!/bin/bash
# analyze-logs.sh - Analyze service logs for errors

SERVICE="nginx"
TIMEFRAME="1 hour ago"

echo "Analyzing $SERVICE logs from $TIMEFRAME..."

# Count errors
errors=$(sudo ./bin/monitor logs $SERVICE --since "$TIMEFRAME" --grep "error" 2>/dev/null | wc -l)

echo "Errors found: $errors"

if [ $errors -gt 10 ]; then
    echo "⚠️  High error rate detected!"
    echo "Recent errors:"
    sudo ./bin/monitor logs $SERVICE --since "$TIMEFRAME" --grep "error" --lines 10
fi
```

### Example 5: Automated Restart on Failure

```bash
#!/bin/bash
# auto-restart.sh - Restart failed services

SERVICES=("nginx" "mysql" "redis")

for service in "${SERVICES[@]}"; do
    status=$(./bin/monitor check $service 2>/dev/null | grep -o "failed")
    
    if [ "$status" = "failed" ]; then
        echo "🔄 Restarting $service..."
        sudo systemctl restart $service
        sleep 5
        
        # Verify restart
        ./bin/monitor check $service
    fi
done
```

### Example 6: Script Logging to Journal

```bash
#!/bin/bash
# monitor-disk.sh - Disk monitoring with journal logging

APP_NAME="disk-monitor"
THRESHOLD=80

usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')

# Log current usage
./bin/monitor write-log \
  --message "Disk usage check: ${usage}% used" \
  --priority info \
  --identifier $APP_NAME

if [ $usage -gt $THRESHOLD ]; then
    # Critical alert
    ./bin/monitor write-log \
      --message "Disk usage critical: ${usage}% used" \
      --priority crit \
      --identifier $APP_NAME
    
    # Send email alert
    echo "Disk usage: ${usage}%" | mail -s "CRITICAL: Disk Space" admin@example.com
else
    # Normal operation
    ./bin/monitor write-log \
      --message "Disk usage normal: ${usage}% used" \
      --priority info \
      --identifier $APP_NAME
fi

# View logs: journalctl -t disk-monitor -n 20
```

---

## ⚙️ Configuration

### Permissions Setup

For non-root users to access logs without sudo:

```bash
# Add user to systemd-journal group
sudo usermod -a -G systemd-journal $USER

# Logout and login again for changes to take effect
# Or use: newgrp systemd-journal
```

### Create Aliases

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Systemd Monitor Aliases
alias sysmon='sudo /path/to/systemd-monitoring/bin/monitor'
alias syslist='sudo /path/to/systemd-monitoring/bin/monitor list'
alias syscheck='sudo /path/to/systemd-monitoring/bin/monitor check'
alias syslog='sudo /path/to/systemd-monitoring/bin/monitor logs'
alias syserr='sudo /path/to/systemd-monitoring/bin/monitor list --status failed'
```

Usage after reload:
```bash
syslist
syscheck nginx
syslog nginx --follow
syserr
```

### Environment Variables

You can set these in your environment:

```bash
# Default sudo usage
export SYSMON_USE_SUDO=true

# Default log file location
export SYSMON_LOG_DIR=/var/log/sysmon

# Default monitoring interval
export SYSMON_INTERVAL=30s
```

---

## 🐛 Troubleshooting

### Common Issues

#### 1. Permission Denied

**Problem:**
```
Error: failed to execute systemctl: permission denied
```

**Solution:**
```bash
# Use sudo
sudo ./bin/monitor list

# Or add --sudo flag
./bin/monitor list --sudo

# Or add user to required groups
sudo usermod -a -G systemd-journal $USER
```

---

#### 2. Command Not Found

**Problem:**
```
bash: monitor: command not found
```

**Solution:**
```bash
# Use full path
/path/to/systemd-monitoring/bin/monitor list

# Or add to PATH
export PATH=$PATH:/path/to/systemd-monitoring/bin

# Or install globally
sudo cp bin/monitor /usr/local/bin/systemd-monitor
```

---

#### 3. Service Not Found

**Problem:**
```
Error: failed to get service status: exit status 4
```

**Solution:**
```bash
# Check if service exists
systemctl list-units --type=service | grep myservice

# Use correct service name
./bin/monitor check nginx.service  # Include .service suffix

# Or let the tool add it automatically
./bin/monitor check nginx
```

---

#### 4. No Logs Displayed

**Problem:**
```
No logs found
```

**Solutions:**
```bash
# Check if service has logs
sudo journalctl -u nginx -n 10

# Use sudo
sudo ./bin/monitor logs nginx

# Check service is running
systemctl status nginx

# Try different time range
sudo ./bin/monitor logs nginx --since "today"
```

---

#### 5. Build Errors

**Problem:**
```
go: module not found
```

**Solution:**
```bash
# Initialize go modules
go mod init github.com/andinianst93/systemd-monitoring
go mod tidy

# Build again
go build -o bin/monitor main.go
```

---

### Debug Mode

For verbose output:

```bash
# Run with go to see detailed errors
go run main.go list

# Check systemctl directly
systemctl list-units --type=service --all --no-pager

# Test journalctl access
journalctl -u nginx -n 10
```

---

## 🛠️ Development

### Project Structure

```
systemd-monitoring/
├── main.go                          # Main entry point
├── internal/
│   ├── models/                      # Data models
│   │   ├── service.go              # Service models
│   │   └── log.go                  # Log models
│   ├── systemd/                     # Systemd interactions
│   │   └── client.go               # Systemd client
│   ├── output/                      # Output formatters
│   │   ├── table.go                # Table formatter
│   │   └── json.go                 # JSON formatter
│   └── logger/                      # File logging
│       └── file_logger.go          # File logger
│       └── journal_logger.go       # Journal logger
├── bin/                             # Compiled binaries
```

### Building from Source

```bash
# Clone repository
git clone https://github.com/andinianst93/systemd-monitoring.git
cd systemd-monitoring

# Install dependencies
go mod download

# Build
go build -o bin/monitor main.go

# Build with optimizations
go build -ldflags="-s -w" -o bin/monitor main.go

# Cross-compile for different architectures
GOOS=linux GOARCH=amd64 go build -o bin/monitor-amd64 main.go
GOOS=linux GOARCH=arm64 go build -o bin/monitor-arm64 main.go
```

### Adding New Features

1. **Create feature branch**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Implement feature** following existing patterns

3. **Test thoroughly**
   ```bash
   go build -o bin/monitor main.go
   ./bin/monitor your-new-command
   ```

4. **Update documentation**

5. **Submit pull request**

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

### Ways to Contribute

1. **Report Bugs** - Open an issue with details
2. **Suggest Features** - Open an issue with your idea
3. **Submit Pull Requests** - Fix bugs or add features
4. **Improve Documentation** - Help make docs better
5. **Share** - Star the repo and share with others

### Contribution Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code of Conduct

- Be respectful and inclusive
- Follow Go best practices
- Write clear commit messages
- Add tests for new features
- Update documentation

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/andinianst93/systemd-monitoring/issues)

---

## 📊 Project Status

**Current Version:** 1.0.0  
**Status:** ✅ Production Ready  
**Last Updated:** December 22, 2025

### Features Status

| Feature | Status |
|---------|--------|
| List Services | ✅ Complete |
| Check Services | ✅ Complete |
| Monitor Services | ✅ Complete |
| View Logs | ✅ Complete |
| Write Logs | ✅ Complete |
| JSON Export | ✅ Complete |
| Real-time Following | ✅ Complete |
| Filtering | ✅ Complete |
| Colorized Output | ✅ Complete |

---

## 🗺️ Roadmap

### Planned Features

- [ ] Web Dashboard UI
- [ ] Email/Slack notifications
- [ ] Service dependency graphs
- [ ] Historical data storage
- [ ] Performance metrics
- [ ] Custom alert rules
- [ ] Configuration file support
- [ ] Docker container support
- [ ] Kubernetes integration
- [ ] Prometheus exporter

---

<div align="center">
⭐ Star this repo if you find it helpful!

</div>
