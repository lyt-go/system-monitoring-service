package store

import (
	"context"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/ingest"
)

func (s *MemoryStore) CreateSample(sa *model.Sample) error {
	return s.CreateSampleContext(context.Background(), sa)
}

func (s *MemoryStore) CreateSampleContext(ctx context.Context, sa *model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[sa.MetricID]; !ok {
		return ErrNotFound
	}
	s.samples[sa.ID] = sa
	if err := ingest.BeforeCommit(ctx); err != nil {
		return err
	}
	return ctx.Err()
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
