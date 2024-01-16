package hashset

import (
	"github.com/zeebo/xxh3"
	"sync"
)

const (
	DefaultSize       = 1 << 7
	DefaultLoadFactor = 0.5
	DefaultGrowFactor = 1 << 2
)

type HashSet struct {
	data       []string
	size       int     // 当前集合数据大小
	loadFactor float32 // 最大负载因子
	growFactor int     // 增长因子：2的指数倍
	threshold  int     // 下一次调整大小的阈值
	conflict   int     //冲突次数,后续可进行删除
}

// NewHashSet 初始化集合
func NewHashSet(size int) *HashSet {
	if size < 1 {
		panic("invalid init size")
	}
	for i := 3; i < 32; i++ {
		if size <= (1 << i) {
			size = 1 << i
			break
		}
	}

	s := &HashSet{
		data:       getSlice(size),
		loadFactor: DefaultLoadFactor,
		growFactor: DefaultGrowFactor,
		size:       0,
	}
	s.setThreshold()
	return s
}

func (s *HashSet) setThreshold() {
	s.threshold = int(float32(len(s.data)) * s.loadFactor)
}

func (s *HashSet) Size() int {
	return s.size
}

func (s *HashSet) Add(value string) {
	if s.add(value) {
		s.size += 1
		if s.size > s.threshold {
			s.resize(len(s.data) * s.growFactor)
		}
	}
}

func (s *HashSet) add(value string) bool {
	//index := int(xxhash.Sum64String(value)) & (len(s.data) - 1)
	index := int(xxh3.HashString(value)) & (len(s.data) - 1)
	for s.data[index] != "" {
		if s.data[index] == value {
			return false
		}
		s.conflict += 1
		index = (index + 1) & (len(s.data) - 1) // 使用开放寻址法解决冲突
	}

	s.data[index] = value
	return true
}

// 检查是否需要调整集合的大小，根据最大负载因子和增长因子进行计算和调整
func (s *HashSet) resize(newCapacity int) {
	src := s.data
	s.data = getSlice(newCapacity)
	s.setThreshold()

	for i := range src {
		if src[i] != "" {
			s.add(src[i])
		}
	}
	putSlice(src)
}

func (s *HashSet) reset() {
	s.size = 0
	s.growFactor = DefaultGrowFactor
	s.loadFactor = DefaultLoadFactor
	s.conflict = 0
	s.data = s.data[:DefaultSize]
	clear(s.data)
	s.setThreshold()
}

func (s *HashSet) Iterator() Iterator {
	return &setIter{
		set: s,
		idx: -1,
	}
}

type Iterator interface {
	// Next advances the iterator and returns true if another value was found.
	Next() bool

	// At returns the value at the current iterator position.
	At() string
}

type setIter struct {
	set *HashSet
	idx int
}

func (i *setIter) Next() bool {
	i.idx += 1
	for i.idx < len(i.set.data) {
		if i.set.data[i.idx] == "" {
			i.idx += 1
		} else {
			return true
		}
	}
	return false
}

func (i *setIter) At() string {
	return i.set.data[i.idx]
}

func GetHashSet() *HashSet {
	v := setPool.Get()
	if v == nil {
		return NewHashSet(DefaultSize)
	}
	s := v.(*HashSet)
	return s
}

func PutHashSet(set *HashSet) {
	set.reset()
	setPool.Put(set)
}

var setPool sync.Pool

func getSlice(sz int) []string {
	o := slicePool.Get()
	if o == nil {
		return make([]string, sz)
	}
	b := o.([]string)
	if cap(b) < sz {
		return make([]string, sz)
	}
	b = b[:sz]
	clear(b)
	return b
}

// putSlice adds a slice to the right Bucket in the slicePool.
func putSlice(s []string) {
	slicePool.Put(s)
}

var slicePool sync.Pool
