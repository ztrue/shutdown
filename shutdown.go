// This package provides convenient interface for working with os.Signal.
//
// Multiple hooks can be applied,
// they will be called simultaneously on app shutdown.
package shutdown

import (
	"math/rand"
	"os"
	"os/signal"
	"sync"
)

var defaultShutdown = &Shutdown{
	hooks: map[string]func(os.Signal){},
	mutex: &sync.Mutex{},
}

// Shutdown is an instance of shutdown handler.
type Shutdown struct {
	hooks map[string]func(os.Signal)
	mutex *sync.Mutex
}

// Add adds a shutdown hook
// and returns hook identificator (key).
func Add(fn func()) string {
	return defaultShutdown.Add(fn)
}

// AddWithKey adds a shutdown hook
// with provided identificator (key).
func AddWithKey(key string, fn func()) {
	defaultShutdown.AddWithKey(key, fn)
}

// AddWithParam adds a shutdown hook with signal parameter
// and returns hook identificator (key).
func AddWithParam(fn func(os.Signal)) string {
	return defaultShutdown.AddWithParam(fn)
}

// AddWithKeyWithParam adds a shutdown hook with signal parameter
// with provided identificator (key).
func AddWithKeyWithParam(key string, fn func(os.Signal)) {
	defaultShutdown.AddWithKeyWithParam(key, fn)
}

// Hooks returns a copy of current hooks.
func Hooks() map[string]func(os.Signal) {
	return defaultShutdown.Hooks()
}

// Listen waits for provided OS signals.
// It will wait for any signal if no signals provided.
func Listen(signals ...os.Signal) {
	defaultShutdown.Listen(signals...)
}

// Remove cancels hook by identificator (key).
func Remove(key string) {
	defaultShutdown.Remove(key)
}

// Add adds a shutdown hook
// and returns hook identificator (key).
func (s *Shutdown) Add(fn func()) string {
	return s.AddWithParam(func(os.Signal) {
		fn()
	})
}

// AddWithKey adds a shutdown hook
// with provided identificator (key).
func (s *Shutdown) AddWithKey(key string, fn func()) {
	s.AddWithKeyWithParam(key, func(os.Signal) {
		fn()
	})
}

// AddWithParam adds a shutdown hook with signal parameter
// and returns hook identificator (key).
func (s *Shutdown) AddWithParam(fn func(os.Signal)) string {
	key := randomKey()
	s.AddWithKeyWithParam(key, fn)
	return key
}

// AddWithKeyWithParam adds a shutdown hook with signal parameter
// with provided identificator (key).
func (s *Shutdown) AddWithKeyWithParam(key string, fn func(os.Signal)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.hooks[key] = fn
}

// Hooks returns a copy of current hooks.
func (s *Shutdown) Hooks() map[string]func(os.Signal) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fns := map[string]func(os.Signal){}
	for key, cb := range s.hooks {
		fns[key] = cb
	}
	return fns
}

// Listen waits for provided OS signals.
// It will wait for any signal if no signals provided.
func (s *Shutdown) Listen(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	sig := <-ch
	var wg sync.WaitGroup
	for _, fn := range s.Hooks() {
		wg.Add(1)
		go func(sig os.Signal, fn func(os.Signal)) {
			defer wg.Done()
			fn(sig)
		}(sig, fn)
	}
	wg.Wait()
}

// Remove cancels hook by identificator (key).
func (s *Shutdown) Remove(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.hooks, key)
}

// randomKey generates a random identificator (key) for hook.
//
// Do not use this identificator for purposes other then to remove a hook
// as long as it's not fairly random without seed.
func randomKey() string {
	runes := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 16)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
