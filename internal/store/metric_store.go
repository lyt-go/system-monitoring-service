package store

import (
	"sysmonitor/internal/model"
)

func (s *MemoryStore) CreateMetric(m *model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.metrics {
		if exist.Name == m.Name {
			return ErrConflict
		}
	}
	s.metrics[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMetric(id string) (*model.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.metrics[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) GetMetricByName(name string) (*model.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, m := range s.metrics {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListMetrics() []*model.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Metric, 0, len(s.metrics))
	for _, m := range s.metrics {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) UpdateMetric(m *model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[m.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.metrics {
		if exist.ID != m.ID && exist.Name == m.Name {
			return ErrConflict
		}
	}
	s.metrics[m.ID] = m
	return nil
}

// HasMetricReferences reports whether any threshold, sample, or alert still
// references the given metric. The caller relies on this to decide whether the
// metric can be deleted without leaving orphaned dependents behind.
func (s *MemoryStore) HasMetricReferences(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.thresholds {
		if t.MetricID == id {
			return true
		}
	}
	for _, sa := range s.samples {
		if sa.MetricID == id {
			return true
		}
	}
	for _, a := range s.alerts {
		if a.MetricID == id {
			return true
		}
	}
	return false
}

func (s *MemoryStore) DeleteMetric(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[id]; !ok {
		return ErrNotFound
	}
	delete(s.metrics, id)
	return nil
}
