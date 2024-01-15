package logstorage

import (
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logstorage/hashset"
	"sync"
	"unicode"
)

var charTokens [128]bool

func init() {
	splits := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	for _, s := range splits {
		charTokens[s] = true
	}
}

// tokenizeStrings extracts word tokens from a, appends them to dst and returns the result.
func tokenizeStrings(dst, a []string) []string {
	set := hashset.GetHashSet()
	for i, s := range a {
		if i > 0 && s == a[i-1] {
			// This string has been already tokenized
			continue
		}
		tokenizeString(set, s)
	}
	//dstLen := len(dst)
	iter := set.Iterator()
	for iter.Next() {
		dst = append(dst, iter.At())
	}
	hashset.PutHashSet(set)
	//putTokenizer(t)

	// Sort tokens with zero memory allocations
	//ss := getStringsSorter(dst[dstLen:])
	//sort.Sort(ss)
	//putStringsSorter(ss)

	return dst
}

type tokenizer struct {
	m map[string]struct{}
}

func (t *tokenizer) reset() {
	m := t.m
	for k := range m {
		delete(m, k)
	}
}

func tokenizeString(dst *hashset.HashSet, s string) {
	start := -1 // valid span start if >= 0
	for end, c := range s {
		if c < 128 && charTokens[c] || unicode.IsLetter(c) {
			if start < 0 {
				start = end
			}
		} else {
			if start >= 0 {
				dst.Add(s[start:end])
				start = ^start
			}
		}
	}
	if start >= 0 {
		dst.Add(s[start:])
	}
}

func isTokenRune(c rune) bool {
	if c > 128 {
		return unicode.IsLetter(c)
	}
	return charTokens[c]
}

func getTokenizer() *tokenizer {
	v := tokenizerPool.Get()
	if v == nil {
		return &tokenizer{
			m: make(map[string]struct{}),
		}
	}
	return v.(*tokenizer)
}

func putTokenizer(t *tokenizer) {
	t.reset()
	tokenizerPool.Put(t)
}

var tokenizerPool sync.Pool

type stringsSorter struct {
	a []string
}

func (ss *stringsSorter) Len() int {
	return len(ss.a)
}
func (ss *stringsSorter) Swap(i, j int) {
	a := ss.a
	a[i], a[j] = a[j], a[i]
}
func (ss *stringsSorter) Less(i, j int) bool {
	a := ss.a
	return a[i] < a[j]
}

func getStringsSorter(a []string) *stringsSorter {
	v := stringsSorterPool.Get()
	if v == nil {
		return &stringsSorter{
			a: a,
		}
	}
	ss := v.(*stringsSorter)
	ss.a = a
	return ss
}

func putStringsSorter(ss *stringsSorter) {
	ss.a = nil
	stringsSorterPool.Put(ss)
}

var stringsSorterPool sync.Pool

type tokensBuf struct {
	A []string
}

func (tb *tokensBuf) reset() {
	a := tb.A
	for i := range a {
		a[i] = ""
	}
	tb.A = a[:0]
}

func getTokensBuf() *tokensBuf {
	v := tokensBufPool.Get()
	if v == nil {
		return &tokensBuf{}
	}
	return v.(*tokensBuf)
}

func putTokensBuf(tb *tokensBuf) {
	tb.reset()
	tokensBufPool.Put(tb)
}

var tokensBufPool sync.Pool
