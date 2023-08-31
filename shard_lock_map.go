/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/31
 * Time: 11:52
 */

package vaedb

import "sync"

const (
	maxShardSize = 1024
)

type ShardMap struct {
	shards []*Shard
	hash   IHash
}

type Shard struct {
	val map[string]interface{}
	mux sync.RWMutex
}

func DefaultShardMap() *ShardMap {
	return NewShardMap(&Fnv32Hash{})
}

func NewShardMap(hash IHash) *ShardMap {
	s := &ShardMap{
		shards: make([]*Shard, maxShardSize),
		hash:   hash,
	}
	for i := 0; i < maxShardSize; i++ {
		s.shards[i] = &Shard{
			val: make(map[string]interface{}),
		}
	}
	return s
}

func (s *ShardMap) get(key string) interface{} {
	shard := s.getShard(key)
	shard.mux.RLock()
	defer shard.mux.RUnlock()
	v, ok := shard.val[key]
	if ok {
		return v
	}
	return nil
}

func (s *ShardMap) set(key string, value interface{}) {
	shard := s.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	shard.val[key] = value
}

func (s *ShardMap) del(key string) {
	shard := s.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	delete(shard.val, key)
}

func (s *ShardMap) getShard(key string) *Shard {
	return s.shards[s.hash.Hash(key)%uint32(maxShardSize)]
}
