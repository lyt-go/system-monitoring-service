package service

import (
	"context"
	"sync"

	"sysmonitor/internal/config"
	"sysmonitor/internal/store"
	"sysmonitor/pkg/logger"
)

type Service struct {
	store          store.Store
	log            *logger.Logger
	cfg            *config.Config
	sampleContexts sampleContextState
}

type sampleContextState struct {
	mu  sync.Mutex
	ctx context.Context
}

func (s *sampleContextState) take(current context.Context) context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx == nil {
		s.ctx = current
	}
	return s.ctx
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
