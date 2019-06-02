package errors

import (
	"sync"
)

// Channel provides errors chan and plumbing for error propogation
type Channel struct {
	mu   sync.RWMutex
	done bool
	once sync.Once

	errChan chan error
	subs    []chan error
}

// Errors get error chan
func (e *Channel) Errors() <-chan error {
	e.once.Do(e.init)

	ch := make(chan error)

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.done {
		close(ch)
	} else {
		e.subs = append(e.subs, ch)
	}
	return ch
}

// Publish propagate error to all subscribed channels
func (e *Channel) Publish(err error) {
	e.once.Do(e.init)

	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.done {
		e.errChan <- err
	}
}

// Close stop accepting new errors
func (e *Channel) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.errChan != nil {
		close(e.errChan)
		e.done = true
	}
}

func (e *Channel) init() {
	e.errChan = make(chan error)
	e.subs = []chan error{}

	go e.propagate()
}

func (e *Channel) propagate() {
	wg := &sync.WaitGroup{}

	for {
		err, ok := <-e.errChan
		if !ok {
			wg.Wait()
			e.mu.RLock()

			for _, ch := range e.subs {
				close(ch)
			}
			e.mu.RUnlock()
			return
		}

		e.mu.RLock()
		for _, ch := range e.subs {
			wg.Add(1)

			go func(ch chan error) {
				ch <- err
				wg.Done()
			}(ch)
		}
		e.mu.RUnlock()
	}
}
