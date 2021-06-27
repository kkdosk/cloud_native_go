package ch04

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func counter() Circuit {
	m := sync.Mutex{}
	count := 0

	return func(ctx context.Context) (string, error) {
		m.Lock()
		count++
		m.Unlock()

		return fmt.Sprintf("%d", count), nil
	}
}

func failAfter(threshold int) Circuit {
	count := 0

	return func(ctx context.Context) (string, error) {
		count++

		if count > threshold {
			return "", errors.New("INTENTIONAL FAIL!")
		}
		return "Success", nil
	}
}

func waitAndContinue() Circuit {
	return func(ctx context.Context) (string, error) {
		time.Sleep(time.Second)

		if rand.Int()%2 == 0 {
			return "success", nil
		}

		return "Failed", fmt.Errorf("foced failure")
	}
}

func TestCircuitBreakerFailAfter(t *testing.T) {
	circuit := failAfter(5)
	ctx := context.Background()

	for count := 1; count <= 5; count++ {
		_, err := circuit(ctx)

		t.Logf("attempt %d: %v", count, err)

		switch {
		case count <= 5 && err != nil:
			t.Error("expected no error; got", err)
		case count > 5 && err == nil:
			t.Error("expected err; got none")
		}
	}
}

func TestCircuitBreaker(t *testing.T) {
	circuit := failAfter(5)

	breaker := Breaker(circuit, 1)

	ctx := context.Background()

	circuitOpen := false
	doesCircuitOpen := false
	doesCircuitReclose := false
	count := 0

	for range time.NewTicker(time.Second).C {
		_, err := breaker(ctx)

		if err != nil {
			if strings.HasPrefix(err.Error(), "service unreachable") {
				if !circuitOpen {
					circuitOpen = true
					doesCircuitOpen = true
					t.Log("circuit has opend")
				}
			} else {
				if circuitOpen {
					circuitOpen = false
					doesCircuitReclose = true

					t.Log("circuit has automatically closed")
				}
			}
		} else {
			t.Log("circuit closed and operational")
		}

		count++
		if count >= 10 {
			break
		}
	}

	if !doesCircuitOpen {
		t.Error("circuit didn't appear to open")
	}

	if !doesCircuitReclose {
		t.Error("circuit didn't appear to close after time")
	}
}

func TestCircuitBreakerDataRace(t *testing.T) {
	ctx := context.Background()
	var threshold time.Time
	t.Log(time.Now().Before(threshold))
	circuit := waitAndContinue()
	breaker := Breaker(circuit, 1)

	wg := sync.WaitGroup{}

	for count := 1; count <= 20; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)

			_, err := breaker(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}
	wg.Wait()
}
