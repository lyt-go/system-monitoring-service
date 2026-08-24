package service

import (
	"sysmonitor/internal/config"
	"sysmonitor/internal/store"
	"sysmonitor/pkg/logger"
	"sysmonitor/pkg/notify"
)

type Service struct {
	store    store.Store
	log      *logger.Logger
	cfg      *config.Config
	notifier notify.Notifier
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

func (s *Service) SetNotifier(notifier notify.Notifier) {
	s.notifier = notifier
}
