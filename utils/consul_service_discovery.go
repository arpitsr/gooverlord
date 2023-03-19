package utils

import (
	"fmt"
	"sync"

	"com.ak.gooverlord/config"
	"github.com/hashicorp/consul/api"
)

type ConsulServiceDiscovery struct{}

var consulServiceDiscovery *ConsulServiceDiscovery
var onceConsulDiscover sync.Once

func (cn *ConsulServiceDiscovery) GetIndexNodes() []string {
	var ips []string
	// Create a new Consul API client
	conf := api.DefaultConfig()
	conf.Address = config.CONSUL_ADDRESS

	client, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}

	// Get a handle to the Catalog API
	catalog := client.Catalog()

	// Get a list of instances for the "my-service" service
	instances, _, err := catalog.Service(config.INDEXER_SERVICE_NAME, "indexer", nil)
	if err != nil {
		panic(err)
	}

	// Print the details of each instance
	for _, instance := range instances {
		ip := fmt.Sprintf("http://%s:%d", instance.Address, instance.ServicePort)
		ips = append(ips, ip)
	}
	return ips
}

func NewConsulDiscovery() *ConsulServiceDiscovery {
	onceConsulDiscover.Do(func() {
		consulServiceDiscovery = &ConsulServiceDiscovery{}
	})
	return consulServiceDiscovery
}
