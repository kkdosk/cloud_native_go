package ch04

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func callsCountFunction(callCounter *int) Effector {
	return func(ctx context.Context) (string, error) {
		*callCounter++
		return fmt.Sprintf("call %d", *callCounter), nil
	}
}

func TestThrottleMax1(t *testing.T) {
	const max uint = 1

	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, max, max, time.Second)

	for i := 0; i < 100; i++ {
		_, err := throttle(ctx)
		if err != nil {
			t.Log(err)
		}
	}

	if callsCounter == 0 {
		t.Error("test is broken; got", callsCounter)
	}

	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

func TestThrottleMax10(t *testing.T) {
	const max uint = 10

	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, max, max, time.Second)

	for i := 0; i < 100; i++ {
		throttle(ctx)
	}

	if callsCounter == 0 {
		t.Error("test is broken; got", callsCounter)
	}
	t.Logf("callsCounte: %d", callsCounter)
	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

func TestThrottleCallFrequency5Seconds(t *testing.T) {
	callsCounter := 0
	effector := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(effector, 1, 1, time.Second)

	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond).C

	for range ticker {
		tickCounts++

		s, e := throttle(ctx)
		if e != nil {
			t.Log("Error:", e)
		} else {
			t.Log("output:", s)
		}

		if tickCounts >= 20 {
			break
		}
	}

	if callsCounter != 5 {
		t.Error("expected 5; got", callsCounter)
	}
}
