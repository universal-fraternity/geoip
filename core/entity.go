// Package core provides data core logics for handling geo ip.
package core

// IPV4Entity Define the entity of IPv4 data.
type IPV4Entity struct {
	ipIndex   uint32
	metaIndex uint32
}

// IPV6Entity Define the entity of IPv6 data.
type IPV6Entity struct {
	ipIndex   uint64
	metaIndex uint32
}

// IPIndex start index
func (i *IPV4Entity) IPIndex() uint32 {
	return i.ipIndex
}

// IPIndex start index
func (i *IPV6Entity) IPIndex() uint64 {
	return i.ipIndex
}
