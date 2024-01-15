package extset

import "unsafe"

//go:noescape
//go:linkname strhash runtime.strhash
func strhash(a unsafe.Pointer, h uintptr) uintptr

// StrHash returns the hash of the given string.
func StrHash(s string) uintptr {
	return strhash(unsafe.Pointer(&s), 0)
}
