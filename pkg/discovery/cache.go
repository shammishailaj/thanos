package discovery

import (
	"sync"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

// Cache is a store for target groups. It provides thread safe updates and a way for obtaining all addresses from
// the stored target groups.
type Cache struct {
	tgs map[string]*targetgroup.Group
	sync.Mutex
}

// NewCache returns a new empty Cache.
func NewCache() *Cache {
	return &Cache{
		tgs: make(map[string]*targetgroup.Group),
	}
}

// Update stores the targets for the given groups.
// Note: targets for a group are replaced entirely on update. If a group with no target is given this is equivalent to
// deleting all the targets for this group.
func (c *Cache) Update(tgs []*targetgroup.Group) {
	c.Lock()
	defer c.Unlock()
	for _, tg := range tgs {
		// Some Discoverers send nil target group so need to check for it to avoid panics.
		if tg == nil {
			continue
		}
		c.tgs[tg.Source] = tg
	}
}

// Addresses returns all the addresses from all target groups present in the Cache.
func (c *Cache) Addresses() []string {
	c.Lock()
	defer c.Unlock()
	var addresses []string
	for _, group := range c.tgs {
		for _, target := range group.Targets {
			addresses = append(addresses, string(target[model.AddressLabel]))
		}
	}
	return addresses
}
