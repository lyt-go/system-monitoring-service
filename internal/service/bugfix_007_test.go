package service

import (
	"errors"
	"testing"

	"sysmonitor/internal/model"
	"sysmonitor/internal/store"
)

func TestMetricDependencyBlocksResourceRemoval(t *testing.T) {
	st := store.NewMemoryStore()
	svc := New(st, nil, nil)
	metric, err := svc.CreateMetric(model.Metric{Name: "cpu", Unit: "%", Type: "gauge"})
	if err != nil { t.Fatalf("metric setup failed: %v", err) }
	threshold, err := svc.CreateThreshold(model.Threshold{MetricID: metric.ID, Operator: "gt", Value: 80})
	if err != nil { t.Fatalf("threshold setup failed: %v", err) }
	sample, err := svc.CreateSample(model.Sample{MetricID: metric.ID, Host: "node-a"})
	if err != nil { t.Fatalf("sample setup failed: %v", err) }
	alert, err := svc.CreateAlert(model.Alert{MetricID: metric.ID, ThresholdID: threshold.ID, Message: "cpu high"})
	if err != nil { t.Fatalf("alert setup failed: %v", err) }

	if err := svc.DeleteMetric(metric.ID); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("referenced metric deletion should conflict, got %v", err)
	}
	if _, err := svc.GetMetric(metric.ID); err != nil { t.Fatalf("metric was deleted: %v", err) }
	if _, err := svc.GetThreshold(threshold.ID); err != nil { t.Fatalf("threshold became orphaned: %v", err) }
	if _, err := svc.GetSample(sample.ID); err != nil { t.Fatalf("sample became orphaned: %v", err) }
	if _, err := svc.GetAlert(alert.ID); err != nil { t.Fatalf("alert became orphaned: %v", err) }

	free, err := svc.CreateMetric(model.Metric{Name: "idle", Unit: "%", Type: "gauge"})
	if err != nil { t.Fatalf("free metric setup failed: %v", err) }
	if err := svc.DeleteMetric(free.ID); err != nil { t.Fatalf("unreferenced metric should delete: %v", err) }
}
