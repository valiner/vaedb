/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/23
 * Time: 16:11
 */

package vaedb

import (
	"crypto/md5"
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
