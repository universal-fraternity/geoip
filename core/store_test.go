package core

import (
	"net"
	"testing"
)

func TestLoad(t *testing.T) {
	st := NewStore()
	t.Log(st.LoadData(Option{
		Files: []string{"testdata/ipv6.txt"},
	}))
}

func TestSearch(t *testing.T) {
	st := NewStore()
	t.Log(st.LoadData(Option{
		Files: []string{"testdata/ipv6.txt", "testdata/ipv4.txt"},
	}))
	t.Log(st.Search(net.ParseIP("1.30.13.90")))
}

func BenchmarkStore_Search(b *testing.B) {
	st := NewStore()
	if err := st.LoadData(Option{
		Files: []string{"testdata/ipv6.txt",
			"testdata/ipv4.txt"},
	}); err != nil {
		b.Error(err)
	}
	addr := net.ParseIP("1.30.13.51")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.Search(addr)
	}
}
