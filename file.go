/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/24
 * Time: 11:38
 */

package vaedb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	vdbSplit = 1000000000
)

type vdbFile struct {
	fileFullName string
	fileIndex    int
	file         *os.File
	offset       int64
	maxSize      int64
}

func newVdbFile(fileName string, maxSize int64) (*vdbFile, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	offset, _ := file.Seek(0, io.SeekEnd)
	fid, _ := getFileId(file.Name())
	return &vdbFile{file: file, maxSize: maxSize, offset: offset, fileFullName: fileName, fileIndex: fid}, nil
}

func (v *vdbFile) NextFile() (*vdbFile, error) {
	fid, _ := getFileId(v.file.Name())
	file, err := newVdbFile(filepath.Join(filepath.Dir(v.fileFullName), getFileStr(fid+1)), v.maxSize)
	if err != nil {
		return nil, err
	}
	file.fileIndex = fid + 1
	return file, nil
}

// 写进磁盘
func (v *vdbFile) WriteEntry(entry []byte) (int64, error) {
	if int64(len(entry))+v.offset > v.maxSize {
		nx, err := v.NextFile()
		if err != nil {
			return 0, err
		}
		*v = *nx
	}
	n, err := v.file.Write(entry)
	if err != nil {
		return 0, err
	}
	v.offset += int64(n)
	return v.offset, err
}

func (v *vdbFile) GetOffset() int64 {
	return v.offset
}

func (v *vdbFile) Close() {
	v.file.Close()
}

type vdbFileNames []string

func (v vdbFileNames) getNeedMergeFiles() vdbFileNames {
	return v[:len(v)-1]
}

func (v vdbFileNames) getNextMergeFile() string {
	for i := len(v) - 1; i >= 0; i-- {
		fid, err := getFileId(v[i])
		if err != nil {
			continue
		}
		if fid < vdbSplit {
			return getFileStr(fid + 1)
		}
	}
	return getFileStr(1)
}

// 不包括活跃file 和 merge file
func (v vdbFileNames) GetOldFiles() vdbFileNames {
	for i := 0; i < len(v); i++ {
		fid, err := getFileId(v[i])
		if err != nil {
			continue
		}
		if fid > vdbSplit {
			return v[i : len(v)-1]
		}
	}
	return v[:0]
}

func (v vdbFileNames) getActiveFiles() vdbFileNames {
	return v[len(v)-2:]
}

type fileIndexWarp struct {
	fileIndex *fileIndex
	key       []byte
}

type fileIndex struct {
	fileId    string
	valueSize int
	valuePos  int64
	timeStamp int64
}

func NewFileIndex(fileId string, valueSize int, valuePos int64, timeStamp int64) *fileIndex {
	return &fileIndex{
		fileId:    fileId,
		valueSize: valueSize,
		valuePos:  valuePos,
		timeStamp: timeStamp,
	}
}

func getFileStr(id int) string {
	return fmt.Sprintf("%010d.%s", id, DataFileSuffix)
}

func getFileId(fileName string) (int, error) {
	fileName = filepath.Base(fileName)
	trimmed := strings.TrimSuffix(fileName, "."+DataFileSuffix)
	return strconv.Atoi(trimmed)
}
