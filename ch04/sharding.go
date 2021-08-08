package ch04

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	m map[string]interface{}
}

type ShardedMap []*Shard

func NewShardedMap(nshards int) ShardedMap {
	shards := make([]*Shard, nshards)

	for i := 0; i < nshards; i++ {
		shard := make(map[string]interface{})
		shards[i] = &Shard{m: shard}
	}
	return shards
}

func (m ShardedMap) getShardIndex(key string) int {
	hash := sha1.Sum([]byte(key))

	return int(hash[17]) % len(m)
}

func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardedMap) Delete(key string) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

func (m ShardedMap) Get(key string) interface{} {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

func (m ShardedMap) Keys() []string {
	keys := make([]string, 0)
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard) {
			s.RLock()

			for key := range s.m {
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}
			s.RUnlock()
			wg.Done()
		}(shard)
	}
	wg.Wait()

	return keys
}
