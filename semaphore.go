package semaphore

import (
	"sync"
)

// ------------------------------------------------------------------------------------------------------------------------
// Struct: CountingSemaphore
// ------------------------------------------------------------------------------------------------------------------------

// CountingSemaphore is a locking mechanism which counts Wait/Signal calls and only locks if Count goes below zero
//
// For example:
// 	CountingSemaphore.Signal() // Count = 1
// 	CountingSemaphore.Signal() // Count = 2
//  CountingSemaphore.Wait() // Count = 1, the method returns immediately
// 	CountingSemaphore.Wait() // Count = 0, the method returns immediately
//  CountingSemaphore.Wait() // Count = -1, this call will block until Signal is called and the count returns to zero
type CountingSemaphore struct {
	
	// Count is the current count of Wait/Signal. Signal increments, Wait decrements
	Count int

	// SyncLock is used to synchronise the methods of this objects
	SyncLock sync.Mutex

	WaitingChan chan bool
}

// Wait decrements the Count by 1 and will pause/block on the underlying lock if the count goes below zero
func (this *CountingSemaphore) Wait() {
	this.SyncLock.Lock()
	this.Count--
	
	if this.Count < 0 {
		this.SyncLock.Unlock()
		<- this.WaitingChan
		this.SyncLock.Lock()
	}

	this.SyncLock.Unlock()
}

// TryAcquire is non blocking. It'll return true and decrement the count if its currently positive, false if Count <= 0
func (this *CountingSemaphore) TryAcquire() bool {
	this.SyncLock.Lock()
	defer this.SyncLock.Unlock()

	if this.Count < 1 {
		return false
	} else {
		this.Count--
		return true
	}
}

// Signal increments the Count by 1 and will unblock a single goroutine blocked on the Wait call
func (this *CountingSemaphore) Signal() {
	this.SyncLock.Lock()
	defer this.SyncLock.Unlock()

	if this.Count < 0 {
		this.WaitingChan <- true
	}

	this.Count++
}

// ------------------------------------------------------------------------------------------------------------------------
// Construction
// ------------------------------------------------------------------------------------------------------------------------

// New creates a new CountingSemaphore
func New() *CountingSemaphore {
	cs := &CountingSemaphore{ WaitingChan: make(chan bool)}
	return cs
}