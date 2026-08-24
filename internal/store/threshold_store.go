package store

import (
	"sysmonitor/internal/model"
)

func (s *MemoryStore) CreateThreshold(t *model.Threshold) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[t.MetricID]; !ok {
		return ErrNotFound
	}
	s.thresholds[t.ID] = t
	return nil
}

func (s *MemoryStore) GetThreshold(id string) (*model.Threshold, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.thresholds[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListThresholds() []*model.Threshold {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Threshold, 0, len(s.thresholds))
	for _, t := range s.thresholds {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateThreshold(t *model.Threshold) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.thresholds[t.ID]; !ok {
		return ErrNotFound
	}
	if _, ok := s.metrics[t.MetricID]; !ok {
		return ErrNotFound
	}
	s.thresholds[t.ID] = t
	return nil
}

func (s *MemoryStore) HasThresholdReferences(id string) bool {
	if id == "" {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.alerts {
		if a.ThresholdID == id {
			return true
		}
	}
	return false
}

func (s *MemoryStore) DeleteThreshold(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.thresholds[id]; !ok {
		return ErrNotFound
	}
	delete(s.thresholds, id)
	return nil
}
