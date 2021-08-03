package ch04

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := SlowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 2)
}

func elapsedCheck(t *testing.T, start time.Time, seconds int) {
	elapsed := int(time.Now().Sub(start).Seconds())

	if seconds != elapsed {
		t.Errorf("expected %d seconds; got %d\n", seconds, elapsed)
	}
}

func TestFutureGetTwice(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := SlowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 2)

	start = time.Now()
	res, err = future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 0)
}

func TestFutureConcurrent(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := SlowFunction(ctx)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			res, err := future.Result()
			if err != nil {
				t.Error(err)
				return
			}
			if !strings.HasPrefix(res, "I slept for") {
				t.Error("unexpected output:", res)
			}

			elapsedCheck(t, start, 2)
		}()
	}
	wg.Wait()
}

func TestFutureTimeout(t *testing.T) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	defer cancel()

	future := SlowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		if !strings.Contains(err.Error(), "deadline") {
			t.Error("received unexpected error maybe:", err)
		}
	}

	if res != "" {
		t.Error("should have an empty result")
	}

	elapsedCheck(t, start, 1)
}

func TestFutureCancel(t *testing.T) {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	future := SlowFunction(ctx)

	res, err := future.Result()
	if err != nil {
		if !strings.Contains(err.Error(), "canceled") {
			t.Error("received unexpected error maybe: ", err)
		}
	}

	if res != "" {
		t.Error("should have an empty result")
	}

	elapsedCheck(t, start, 1)
}
