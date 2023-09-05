/**
 * User: coder.sdp@gmail.com
 * Date: 2023/9/5
 * Time: 12:13
 */

package vaedb

import (
	"container/list"
	"sync"
)

const DefaultSize = 1 << 16

type LruCache struct {
	cache map[string]*list.Element
	size  int
	list  *list.List
	mux   sync.Mutex
}

type CacheItem struct {
	key string
	val []byte
}

func DefaultLruCache() *LruCache {
	return NewLruCache(DefaultSize)
}

// size = 0 表示无限制
func NewLruCache(size int) *LruCache {
	return &LruCache{
		cache: make(map[string]*list.Element),
		size:  size,
		list:  list.New(),
	}
}

func (l *LruCache) isFull() bool {
	if l.size == 0 {
		return false
	}
	return l.list.Len() >= l.size
}

func (l *LruCache) Set(key string, val []byte) {
	l.mux.Lock()
	defer l.mux.Unlock()
	e, ok := l.cache[key]
	if !ok {
		if l.isFull() {
			l.list.Remove(l.list.Back())
		}
	} else {
		l.list.Remove(e)
	}
	newEle := l.list.PushFront(&CacheItem{
		key: key,
		val: val,
	})
	l.cache[key] = newEle
}

func (l *LruCache) Get(key string) *CacheItem {
	l.mux.Lock()
	defer l.mux.Unlock()
	e, ok := l.cache[key]
	if ok {
		l.list.MoveToFront(e)
		return e.Value.(*CacheItem)
	}
	return nil
}
