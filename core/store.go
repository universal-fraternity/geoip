// Package core provides data core logics for handling geo ip.
package core

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"

	"github.com/universal-fraternity/geoip/utils"
)

// Store Define the storage area for storing IP location data.
type Store struct {
	ipv4EntityList []*IPV4Entity
	ipv6EntityList []*IPV6Entity
	metaList       []*Meta
	tmpMetaList    []*Meta // Meta list of intermediate states
	opt            Option  // configuration parameter
}

// WithMetaList Set tuple list
func (s *Store) WithMetaList(ml []*Meta) {
	s.metaList = ml
}

// NewStore returns a new store.
func NewStore() *Store {
	return &Store{}
}

// WithDataFiles Set data file
func (s *Store) WithDataFiles(fs []string) {
	if s != nil && len(fs) > 0 {
		s.opt.Files = fs
	}
}

// IPV4EntityCount Number of IPv4 instances
func (s *Store) IPV4EntityCount() int {
	if s != nil {
		return len(s.ipv4EntityList)
	}
	return 0
}

// IPV6EntityCount Number of IPv6 instances
func (s *Store) IPV6EntityCount() int {
	if s != nil {
		return len(s.ipv6EntityList)
	}
	return 0
}

// IPV4Entity Return the index IPV4Entity pointing to the entity list.
func (s *Store) IPV4Entity(i int) *IPV4Entity {
	if s != nil {
		if i < len(s.ipv4EntityList) && i >= 0 {
			return s.ipv4EntityList[i]
		}
	}
	return nil
}

// IPV6Entity Return the index IPV6Entity pointing to the entity list.
func (s *Store) IPV6Entity(i int) *IPV6Entity {
	if s != nil && i < s.IPV6EntityCount() && i >= 0 {
		return s.ipv6EntityList[i]
	}
	return nil
}

// Search search for IP addresses
func (s *Store) Search(addr net.IP) *Meta {
	if addr == nil {
		return nil
	}
	if utils.IsIPv4(addr.String()) {
		// IPv4
		ipIndex := binary.BigEndian.Uint32(addr.To4())
		if index := sort.Search(s.IPV4EntityCount(), func(i int) bool {
			return s.ipv4EntityList[i].IPIndex() >= ipIndex
		}); s.IPV4Entity(index) != nil {
			if s.IPV4Entity(index).ipIndex != ipIndex {
				index -= 1
			}
			mi := s.ipv4EntityList[index].metaIndex
			return s.metaList[mi]
		}
	} else if utils.IsIPv6(addr.String()) {
		// IPV6
		ipIndex := binary.BigEndian.Uint64(addr.To16())
		if index := sort.Search(s.IPV6EntityCount(), func(i int) bool {
			return s.ipv6EntityList[i].ipIndex >= ipIndex
		}); s.IPV6Entity(index) != nil {
			if s.IPV6Entity(index).IPIndex() != ipIndex {
				index -= 1
			}
			mi := s.ipv6EntityList[index].metaIndex
			return s.metaList[mi]
		}
	}
	return nil
}

// UnmarshalFrom Decompose and store from raeder.
func (s *Store) UnmarshalFrom(reader io.Reader) error {
	var err error
	ipv4List := make([]*IPV4Entity, 0)
	ipv6List := make([]*IPV6Entity, 0)
	metaTable := make(map[string]uint32)
	iReader := bufio.NewReader(reader)
	for {
		var line []byte
		if line, err = iReader.ReadBytes('\n'); err != nil {
			break
		}
		rowMeta := &RowMeta{}
		if err = rowMeta.Unmarshal(line); err != nil {
			_, _ = fmt.Fprint(os.Stderr, "meta unmarshal error, ", err.Error(), string(line))
			continue
		}
		fp := rowMeta.Hash()
		if fp == "" {
			// Fingerprint calculation error
			_, _ = fmt.Fprint(os.Stderr, "Fingerprint calculation error")
			continue
		}
		var index uint32
		var ok bool
		if index, ok = metaTable[fp]; !ok {
			meta := &Meta{
				Country:     rowMeta.Country,
				Province:    rowMeta.Province,
				City:        rowMeta.City,
				Region:      rowMeta.Region,
				FrontISP:    rowMeta.FrontISP,
				BackboneISP: rowMeta.BackboneISP,
				AsID:        rowMeta.AsID,
				Comment:     rowMeta.Comment,
				Type:        rowMeta.Type,
			}
			if s.opt.CB != nil {
				meta.Extends = s.opt.CB(meta)
			}
			index = uint32(len(s.tmpMetaList))
			s.tmpMetaList = append(s.tmpMetaList, meta)
			metaTable[fp] = index
		}

		ipObj := rowMeta.StartIPObj()
		if m := rowMeta.Mode(); m == IPV4 {
			entity := &IPV4Entity{
				ipIndex:   binary.BigEndian.Uint32(ipObj.To4()),
				metaIndex: index,
			}
			ipv4List = append(ipv4List, entity)
		} else if m == IPV6 {
			entity := &IPV6Entity{
				ipIndex:   binary.BigEndian.Uint64(ipObj.To16()),
				metaIndex: index,
			}
			ipv6List = append(ipv6List, entity)
		} else {
			// Bad IP metadata
			_, _ = fmt.Fprint(os.Stderr, "Bad IP metadata, start_ip=", rowMeta.StartIP)
			continue
		}
	}
	if err != io.EOF {
		return fmt.Errorf("unmarshal entity list error, %s", err)
	}

	if len(ipv4List) > 0 {
		s.ipv4EntityList = ipv4List
	}
	if len(ipv6List) > 0 {
		s.ipv6EntityList = ipv6List
	}
	return nil
}

// LoadData load data
func (s *Store) LoadData(opt Option) error {
	if len(opt.Files) <= 0 {
		return errors.New("no incoming data file")
	}
	s.opt = opt
	return s.update()
}

// Update update data
func (s *Store) Update() error {
	return s.update()
}

func (s *Store) update() error {
	s.tmpMetaList = []*Meta{}
	var err error

	for _, fn := range s.opt.Files {
		// open file by fiilename
		var fReader *os.File
		if fReader, err = os.Open(fn); err != nil {
			return err
		}
		err = s.UnmarshalFrom(fReader)
		if err != nil {
			return err
		}

		_ = fReader.Close()
	}

	// Update Meta List
	s.WithMetaList(s.tmpMetaList)
	s.tmpMetaList = []*Meta{}
	return nil
}
