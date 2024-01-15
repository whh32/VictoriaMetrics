package extset

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSet_Insert(t *testing.T) {
	set := GetSet()
	size := 500
	for i := 0; i < size; i++ {
		set.Add(strconv.Itoa(i) + "-test")
	}
	t.Log("Data len:", len(set.Data))
	t.Log("dataSize:", set.dataSize)
	t.Log("nextResize:", set.nextResize)
	t.Log("conflict:", set.conflict)

	s := make([]string, 0, size)
	iter := set.Iterator()
	for iter.Next() {
		s = append(s, iter.At())
	}
	assert.Equal(t, size, len(s))
}

func BenchmarkTest(b *testing.B) {
	set := GetSet()
	for i := 0; i < b.N; i++ {
		set.Add(strconv.Itoa(i) + "-test")
	}

	iter := set.Iterator()
	for iter.Next() {
		iter.At()
	}
	b.Log("Data len:", len(set.Data))
	b.Log("dataSize:", set.dataSize)
	b.Log("nextResize:", set.nextResize)
	b.Log("conflict:", set.conflict)
}

func TestHash(t *testing.T) {
	size := 10
	for i := 0; i < size; i++ {
		t.Log(StrHash(strconv.Itoa(i) + "-test"))
		//set.Add(strconv.Itoa(i) + "-test")
	}
}
