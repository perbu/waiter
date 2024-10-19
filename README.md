# waiter

Package waiter provides a wrapper around `sync.WaitGroup` that allows to 
figure out which workers are still waiting.

This can make it easier to debug issues when there is one worker process that isn't ending.

## Usage

```go
import (
	"github.com/perbu/waiter"
)

wg := waiter.NewWaitGroupWithIDs()

// Add a worker with a unique ID
wg.Add("worker1") 

go func() {
    defer wg.Done("worker1")
    // ... do some work ...
}()

// List the IDs of workers that are still waiting
waiters := wg.ListWaiters() 

// Wait for all workers to finish
wg.Wait()
```