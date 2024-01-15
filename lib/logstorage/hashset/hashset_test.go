package hashset

import (
	"github.com/cespare/xxhash/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/xxh3"
	"hash/maphash"
	"strconv"
	"testing"
)

func TestSet_Insert(t *testing.T) {
	set := GetHashSet()
	size := 30
	for i := 0; i < size; i++ {
		set.Add(strconv.Itoa(i) + "-test")
	}
	t.Log("Data len:", len(set.data))
	t.Log("size:", set.size)
	t.Log("threshold:", set.threshold)
	t.Log("conflict:", set.conflict)

	s := make([]string, 0, size)
	iter := set.Iterator()
	for iter.Next() {
		s = append(s, iter.At())
	}
	assert.Equal(t, size, len(s))
}

func BenchmarkTest(b *testing.B) {
	set := GetHashSet()
	for i := 0; i < b.N; i++ {
		set.Add(strconv.Itoa(i) + "-test")
	}

	iter := set.Iterator()
	for iter.Next() {
		iter.At()
	}
	b.Log("Data len:", len(set.data))
	b.Log("size:", set.size)
	b.Log("threshold:", set.threshold)
	b.Log("conflict:", set.conflict)
}

func TestHash(t *testing.T) {
	size := 10
	for i := 0; i < size; i++ {
		t.Log(StrHash(strconv.Itoa(i) + "-test"))
		//set.Add(strconv.Itoa(i) + "-test")
	}
}

func BenchmarkHashTest(b *testing.B) {
	var h maphash.Hash
	// Reset discards all data previously added to the Hash, without
	// changing its seed.

	str := "10000000-test"
	for i := 0; i < b.N; i++ {
		h.WriteString(str)
		h.Sum64()
		h.Reset()
	}
}

func Benchmark_xxhash2(b *testing.B) {
	// Reset discards all data previously added to the Hash, without
	// changing its seed.
	str := "10000000-test"
	for i := 0; i < b.N; i++ {
		xxhash.Sum64String(str)
	}
}

func Benchmark_xxhash3(b *testing.B) {
	// Reset discards all data previously added to the Hash, without
	// changing its seed.
	str := "10000000-test"
	for i := 0; i < b.N; i++ {
		_ = xxh3.HashString(str)
	}
}
