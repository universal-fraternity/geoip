package geoip

import "testing"

func TestSearch(t *testing.T) {
	if err := Load(Option{
		Files: []string{
			"store/testdata/ipv6.txt",
			"store/testdata/ipv4.txt",
		},
		CB: nil,
	}); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(Search("1.30.13.242"))
}
