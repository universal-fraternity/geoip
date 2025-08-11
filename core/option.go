// Package core provides data core logics for handling geo ip.
package core

// CallBackFunc Callback function format definition
type CallBackFunc func(meta *Meta) interface{}

// Option config option
type Option struct {
	Files []string
	CB    CallBackFunc
}
