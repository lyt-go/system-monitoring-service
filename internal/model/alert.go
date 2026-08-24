package model

import (
	"strings"
	"time"
)

const (
	AlertLevelWarn     = "warn"
	AlertLevelCritical = "critical"

	AlertStatusOpen         = "open"
	AlertStatusAcknowledged = "acknowledged"
	AlertStatusResolved     = "resolved"
)

type Alert struct {
	ID                string    `json:"id"`
	MetricID          string    `json:"metric_id"`
	ThresholdID       string    `json:"threshold_id"`
	Level             string    `json:"level"`
	Message           string    `json:"message"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	DeliveryAttempts  int       `json:"delivery_attempts"`
	LastDeliveryError string    `json:"last_delivery_error,omitempty"`
}

func (a *Alert) Clone() *Alert {
	if a == nil {
		return nil
	}
	clone := *a
	return &clone
}

func (a *Alert) MarkDeliveryFailure(err error) {
	a.Status = AlertStatusAcknowledged
	a.LastDeliveryError = err.Error()
}

func (a *Alert) Validate() error {
	a.MetricID = strings.TrimSpace(a.MetricID)
	a.ThresholdID = strings.TrimSpace(a.ThresholdID)
	a.Level = strings.TrimSpace(a.Level)
	a.Message = strings.TrimSpace(a.Message)
	if a.MetricID == "" {
		return NewValidationError("metric_id", "关联指标不能为空")
	}
	if a.ThresholdID == "" {
		return NewValidationError("threshold_id", "关联阈值不能为空")
	}
	if a.Message == "" {
		return NewValidationError("message", "告警消息不能为空")
	}
	if a.Level == "" {
		a.Level = AlertLevelWarn
	}
	if a.Level != AlertLevelWarn && a.Level != AlertLevelCritical {
		return NewValidationError("level", "告警级别不合法")
	}
	if a.Status == "" {
		a.Status = AlertStatusOpen
	}
	if a.Status != AlertStatusOpen && a.Status != AlertStatusAcknowledged && a.Status != AlertStatusResolved {
		return NewValidationError("status", "告警状态不合法")
	}
	return nil
}

var alertTransitions = map[string]map[string]bool{
	AlertStatusOpen:         {AlertStatusAcknowledged: true, AlertStatusResolved: true},
	AlertStatusAcknowledged: {AlertStatusResolved: true},
	AlertStatusResolved:     {},
}

func AlertCanTransition(from, to string) bool {
	if m, ok := alertTransitions[from]; ok {
		return m[to]
	}
	return false
}

type AlertFilter struct {
	MetricID string
	Level    string
	Status   string
	Keyword  string
}

func (f AlertFilter) Match(a *Alert) bool {
	if f.MetricID != "" && a.MetricID != f.MetricID {
		return false
	}
	if f.Level != "" && a.Level != f.Level {
		return false
	}
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Message), k) {
			return false
		}
	}
	return true
}
