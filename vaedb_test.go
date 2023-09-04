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
)

var cnt = 100000

const keyPrefix = "test"

var db, _ = NewVaeDB(".")

func TestBigSet(t *testing.T) {
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		fmt.Println(db.Set(k, []byte(k+"tv")))
	}
}

func TestBigGet(t *testing.T) {
	for i := 0; i < cnt; i++ {
		k := keyPrefix + strconv.Itoa(i)
		v := string(db.Get(k))
		if v != k+"tv" {
			t.Errorf("err k:%s v%s", k, v)
		}
	}
}

func TestSetAndGet(t *testing.T) {
	t.Run("c_set", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < cnt; i++ {
			k := keyPrefix + strconv.Itoa(i)
			db.Set(k, []byte(k+"tv"))
		}
	})

	t.Run("c_get", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < cnt; i++ {
			k := keyPrefix + strconv.Itoa(i)
			v := string(db.Get(k))
			if v != "" && v != k+"tv" {
				t.Errorf("err k:%s v%s", k, v)
			}
		}
	})
}
