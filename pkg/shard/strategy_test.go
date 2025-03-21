package shard

import (
	"testing"
)

func TestConsistentHashStrategy_GetShardID(t *testing.T) {
	//Setyp test data
	shardIDs := []string{"shard1", "shard2", "shard3"}
	strategy := NewConsistentHashStrategy(shardIDs, 10)

	//Test key consistency

	key := "user:1001"

	shardID, err := strategy.GetShardID(key)

	//check the results
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if shardID == "" {
		t.Error("expected non-empty shard ID")
	}

	// Verify key always maps to the same shard
	for i := 0; i < 10; i++ {
		nextShardID, err := strategy.GetShardID(key)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if nextShardID != shardID {
			t.Errorf("expected consistent shard mapping, got %s then %s",
				shardID, nextShardID)
		}
	}

}
