package example

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// TestRedisPool 测试redis连接池
// maxidle 最大空闲数 一旦连接会一直保持tcp连接 处于ESTABLISHED状态,被连接池复用
// maxactive 最大活跃数 使用完会被回收，tcp连接被关闭，处于TIME_WAIT状态，等待tcp真正关闭
// 所以maxactive>=maxidle，maxactive决定了应用操作redis的最大并发能力，maxidle决定了连接池复用的能力
// windows下查看连接数 netstat -ano | findstr "127.0.0.1:6379"
// redis客户端通讯协议 https://redis.io/topics/protocol
// redis下的编码问题
// window下使用cmd进行telnet，写入中文，[set hi1 hello中国],value占9个字节，因为cmd是gbk编码，中文占两个字节
// 使用redigo执行[set hi2 hello中国]，value占11个字节，因为golang中使用utf8编码，中文占3个字节
// 此时在telnet中执行get hi1，返回正常，[$9\r\nhello中国\r\n]
// 在telnet中执行get hi2，返回乱码，[$11\r\nhello涓浗\r\n]
// 在redigo中执行get hi，返回乱码，[$9\r\nhello�й�\r\n]
// .net redis驱动ServiceStack.Redis中，也可以看到字符串以utf8方式进行编码
//	public void SetEntry(string key, string value)
//	{
//	    byte[] buffer = (value != null) ? value.ToUtf8Bytes() : null;
//	    base.Set(key, buffer);
//	}
func TestRedisPool(t *testing.T) {
	fmt.Println(len("hello中国"))
	//callTime(redisTestReadRepeat)
	//callTime(redisTestReadWithConnRepeat)
	//callTime(redisTestReadWithConnRepeatConcurrency)
	//callTime(redisTestReadRepeatConcurrency)
	//callTime(redisTestChinese)

	for i := 0; i < 1500; i++ {
		redisPool()
	}
	time.Sleep(time.Minute * 5)
}

func redisPool() {
	p := newPool("127.0.0.1:6379")
	for i := 0; i < 10; i++ {
		go func() {
			c := p.Get()
			str, _ := redis.String(c.Do("GET", "hello_chinese"))
			fmt.Println(str)
			c.Close()
		}()
	}
	time.Sleep(11 * time.Second)
	c := p.Get()
	str, _ := redis.String(c.Do("GET", "hello_chinese"))
	fmt.Println(str)
	c.Close()
	//p.Close()
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     2,
		MaxActive:   4,
		IdleTimeout: 10 * time.Second, // 触发式超时释放，在从连接池中拿取连接的时候才会进行idle释放
		//Wait:true,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func callTime(f func()) {
	now := time.Now().UnixNano()
	f()
	end := time.Now().UnixNano()
	fmt.Printf("call [%s] finish span [%d]\r\n", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), (end-now)/1000000)
}

func redisTestChinese() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialReadTimeout(time.Second*1000), redis.DialWriteTimeout(time.Second*1000), redis.DialDatabase(0))
	if err != nil {
		fmt.Printf("dial err %s \r\n", err)
		return
	}
	conn.Do("set", "hello_chinese", "hello中国")
	str, err := redis.String(conn.Do("get", "hi"))
	if err != nil {
		fmt.Printf("do  err %s \r\n", err)
		return
	}
	_ = str
	fmt.Println("do finish " + str)
	fmt.Println("redis done !")
}

func redisTestReadRepeat() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialReadTimeout(time.Second), redis.DialWriteTimeout(time.Second), redis.DialDatabase(12))
	if err != nil {
		fmt.Printf("dial err %s \r\n", err)
		return
	}
	conn.Do("set", "hello_redis", "Redis（Remote Dictionary Server )")
	for i := 0; i < 3500; i++ {
		str, err := redis.String(conn.Do("get", "hello_redis"))
		if err != nil {
			fmt.Printf("do [%d] err %s \r\n", i, err)
			return
		}
		_ = str
		fmt.Println("do[" + strconv.Itoa(i) + "] " + str)
	}
	fmt.Println("redis done !")
}

func redisTestReadWithConnRepeat() {
	for i := 0; i < 5000; i++ {
		//conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialReadTimeout(time.Second), redis.DialWriteTimeout(time.Second), redis.DialDatabase(12))
		conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(12))
		if err != nil {
			fmt.Printf("dial [%d] err %s \r\n", i, err)
			return
		}
		str, err := redis.String(conn.Do("get", "hello_redis"))
		if err != nil {
			fmt.Printf("do [%d] err %s \r\n", i, err)
			return
		}
		_ = str
		fmt.Println("do[" + strconv.Itoa(i) + "] " + str)
	}
	fmt.Println("redis done !")
}

func redisTestReadWithConnRepeatConcurrency() {
	for i := 0; i < 100; i++ {
		go func(num int) {
			// read timeout 高并发下延迟较大
			// write timeout 高并发下延迟较小
			conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialReadTimeout(time.Second*5), redis.DialWriteTimeout(time.Second*5), redis.DialDatabase(12))
			//conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(12))

			if err != nil {
				fmt.Printf("dial [%d] err %s \r\n", num, err)
				return
			}
			str, err := redis.String(conn.Do("get", "hello_redis"))
			conn.Close()
			if err != nil {
				fmt.Printf("do [%d] err %s \r\n", num, err)
				return
			}
			_ = str
			fmt.Println("do[" + strconv.Itoa(num) + "] " + str)
		}(i)
	}
	fmt.Println("redis done !")
}

func redisTestReadRepeatConcurrency() {
	// conn.Do 同一个conn不能并发操作Do
	conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialReadTimeout(time.Second*5), redis.DialWriteTimeout(time.Second*5), redis.DialDatabase(12))
	if err != nil {
		fmt.Printf("dial err %s \r\n", err)
		return
	}
	//conn.Do("set", "hello_redis", "Redis（Remote Dictionary Server )")
	for i := 0; i < 5; i++ {
		go func(num int) {
			str, err := redis.String(conn.Do("get", "hello_redis"))
			if err != nil {
				fmt.Printf("do [%d] err %s \r\n", num, err)
				return
			}
			_ = str
			fmt.Println("do[" + strconv.Itoa(num) + "] " + str)
		}(i)
	}
	time.Sleep(time.Second * 10)
	fmt.Println("redis done !")
}
