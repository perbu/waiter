// Package waiter provides a wrapper around sync.WaitGroup that allows to figure out which workers are still waiting.
package waiter

import (
	"sync"
)

type WaitGroupWithIDs struct {
	mu  sync.Mutex
	ids map[string]struct{}
	wg  sync.WaitGroup
}

func (wg *WaitGroupWithIDs) Add(id string) {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	if wg.ids == nil {
		wg.ids = make(map[string]struct{})
	}
	// check if the id is already in the map
	if _, ok := wg.ids[id]; ok {
		panic("id already exists")
	}
	wg.ids[id] = struct{}{}
	wg.wg.Add(1)
}

func (wg *WaitGroupWithIDs) Done(id string) {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	if _, ok := wg.ids[id]; !ok {
		panic("id does not exist")
	}
	delete(wg.ids, id)
	wg.wg.Done()
}

func (wg *WaitGroupWithIDs) ListWaiters() []string {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	waiters := make([]string, 0, len(wg.ids))
	for id := range wg.ids {
		waiters = append(waiters, id)
	}
	return waiters
}

func (wg *WaitGroupWithIDs) Wait() {
	wg.wg.Wait()
}
