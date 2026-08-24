package model

import (
	"strings"
	"time"
)

type Sample struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	Value     float64   `json:"value"`
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Sample) Clone() *Sample {
	return s
}

func (s *Sample) Validate() error {
	s.MetricID = strings.TrimSpace(s.MetricID)
	s.Host = strings.TrimSpace(s.Host)
	if s.MetricID == "" {
		return NewValidationError("metric_id", "关联指标不能为空")
	}
	if s.Host == "" {
		return NewValidationError("host", "主机不能为空")
	}
	if s.Timestamp.IsZero() {
		s.Timestamp = time.Now()
	}
	return nil
}

type SampleFilter struct {
	MetricID string
	Host     string
}

func (f SampleFilter) Match(s *Sample) bool {
	if f.MetricID != "" && s.MetricID != f.MetricID {
		return false
	}
	if f.Host != "" && s.Host != f.Host {
		return false
	}
	return true
}
