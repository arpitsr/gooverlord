package metastore

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *redis.Client
)

func newClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func getClient() *redis.Client {
	once.Do(func() {
		client = newClient()
	})
	return client
}

func Set(key string, value interface{}) error {
	client := getClient()
	return client.Set(context.Background(), key, value, 0).Err()
}

func Get(key string) (string, error) {
	conn := getClient()
	return conn.Get(context.Background(), key).Result()
}

// GetNodes for meili nodes from redis - this is part of alternate approach
// func GetNodes() ([]string, error) {
// 	var cursor uint64
// 	var results []string
// 	for {
// 		var keys []string
// 		var err error
// 		keys, cursor, err = getClient().SScan(context.TODO(), partitioner.MEILI_NODES_LIST, cursor, "*", 10).Result()
// 		if err != nil {
// 			panic(err)
// 		}
// 		results = append(results, keys...)
// 		if cursor == 0 {
// 			break
// 		}
// 	}
// 	return results, nil
// }
