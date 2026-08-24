// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"sysmonitor/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Metric
	CreateMetric(m *model.Metric) error
	GetMetric(id string) (*model.Metric, error)
	GetMetricByName(name string) (*model.Metric, error)
	ListMetrics() []*model.Metric
	UpdateMetric(m *model.Metric) error
	DeleteMetric(id string) error

	// Collector
	CreateCollector(c *model.Collector) error
	GetCollector(id string) (*model.Collector, error)
	GetCollectorByName(name string) (*model.Collector, error)
	ListCollectors() []*model.Collector
	UpdateCollector(c *model.Collector) error
	DeleteCollector(id string) error

	// Threshold
	CreateThreshold(t *model.Threshold) error
	GetThreshold(id string) (*model.Threshold, error)
	ListThresholds() []*model.Threshold
	UpdateThreshold(t *model.Threshold) error
	DeleteThreshold(id string) error

	// Sample
	CreateSample(s *model.Sample) error
	CommitSampleBatch(samples []*model.Sample) error
	GetSample(id string) (*model.Sample, error)
	ListSamples() []*model.Sample
	DeleteSample(id string) error
	BatchDeleteSamples(ids []string) error

	// Alert
	CreateAlert(a *model.Alert) error
	GetAlert(id string) (*model.Alert, error)
	ListAlerts() []*model.Alert
	UpdateAlert(a *model.Alert) error
	BatchUpdateAlertStatus(ids []string, status string) error
	DeleteAlert(id string) error
}
