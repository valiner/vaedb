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

type myDir struct {
	dir  *os.File
	path string
}

func newMyDir(path string) (*myDir, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &myDir{dir: dir, path: path}, nil
}

func (d *myDir) Close() error {
	return d.dir.Close()
}

func (d *myDir) getVdbs() vdbFiles {
	dirEntry, _ := d.dir.ReadDir(0)
	vdbs := make([]string, 0)
	for _, item := range dirEntry {
		if strings.HasSuffix(item.Name(), DataFileSuffix) {
			vdbs = append(vdbs, item.Name())
		}
	}
	sort.Strings(vdbs)
	return vdbs
}

func (d *myDir) readFile(fileName string, f func(*entry)) error {
	fd, err := os.Open(filepath.Join(d.path, fileName))
	defer fd.Close()
	if err != nil {
		return err
	}
	readEntryFormFile(fd, f)
	return nil
}
