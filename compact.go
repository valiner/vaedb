/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/25
 * Time: 09:40
 */

package vaedb

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	CompactFileMaxSize = FileMaxSize * 2
	MinCompactNum      = 10
	DefaultInterval    = time.Second * 10000
)

type compacter interface {
	run()
}

type compactness struct {
	interval            time.Duration
	path                string
	maxSize             int64 //合并文件的最大大小
	minMinCompactNumNum int   //老文件大于这个数量，才开始进行合并操作
	msgCh               chan *fileIndexWarp
	entryBuffer         []byte
	hash                Hasher
}

func defaultCompactness(path string, msgCh chan *fileIndexWarp) *compactness {
	return newCompactness(DefaultInterval, path, CompactFileMaxSize, MinCompactNum, msgCh)
}

func newCompactness(interval time.Duration, path string, maxSize int64, minMinCompactNumNum int, msgCh chan *fileIndexWarp) *compactness {
	return &compactness{interval: interval, path: path, maxSize: maxSize, minMinCompactNumNum: minMinCompactNumNum, msgCh: msgCh, hash: NewMd5Hash(), entryBuffer: make([]byte, 1024)}
}

func (c *compactness) run() {
	ticker := time.NewTicker(c.interval)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	for {
		select {
		case <-ticker.C:
			c.compact()
		}
	}
}

func (c *compactness) compact() {
	dir, err := newMyDir(c.path)
	if err != nil {
		fmt.Println("newMyDir error:", err)
		return
	}
	defer dir.Close()
	vdbs := dir.getVdbs()
	if len(vdbs.GetOldFiles()) < c.minMinCompactNumNum {
		return
	}
	nmf := vdbs.getNeedMergeFiles()
	filter := make(map[string]entry)
	for _, file := range nmf {
		dir.readFile(file, func(e *entry) {
			//Process more data？
			filter[string(e.key)] = e.Clone()
			//e中含有引用类型，下列写法错误！
			//filter[string(e.key)] = *e
		})
	}
	newFileName := vdbs.getNextMergeFile()
	mergeFile, err := newVdbFile(newFileName, c.maxSize)
	if err != nil {
		fmt.Println("newVdbFile error:", err)
		return
	}
	for _, e := range filter {
		fmt.Println(string(e.key))
		offset := mergeFile.GetOffset()
		entry := wrapEntry(int64(e.timeStamp), string(e.key), e.value, &c.entryBuffer, c.hash)
		_, err := mergeFile.WriteEntry(entry)
		if err != nil {
			fmt.Println("mergeFile.WriteEntry error:", err)
			return
		}
		m := &fileIndexWarp{
			fileIndex: NewFileIndex(getFileStr(mergeFile.fileIndex), int(e.length), offset, int64(e.timeStamp)),
			key:       e.key,
		}
		c.msgCh <- m
	}

	for i := 0; i < len(nmf); i++ {
		os.Remove(filepath.Join(c.path, nmf[i]))
	}
}
