package service

import (
	"sort"
	"time"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/batch"
	"sysmonitor/pkg/idgen"
)

func (s *Service) CreateSample(input model.Sample) (*model.Sample, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMetric(input.MetricID); err != nil {
		return nil, model.NewValidationError("metric_id", "关联指标不存在")
	}
	sa := &model.Sample{
		ID:        idgen.Hex(),
		MetricID:  input.MetricID,
		Value:     input.Value,
		Host:      input.Host,
		Timestamp: input.Timestamp,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateSample(sa); err != nil {
		return nil, err
	}
	return sa, nil
}

func (s *Service) StageSampleBatch(samples []*model.Sample) error {
	if len(samples) == 0 {
		return model.NewValidationError("samples", "样本批次不能为空")
	}
	for _, sample := range samples {
		if sample == nil {
			return model.NewValidationError("samples", "样本不能为空")
		}
		if err := sample.Validate(); err != nil {
			return err
		}
	}
	s.batchMu.Lock()
	s.pendingSamples = samples
	s.batchMu.Unlock()
	return nil
}

func (s *Service) CommitSampleBatch() error {
	s.batchMu.Lock()
	pending := batch.SnapshotSamples(s.pendingSamples)
	s.pendingSamples = nil
	s.batchMu.Unlock()
	return s.store.CommitSampleBatch(pending)
}

func (s *Service) GetSample(id string) (*model.Sample, error) {
	return s.store.GetSample(id)
}

func (s *Service) ListSamples(filter model.SampleFilter, page, size int) ([]*model.Sample, int, error) {
	all := s.store.ListSamples()
	matched := make([]*model.Sample, 0, len(all))
	for _, sa := range all {
		if filter.Match(sa) {
			matched = append(matched, sa)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.After(matched[j].Timestamp)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Sample{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteSample(id string) error {
	return s.store.DeleteSample(id)
}

func (s *Service) BatchDeleteSamples(ids []string) error {
	if len(ids) == 0 {
		return model.NewValidationError("ids", "ID 列表不能为空")
	}
	return s.store.BatchDeleteSamples(ids)
}
