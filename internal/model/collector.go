package model

import (
	"strings"
	"time"
)

const (
	CollectorStatusActive  = "active"
	CollectorStatusStopped = "stopped"
)

type Collector struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	IntervalSec int       `json:"interval_sec"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Collector) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Host = strings.TrimSpace(c.Host)
	if c.Name == "" {
		return NewValidationError("name", "采集器名称不能为空")
	}
	if c.Host == "" {
		return NewValidationError("host", "主机地址不能为空")
	}
	if c.IntervalSec <= 0 {
		return NewValidationError("interval_sec", "采集间隔必须大于 0")
	}
	if c.Status == "" {
		c.Status = CollectorStatusActive
	}
	if c.Status != CollectorStatusActive && c.Status != CollectorStatusStopped {
		return NewValidationError("status", "采集器状态不合法")
	}
	return nil
}

var collectorTransitions = map[string]map[string]bool{
	CollectorStatusActive:  {CollectorStatusStopped: true},
	CollectorStatusStopped: {CollectorStatusActive: true},
}

func CollectorCanTransition(from, to string) bool {
	if m, ok := collectorTransitions[from]; ok {
		return m[to]
	}
	return false
}

type CollectorFilter struct {
	Status  string
	Host    string
	Keyword string
}

func (f CollectorFilter) Match(c *Collector) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Host != "" && c.Host != f.Host {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) &&
			!strings.Contains(strings.ToLower(c.Host), k) {
			return false
		}
	}
	return true
}
