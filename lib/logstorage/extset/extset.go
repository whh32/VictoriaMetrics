package extset

import (
	"hash/crc32"
	"sync"
	"unsafe"
)

// 帮我做下整体的实现. 用slice实现一个set,并且实现set的insert,delete,遍历等方法，要求底层实现为slice，且类型为string,hash冲突时使用开放寻址法解决
const (
	DefaultSize = 128
	DefaultLoad = 0.5
	DefaultGrow = 4
)

type Set struct {
	Data       []string
	dataSize   int     // 当前集合数据大小
	loadFactor float32 // 最大负载因子
	growFactor int     // 增长因子：2的指数倍
	nextResize int     // 下一次调整大小的阈值
	conflict   int     //冲突次数,后续可进行删除
}

// NewSet 初始化集合
func NewSet(size int, loadFactor float32, growth int) *Set {
	s := &Set{
		Data:       getSlice(size),
		loadFactor: loadFactor,
		growFactor: growth,
		dataSize:   0,
	}
	s.setNextResize()
	return s
}

func (s *Set) setNextResize() {
	s.nextResize = int(float32(len(s.Data)) * s.loadFactor)
}

// Add 插入一个元素到集合中
func (s *Set) Add(value string) {
	if s.add(value) {
		s.dataSize += 1
		s.checkResize() // 检查是否需要调整大小
	}
}

//func YoloBytes(s string) []byte {
//	return
//}

func (s *Set) add(value string) bool {
	//hash:=crc32.ChecksumIEEE(*(*[]byte)(unsafe.Pointer(&value)))

	//index := int(xxhash.Sum64String(value)) & (len(s.Data) - 1)
	index := int(crc32.ChecksumIEEE(*(*[]byte)(unsafe.Pointer(&value)))) & (len(s.Data) - 1)
	for s.Data[index] != "" {
		if s.Data[index] == value {
			return false
		}
		s.conflict += 1
		index = (index + 1) & (len(s.Data) - 1) // 使用开放寻址法解决冲突
	}

	s.Data[index] = value
	return true
}

func (s *Set) Iterator() Iterator {
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
	set *Set
	idx int
}

func (i *setIter) Next() bool {
	i.idx += 1
	for i.idx < len(i.set.Data) {
		if i.set.Data[i.idx] == "" {
			i.idx += 1
		} else {
			return true
		}
	}
	return false
}

func (i *setIter) At() string {
	return i.set.Data[i.idx]
}

// 检查是否需要调整集合的大小，根据最大负载因子和增长因子进行计算和调整
func (s *Set) checkResize() {
	if s.dataSize > s.nextResize {
		src := s.Data
		s.Data = getSlice(len(s.Data) * s.growFactor) // 创建新的底层数据结构切片
		s.setNextResize()

		for i := range src {
			if src[i] != "" {
				s.add(src[i])
			}
		}
		putSlice(src)
	}
}

func (s *Set) reset() {
	s.dataSize = 0
	s.growFactor = DefaultGrow
	s.loadFactor = DefaultLoad
	s.conflict = 0
	s.Data = s.Data[:DefaultSize]
	clear(s.Data)
	s.setNextResize()
}

func GetSet() *Set {
	v := setPool.Get()
	if v == nil {
		return NewSet(DefaultSize, DefaultLoad, DefaultGrow)
	}
	s := v.(*Set)
	return s
}

func PutSet(set *Set) {
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
