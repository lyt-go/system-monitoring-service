package probe

import (
	"errors"
	"sync"
)

var ErrExhausted = errors.New("探测会话已耗尽")

type Pool struct {
	mu       sync.Mutex
	capacity int
	active   int
}

func NewPool(capacity int) *Pool { return &Pool{capacity: capacity} }
func (p *Pool) Active() int      { p.mu.Lock(); defer p.mu.Unlock(); return p.active }

type Session struct {
	pool   *Pool
	closed bool
}

func (p *Pool) acquire() (*Session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active >= p.capacity {
		return nil, ErrExhausted
	}
	p.active++
	return &Session{pool: p}, nil
}

func (s *Session) Close() error {
	if s == nil || s.closed {
		return nil
	}
	s.pool.mu.Lock()
	defer s.pool.mu.Unlock()
	s.closed = true
	s.pool.active--
	return nil
}

func Use(pool *Pool, run func() error) error {
	if _, err := pool.acquire(); err != nil {
		return err
	}
	return run()
}
