package service

import (
	"testing"
	"sysmonitor/internal/model"
	"sysmonitor/internal/store"
)

func TestTopAlertMetricsRejectsNonPositiveLimit(t *testing.T) {
	st := store.NewMemoryStore(); svc := New(st, nil, nil)
	metric, err := svc.CreateMetric(model.Metric{Name:"cpu", Unit:"%", Type:"gauge"}); if err != nil { t.Fatalf("metric setup failed: %v", err) }
	threshold, err := svc.CreateThreshold(model.Threshold{MetricID:metric.ID, Operator:"gt", Value:80}); if err != nil { t.Fatalf("threshold setup failed: %v", err) }
	if _, err := svc.CreateAlert(model.Alert{MetricID:metric.ID, ThresholdID:threshold.ID, Message:"high"}); err != nil { t.Fatalf("alert setup failed: %v", err) }
	if _, err := svc.GetTopAlertMetrics(0); !model.IsValidationError(err) { t.Fatalf("zero limit should be rejected, got %v", err) }
	if _, err := svc.GetTopAlertMetrics(-1); !model.IsValidationError(err) { t.Fatalf("negative limit should be rejected, got %v", err) }
	items, err := svc.GetTopAlertMetrics(1); if err != nil || len(items) != 1 { t.Fatalf("positive limit changed: items=%v err=%v", items, err) }
}
