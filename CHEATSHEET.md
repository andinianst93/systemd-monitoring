# 🚀 Systemd Monitor - Quick Reference Cheatsheet

Fast reference for the most commonly used commands.

---

## 📋 Basic Commands

```bash
# List all services
sudo ./bin/monitor list

# Check a service
sudo ./bin/monitor check nginx

# View logs
sudo ./bin/monitor logs nginx

# Monitor continuously
sudo ./bin/monitor monitor --services nginx --interval 30s

# Write logs to journal
./bin/monitor write-log --message "Application started"
```

---

## 📊 LIST Command

```bash
# All services (default)
sudo ./bin/monitor list

# Running services only
sudo ./bin/monitor list --status running

# Failed services only
sudo ./bin/monitor list --status failed

# Stopped services only
sudo ./bin/monitor list --status stopped

# Export as JSON
sudo ./bin/monitor list --output json

# JSON + filter
sudo ./bin/monitor list --status running --output json > running.json
```

---

## 🔍 CHECK Command

```bash
# Single service
sudo ./bin/monitor check nginx

# Multiple services
sudo ./bin/monitor check nginx mysql redis

# With .service suffix (optional)
sudo ./bin/monitor check nginx.service

# Check and pipe to file
sudo ./bin/monitor check nginx > status.txt
```

**Output includes:**
- Status icon (✅/❌/⏸️)
- Service name
- PID
- Memory usage
- Uptime

---

## 📡 MONITOR Command

```bash
# Basic monitoring (30s default)
sudo ./bin/monitor monitor --services nginx

# Custom interval
sudo ./bin/monitor monitor --services nginx --interval 10s

# Multiple services
sudo ./bin/monitor monitor --services nginx,mysql,redis --interval 1m

# Custom log file
sudo ./bin/monitor monitor --services nginx --interval 30s --log-file /var/log/monitor.log

# Quick intervals
sudo ./bin/monitor monitor --services redis --interval 5s
```

**Interval formats:**
- `5s` = 5 seconds
- `1m` = 1 minute
- `5m` = 5 minutes
- `1h` = 1 hour

**Stop:** Press `Ctrl+C`

---

## 📝 LOGS Command

### Basic Usage

```bash
# Last 50 lines (default)
sudo ./bin/monitor logs nginx

# Last 100 lines
sudo ./bin/monitor logs nginx --lines 100

# Last 20 lines
sudo ./bin/monitor logs nginx --lines 20
```

### Follow Logs (Real-time)

```bash
# Follow logs (like tail -f)
sudo ./bin/monitor logs nginx --follow

# Follow with initial context
sudo ./bin/monitor logs nginx --follow --lines 20
```

### Time-based Filtering

```bash
# Last hour
sudo ./bin/monitor logs nginx --since "1 hour ago"

# Last 2 hours
sudo ./bin/monitor logs nginx --since "2 hours ago"

# Today's logs
sudo ./bin/monitor logs nginx --since "today"

# Yesterday's logs
sudo ./bin/monitor logs nginx --since "yesterday"

# Specific date
sudo ./bin/monitor logs nginx --since "2024-12-22"

# Specific time range
sudo ./bin/monitor logs nginx --since "2024-12-22 00:00" --until "2024-12-22 12:00"

# Last 10 minutes
sudo ./bin/monitor logs nginx --since "10 minutes ago"
```

### Priority Filtering

```bash
# Only errors
sudo ./bin/monitor logs nginx --priority err

# Only warnings and above
sudo ./bin/monitor logs nginx --priority warning

# Only critical
sudo ./bin/monitor logs nginx --priority crit

# Info and above
sudo ./bin/monitor logs nginx --priority info
```

**Priority levels:** `emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`

### Search/Grep

```bash
# Find "error" in logs
sudo ./bin/monitor logs nginx --grep "error"

# Find "timeout"
sudo ./bin/monitor logs nginx --grep "timeout"

# Find "configuration"
sudo ./bin/monitor logs nginx --grep "configuration"

# Case-insensitive (automatic)
sudo ./bin/monitor logs nginx --grep "ERROR"
```

### Combined Filters

```bash
# Last hour's errors
sudo ./bin/monitor logs nginx --since "1 hour ago" --priority err

# Follow errors only
sudo ./bin/monitor logs nginx --follow --priority err

# Search in recent logs
sudo ./bin/monitor logs nginx --since "today" --grep "connection"

# Complex query
sudo ./bin/monitor logs nginx --lines 100 --since "1 hour ago" --grep "error" --priority err
```

---

## 📝 WRITE-LOG Command

### Basic Usage

```bash
# Write info message (default)
./bin/monitor write-log --message "Service started"

# Write with priority
./bin/monitor write-log --message "Error occurred" --priority err

# Write with custom identifier
./bin/monitor write-log --message "Deployment complete" --identifier my-app
```

### Priority Levels

```bash
# Info (default)
./bin/monitor write-log --message "Normal operation"

# Warning
./bin/monitor write-log --message "High memory" --priority warning

# Error
./bin/monitor write-log --message "Connection failed" --priority err

# Critical
./bin/monitor write-log --message "System critical" --priority crit

# Debug
./bin/monitor write-log --message "Debug info" --priority debug
```

### View Written Logs

```bash
# View by identifier
journalctl -t systemd-monitor -n 20

# View custom identifier
journalctl -t my-app -n 50

# Follow logs
journalctl -t systemd-monitor -f

# Filter by priority
journalctl -t systemd-monitor -p err
```

### Script Integration

