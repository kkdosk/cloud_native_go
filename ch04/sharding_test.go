package ch04

import (
	"testing"
)

func TestShardingGetShardIndex(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)
	counts := make([]int, BUCKETS)

	keys := []string{"A", "B", "C", "D", "E"}

	for _, key := range keys {
		idx := sMap.getShardIndex(key)
		counts[idx]++
		t.Log(key, idx)
		if counts[idx] > 1 {
			t.Error("Hash collision")
		}
	}
}

func TestShardingSetAndGet(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Set(k, v)
	}

	for k, v := range truthMap {
		got := sMap.Get(k)

		if got != v {
			t.Errorf("Key mismatch on %s: expected %d, got %d", k, v, got)
		}
	}
}

func TestShardingKeys(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Set(k, v)
	}

	keys := sMap.Keys()

	if len(truthMap) != len(keys) {
		t.Error("Map/keys mismach")
	}

	for _, key := range sMap.Keys() {
		if _, ok := truthMap[key]; !ok {
			t.Error("Key", key, "not in truthMap")
		}

		delete(truthMap, key)
	}

	if len(truthMap) != 0 {
		t.Error("Key mismatch")
	}
}

func TestShardingDelete(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Set(k, v)
	}

	keys := sMap.Keys()
	for _, key := range keys {
		sMap.Delete(key)
	}

	if len(sMap.Keys()) != 0 {
		t.Error("Deletion failure")
	}

}
