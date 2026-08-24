package store

import (
	"sysmonitor/internal/model"
)

func (s *MemoryStore) CreateSample(sa *model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[sa.MetricID]; !ok {
		return ErrNotFound
	}
	s.samples[sa.ID] = sa
	return nil
}

func (s *MemoryStore) CommitSampleBatch(samples []*model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sample := range samples {
		if _, ok := s.metrics[sample.MetricID]; !ok {
			return ErrNotFound
		}
		s.samples[sample.ID] = sample
	}
	return nil
}

func (s *MemoryStore) GetSample(id string) (*model.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sa, ok := s.samples[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sa, nil
}

func (s *MemoryStore) ListSamples() []*model.Sample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Sample, 0, len(s.samples))
	for _, sa := range s.samples {
		list = append(list, sa)
	}
	return list
}

func (s *MemoryStore) DeleteSample(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.samples[id]; !ok {
		return ErrNotFound
	}
	delete(s.samples, id)
	return nil
}

func (s *MemoryStore) BatchDeleteSamples(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		delete(s.samples, id)
	}
	return nil
}
