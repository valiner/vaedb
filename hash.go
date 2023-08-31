/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/23
 * Time: 16:11
 */

package vaedb

import (
	"crypto/md5"
)

const (
	Prime   = 16777619
	HashVal = 2166136261
)

type Hasher interface {
	Hash([]byte) []byte
}

func NewMd5Hash() *Md5Hash {
	return &Md5Hash{}
}

type Md5Hash struct {
}

// Hash returns the first 4 bytes of the md5 checksum of the byte slice.
func (m *Md5Hash) Hash(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:4]
}

type IHash interface {
	Hash(string) uint32
}

type Fnv32Hash struct{}

func NewFnv32Hash() *Fnv32Hash {
	return &Fnv32Hash{}
}

// fnv32 algorithm
func (f *Fnv32Hash) Hash(key string) uint32 {
	hashVal := uint32(HashVal)
	prime := uint32(Prime)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hashVal *= prime
		hashVal ^= uint32(key[i])
	}
	return hashVal
}
