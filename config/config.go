package config

import (
	"errors"
	"sync"
	"time"

	"github.com/auburnhacks/lockd/lock"
	"github.com/golang/glog"
)

var NotAvailable = errors.New("lock for service name not found")

// Config is a struct that manages the global config of the lockd application
type Config struct {
	mu sync.RWMutex
	// Locks is a map that maps a service name to a lock
	locks map[string]*lock.Lock
}

type Broker struct {
	EventChan      chan *Event
	NewClient      chan chan chan *Event
	ClosingClients chan chan *Event
	clients        map[chan *Event]bool
}

type Event interface {
	GetEventType() string
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

// Cleanup is function that runs in its own go routine and it cleans the map that
// holds all the locks
func (c *Config) Cleanup(d time.Duration) {
	for {
		glog.Infof("running cleanup...")
		c.mu.Lock()
		for _, l := range c.locks {
			// check if the lock has an expired status and then delete it from the map
			if l.IsTerminated {
				delete(c.locks, l.ServiceName)
			}
		}
		c.mu.Unlock()
		glog.Infof("finished cleanup...")
		glog.Infof("%v", c.locks)
		<-time.Tick(d)
	}
}
