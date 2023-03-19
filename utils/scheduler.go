package utils

import (
	"time"

	"com.ak.gooverlord/indexer"
	"com.ak.gooverlord/partitioner"
)

func init() {
	schRun()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for _ = range ticker.C {
			schRun()
		}
	}()
}

func schRun() {
	// ips := k8sDiscovery.GetIndexNodes()
	ips := consulServiceDiscovery.GetIndexNodes()
	go func(ips []string) {
		newCHR := partitioner.NewConsistenHashRing()
		for _, ip := range ips {
			newCHR.AddNode(ip)
		}
		partitioner.GetConsistentHashRing().RWLock.Lock()
		defer partitioner.GetConsistentHashRing().RWLock.Unlock()
		partitioner.CHR = newCHR
		indexer.UpdateInstance()
	}(ips)
}
