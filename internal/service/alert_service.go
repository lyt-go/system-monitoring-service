package service

import (
	"sort"
	"time"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/idgen"
	"sysmonitor/pkg/notify"
)

func (s *Service) CreateAlert(input model.Alert) (*model.Alert, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMetric(input.MetricID); err != nil {
		return nil, model.NewValidationError("metric_id", "关联指标不存在")
	}
	if _, err := s.store.GetThreshold(input.ThresholdID); err != nil {
		return nil, model.NewValidationError("threshold_id", "关联阈值不存在")
	}
	a := &model.Alert{
		ID:          idgen.Hex(),
		MetricID:    input.MetricID,
		ThresholdID: input.ThresholdID,
		Level:       input.Level,
		Message:     input.Message,
		Status:      input.Status,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateAlert(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DispatchAlert(id string) (err error) {
	defer func() {
		if recover() != nil {
			err = nil
		}
	}()
	alert, err := s.store.GetAlertForDelivery(id)
	if err != nil {
		return err
	}
	alert.DeliveryAttempts++
	if !notify.Available(s.notifier) {
		return nil
	}
	if err := s.notifier.Send(alert.Message); err != nil {
		alert.MarkDeliveryFailure(err)
		_ = s.store.UpdateAlert(alert)
		return nil
	}
	alert.Status = model.AlertStatusAcknowledged
	alert.LastDeliveryError = ""
	return s.store.UpdateAlert(alert)
}

func (s *Service) GetAlert(id string) (*model.Alert, error) {
	return s.store.GetAlert(id)
}

func (s *Service) ListAlerts(filter model.AlertFilter, page, size int) ([]*model.Alert, int, error) {
	all := s.store.ListAlerts()
	matched := make([]*model.Alert, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Alert{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlert(id string, input model.Alert) (*model.Alert, error) {
	a, err := s.store.GetAlert(id)
	if err != nil {
		return nil, err
	}
	if input.Message != "" {
		a.Message = input.Message
	}
	if input.Level != "" {
		a.Level = input.Level
	}
	if input.Status != "" {
		if !model.AlertCanTransition(a.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态流转不合法")
		}
		a.Status = input.Status
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAlert(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAlert(id string) error {
	return s.store.DeleteAlert(id)
}

func (s *Service) BatchUpdateAlertStatus(ids []string, status string) error {
	if len(ids) == 0 {
		return model.NewValidationError("ids", "ID 列表不能为空")
	}
	validStatus := map[string]bool{model.AlertStatusOpen: true, model.AlertStatusAcknowledged: true, model.AlertStatusResolved: true}
	if !validStatus[status] {
		return model.NewValidationError("status", "告警状态不合法")
	}
	return s.store.BatchUpdateAlertStatus(ids, status)
}
