package infrastructure

import "github.com/go-redsync/redsync/v4"

type MutexProvider interface {
	Acquire(key string) (mutex *redsync.Mutex, err error)
	Release(mutex *redsync.Mutex) (err error)
}

type mutexProvider struct {
	provider redsync.Redsync
}

func NewMutexProvider(provider redsync.Redsync) MutexProvider {
	return &mutexProvider{
		provider: provider,
	}
}

func (mutex *mutexProvider) Acquire(key string) (res *redsync.Mutex, err error) {
	activityMutex := mutex.provider.NewMutex(key)

	if err := activityMutex.Lock(); err != nil {
		return nil, err
	}

	return activityMutex, nil
}

func (mutex *mutexProvider) Release(_mutex *redsync.Mutex) (err error) {
	if ok, err := _mutex.Unlock(); !ok || err != nil {
		return err
	}
	return nil
}
