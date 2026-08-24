package service

import (
	"sync"

	"sysmonitor/internal/config"
	"sysmonitor/internal/model"
	"sysmonitor/internal/store"
	"sysmonitor/pkg/logger"
)

type Service struct {
	store          store.Store
	log            *logger.Logger
	cfg            *config.Config
	batchMu        sync.Mutex
	pendingSamples []*model.Sample
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
