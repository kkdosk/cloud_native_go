package ch04

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var tokens = max
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						m.Lock()
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
						m.Unlock()
					}
				}
			}()
		})

		m.Lock()
		defer m.Unlock()

		if tokens <= 0 {

			return "", fmt.Errorf("too many calls")
		}

		tokens--

		return e(ctx)
	}
}

func TestThrottleVariableRefill(t *testing.T) {
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, 4, 2, 500*time.Millisecond)

	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond)
	timer := time.NewTimer(2 * time.Second)

time:
	for {
		select {
		case <-ticker.C:
			tickCounts++

			s, e := throttle(ctx)
			if e != nil {
				t.Log("Error:", e)
			} else {
				t.Log("output:", s)
			}
		case <-timer.C:
			break time
		}
	}

	if callsCounter != 8 {
		t.Error("expected 8; got", callsCounter)
	}
}

func TestThrottleContextTimeout(t *testing.T) {
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Microsecond)
	defer cancel()

	throttle := Throttle(effector, 1, 1, time.Second)

	s, e := throttle(ctx)
	if e != nil {
		t.Error("unexpected error:", e)
	} else {
		t.Log("output:", s)
	}

	time.Sleep(300 * time.Millisecond)

	_, e = throttle(ctx)
	if e != nil {
		t.Log("got expected error:", e)
	} else {
		t.Error("didn't get expected error")
	}
}
