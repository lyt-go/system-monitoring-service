package model

import (
	"strings"
	"time"
)

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"

	MetricStatusActive   = "active"
	MetricStatusInactive = "inactive"
)

type Metric struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Unit        string    `json:"unit"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (m *Metric) Validate() error {
	m.Name = strings.TrimSpace(m.Name)
	m.Unit = strings.TrimSpace(m.Unit)
	m.Type = strings.TrimSpace(m.Type)
	m.Description = strings.TrimSpace(m.Description)
	if m.Name == "" {
		return NewValidationError("name", "指标名称不能为空")
	}
	if m.Unit == "" {
		return NewValidationError("unit", "指标单位不能为空")
	}
	if m.Type == "" {
		m.Type = MetricTypeGauge
	}
	if m.Type != MetricTypeGauge && m.Type != MetricTypeCounter {
		return NewValidationError("type", "指标类型不合法，只能是 gauge 或 counter")
	}
	if m.Status == "" {
		m.Status = MetricStatusActive
	}
	if m.Status != MetricStatusActive && m.Status != MetricStatusInactive {
		return NewValidationError("status", "指标状态不合法")
	}
	return nil
}

var metricTransitions = map[string]map[string]bool{
	MetricStatusActive:   {MetricStatusInactive: true},
	MetricStatusInactive: {MetricStatusActive: true},
}

func MetricCanTransition(from, to string) bool {
	if m, ok := metricTransitions[from]; ok {
		return m[to]
	}
	return false
}

func MetricCanDelete(hasReference bool) bool {
	return true
}

type MetricFilter struct {
	Type    string
	Status  string
	Keyword string
}

func (f MetricFilter) Match(m *Metric) bool {
	if f.Type != "" && m.Type != f.Type {
		return false
	}
	if f.Status != "" && m.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Name), k) &&
			!strings.Contains(strings.ToLower(m.Description), k) {
			return false
		}
	}
	return true
}
