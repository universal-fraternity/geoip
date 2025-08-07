// Package geoip provides location information query through IP
package geoip

import (
	"net"
	"sync"

	"github.com/universal-fraternity/geoip/store"
)

var (
	defaultStore *store.Store
	mu           sync.RWMutex
)

// Store output store.Store.
type Store = store.Store

// Meta output store.Meta。
type Meta = store.Meta

// Option output store.Option。
type Option = store.Option

// Load load data
func Load(opt Option) error {
	mu.Lock()
	defer mu.Unlock()
	if defaultStore == nil {
		defaultStore = store.NewStore()
	}
	return defaultStore.LoadData(opt)
}

// Update update data
func Update(fs ...string) error {
	return update()
}

func update(fs ...string) error {
	mu.Lock()
	defer mu.Unlock()
	if len(fs) > 0 {
		defaultStore.WithDataFiles(fs)
	}
	return defaultStore.Update()
}

// Search meta by address .
func Search(addr string) *Meta {
	mu.RLock()
	defer mu.RUnlock()
	return defaultStore.Search(net.ParseIP(addr))
}
