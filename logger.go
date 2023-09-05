/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/24
 * Time: 10:25
 */

package vaedb

import (
	"log"
)

type Logger interface {
	Println(v ...interface{})
}

var _ Logger = &log.Logger{}

func DefaultLogger() *log.Logger {
	return log.Default()
}

func newLogger(custom Logger) Logger {
	if custom != nil {
		return custom
	}
	return DefaultLogger()
}
