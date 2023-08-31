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
```

## todo
- [x] 更高性能的hash存储结构 => 分片锁
- [ ] 增加hit文件，提高恢复速度
- [ ] 提高文件读取速度



