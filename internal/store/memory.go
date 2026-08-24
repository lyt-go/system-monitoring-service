package store

import (
	"sync"

	"sysmonitor/internal/model"
)

type MemoryStore struct {
	mu          sync.RWMutex
	metrics     map[string]*model.Metric
	collectors  map[string]*model.Collector
	thresholds  map[string]*model.Threshold
	samples     map[string]*model.Sample
	alerts      map[string]*model.Alert
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics:    make(map[string]*model.Metric),
		collectors: make(map[string]*model.Collector),
		thresholds: make(map[string]*model.Threshold),
		samples:    make(map[string]*model.Sample),
		alerts:     make(map[string]*model.Alert),
	}
}

var _ Store = (*MemoryStore)(nil)
