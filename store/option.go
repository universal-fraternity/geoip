// Package store provides data store definitions
package store

// CallBackFunc Callback function format definition
type CallBackFunc func(meta *Meta) interface{}

// Option config option
type Option struct {
	Files []string
	CB    CallBackFunc
}
