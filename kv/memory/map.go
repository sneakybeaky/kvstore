package memory

import "sync"

// Store is like a Go map[string]string but is safe for concurrent use by multiple goroutines
type Store struct {
	m     map[string]string
	mutex sync.RWMutex
}

func NewStore() *Store {
	return &Store{m: map[string]string{}}
}

func (s *Store) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.m[key]
	return value, ok
}
