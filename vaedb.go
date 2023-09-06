/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/23
 * Time: 11:16
 */

package vaedb

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DataFileSuffix = "vdb"
	FileMaxSize    = 1024 * 1024 * 1024 * 1 // 1MB
	FirstVdbName   = "1000000001.vdb"
)

type entryHandle func(*entry)

type VaeDB struct {
	//todo bigMap Optimize https://github.com/golang/go/issues/9477
	keys           *ShardMap
	dataDir        *dataDir
	path           string
	activeFile     *vdbFile
	mux            sync.RWMutex
	hash           Hasher
	entryBuffer    []byte
	logger         Logger
	compacter      compacter
	msgCh          chan *fileIndexWarp
	lruCache       *LruCache
	getInterceptor []Interceptor
	setInterceptor []Interceptor
}

type Option func(db *VaeDB)

func NewVaeDB(path string, opts ...Option) (v *VaeDB, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	dir, err := newMyDir(path)
	if err != nil {
		return
	}
	msgCh := make(chan *fileIndexWarp, 10)
	v = &VaeDB{
		keys:        DefaultShardMap(),
		dataDir:     dir,
		hash:        NewMd5Hash(),
		entryBuffer: make([]byte, 1024),
		path:        path,
		logger:      DefaultLogger(),
		msgCh:       msgCh,
		compacter:   defaultCompactness(path, msgCh),
		lruCache:    DefaultLruCache(),
	}

	for _, opt := range opts {
		opt(v)
	}

	if err = v.loadData(); err != nil {
		return
	}
	go v.compacter.run()
	go v.mergeKeys()
	return
}

func SetGetInterceptor(interceptor ...Interceptor) Option {
	return func(v *VaeDB) {
		v.getInterceptor = interceptor
	}
}

func SetSetInterceptor(interceptor ...Interceptor) Option {
	return func(v *VaeDB) {
		v.setInterceptor = interceptor
	}
}

// loadData 从磁盘中恢复数据
func (v *VaeDB) loadData() error {
	vdbs := v.dataDir.getVdbFileNames()
	if len(vdbs) == 0 {
		file, err := newVdbFile(filepath.Join(v.path, FirstVdbName), FileMaxSize)
		if err != nil {
			return err
		}
		v.activeFile = file
		return nil
	}
	for _, fileName := range vdbs {
		var offset int64
		v.dataDir.readFile(fileName, func(e *entry) {
			k := string(e.key)
			v.keys.set(k, NewFileIndex(fileName, int(e.length), offset, int64(e.timeStamp)))
			offset += e.length
		})
	}
	activeFileName := vdbs[len(vdbs)-1]
	file, err := newVdbFile(filepath.Join(v.path, activeFileName), FileMaxSize)
	if err != nil {
		return err
	}
	v.activeFile = file
	return nil
}

// 	合并数据，打包旧数据
func (v *VaeDB) mergeKeys() {
	for newFileIndexWarp := range v.msgCh {
		fIndex := v.keys.get(string(newFileIndexWarp.key))
		if fIndex == nil {
			v.logger.Println("not found Key", string(newFileIndexWarp.key))
			continue
		}
		fileIndex := fIndex.(*fileIndex)
		//活跃文件的key不操作
		if fileIndex.fileId != getFileStr(v.activeFile.fileIndex) {
			if newFileIndexWarp.fileIndex.timeStamp == 0 {
				v.keys.del(string(newFileIndexWarp.key))
				continue
			}
			v.keys.set(string(newFileIndexWarp.key), newFileIndexWarp.fileIndex)
		}
	}
}

//empty or delete return nil
func (v *VaeDB) Get(key string) (val []byte) {
	chain := NewChain(key, []byte(""))
	chain.AddInterceptor(v.getInterceptor...)
	chain.AddInterceptor(v.get)
	chain.Exec()
	return chain.Val
}

func (v *VaeDB) get(chain *Chain) {
	key := chain.Key
	cacheItem := v.lruCache.Get(key)
	if cacheItem != nil {
		chain.Val = cacheItem.val
	}
	fIndex := v.keys.get(key)
	if fIndex == nil {
		return
	}
	fileIndex := fIndex.(*fileIndex)
	file, err := os.Open(filepath.Join(v.path, fileIndex.fileId))
	if err != nil {
		return
	}
	defer file.Close()
	e, err := readEntryByPos(file, fileIndex.valueSize, fileIndex.valuePos)
	if err != nil {
		return
	}
	//deleted
	if e.timeStamp == 0 {
		return
	}
	v.lruCache.Set(key, e.value)
	chain.Val = e.value
}

func (v *VaeDB) Set(key string, val []byte) (err error) {
	chain := NewChain(key, val)
	chain.AddInterceptor(v.setInterceptor...)
	chain.AddInterceptor(v.setLru)
	chain.Exec()
	return chain.Err
}

func (v *VaeDB) setLru(chain *Chain) {
	val := chain.Val
	key := chain.Key
	v.lruCache.Set(key, val)
	chain.Err = v.set(key, val, time.Now().UnixNano())
}

func (v *VaeDB) set(key string, val []byte, ts int64) error {
	v.mux.Lock()
	defer v.mux.Unlock()
	entry := wrapEntry(ts, key, val, &v.entryBuffer, v.hash)
	offset := v.activeFile.GetOffset()
	_, err := v.activeFile.WriteEntry(entry)
	if err != nil {
		return err
	}
	v.keys.set(key, NewFileIndex(getFileStr(v.activeFile.fileIndex), len(entry), offset, ts))
	return nil
}

//del 追加一条数据为空的纪录,且ts==0
func (v *VaeDB) Del(key string) {
	v.set(key, []byte{}, 0)
}
