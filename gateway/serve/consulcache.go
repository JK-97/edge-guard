package serve

import (
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/patrickmn/go-cache"
)

type CousulCache struct {
	cache  *cache.Cache
	mux    *sync.Mutex
	client *consulapi.Client
}

func (c *CousulCache) Get(name string) (interface{}, bool) {
	return c.cache.Get(name)
}
func (c *CousulCache) Delete(name string) {
	c.cache.Delete(name)
}

func (c *CousulCache) Set(name string) bool {
	return true
}

func (c *CousulCache) syncRemoteService() {
	catalog := c.client.Catalog()

	catalog.Service("", "ai_service", &consulapi.QueryOptions{
		UseCache: true,
		MaxAge:   3 * time.Hour,
	})
}
