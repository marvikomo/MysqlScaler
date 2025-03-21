package shard

import (
	"fmt"
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
