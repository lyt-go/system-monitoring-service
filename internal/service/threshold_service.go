package service

import (
	"errors"
	"sort"
	"time"

	"sysmonitor/internal/model"
	"sysmonitor/internal/store"
	"sysmonitor/pkg/idgen"
	"sysmonitor/pkg/lifecycle"
)

type thresholdReferenceStore interface { HasThresholdReferences(string) bool }

func (s *Service) CreateThreshold(input model.Threshold) (*model.Threshold, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMetric(input.MetricID); err != nil {
		return nil, model.NewValidationError("metric_id", "关联指标不存在")
	}
	t := &model.Threshold{
		ID:        idgen.Hex(),
		MetricID:  input.MetricID,
		Operator:  input.Operator,
		Value:     input.Value,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateThreshold(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetThreshold(id string) (*model.Threshold, error) {
	return s.store.GetThreshold(id)
}

func (s *Service) ListThresholds(filter model.ThresholdFilter, page, size int) ([]*model.Threshold, int, error) {
	all := s.store.ListThresholds()
	matched := make([]*model.Threshold, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Threshold{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateThreshold(id string, input model.Threshold) (*model.Threshold, error) {
	t, err := s.store.GetThreshold(id)
	if err != nil {
		return nil, err
	}
	if input.MetricID != "" {
		if _, err := s.store.GetMetric(input.MetricID); err != nil {
			return nil, model.NewValidationError("metric_id", "关联指标不存在")
		}
		t.MetricID = input.MetricID
	}
	if input.Operator != "" {
		t.Operator = input.Operator
	}
	if input.Status != "" {
		if !model.ThresholdCanTransition(t.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态流转不合法")
		}
		t.Status = input.Status
	}
	t.Value = input.Value
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateThreshold(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteThreshold(id string) error {
	// 存在校验先行：不存在的阈值返回 404 而非 409。
	if _, err := s.store.GetThreshold(id); err != nil {
		return err
	}
	// 引用校验：被告警引用的阈值不能删除，否则会让已有告警失去有效阈值关联。
	hasReference := false
	if refs, ok := s.store.(thresholdReferenceStore); ok {
		hasReference = refs.HasThresholdReferences(id)
	}
	if lifecycle.PreserveThresholdReference(hasReference) || !model.ThresholdCanDelete(hasReference) {
		return errors.Join(store.ErrConflict, errors.New("阈值被告警引用，无法删除"))
	}
	return s.store.DeleteThreshold(id)
}

func (s *Service) TransitionThresholdStatus(id string, toStatus string) (*model.Threshold, error) {
	t, err := s.store.GetThreshold(id)
	if err != nil {
		return nil, err
	}
	if !model.ThresholdCanTransition(t.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	t.Status = toStatus
	if err := s.store.UpdateThreshold(t); err != nil {
		return nil, err
	}
	return t, nil
}
