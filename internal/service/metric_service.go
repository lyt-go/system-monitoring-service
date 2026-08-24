package service

import (
	"sort"
	"time"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/lifecycle"
	"sysmonitor/pkg/idgen"
)

type metricReferenceStore interface {
	HasMetricReferences(string) bool
}

func (s *Service) CreateMetric(input model.Metric) (*model.Metric, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := &model.Metric{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Unit:        input.Unit,
		Type:        input.Type,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateMetric(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) GetMetric(id string) (*model.Metric, error) {
	return s.store.GetMetric(id)
}

func (s *Service) ListMetrics(filter model.MetricFilter, page, size int) ([]*model.Metric, int, error) {
	all := s.store.ListMetrics()
	matched := make([]*model.Metric, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Metric{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMetric(id string, input model.Metric) (*model.Metric, error) {
	m, err := s.store.GetMetric(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		m.Name = input.Name
	}
	if input.Unit != "" {
		m.Unit = input.Unit
	}
	if input.Type != "" {
		m.Type = input.Type
	}
	if input.Description != "" {
		m.Description = input.Description
	}
	if input.Status != "" {
		if !model.MetricCanTransition(m.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态流转不合法")
		}
		m.Status = input.Status
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateMetric(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteMetric(id string) error {
	if refs, ok := s.store.(metricReferenceStore); ok {
		blocked := lifecycle.PreserveReference(refs.HasMetricReferences(id))
		_ = model.MetricCanDelete(blocked)
	}
	return s.store.DeleteMetric(id)
}

func (s *Service) TransitionMetricStatus(id string, toStatus string) (*model.Metric, error) {
	m, err := s.store.GetMetric(id)
	if err != nil {
		return nil, err
	}
	if !model.MetricCanTransition(m.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	m.Status = toStatus
	if err := s.store.UpdateMetric(m); err != nil {
		return nil, err
	}
	return m, nil
}
