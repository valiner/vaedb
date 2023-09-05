# vaedb

> 一个简单，高效的 key-value数据库，基于bitcask

## bitcask
[bitcask](https://blog.csdn.net/Z_Stand/article/details/115606758?ydreferer=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8%3D?ydreferer=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8%3D)

## examples
```go
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


//拦截器
//get拦截器测试
func TestSetGetInterceptor(t *testing.T) {
    getInter := func (chain *Chain) {
		//get执行前的操作 改变chain的key,val会影响之后操作
        key := chain.Key
        chain.Key = "test2"
        fmt.Printf("获取到原key %s,并更改Key %s", key, chain.Key)
        chain.Next()
		//get获取到值后的操作
        fmt.Printf(("获取到key %s,val %s", chain.Key, chain.Val)
    }
    db, _ = NewVaeDB(".", SetGetInterceptor(getInter))
    db.Get("test1")
}

//set拦截器测试
func TestSetSetInterceptor(t *testing.T) {
    setInter := func(chain *Chain) {
        //set执行前的操作  改变chain的key,val会影响之后操作
        key := chain.Key
        chain.Key = "test2"
        t.Logf("获取到原key %s,并更改Key %s", key, chain.Key)
        chain.Next()
        //set执行后的操作
    }
    db, _ = NewVaeDB(".", SetSetInterceptor(setInter))
    db.Set("test1", []byte("test1"))
    fmt.Printf("获取 test1 %s", db.Get("test1"))
	fmt.Printf("获取 test2 %s", db.Get("test2"))
}

```

## todo
- [x] 更高性能的hash存储结构 => 分片锁
- [ ] 增加hit文件，提高恢复速度
- [x] get方法 增加缓存 
- [x] get / set 拦截器
