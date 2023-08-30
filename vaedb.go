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
	keys          map[string]*fileIndex
	dir           *myDir
	path          string
	activeFile    *vdbFile
	mux           sync.RWMutex
	activeFileId  int
	curFileOffset int64
	hash          Hasher
	entryBuffer   []byte
	logger        Logger
	compacter     compacter
	msgCh         chan *fileIndexWarp
}

func NewVaeDB(path string) (v *VaeDB, err error) {
	dir, err := newMyDir(path)
	defer dir.dir.Close()
	if err != nil {
		return
	}
	msgCh := make(chan *fileIndexWarp, 10)
	v = &VaeDB{
		keys:        make(map[string]*fileIndex),
		dir:         dir,
		hash:        NewMd5Hash(),
		entryBuffer: make([]byte, 1024),
		path:        path,
		logger:      DefaultLogger(),
		msgCh:       msgCh,
		compacter:   defaultCompactness(path, msgCh),
	}
	err = v.loadData()
	go v.compacter.run()
	go v.mergeKeys()
	return
}

// loadData 从磁盘中恢复数据
func (v *VaeDB) loadData() error {
	vdbs := v.dir.getVdbs()
	if len(vdbs) == 0 {
		file, err := newVdbFile(filepath.Join(v.path, FirstVdbName), FileMaxSize)
		if err != nil {
			return err
		}
		v.activeFile = file
	} else {
		for _, fileName := range vdbs {
			var offset int64
			v.dir.readFile(fileName, func(e *entry) {
				k := string(e.key)
				v.keys[k] = NewFileIndex(fileName, int(e.length), offset, int64(e.timeStamp))
				offset += e.length
			})
		}
		activeFileName := vdbs[len(vdbs)-1]
		file, err := newVdbFile(filepath.Join(v.path, activeFileName), FileMaxSize)
		if err != nil {
			return err
		}
		v.activeFile = file
	}
	return nil
}

// 	合并数据，打包旧数据
func (v *VaeDB) mergeKeys() {
	for newFileIndexWarp := range v.msgCh {
		v.mux.Lock()
		fileIndex, ok := v.keys[string(newFileIndexWarp.key)]
		if !ok {
			v.logger.Println("not found key", string(newFileIndexWarp.key))
			continue
		}
		//活跃文件的key不操作
		if fileIndex.fileId != getFileStr(v.activeFile.fileIndex) {
			if newFileIndexWarp.fileIndex.timeStamp == 0 {
				delete(v.keys, string(newFileIndexWarp.key))
				continue
			}
			v.keys[string(newFileIndexWarp.key)] = newFileIndexWarp.fileIndex
		}
		v.mux.Unlock()
	}
}

func (v *VaeDB) Get(key string) []byte {
	v.mux.Lock()
	defer v.mux.Unlock()
	val := make([]byte, 0)
	fileIndex, ok := v.keys[key]
	if !ok {
		return val
	}
	file, err := os.Open(filepath.Join(v.path, fileIndex.fileId))
	if err != nil {
		return val
	}
	defer file.Close()
	e, err := readEntryByPos(file, fileIndex.valueSize, fileIndex.valuePos)
	if err != nil {
		return val
	}
	//deleted
	if e.timeStamp == 0 {
		//fmt.Println("删除了")
	}
	return e.value
}

func (v *VaeDB) Set(key string, val []byte) (err error) {
	return v.set(key, val, time.Now().UnixNano())
}

func (v *VaeDB) set(key string, val []byte, ts int64) (err error) {
	v.mux.Lock()
	defer v.mux.Unlock()
	entry := wrapEntry(ts, key, val, &v.entryBuffer, v.hash)
	offset := v.activeFile.GetOffset()
	_, err = v.activeFile.WriteEntry(entry)
	if err != nil {
		return
	}
	v.keys[key] = NewFileIndex(getFileStr(v.activeFile.fileIndex), len(entry), offset, ts)
	return
}

//Del 追加一条数据为空的纪录,且ts==0
func (v *VaeDB) Del(key string) {
	v.set(key, []byte{}, 0)
}
