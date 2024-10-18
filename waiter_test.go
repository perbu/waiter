package waiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWaitGroupWithIDs(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		wg := &WaitGroupWithIDs{}
		wg.Add("worker1")
		wg.Add("worker2")

		go func() {
			time.Sleep(time.Millisecond)
			wg.Done("worker1")
		}()

		go func() {
			time.Sleep(time.Millisecond)
			wg.Done("worker2")
		}()

		wg.Wait()
	})

	t.Run("Wait with remaining workers", func(t *testing.T) {
		iWg := sync.WaitGroup{}
		iWg.Add(1)
		wg := &WaitGroupWithIDs{}
		wg.Add("worker1")
		wg.Add("worker2")

		go func() {
			time.Sleep(time.Millisecond)
			wg.Done("worker1")
			iWg.Done() // release
		}()

		go func() {
			time.Sleep(time.Millisecond * 10)
			wg.Done("worker2")
		}()
		go func() {
			wg.Wait()
		}()
		iWg.Wait()
		l1 := wg.ListWaiters()
		time.Sleep(time.Millisecond * 20)
		l2 := wg.ListWaiters()
		if len(l1) != 1 {
			t.Fatalf("Expected 1 waiter, got %d", len(l1))
		}
		if len(l2) != 0 {
			t.Fatalf("Expected 0 waiters, got %d", len(l2))
		}

	})

	t.Run("Concurrent access", func(t *testing.T) {
		wg := &WaitGroupWithIDs{}
		const numWorkers = 100

		for i := 0; i < numWorkers; i++ {
			wg.Add(fmt.Sprintf("worker%d", i))
		}

		var wg2 sync.WaitGroup
		wg2.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go func(id string) {
				defer wg2.Done()
				wg.Done(id)
			}(fmt.Sprintf("worker%d", i))
		}
		wg2.Wait()

		wg.Wait()
	})
}
