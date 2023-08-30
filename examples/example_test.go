/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/23
 * Time: 11:18
 */

package examples

import (
	"fmt"
	"github.com/valiner/vaedb"
)

func Example() {
	//打开对应目录，会在这个目录下生成数据文件 .vdb结尾
	db, err := vaedb.NewVaeDB("./")
	if err != nil {
		panic(err)
	}
	_ = db.Set("test", []byte("vdb"))
	val := db.Get("test")
	db.Del("test")
	val1 := db.Get("test")
	fmt.Println(string(val), string(val1))
	// Output: value
}
