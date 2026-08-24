package service

import (
	"sort"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/ranking"
)

// StatsOverview 综合概览。
type StatsOverview struct {
	MetricCount     int `json:"metric_count"`
	CollectorCount  int `json:"collector_count"`
	ThresholdCount  int `json:"threshold_count"`
	SampleCount     int `json:"sample_count"`
	AlertCount      int `json:"alert_count"`
	OpenAlertCount  int `json:"open_alert_count"`
}

func (s *Service) GetStatsOverview() (*StatsOverview, error) {
	metrics := s.store.ListMetrics()
	collectors := s.store.ListCollectors()
	thresholds := s.store.ListThresholds()
	samples := s.store.ListSamples()
	alerts := s.store.ListAlerts()
	openAlerts := 0
	for _, a := range alerts {
		if a.Status == "open" {
			openAlerts++
		}
	}
	return &StatsOverview{
		MetricCount:    len(metrics),
		CollectorCount: len(collectors),
		ThresholdCount: len(thresholds),
		SampleCount:    len(samples),
		AlertCount:     len(alerts),
		OpenAlertCount: openAlerts,
	}, nil
}

// SampleCountByMetric 按指标统计样本数。
type SampleCountByMetric struct {
	MetricID string `json:"metric_id"`
	Count    int    `json:"count"`
}

func (s *Service) GetSampleCountByMetric() ([]SampleCountByMetric, error) {
	samples := s.store.ListSamples()
	counts := make(map[string]int)
	for _, sa := range samples {
		counts[sa.MetricID]++
	}
	result := make([]SampleCountByMetric, 0, len(counts))
	for mid, c := range counts {
		result = append(result, SampleCountByMetric{MetricID: mid, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

// SampleCountByHost 按主机统计样本数。
type SampleCountByHost struct {
	Host  string `json:"host"`
	Count int    `json:"count"`
}

func (s *Service) GetSampleCountByHost() ([]SampleCountByHost, error) {
	samples := s.store.ListSamples()
	counts := make(map[string]int)
	for _, sa := range samples {
		counts[sa.Host]++
	}
	result := make([]SampleCountByHost, 0, len(counts))
	for h, c := range counts {
		result = append(result, SampleCountByHost{Host: h, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

// TopAlertMetric TOP N 告警最多指标。
type TopAlertMetric struct {
	MetricID string `json:"metric_id"`
	Count    int    `json:"count"`
}

func (s *Service) GetTopAlertMetrics(n int) ([]TopAlertMetric, error) {
	if !ranking.ValidLimit(n) {
		return nil, model.NewValidationError("limit", "limit 必须为正整数")
	}
	alerts := s.store.ListAlerts()
	counts := make(map[string]int)
	for _, a := range alerts {
		counts[a.MetricID]++
	}
	result := make([]TopAlertMetric, 0, len(counts))
	for mid, c := range counts {
		result = append(result, TopAlertMetric{MetricID: mid, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if n < len(result) {
		result = result[:n]
	}
	return result, nil
}
