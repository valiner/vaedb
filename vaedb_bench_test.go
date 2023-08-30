/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/24
 * Time: 15:20
 */

package vaedb

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var value = bytes.Repeat([]byte("a"), 1024)

//BenchmarkVaeDB_Set-8   	  150385	      7754 ns/op	     148 B/op	       6 allocs/op
//BenchmarkVaeDB_Set-8   	  163815	      7918 ns/op	     135 B/op	       6 allocs/op
//BenchmarkVaeDB_Set-8   	  162955	      8022 ns/op	     135 B/op	       6 allocs/op
func BenchmarkVaeDB_Set(b *testing.B) {
	//3.5G 耗时1s多 带bufio速度提升3倍
	db, _ := NewVaeDB("./")
	b.ResetTimer()
	rand.Seed(time.Now().UnixNano())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := fmt.Sprintf("test:%d", rand.Int())
			_ = db.Set(key, value)
		}
	})
}

// BenchmarkVaeDB_Get-8   	  532131	      3004 ns/op	     404 B/op	       7 allocs/op
// BenchmarkVaeDB_Get-8   	  501195	      2763 ns/op	     408 B/op	       7 allocs/op
func BenchmarkVaeDB_Get(b *testing.B) {
	db, _ := NewVaeDB("./")
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("test:%d", i)
		if err := db.Set(key, value); err != nil {
			panic(err)
		}
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		for pb.Next() {
			key := fmt.Sprintf("test:%d", rand.Intn(b.N))
			_ = db.Get(key)
		}
	})
}
