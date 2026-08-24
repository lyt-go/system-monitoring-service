package model

import (
	"strings"
	"time"
)

const (
	ThresholdOpGT  = "gt"
	ThresholdOpLT  = "lt"
	ThresholdOpGTE = "gte"
	ThresholdOpLTE = "lte"
	ThresholdOpEQ  = "eq"

	ThresholdStatusActive   = "active"
	ThresholdStatusDisabled = "disabled"
)

type Threshold struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	Operator  string    `json:"operator"`
	Value     float64   `json:"value"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Threshold) Validate() error {
	t.MetricID = strings.TrimSpace(t.MetricID)
	t.Operator = strings.TrimSpace(t.Operator)
	if t.MetricID == "" {
		return NewValidationError("metric_id", "关联指标不能为空")
	}
	if t.Operator == "" {
		return NewValidationError("operator", "操作符不能为空")
	}
	validOps := map[string]bool{ThresholdOpGT: true, ThresholdOpLT: true, ThresholdOpGTE: true, ThresholdOpLTE: true, ThresholdOpEQ: true}
	if !validOps[t.Operator] {
		return NewValidationError("operator", "操作符不合法")
	}
	if t.Status == "" {
		t.Status = ThresholdStatusActive
	}
	if t.Status != ThresholdStatusActive && t.Status != ThresholdStatusDisabled {
		return NewValidationError("status", "阈值状态不合法")
	}
	return nil
}

var thresholdTransitions = map[string]map[string]bool{
	ThresholdStatusActive:   {ThresholdStatusDisabled: true},
	ThresholdStatusDisabled: {ThresholdStatusActive: true},
}

func ThresholdCanTransition(from, to string) bool {
	if m, ok := thresholdTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ThresholdFilter struct {
	MetricID string
	Operator string
	Status   string
}

func (f ThresholdFilter) Match(t *Threshold) bool {
	if f.MetricID != "" && t.MetricID != f.MetricID {
		return false
	}
	if f.Operator != "" && t.Operator != f.Operator {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	return true
}
