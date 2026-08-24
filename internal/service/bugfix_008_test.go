package service

import (
	"errors"
	"testing"
	"sysmonitor/internal/model"
	"sysmonitor/internal/store"
)

func TestActiveAgentSurvivesPrematureRemoval(t *testing.T) {
	st := store.NewMemoryStore(); svc := New(st, nil, nil)
	collector, err := svc.CreateCollector(model.Collector{Name:"node-a", Host:"10.0.0.1", IntervalSec:10})
	if err != nil { t.Fatalf("collector setup failed: %v", err) }
	if err := svc.DeleteCollector(collector.ID); !errors.Is(err, store.ErrConflict) { t.Fatalf("active collector deletion should conflict, got %v", err) }
	if _, err := svc.GetCollector(collector.ID); err != nil { t.Fatalf("active collector disappeared: %v", err) }
	if _, err := svc.TransitionCollectorStatus(collector.ID, model.CollectorStatusStopped); err != nil { t.Fatalf("stopping collector failed: %v", err) }
	if err := svc.DeleteCollector(collector.ID); err != nil { t.Fatalf("stopped collector should delete: %v", err) }
}
