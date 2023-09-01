/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/30
 * Time: 10:03
 */

package vaedb

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type dataDir struct {
	path string
}

func newMyDir(path string) (*dataDir, error) {
	return &dataDir{path: path}, nil
}

// 获取当前的数据文件名
func (d *dataDir) getVdbFileNames() vdbFileNames {
	dirEntry, _ := os.ReadDir(d.path)
	vdbs := make([]string, 0)
	for _, item := range dirEntry {
		if strings.HasSuffix(item.Name(), DataFileSuffix) {
			vdbs = append(vdbs, item.Name())
		}
	}
	sort.Strings(vdbs)
	return vdbs
}

func (d *dataDir) readFile(fileName string, f func(*entry)) error {
	fd, err := os.Open(filepath.Join(d.path, fileName))
	defer fd.Close()
	if err != nil {
		return err
	}
	readEntryFormFile(fd, f)
	return nil
}
