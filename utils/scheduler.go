package utils

import (
	"fmt"
	"time"

	"com.ak.gooverlord/indexer"
	"com.ak.gooverlord/partitioner"
)

func init() {
	schRun()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for _ = range ticker.C {
			schRun()
		}
	}()
}

func schRun() {
	ips := k8sDiscovery.GetIndexNodes()
	go func(ips []string) {
		newCHR := partitioner.NewConsistenHashRing()
		for _, ip := range ips {
			ip = fmt.Sprintf("http://%s:7700", ip)
			newCHR.AddNode(ip)
		}
		partitioner.GetConsistentHashRing().RWLock.Lock()
		defer partitioner.GetConsistentHashRing().RWLock.Unlock()
		partitioner.CHR = newCHR
		indexer.UpdateInstance()
	}(ips)
}
