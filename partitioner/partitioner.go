package partitioner

import (
	"errors"
	"hash/fnv"
	"log"
	"sort"
	"strconv"
	"sync"
)

type virtualNode struct {
	nodeID string
	hash   uint32
}

type ConsistentHashRing struct {
	virtualNodes         []virtualNode
	virtualToRealNodeMap map[uint32]string
	RealNodesSet         map[string]struct{}
	RWLock               sync.RWMutex
}

var (
	once sync.Once
	CHR  *ConsistentHashRing
)

func GetConsistentHashRing() *ConsistentHashRing {
	once.Do(func() {
		CHR = NewConsistenHashRing()
	})
	return CHR
}

func (c *ConsistentHashRing) GetHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (c *ConsistentHashRing) AddNode(nodeID string) {
	for i := 0; i < 3; i++ {
		hash := c.GetHash(nodeID + strconv.Itoa(i))
		c.virtualNodes = append(c.virtualNodes, virtualNode{nodeID, hash})
		c.virtualToRealNodeMap[hash] = nodeID
		c.RealNodesSet[nodeID] = struct{}{}
	}
	sort.Slice(c.virtualNodes, func(i, j int) bool {
		return c.virtualNodes[i].hash < c.virtualNodes[j].hash
	})
}

func (c *ConsistentHashRing) RemoveNode(nodeID string) {
	var newVirtualNodes []virtualNode
	for _, virtualNode := range c.virtualNodes {
		if virtualNode.nodeID != nodeID {
			newVirtualNodes = append(newVirtualNodes, virtualNode)
		} else {
			delete(c.RealNodesSet, c.virtualToRealNodeMap[virtualNode.hash])
			delete(c.virtualToRealNodeMap, virtualNode.hash)
		}
	}
	c.virtualNodes = newVirtualNodes
}

func (c *ConsistentHashRing) GetNode(key string) (string, error) {
	log.Printf("ConsistentHashRing %v\n", c.RealNodesSet)
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()
	hash := c.GetHash(key)
	i := sort.Search(len(c.virtualNodes), func(i int) bool {
		return c.virtualNodes[i].hash >= hash
	})
	if i == len(c.virtualNodes) {
		i = 0
	}
	if len(c.virtualNodes) == 0 {
		return "", errors.New("no nodes up")
	}
	return c.virtualToRealNodeMap[c.virtualNodes[i].hash], nil
}

func NewConsistenHashRing() *ConsistentHashRing {
	return &ConsistentHashRing{
		virtualNodes:         make([]virtualNode, 0),
		virtualToRealNodeMap: make(map[uint32]string, 0),
		RealNodesSet:         make(map[string]struct{}, 0),
	}
}
