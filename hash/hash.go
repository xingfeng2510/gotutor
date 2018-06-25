package hash

import "hash"

type (
	sum32 uint32
)

const (
	offset32 = 2166136261
	prime32  = 16777619
)

func NewFNV32() hash.Hash {
	var s sum32 = offset32
	return &s
}

func (s *sum32) Reset() { *s = offset32 }

func (s *sum32) Sum32() uint32 { return uint32(*s) }

func (s *sum32) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash *= prime32
		hash ^= sum32(c)
	}
	*s = hash
	return len(data), nil
}

func (s *sum32) Size() int { return 4 }

func (s *sum32) BlockSize() int { return 1 }

func (s *sum32) Sum(b []byte) []byte {
	v := uint32(*s)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
