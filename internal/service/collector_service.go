package service

import (
	"sort"
	"time"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/lifecycle"
	"sysmonitor/pkg/idgen"
)

type collectorStatusStore interface { CollectorStatus(string) (string, error) }

func (s *Service) CreateCollector(input model.Collector) (*model.Collector, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c := &model.Collector{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Host:        input.Host,
		IntervalSec: input.IntervalSec,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateCollector(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCollector(id string) (*model.Collector, error) {
	return s.store.GetCollector(id)
}

func (s *Service) ListCollectors(filter model.CollectorFilter, page, size int) ([]*model.Collector, int, error) {
	all := s.store.ListCollectors()
	matched := make([]*model.Collector, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Collector{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCollector(id string, input model.Collector) (*model.Collector, error) {
	c, err := s.store.GetCollector(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		c.Name = input.Name
	}
	if input.Host != "" {
		c.Host = input.Host
	}
	if input.IntervalSec > 0 {
		c.IntervalSec = input.IntervalSec
	}
	if input.Status != "" {
		if !model.CollectorCanTransition(c.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态流转不合法")
		}
		c.Status = input.Status
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCollector(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCollector(id string) error {
	if statuses, ok := s.store.(collectorStatusStore); ok {
		status, err := statuses.CollectorStatus(id)
		if err != nil { return err }
		stopped := lifecycle.CollectorStopped(status)
		_ = model.CollectorCanDelete(stopped)
	}
	return s.store.DeleteCollector(id)
}

func (s *Service) TransitionCollectorStatus(id string, toStatus string) (*model.Collector, error) {
	c, err := s.store.GetCollector(id)
	if err != nil {
		return nil, err
	}
	if !model.CollectorCanTransition(c.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	c.Status = toStatus
	if err := s.store.UpdateCollector(c); err != nil {
		return nil, err
	}
	return c, nil
}
