package shard

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

// ShardingStrategy defines the interface for sharding strategies

type ShardingStrategy interface {
	// GetShardID determines which shard a data item belongs to based on sharding key
	GetShardID(Key string) (string, error)
}

type ConsistentHashStrategy struct {
	mu           sync.RWMutex
	ring         []uint32
	shardMap     map[uint32]string // Maps position to shard ID
	virtualNodes int               // Number of virtual nodes per shard
}

// NewConsistentHashStrategy creates a new consistent hashing strategy
func NewConsistentHashStrategy(shardIDs []string, virtualNodes int) *ConsistentHashStrategy {
	c := &ConsistentHashStrategy{
		shardMap:     make(map[uint32]string),
		virtualNodes: virtualNodes,
		ring:         make([]uint32, 0),
	}

	// Add all shards to the hash ring
	for _, shardID := range shardIDs {
		c.AddShard(shardID)
	}

	return c
}

func (c *ConsistentHashStrategy) AddShard(shardID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < c.virtualNodes; i++ {
		virtualNodeKey := fmt.Sprintf("%s-%d", shardID, i)
		hash := c.hashKey(virtualNodeKey)
		c.shardMap[hash] = shardID
		c.ring = append(c.ring, hash)

	}

	//sort ring

	sort.Slice(c.ring, func(i, j int) bool {
		return c.ring[i] < c.ring[j]
	})

}

// hashKey hashes a key to a uint32 value
func (c *ConsistentHashStrategy) hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// GetShardID implements ShardingStrategy.GetShardID
func (c *ConsistentHashStrategy) GetShardID(key string) (string, error) {
	if len(c.ring) == 0 {
		return "", fmt.Errorf("no shards available")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	hash := c.hashKey(key)
	idx := c.search(hash)

	// If the hash is greater than the largest hash in the ring, wrap around
	if idx == len(c.ring) {
		idx = 0
	}

	return c.shardMap[c.ring[idx]], nil
}

// search finds the first index in the ring with a hash >= the given hash
func (c *ConsistentHashStrategy) search(hash uint32) int {
	// Binary search the ring
	return sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hash
	})
}
