package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 定义一个全局pool
var pool *redis.Pool

// 当启动程序时，就初始化连接池
func init() {
	pool = &redis.Pool{
		MaxIdle:   8, // 最大空闲链接数
		MaxActive: 0, // 表示数据库的最大链接数，0表示没有限制
		Dial: func() (redis.Conn, error) { // 初始化代码，链接那个ip的redis
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func main() {
	// 先从pool取出一个链接
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("Set", "name", "tom")
	if err != nil {
		fmt.Println("conn.Do err=", err)
		return
	}

	// 取出
	r, err := redis.String(conn.Do("Get", "name"))
	if err != nil {
		fmt.Println("conn.Do err=", err)
		return
	}
	fmt.Println("r=", r)
	/*
		// 通过go向redis写入和读取数据
		// 链接到redis
		conn, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("redis.Dial err=", err)
			return
		}
		// 2.通过go向redis写入数据 string [key-val]
		_, err = conn.Do("set", "name", "tomjerry")
		if err != nil {
			fmt.Println("redis.Dial err=", err)
			return
		}
		defer conn.Close() // 关闭
		fmt.Println("操作成功")

		// 通过go向redis读取数据string [key-val]
		//r, err := redis.String(conn.Do("Get", "name"))
		r, err := redis.String(conn.Do("HMGet", "name")) // 读取多个
		if err != nil {
			fmt.Println("set err=", err)
			return
		}
		// 因为返回r是interface{}
		// 因为name对应的值是string,因此我们需要转换
		//nameString := r.(string)
		fmt.Println("操作ok", r)

		// 批量操作

	*/

}
