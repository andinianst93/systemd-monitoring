package models

import (
	"fmt"
	"time"
)

type ServiceStatus string

const (
	StatusRunning ServiceStatus = "running"
	StatusStopped ServiceStatus = "stopped"
	StatusFailed  ServiceStatus = "failed"
	StatusUnknown ServiceStatus = "unknown"
)

type ServiceInfo struct {
	Name        string
	Status      ServiceStatus
	ActiveState string
	SubState    string
	Uptime      time.Duration
	PID         int
	MemoryUsage string
	CheckedAt   time.Time
}

func NewServiceInfo(name string) *ServiceInfo {
	return &ServiceInfo{
		Name:        name,
		Status:      StatusUnknown,
		ActiveState: "",
		SubState:    "",
		Uptime:      0,
		PID:         0,
		MemoryUsage: "",
		CheckedAt:   time.Now(),
	}
}

func (s *ServiceInfo) IsRunning() bool {
	return s.Status == StatusRunning
}

func (s *ServiceInfo) IsFailed() bool {
	return s.Status == StatusFailed
}

func (s *ServiceInfo) GetUptimeString() string {
	// TODO: Format uptime ke human-readable
	// Examples: "2h 15m", "5d 3h", "45m"
	if s.Uptime == 0 {
		return "0s"
	}
	hours := s.Uptime / time.Hour
	minutes := (s.Uptime % time.Hour) / time.Minute
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func (s *ServiceInfo) GetStatusIcon() string {
	// TODO: Return emoji based on status
	// ✅ for running, ❌ for failed, ⏸️ for stopped
	switch s.Status {
	case StatusRunning:
		return "✅"
	case StatusFailed:
		return "❌"
	case StatusStopped:
		return "⏸️"
	default:
		return "❓"
	}
}

type ServiceList struct {
	Services  []*ServiceInfo
	Timestamp time.Time
	Total     int
	Running   int
	Failed    int
	Stopped   int
}

func NewServiceList() *ServiceList {
	// TODO: Initialize empty list
	return &ServiceList{
		Services:  []*ServiceInfo{},
		Timestamp: time.Now(),
		Total:     0,
		Running:   0,
		Failed:    0,
		Stopped:   0,
	}
}

func (sl *ServiceList) AddService(service *ServiceInfo) {
	// TODO: Add service and update counters
	sl.Services = append(sl.Services, service)
	sl.Total++
	switch service.Status {
	case StatusRunning:
		sl.Running++
	case StatusFailed:
		sl.Failed++
	case StatusStopped:
		sl.Stopped++
	}
}

func (sl *ServiceList) GetByStatus(status ServiceStatus) []*ServiceInfo {
	// TODO: Filter services by status
	var filtered []*ServiceInfo
	for _, service := range sl.Services {
		if service.Status == status {
			filtered = append(filtered, service)
		}
	}
	return filtered
}

func (sl *ServiceList) HasFailures() bool {
	return sl.Failed > 0
}
