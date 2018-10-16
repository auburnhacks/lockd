package config

import (
	"errors"
	"sync"

	"github.com/auburnhacks/lockd/lock"
)

var NotAvailable = errors.New("lock for service name not found")

// Config is a struct that manages the global config of the lockd application
type Config struct {
	mu sync.RWMutex
	// Locks is a map that maps a service name to a lock
	locks map[string]*lock.Lock
}

func New() *Config {
	return &Config{
		mu:    sync.RWMutex{},
		locks: map[string]*lock.Lock{},
	}
}

func (c *Config) GetLockWithServiceName(serviceName string) (*lock.Lock, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	l, ok := c.locks[serviceName]
	if !ok {
		return nil, NotAvailable
	}

	return l, nil
}

func (c *Config) SetLock(lock *lock.Lock) {
	c.locks[lock.ServiceName] = lock
}

func (c *Config) DeleteLock(lock *lock.Lock) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.locks, lock.ServiceName)
}
