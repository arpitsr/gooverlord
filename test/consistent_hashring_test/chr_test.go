package chr_test

import (
	"testing"

	"com.ak.gooverlord/partitioner"
)

func TestConsistentHashing(t *testing.T) {
	nodes := []string{"10.0.0.1", "10.0.0.2"}
	for _, node := range nodes {
		partitioner.GetConsistentHashRing().AddNode(node)
	}
	key1 := "key1"
	key2 := "newkey1"

	val, err := partitioner.GetConsistentHashRing().GetNode(key1)
	if err != nil {
		t.Fatalf("%s", err)
	}

	val2, err := partitioner.GetConsistentHashRing().GetNode(key2)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if val != nodes[0] {
		t.Fatalf("Expected Node: %s, Result Node: %s", nodes[0], val)
	}
	if val2 != nodes[1] {
		t.Fatalf("Expected Node: %s, Result Node: %s", nodes[1], val2)
	}
}
