/**
 * User: coder.sdp@gmail.com
 * Date: 2023/9/5
 * Time: 15:47
 */

package vaedb

//责任链模式，为get / set 设置拦截器
type IChain interface {
	AddInterceptor(interceptor ...Interceptor)
	Next()
}

type Interceptor func(c *Chain)

type Chain struct {
	pos          int
	interceptors []Interceptor
	Key          string
	Val          []byte
	Err          error
}

func NewChain(k string, v []byte) *Chain {
	return &Chain{Key: k, Val: v}
}

func (c *Chain) AddInterceptor(interceptor ...Interceptor) {
	c.interceptors = append(c.interceptors, interceptor...)
}

func (c *Chain) Next() {
	if c.pos < len(c.interceptors) {
		index := c.pos
		c.pos += 1
		c.interceptors[index](c)
	}
	return
}

func (c *Chain) Exec() {
	c.Next()
}
