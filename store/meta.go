// Package store provides data store definitions
package store

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/universal-fraternity/geoip/utils"
)

const (
	Unknown = iota
	IPV4
	IPV6
)

// RowMeta define row metadata
type RowMeta struct {
	StartIP     string      // The starting point of the IP segment is divided into IP addresses
	EndIP       string      // The termination point of the IP segment is divided into IP addresses
	Netmask     string      // Subnet Mask
	Country     string      // Location Information - Country
	Province    string      // Location Information - Province
	City        string      // Location Information - City
	Region      string      // Location Information - Region
	FrontISP    string      // Export ISP
	BackboneISP string      // Backbone ISP
	AsID        int         // AS ID number
	Comment     *string     // Other remarks information
	Type        *string     // Network type
	Extends     interface{} // Extended Information
	startIPObj  net.IP
}

// Meta Definition of metadata results (location information)
type Meta struct {
	Country     string      // Location Information - Country
	Province    string      // Location Information - Province
	City        string      // Location Information - City
	Region      string      // Location Information - Region
	FrontISP    string      // Export ISP
	BackboneISP string      // Backbone ISP
	AsID        int         // AS ID number
	Comment     *string     // Other remarks information
	Type        *string     // Network type
	Extends     interface{} // Extended Information
}

// NewMeta Return a new meta
func NewMeta() *Meta {
	return &Meta{}
}

// WithExtends Set extension information
func (m *Meta) WithExtends(en interface{}) {
	m.Extends = en
}

// IsEmpty If all fields are empty, return true.
func (m Meta) IsEmpty() bool {
	return m.Country == "" &&
		m.Province == "" &&
		m.City == "" &&
		m.Region == "" &&
		m.FrontISP == "" &&
		m.BackboneISP == "" &&
		m.AsID == 0
}

// IsEmpty If all fields are empty, return true.
func (m *RowMeta) IsEmpty() bool {
	return m.StartIP == "" &&
		m.EndIP == "" &&
		m.Netmask == "" &&
		m.Country == "" &&
		m.Province == "" &&
		m.City == "" &&
		m.Region == "" &&
		m.FrontISP == "" &&
		m.BackboneISP == "" &&
		m.AsID == 0
}

// UnmarshalString Parse detailed information into meta format.
func (m *RowMeta) UnmarshalString(row string) error {
	var e error
	if m != nil {
		for i, item := range strings.Split(strings.TrimSuffix(row, "\n"), "\t") {
			switch i {
			case 0:
				m.StartIP = item
				m.startIPObj = net.ParseIP(item)
			case 1:
				m.EndIP = item
			case 2:
				m.Netmask = item
			case 3:
				m.Country = item
			case 4:
				m.Province = item
			case 5:
				m.City = item
			case 6:
				m.Region = item
			case 7:
				m.FrontISP = item
			case 8:
				m.BackboneISP = item
			case 9:
				if m.AsID, e = utils.String2Int(item); e != nil {
					return e
				}
			case 10:
				if item != "NULL" {
					m.Comment = &item
				}
			case 11:
				if item != "NULL" {
					m.Type = &item
				}
			}
		}
	} else {
		return errors.New("meta is null")
	}
	return e
}

// Unmarshal Parse detailed information into meta format.
func (m *RowMeta) Unmarshal(buffer []byte) error {
	return m.UnmarshalString(string(buffer))
}

// Hash hash calculation
func (m *RowMeta) Hash() string {
	if m.IsEmpty() {
		return ""
	}
	h1 := sha1.New()
	source := m.Country + m.Province + m.City + m.Region + m.FrontISP + m.BackboneISP + strconv.Itoa(m.AsID)
	_, _ = io.WriteString(h1, source)
	return string(h1.Sum(nil))
}

// Mode Type judgment
func (m *RowMeta) Mode() int {
	if m.startIPObj.To4() != nil {
		return IPV4
	} else if m.startIPObj.To16() != nil {
		return IPV6
	}
	return Unknown
}

// StartIPObj Starting IP
func (m *RowMeta) StartIPObj() net.IP {
	return m.startIPObj
}

// String  Format output
func (m *Meta) String() string {
	if m != nil {
		return fmt.Sprintf("country:%q province:%q city:%q region:%q "+
			"ISP:%q backboneISP:%q asid:%d",
			m.Country, m.Province, m.City, m.Region,
			m.FrontISP, m.BackboneISP, m.AsID)
	}
	return ""
}

// MarshalString Serialize data entities into strings.
func (m *Meta) MarshalString() (string, error) {
	var line string
	if m != nil {
		var comment, mType string
		if m.Comment == nil {
			comment = "NULL"
		} else {
			comment = *m.Comment
		}
		if m.Type == nil {
			mType = "NULL"
		} else {
			mType = *m.Type
		}
		line = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s",
			m.Country, m.Province, m.City, m.Region,
			m.FrontISP, m.BackboneISP, m.AsID, comment, mType)
	}
	return line, nil
}

// Marshal Serialize data entities into strings.
func (m *Meta) Marshal() ([]byte, error) {
	line, err := m.MarshalString()
	return []byte(line), err
}
