package service

import (
	"sysmonitor/internal/config"
	"sysmonitor/internal/store"
	"sysmonitor/pkg/logger"
)

type Service struct {
	store          store.Store
	log            *logger.Logger
	cfg            *config.Config
	collectorProbe func(string) error
}

func (s *Service) SetCollectorProbeRunner(run func(string) error) {
	s.collectorProbe = run
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
