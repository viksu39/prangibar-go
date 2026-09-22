package middleware

import (
	"sync"
	"time"

	"prangibar-go/config"
)

type pinAttempt struct {
	count    int
	lastSeen time.Time
	lockedAt time.Time
}

var (
	pinAttempts = make(map[string]*pinAttempt)
	pinMu       sync.Mutex
	pinInit     = false
)

func ensurePinInit() {
	if !pinInit {
		pinInit = true
		go pinCleanupRoutine()
	}
}

func pinCleanupRoutine() {
	for {
		time.Sleep(5 * time.Minute)
		lockoutDur := time.Duration(config.App.Security.PinLockout.LockoutMinutes) * time.Minute
		window := time.Duration(config.App.Security.PinLockout.WindowMinutes) * time.Minute
		pinMu.Lock()
		for key, a := range pinAttempts {
			if time.Since(a.lastSeen) > window && time.Since(a.lockedAt) > lockoutDur {
				delete(pinAttempts, key)
			}
		}
		pinMu.Unlock()
	}
}

func RecordPinFailure(key string) bool {
	ensurePinInit()
	pinMu.Lock()
	defer pinMu.Unlock()

	attempt, exists := pinAttempts[key]
	maxAttempts := config.App.Security.PinLockout.MaxAttempts
	window := time.Duration(config.App.Security.PinLockout.WindowMinutes) * time.Minute

	if !exists {
		pinAttempts[key] = &pinAttempt{
			count:    1,
			lastSeen: time.Now(),
		}
		return false
	}

	if time.Since(attempt.lastSeen) > window {
		attempt.count = 1
		attempt.lastSeen = time.Now()
		return false
	}

	attempt.count++
	attempt.lastSeen = time.Now()

	if attempt.count >= maxAttempts {
		attempt.lockedAt = time.Now()
		return true
	}
	return false
}

func ResetPinAttempts(key string) {
	pinMu.Lock()
	defer pinMu.Unlock()
	delete(pinAttempts, key)
}

func IsPinLocked(key string) (bool, time.Duration) {
	ensurePinInit()
	pinMu.Lock()
	defer pinMu.Unlock()

	attempt, exists := pinAttempts[key]
	if !exists {
		return false, 0
	}

	maxAttempts := config.App.Security.PinLockout.MaxAttempts
	lockoutDur := time.Duration(config.App.Security.PinLockout.LockoutMinutes) * time.Minute

	if attempt.count >= maxAttempts {
		remaining := lockoutDur - time.Since(attempt.lockedAt)
		if remaining > 0 {
			return true, remaining
		}
		attempt.count = 0
		attempt.lockedAt = time.Time{}
	}
	return false, 0
}
