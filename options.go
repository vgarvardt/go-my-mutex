package mymutex

// Option is the configuration options type for MyMutex
type Option func(mutex *MyMutex)

// WithTimeout returns option that sets MyMutex timeout
func WithTimeout(timeout int) Option {
	return func(mutex *MyMutex) {
		mutex.timeout = timeout
	}
}
