// Package geoip provides location information query through IP
package geoip

import (
	"net"
	"sync"

	"github.com/universal-fraternity/geoip/core"
)

var (
	defaultStore *core.Store
	mu           sync.RWMutex
)

// Store output core.Store.
type Store = core.Store

// Meta output core.Meta。
type Meta = core.Meta

// Option output core.Option。
type Option = core.Option

// Init init store and load data
func Init(opt Option) error {
	mu.Lock()
	defer mu.Unlock()
	if defaultStore == nil {
		defaultStore = core.NewStore()
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
