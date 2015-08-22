package semaphore

import (
	"testing"
	"time"
	"sync"
	//"fmt"
)

var (
	m = sync.Mutex { }
)

func TestSemaphore(t *testing.T) {
	s := New()
	s.Signal()
	s.Signal()

	// Acquire should return true twice as we've had 2 signal calls (Count = 2)
	if !(s.TryAcquire() && s.TryAcquire()) {
		t.Error("Should have been able to acquire twice")
	}

	// Count = 0
	s.Signal()	// Count = 1
	s.Signal()	// Count = 2

	// Wait should return immediately twice (we'll see a deadlock panic otherwise)
	s.Wait()	// Count = 1
	s.Wait()	// Count = 0
	
	go DelayedSignal(s, 2)
	s.Wait()

	// Delay so 
	m.Lock()
	defer m.Unlock()

	// Test that we can acquire here again as we called signal twice
	if !(s.TryAcquire()) {
		t.Error("Should have been able to acquire here")
	}
}

func DelayedSignal(sem *CountingSemaphore, signalCount int) {
	time.Sleep(100 * time.Millisecond)

	// Lock so all signals complete before control returns to TestSemaphore
	m.Lock()
	defer m.Unlock()

	for i := 0; i < signalCount; i++ {
		sem.Signal()
	}
}

