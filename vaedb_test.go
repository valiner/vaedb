/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/24
 * Time: 12:00
 */

package vaedb

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var cnt = 100000

const keyPrefix = "test"

func TestBigSet(t *testing.T) {
	db, err := NewVaeDB("./")
	if err != nil {
		panic(err)
	}
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		fmt.Println(db.Set(k, []byte(k+"tv")))
	}
}

func TestBigGet(t *testing.T) {
	db, err := NewVaeDB("./")
	if err != nil {
		panic(err)
	}
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		v := string(db.Get(k))
		fmt.Println(v)
	}
	time.Sleep(time.Second * 3)
}

func TestSetAndGet(t *testing.T) {
	db, err := NewVaeDB("./")
	if err != nil {
		panic(err)
	}
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		fmt.Println(db.Set(k, []byte(k+"tv")))
	}
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		v := string(db.Get(k))
		fmt.Println(v)
	}
	time.Sleep(20 * time.Second)
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		v := string(db.Get(k))
		fmt.Println(v)
	}
}