```bash
#!/bin/bash
# Example: Backup script with logging

APP="backup-service"

./bin/monitor write-log \
  --message "Backup started" \
  --priority info \
  --identifier $APP

if /usr/bin/backup.sh; then
    ./bin/monitor write-log \
      --message "Backup completed" \
      --priority info \
      --identifier $APP
else
    ./bin/monitor write-log \
      --message "Backup failed" \
      --priority err \
      --identifier $APP
fi
```

---

## 🎨 Output Colors & Icons

### Service Status

| Icon | Status | Color |
|------|--------|-------|
| ✅ | Running | Green |
| ❌ | Failed | Red |
| ⏸️ | Stopped | Yellow |
| ❓ | Unknown | White |

### Log Levels

| Icon | Level | Color |
|------|-------|-------|
| ✅ | INFO | Green |
| 🔍 | DEBUG | White |
| ℹ️ | NOTICE | Cyan |
| ⚠️ | WARNING | Yellow |
| ❌ | ERROR | Red |
| ❌ | CRITICAL | Red |

---

## 🔑 Common Flags

```bash
--sudo              # Use sudo for operations
--output json       # JSON output (list command)
--status <filter>   # Filter by status (list command)
--lines <n>         # Number of log lines (logs command)
--follow            # Follow logs in real-time (logs command)
--since <time>      # Time filter (logs command)
--until <time>      # Time until (logs command)
--priority <level>  # Priority filter (logs/write-log command)
--grep <pattern>    # Search pattern (logs command)
--interval <dur>    # Check interval (monitor command)
--services <list>   # Service list (monitor command)
--log-file <path>   # Log file path (monitor command)
--message <text>    # Message to write (write-log command)
--identifier <name> # Application identifier (write-log command)
```

---

## 💾 Export Examples

```bash
# Export all services to JSON
sudo ./bin/monitor list --output json > services.json

# Export only running services
sudo ./bin/monitor list --status running --output json > running-services.json

# Export failed services
sudo ./bin/monitor list --status failed --output json > failed-services.json

# Save logs to file
sudo ./bin/monitor logs nginx --lines 1000 > nginx-logs.txt

# Save error logs
sudo ./bin/monitor logs nginx --priority err --since "today" > nginx-errors.txt
```

---

## 🔧 Installation Quick Start

```bash
# Clone & Build
git clone https://github.com/yourusername/systemd-monitoring.git
cd systemd-monitoring
go build -o bin/monitor main.go

# Optional: Install globally
sudo cp bin/monitor /usr/local/bin/systemd-monitor

# Test
sudo ./bin/monitor list
```

---

## 🎯 Common Use Cases

### Daily Health Check
```bash
sudo ./bin/monitor check nginx mysql redis
sudo ./bin/monitor list --status failed
```

### Debug Service Issues
```bash
sudo ./bin/monitor check nginx
sudo ./bin/monitor logs nginx --since "1 hour ago" --priority err
```

### Monitor Critical Service
```bash
sudo ./bin/monitor monitor --services nginx --interval 10s --log-file logs/nginx-monitor.log
```

### Find Recent Errors
```bash
sudo ./bin/monitor logs nginx --since "1 hour ago" --grep "error"
```

### Follow Logs for Debugging
```bash
sudo ./bin/monitor logs nginx --follow --priority err
```

### Export Service Inventory
```bash
sudo ./bin/monitor list --output json > inventory-$(date +%Y%m%d).json
```

### Write Custom Logs
```bash
# Log deployment event
./bin/monitor write-log --message "Deployment v1.2.3 started" --priority info --identifier deployment

# Log monitoring alert
./bin/monitor write-log --message "CPU usage exceeded 90%" --priority warning --identifier monitoring
```

---

## 🐛 Quick Troubleshooting

### Permission Denied
```bash
# Use sudo
sudo ./bin/monitor list

# Or add to group
sudo usermod -a -G systemd-journal $USER
```

### Service Not Found
```bash
# Check service exists
systemctl list-units --type=service | grep myservice

# Use correct name
sudo ./bin/monitor check nginx.service
```

### No Logs
```bash
# Use sudo
sudo ./bin/monitor logs nginx

# Try without filters
sudo ./bin/monitor logs nginx --lines 10
```

---

## 📦 Useful Aliases

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Basic aliases
alias sysmon='sudo /path/to/bin/monitor'
alias syslist='sudo /path/to/bin/monitor list'
alias syscheck='sudo /path/to/bin/monitor check'
alias syslog='sudo /path/to/bin/monitor logs'
alias syswrite='/path/to/bin/monitor write-log'

# Quick commands
alias syserr='sudo /path/to/bin/monitor list --status failed'
alias sysrun='sudo /path/to/bin/monitor list --status running'
alias sysfollow='sudo /path/to/bin/monitor logs'
```

Usage:
```bash
syslist
syscheck nginx
syslog nginx --follow
syswrite --message "Event occurred" --priority info
syserr
```

---

## 🔗 More Information

- **Full Documentation:** [README.md](README.md)
- **Complete Guide:** [SYSTEMD_MONITOR_GUIDE.md](SYSTEMD_MONITOR_GUIDE.md)
- **Testing Guide:** [TESTING.md](TESTING.md)
- **Log Viewing:** [LOGS_FEATURE.md](LOGS_FEATURE.md)

---

## 📞 Quick Help

```bash
# Show help
./bin/monitor --help
./bin/monitor list --help
./bin/monitor check --help
./bin/monitor monitor --help
./bin/monitor logs --help
```

---

**Version:** 1.0.0  
**Last Updated:** December 22, 2024  
**Status:** ✅ Production Ready

---

Made with ❤️ - [Star on GitHub](https://github.com/yourusername/systemd-monitoring) ⭐