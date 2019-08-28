package mymutex

import (
	"errors"

	"github.com/vgarvardt/go-pg-adapter"
)

// MyMutex is the mutex lock based on MySQL user-level locks for Go
type MyMutex struct {
	adapter pgadapter.Adapter

	timeout int
}

// New instantiates and prepares MyMutex instance
func New(adapter pgadapter.Adapter, options ...Option) (*MyMutex, error) {
	instance := &MyMutex{
		adapter: adapter,
		timeout: -1,
	}

	for _, o := range options {
		o(instance)
	}

	return instance, nil
}

// Lock puts a named lock or waits until the the resource becomes available
func (m *MyMutex) Lock(name string) error {
	return m.adapter.Exec("SELECT GET_LOCK(?, ?)", name, m.timeout)
}

// Unlock releases a previously-acquired exclusive session level advisory lock
func (m *MyMutex) Unlock(name string) error {
	result := struct {
		Success int `db:"success"`
	}{}
	if err := m.adapter.SelectOne(&result, "SELECT RELEASE_LOCK(?) as success;", name); err != nil {
		return err
	}

	if result.Success != 1 {
		return errors.New("could not release a lock")
	}

	return nil
}

// TryLock puts a named lock immediately and returns true or returns false if the lock can not be acquired immediately.
// In some cases it may act as Lock() because implementation is not atomic.
func (m *MyMutex) TryLock(name string) (bool, error) {
	result := struct {
		IsFree int `db:"is_free"`
	}{}
	if err := m.adapter.SelectOne(&result, "SELECT IS_FREE_LOCK(?) as is_free;", name); err != nil {
		return false, err
	}

	if result.IsFree != 1 {
		return false, nil
	}

	return true, m.Lock(name)
}
