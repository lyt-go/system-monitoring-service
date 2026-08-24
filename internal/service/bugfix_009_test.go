package service

import (
	"errors"
	"testing"
	"sysmonitor/internal/model"
	"sysmonitor/internal/store"
)

func TestRuleRemovalKeepsAlertDependencyIntact(t *testing.T) {
	st := store.NewMemoryStore(); svc := New(st, nil, nil)
	metric, err := svc.CreateMetric(model.Metric{Name:"cpu", Unit:"%", Type:"gauge"}); if err != nil { t.Fatalf("metric setup failed: %v", err) }
	threshold, err := svc.CreateThreshold(model.Threshold{MetricID:metric.ID, Operator:"gt", Value:80}); if err != nil { t.Fatalf("threshold setup failed: %v", err) }
	alert, err := svc.CreateAlert(model.Alert{MetricID:metric.ID, ThresholdID:threshold.ID, Message:"high"}); if err != nil { t.Fatalf("alert setup failed: %v", err) }
	if err := svc.DeleteThreshold(threshold.ID); !errors.Is(err, store.ErrConflict) { t.Fatalf("referenced threshold deletion should conflict, got %v", err) }
	if _, err := svc.GetThreshold(threshold.ID); err != nil { t.Fatalf("threshold was deleted: %v", err) }
	if _, err := svc.GetAlert(alert.ID); err != nil { t.Fatalf("alert lost its threshold: %v", err) }
	free, err := svc.CreateThreshold(model.Threshold{MetricID:metric.ID, Operator:"lt", Value:10}); if err != nil { t.Fatalf("free threshold setup failed: %v", err) }
	if err := svc.DeleteThreshold(free.ID); err != nil { t.Fatalf("unreferenced threshold should delete: %v", err) }
}
