package lock

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

const (
	unlocked uint32 = iota
	locked

	LockAcquired string = "acquired"
	LockReleased string = "released"
	LockStandby  string = "standby"
)

var (
	errorLocked = errors.New("lock already aquired")
)

type Lock struct {
	ServiceName  string        `json:"service_name"`
	TTL          time.Duration `json:"ttl"`
	IsTerminated bool          `json:"is_terminated"`
	locker       uint32
}

func NewLock(serviceName string, ttl time.Duration) *Lock {
	return &Lock{
		ServiceName:  serviceName,
		TTL:          ttl,
		IsTerminated: false,
	}
}

func (l *Lock) Lock() error {
	if !atomic.CompareAndSwapUint32(&l.locker, unlocked, locked) {
		return errorLocked
	}
	// run a go routine in the background and  unlock after a certain duration
	go l.notify(l.TTL)
	return nil
}

func (l *Lock) Unlock() {
	defer atomic.StoreUint32(&l.locker, unlocked)
}

func (l *Lock) notify(d time.Duration) {
	select {
	// TODO: will have to change this to a better notification
	case <-time.Tick(d):
		glog.Infof("lock %s expired", l.ServiceName)
		l.IsTerminated = true
		l.Unlock()
	}
}

type LockEvent struct {
	EventType string
	Lock      *Lock
}

func (le *LockEvent) GetEventType() {
	return le.EventType
}
