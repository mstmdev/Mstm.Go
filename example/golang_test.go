package example

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println(os.Args[0])

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

func TestExit(t *testing.T) {
	defer fmt.Println("test") //因为程序已退出  所以永远不会执行
	os.Exit(3)
}

func TestOtherProcess(t *testing.T) {
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"ls", "-a", "-l", "-h"}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
func TestSpawningProcess(t *testing.T) {
	//启动计算器
	exec.Command("calc").Start()

	//date
	dateCmd := exec.Command("date")
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))

	//grep
	grepCmd := exec.Command("grep", "hello")
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep \ngoodbye grep"))
	grepIn.Close()
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()
	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	//bash
	//用bash调用ls遍历当前目录下的文件信息
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))

	//用bash打开记事本
	exec.Command("bash", "-c", "notepad").Output()

}

// TestEnvironmentVariable 环境变量操作
func TestEnvironmentVariable(t *testing.T) {
	//获取指定的环境变量
	fmt.Println(os.Getenv("OS"))

	//设置环境变量  windows7下测试不起作用
	// os.Setenv("testpan","valuetest2")

	fmt.Println(os.Getenv("testpan"))

	// fmt.Println(os.Environ())

	//遍历所有环境变量
	for _, item := range os.Environ() {
		// fmt.Println(item)
		key := strings.Split(item, "=")[0]
		value := strings.Split(item, "=")[1]
		fmt.Println(key)
		fmt.Println(value)
	}
}

// TestCmdLineArgFlag
//go run test.go -strK=testString -intK=99 -boolK=false -svar=flag a b c
//前面的参数将会进行标记的匹配  最后的参数会在Args中输出
//-h 可以显示需要输入的参数说明
func TestCmdLineArgFlag(t *testing.T) {
	//键  值  说明
	wordPtr := flag.String("strK", "texValue", "this is String")

	numPtr := flag.Int("intK", 23, "this is Int")

	boolPtr := flag.Bool("boolK", true, "this bool")

	var svar string
	flag.StringVar(&svar, "svar", "varValue", "a string var")

	flag.Parse()

	fmt.Println(*wordPtr)
	fmt.Println(*numPtr)
	fmt.Println(*boolPtr)
	fmt.Println(svar)

	fmt.Println(flag.Args())

}

func TestCommandLineArguement(t *testing.T) {
	args := os.Args
	fmt.Println(args)
	//参数第一项为当前应用的路径  这里主输出命令行参数
	argsWithoutPath := os.Args[1:]
	fmt.Println(argsWithoutPath)

	//输出第二个命令行参数
	fmt.Println(os.Args[2])
}

// TestConsoleScan 控制台输入
func TestConsoleScan(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		ucl := strings.ToUpper(scanner.Text())
		fmt.Println(ucl)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func TestWriteFile(t *testing.T) {
	//将字符串写入文件
	file := "write.txt"
	data := []byte("write a string to file")
	err := ioutil.WriteFile(file, data, 0644)
	check(err)

	f, err := os.Create("new.txt")
	check(err)
	defer f.Close()

	//将字节数组写入文件
	data2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(data2)
	check(err)
	fmt.Printf("write %d bytes \n", n2)

	//直接将字符串写入文件
	n3, err := f.WriteString("通过字符串写入文件\n")
	fmt.Printf("write %d bytes \n", n3)
	f.Sync()

	//带缓存的文件写入
	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffer\n")
	fmt.Printf("Write bytes: %d", n4)
	w.Flush()
}

func TestReadFile(t *testing.T) {
	file := "file.txt" //UTF-8 无BOM编码文件
	//读取文件
	data, err := ioutil.ReadFile(file)
	check(err)
	fmt.Println(string(data))

	f, err := os.Open(file)
	//将会延迟到程序的末尾执行  关闭文件流的读取
	defer f.Close()
	check(err)

	//读取指定的前5个字节
	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1))

	//从指定位置读取字节
	o2, err := f.Seek(6, 0)
	check(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	check(err)
	fmt.Printf("%d bytes:@%d %s\n", n2, o2, string(b2))

	//从指定位置读取字节
	o3, err := f.Seek(6, 0)
	check(err)
	b3 := make([]byte, 2)
	//至少读取2个字节  否则会报错
	n3, err := io.ReadAtLeast(f, b3, 2)
	check(err)
	fmt.Printf("%d bytes @ %d : %s\n", n3, o3, string(b3))

	//将指针移动到起始位置
	_, err = f.Seek(0, 0)
	check(err)

	//读取前5个字节  使用缓存
	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(5)
	check(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestBase64(t *testing.T) {
	data := "abc123!?$*&()'-=@~"
	//Base64加密
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc) //YWJjMTIzIT8kKiYoKSctPUB+

	//Base64解密
	sDec, _ := base64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))

	//Url Base64加密
	uEnc := base64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc) //YWJjMTIzIT8kKiYoKSctPUB-

	//URL Base64解密
	uDec, _ := base64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}

// TestMD5 32位MD5值加密
func TestMD5(t *testing.T) {
	str := "md5test string"
	h := md5.New()
	h.Write([]byte(str))
	res := h.Sum(nil)
	fmt.Printf("%x", res)
}

// TestSHA1Hashes 哈希值计算 40位
func TestSHA1Hashes(t *testing.T) {
	str := "sha1 hashes test string"
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	fmt.Println(bs)
	fmt.Printf("%x\n", bs)
}

func TestUrlParsing(t *testing.T) {
	s := "postgres://user:pass@host.com:5432/path?key=v&test=testV#ff"

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	//解析协议
	fmt.Println(u.Scheme)

	//解析用户信息
	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	fmt.Println(u.User.Password())

	//解析主机信息
	fmt.Println(u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	fmt.Println(host)
	fmt.Println(port)

	fmt.Println(u.Path)
	fmt.Println(u.Fragment)

	//解析查询参数
	fmt.Println(u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["key"][0])
}

func TestStrConv(t *testing.T) {
	fmt.Println(strconv.ParseFloat("1.2356", 64))
	fmt.Println(strconv.ParseInt("1234", 0, 64))
	fmt.Println(strconv.ParseInt("0x212", 0, 64))
	fmt.Println(strconv.ParseUint("1223", 0, 64))
	fmt.Println(strconv.Atoi("-123"))
	fmt.Println(strconv.Atoi("123"))
	//fmt.Println(strconv.Atoi("test"))
}

func TestRandom(t *testing.T) {
	//默认情况下 随机数是确定的
	fmt.Print(rand.Intn(100), ",")
	fmt.Println(rand.Float64())

	//指定一个可变的Seed  随机输出数字
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	fmt.Println(r1.Intn(100))

}

func TestTimeFormatting(t *testing.T) {
	fmt.Println(time.Now())
	fmt.Println(time.Now().Format(time.RFC3339))
	fmt.Println(time.Parse(time.RFC3339, "2012-08-08T15:04:05.999999-07:00"))

	p := fmt.Println
	now := time.Now()
	p(".........")
	p(now.Format("3:04PM"))
	p(now.Format("Mon Jan _2 15:04:05 2006"))
	p(now.Format("2006-01-02T15:04:05.999999-07:00"))
	form := "3 04 PM"
	t2, _ := time.Parse(form, "8 41 PM")
	p(t2)
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	ansic := "Mon Jan _2 15:04:05 2006"
	_, e := time.Parse(ansic, "8:41PM")
	p(e)
}

func TestUnixEpoch(t *testing.T) {
	fmt.Println(time.Now())
	//Unix时间戳
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano())
	millis := time.Now().UnixNano() / 1000000
	fmt.Println(millis)
	fmt.Println(time.Unix(time.Now().Unix(), 0))
	fmt.Println(time.Unix(0, time.Now().UnixNano()))
}

func TestTimeNow(t *testing.T) {
	p := fmt.Println
	fmt.Println(time.Now())

	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	p(then)

	p("then.Year()", then.Year())
	p("then.Month()", then.Month())
	p("then.Day()", then.Day())
	p("then.Hour()", then.Hour())
	p("then.Minute()", then.Minute())
	p("then.Second()", then.Second())
	p("then.Nanosecond()", then.Nanosecond())
	p("then.Location()", then.Location())
	p("then.Weekday()", then.Weekday())

	p(then.Before(time.Now()))
	p(then.After(time.Now()))
	p(then.Equal(time.Now()))

	p(time.Now().Sub(then))
	p(time.Now().Sub(then).Hours())
	p(time.Now().Sub(then).Minutes())
	p(time.Now().Sub(then).Nanoseconds())

	p(then.Add(time.Now().Sub(then)))
	p(then.Add(-time.Now().Sub(then)))
}

func TestJson(t *testing.T) {
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(100)
	fmt.Println(string(intB))

	perJson, _ := json.Marshal(PersonalJson{name: "zhangsan", age: 16, address: "beijing"})
	fmt.Println(string(perJson))

	floatJson, _ := json.Marshal(1.34)
	fmt.Println(string(floatJson))

	strJson, _ := json.Marshal("json test hahahahaha")
	fmt.Println(string(strJson))

	slcD := []string{"apple", "test"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	res1D := &Response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}

	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	res2D := &Response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	num := dat["num"].(float64)
	fmt.Println(num)

	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := &Response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)

}

type Response1 struct {
	Page   int
	Fruits []string
}

type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

type PersonalJson struct {
	name    string `json:"name"`
	age     int    `json:"age"`
	address string `json:"address"`
}

func TestRegex(t *testing.T) {

	//判断字符串是否符合正则
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println(match)

	//判断字符串是否符合正则
	r, _ := regexp.Compile("p([a-z]+)ch")
	fmt.Println(r.MatchString("peach"))

	//输出符合正则的字符串
	fmt.Println(r.FindString("peach punch"))

	//输出符合正则的字符串索引集合
	fmt.Println(r.FindStringIndex("peach punch"))

	//输出符合正则的 和符合（）里的正则的字符串
	fmt.Println(r.FindStringSubmatch("peach punch"))

	fmt.Println(r.FindStringSubmatchIndex("peach punch"))

	fmt.Println(r.FindAllString("peach punch pinch", -1))

	fmt.Println(r.FindAllStringSubmatchIndex("peach punch pinch", -1))

	fmt.Println(r.FindAllString("peach punch pinch", 2))

	fmt.Println(r.Match([]byte("peach")))

	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println(r)

	fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)
	fmt.Println(string(out))
}

func TestStringFormatting(t *testing.T) {

	p := Point{5, 7}
	fmt.Println(p.x)
	fmt.Println(p.y)
	fmt.Println(p)

	fmt.Printf("%v\n", p)
	fmt.Printf("%+v\n", p)
	fmt.Printf("%#v\n", p)
	//输出p的类型
	fmt.Printf("%T\n", p)

	fmt.Printf("%t\n", true)
	fmt.Printf("%d\n", 123)
	fmt.Printf("%b\n", 14)
	//将ASCII码转换为字符
	fmt.Printf("%c\n", 33)
	//输出十六进制
	fmt.Printf("%x\n", 456)
	//decimal
	fmt.Printf("%f\n", 78.9)

	//科学计算
	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)

	//字符转义后输出
	fmt.Printf("%s\n", "\"string\"")
	//原样输出字符
	fmt.Printf("%q\n", "\"string\"")

	fmt.Printf("%p\n", &p)
	//空格
	fmt.Printf("|%6d|%6d|\n", 12, 345)

	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)
	//在字符前添加指定的空格
	fmt.Printf("|%6s|%6s|\n", "foo", "b")
	//在字符后添加指定的空格
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)
	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}

type Point struct {
	x, y int
}

func TestStrings(t *testing.T) {
	p := fmt.Println
	//查找字符串中是否包含某个字符串
	p("Contains:", strings.Contains("test", "te"))
	//计算字符串中某个字符串出现的次数
	p("Count:", strings.Count("testcard", "t"))
	//字符串是否以指定的前缀开始
	p("HasPrefix:", strings.HasPrefix("testcard", "tes"))
	//字符串是否以指定的后缀开始
	p("HasSuffix:", strings.HasSuffix("testcard", "rd"))
	//指定字符串中某个字符串所在的索引位置
	p("Index:", strings.Index("testcard", "tc"))
	//将字符数组以某个字符串拼接起来
	p("Join:", strings.Join([]string{"test", "card", "new"}, "-->"))
	//将指定的字符串重复指定的次数返回
	p("Repeat:", strings.Repeat("test", 5))
	//替换指定的字符串 指定替换的此处 -1为替换全部
	p("Replace:", strings.Replace("taaestcaresd", "es", "oo", -1))
	//根据制定的字符串分割字符串 返回一个字符数组
	p("Split:", strings.Split("test->card->new", "->"))
	//转换为大写
	p("ToUpper:", strings.ToUpper("teStCard"))
	//转换为小写
	p("ToLower:", strings.ToLower("TestCard"))
	p("len:", len("testcard"))
	p("Index->", "testCard"[1])

}

func TestFileOpera(t *testing.T) {
	f, err := os.Create("defer.txt")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(f, "testdata")
	f.Close()
}

func TestDefer(t *testing.T) {
	fmt.Println("opera1")
	defer deferFunc("def1")
	defer deferFunc("def2")
	fmt.Println("opera2")
	defer deferFunc("def3")
}

func deferFunc(str string) {
	fmt.Println("defer func invoked", str)
}

func TestPanic(t *testing.T) {
	panic("this is a problem")

	_, err := os.Create("/file.test")
	if err != nil {
		panic(err)
		// fmt.Println(err)
	}
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}

func (s ByLength) Swap(left, right int) {
	s[left], s[right] = s[right], s[left]
}

func (s ByLength) Less(left, right int) bool {
	return len(s[left]) < len(s[right])
}

//自定义排序 实现排序接口中的Len Swap Less函数
func TestUserSort(t *testing.T) {
	strs := []string{"c", "aa", "b", "cccc"}
	sort.Sort(ByLength(strs))
	fmt.Println(strs)
}

func TestSort(t *testing.T) {
	strs := []string{"b", "c", "d", "a"}
	fmt.Println("strs is Sorted?-->", sort.StringsAreSorted(strs))
	sort.Strings(strs)
	fmt.Println(strs)
	fmt.Println("strs is Sorted?-->", sort.StringsAreSorted(strs))

	ints := []int{4, 6, 3, 22, 12, 3, 45}
	sort.Ints(ints)
	fmt.Println(ints)

}

func TestStatefulGoroutines2(t *testing.T) {
	var ops int64 = 0
	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := &readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddInt64(&ops, 1)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := &writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddInt64(&ops, 1)
			}
		}()
	}

	time.Sleep(time.Second)
	opsFinal := atomic.LoadInt64(&ops)
	fmt.Println("ops:", opsFinal)

}

type readOp struct {
	key  int
	resp chan int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func TestStatefulGoroutines(t *testing.T) {
	var ops int64 = 0
	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	//定义一个Goroutine进行读写的同步操作
	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
			fmt.Println(atomic.LoadInt64(&ops))
		}
	}()

	//开启100个Goroutine进行读操作
	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := &readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddInt64(&ops, 1)
			}
		}()
	}

	//开启10个Goroutine进行写操作
	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := &writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddInt64(&ops, 1)
			}
		}()
	}

	//主线程休眠1秒 保持所有Goroutine执行1秒钟
	time.Sleep(time.Second)
	//读取最后的读写计数ops
	opsFinal := atomic.LoadInt64(&ops)
	//打印读写总次数
	fmt.Println("ops:", opsFinal)

}

func TestMutexTest2(t *testing.T) {
	//进行读写锁定时读写的内容都得放在
	//mutex.Lock与mutex.Unlock之间

	tag := 9
	mutex := &sync.Mutex{}
	fmt.Println("Init -->", 9)
	go func() {
		//写锁定  其他地方无法读取
		mutex.Lock()
		tag += 1
		time.Sleep(time.Second * 5)
		mutex.Unlock()

	}()

	time.Sleep(time.Second * 2)
	//等写锁定解锁之后才可以读取
	mutex.Lock()
	fmt.Println(tag)
	mutex.Unlock()
}

func TestRandom2(t *testing.T) {
	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 200)
		fmt.Println(rand.Intn(100))
	}
}

func TestMutexes(t *testing.T) {
	var state = make(map[int]int)

	var mutex = &sync.Mutex{}

	var ops int64 = 0

	for r := 0; r < 100; r++ {
		go func() {
			total := 0
			for {
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				// fmt.Println(state[key])
				mutex.Unlock()
				atomic.AddInt64(&ops, 1)
				runtime.Gosched()
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddInt64(&ops, 1)
				runtime.Gosched()
			}
		}()
	}

	time.Sleep(time.Second)

	opsFinal := atomic.LoadInt64(&ops)
	fmt.Println("ops:", opsFinal)

	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}

func TestFor2(t *testing.T) {
	fmt.Println("start")
	for i := 0; i <= 10; i++ {
		go func() {
			for {
				fmt.Println("working......")
			}
		}()
		fmt.Println("执行-->", i)
	}
	time.Sleep(time.Second)
	fmt.Println("end")
}

func TestAtomicCounters(t *testing.T) {
	var ops uint64 = 0

	for i := 0; i < 50; i++ {
		go func() {
			for {
				atomic.AddUint64(&ops, 1)
				runtime.Gosched()
			}
		}()
	}

	time.Sleep(time.Second)

	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops", opsFinal)
}

// TestRateLimiting 控制任务的执行速率
func TestRateLimiting(t *testing.T) {

	//1 相同频率
	//每隔200毫秒执行一次任务
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}

	close(requests)

	limiter := time.Tick(time.Millisecond * 200)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	//2  不同频率
	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Millisecond * 200) {
			burstyLimiter <- t
			fmt.Println("--->>>>>")
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}

	close(burstyRequests)

	for req := range burstyRequests {
		//前三个chan已经接受不会阻塞 之后2个每隔200毫秒才会接受  会产生阻塞
		<-burstyLimiter
		fmt.Println("request[2]", req, time.Now())
	}
}

func TestWorker1(t *testing.T) {
	input := make(chan int, 100)
	output := make(chan int, 100)
	for i := 1; i <= 3; i++ {
		go worker1(i, input, output)
	}

	for i := 1; i <= 16; i++ {
		input <- i
	}

	close(input)

	for i := 1; i <= 16; i++ {
		<-output
	}

	fmt.Println("the end")
}

func worker1(workId int, input <-chan int, output chan<- int) {
	for item := range input {
		fmt.Println("workId-->", workId, "jobs->", item)
		time.Sleep(time.Second * 2)
		output <- item
	}
}

// TestTicket2 使用ticker每个0.5秒执行一次
func TestTicket2(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for item := range ticker.C {
			fmt.Println(item)
		}
	}()
	time.Sleep(time.Second * 3)
	ticker.Stop()
	fmt.Println("ticker Stop")
}

// TestTicket 使用ticker每个0.5秒执行一次 三秒后停止
func TestTicket(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for item := range ticker.C {
		fmt.Println(item)
	}
	time.Sleep(time.Second * 3)
	ticker.Stop()
	fmt.Println("ticker Stop")
}

func TestTimer2(t *testing.T) {
	timer := time.NewTimer(time.Second * 5)
	go func() {
		<-timer.C
		fmt.Println("time received")
	}()

	res := timer.Stop()
	if res {
		fmt.Println(res)
	}
}

// TestTimer1 使用timer在5秒后执行
func TestTimer1(t *testing.T) {
	timer1 := time.NewTimer(time.Second * 5)
	<-timer1.C
	fmt.Println("5S After")
}

func TestRangeChanell2(t *testing.T) {
	chs := make(chan int, 3)
	chs <- 1
	chs <- 2
	chs <- 3
	close(chs)
	for item := range chs {
		fmt.Println(item)
	}
}

func TestCloserChannel(t *testing.T) {
	jobs := make(chan string, 8)
	done := make(chan bool)

	go func() {
		for {
			jobs, hasMore := <-jobs
			if hasMore {
				fmt.Println("isReceived  ", jobs)
			} else {
				fmt.Println("received all")
				done <- true
				return
			}
		}
	}()

	for i := 0; i <= 4; i++ {
		jobs <- strconv.Itoa(i)
		fmt.Println("send ", i)
	}
	close(jobs)
	fmt.Println("send All End")
	<-done
}

// TestNonBlockingChannel 非阻塞Channel
func TestNonBlockingChannel(t *testing.T) {
	msg := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 3)
		msg <- "test"
	}()

	select {
	case receive := <-msg:
		fmt.Println("msg Received", receive)
	default:
		fmt.Println("exec default")
	}
}

func TestChannel01(t *testing.T) {
	chan1 := make(chan int, 8)
	chan1 <- 2
	chan1 <- 4
	chan1 <- 6
	chan1 <- 8
	chan1 <- 10
	fmt.Println(chan1)
	fmt.Println(&chan1)
	fmt.Println("len=", len(chan1))
	length := len(chan1)
	for i := 0; i < length; i++ {
		fmt.Println(<-chan1, "len->", len(chan1))
	}
}

func TestWhatSelect(t *testing.T) {
	a := make(chan string, 1)
	b := make(chan string, 1)
	c := make(chan string, 1)
	for {
		//会随机执行一个case
		select {
		case a <- "a":
			fmt.Println("a")
		case b <- "b":
			fmt.Println("b")
		case c <- "c":
			fmt.Println("c")
			return
		}
	}

}

func TestTimeOut2(t *testing.T) {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 5)
		c1 <- "result"
	}()

	//输出时会等待上述匿名函数执行完毕  完成c1的接收
	// fmt.Println(<-c1)
	fmt.Println("ending")
}

func TestTimeOut(t *testing.T) {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "result 1"
	}()

	//第一个case等待就绪的时间较长 所以会执行第二个case 输出 timeout 1
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "result 2"
	}()

	//第二个case等待就绪的时间较长 所以会执行第一个case 输出 result 2
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(time.Second * 3):
		fmt.Println("timeout 2")
	}
}

func TestSelect(t *testing.T) {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "one"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

func TestChannelDirections(t *testing.T) {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	// chan<-传入     <-chan传出
	//将数据传入pings中
	ping(pings, "接收")
	//将pings中的数据传出 并传入到pongs中
	pong(pings, pongs)
	fmt.Println(<-pongs)
}

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func TestSyncChannel(t *testing.T) {
	done := make(chan bool, 1)
	go syncChannel(done)
	<-done
}

// syncChannel 同步Channel
func syncChannel(done chan bool) {
	fmt.Println("working......")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}

func TestChannelBuffer(t *testing.T) {
	msg := make(chan string, 2)
	msg <- "test1"
	msg <- "test2"
	//有点类似于队列 先进先出
	fmt.Println(<-msg)
	fmt.Println(<-msg)
	msg <- "test3"
	fmt.Println(<-msg)
}

func TestChannel2(t *testing.T) {
	variable1 := make(chan string)
	go func() { variable1 <- "value1" }()
	variable2 := <-variable1
	fmt.Println(variable2)
}

func TestChannel(t *testing.T) {
	message := make(chan string)
	go func() { message <- "test" }()
	msg := <-message
	fmt.Println(msg)
}

func TestScanner(t *testing.T) {
	var inputStr string
	fmt.Scanln(&inputStr)
	fmt.Println("Input End")
	fmt.Println(inputStr)
}

func TestGoroutines(t *testing.T) {
	printNumString("Usual")
	fmt.Println("...........")
	go printNumString("go-->")
	fmt.Println("end....................")
}

func printNumString(mark string) {
	for i := 0; i <= 10; i++ {
		fmt.Println(mark, i)
	}
}

func TestCustomerError(t *testing.T) {
	_, err := CustomerError()
	fmt.Println(err)
	//type assertion
	if ae, ok := err.(*ErrorInfo); ok {
		fmt.Println(ae.message)
	}
}

type ErrorInfo struct {
	result  int
	message string
}

//自定义实现Error接口
func (err *ErrorInfo) Error() string {
	return fmt.Sprintf("%d  -  %s", err.result, err.message)
}

func CustomerError() (int, error) {
	return 98, &ErrorInfo{result: 99, message: "find error"}
}

func ErrorTest() (string, error) {
	return "-1", errors.New("this is error")
}

func TestInterface(t *testing.T) {
	names := []string{"test1", "test2", "test3"}
	inters := make([]interface{}, len(names))
	for index, value := range names {
		inters[index] = value
	}
	printAll(inters)
}

func printAll(inters []interface{}) {
	for _, item := range inters {
		fmt.Println(item)
	}
}

func TestUserInfo(t *testing.T) {
	userinfo := UserInfo{UserName: "ZN", Age: 23, Password: "123456789", Address: "bj"}
	UserInfoOpera(userinfo)
}

//传递的实例必须实现的IUserInfo接口上的所有方法
func UserInfoOpera(usr IUserInfo) {
	fmt.Println(".............UserInfoOpera Start.............")
	fmt.Println(usr.GetUserName())
	fmt.Println(usr.StartLogin())
	fmt.Println(".............UserInfoOpera End.............")
}

func (usr UserInfo) GetUserName() string {
	return usr.UserName
}

func (usr UserInfo) StartLogin() bool {
	return true
}

type IUserInfo interface {
	GetUserName() string
	StartLogin() bool
}

type IHttpClient interface {
	HttpGet() string
	HttpPost() string
	HttpHead() string
	HttpDelete() string
	HttpTrace() string
	HttpPut() string
}

type UserInfo struct {
	UserName, Address, Password string
	Age                         int
}

func (user UserInfo) SayHello() {
	fmt.Println(user.UserName)
}

func TestStruct(t *testing.T) {
	person := Person2{Name: "Zhang", Age: 23, Address: "BJ"}
	fmt.Println(person)
	fmt.Println(person.Name)
	fmt.Println(&person)
	point := &person
	fmt.Println(point.Name)
	point.Name = "Change"
	fmt.Println(person.Name)
}

//结构Struct
type Person2 struct {
	Name    string
	Age     int
	Address string
}

func TestPoint(t *testing.T) {
	num1 := 3
	num2 := 5
	Pointers(&num1)
	Values(num2)
	fmt.Println(num1) //9
	fmt.Println(num2) //5
	fmt.Println("num1的地址为：", &num1)
	fmt.Println("num2的地址为：", &num2)
}

//指针使用
//传递引用
func Pointers(ptr *int) {
	*ptr = 9
	fmt.Println("Pointers中的地址", ptr)
}

//拷贝值
func Values(num int) {
	num = 99
	fmt.Println("Values中的地址", &num)
}

//递归
func Recursion(num int) int {
	if num == 0 {
		return 1
	} else {
		return num * Recursion(num-1)
	}
}

func TestClosures(t *testing.T) {
	//call Closures
	func1 := Closures()
	func1()
	func1()
	func1()
}

//闭包
func Closures() func() int {
	i := 0
	return func() int {
		i = i + 1
		fmt.Println(i)
		return i
	}
}

// TestVariadicParams 可变参数测试
func TestVariadicParams(t *testing.T) {
	pars := []int{1, 4, 6, 8, 89, 5, 2}
	sum := VariadicParm(pars...)
	fmt.Println(sum)
}

func VariadicParm(nums ...int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

func MutilOutput() (string, string, string, string) {
	return "i1", "i2", "i3", "i4"
}

func Calc(left, right int) (int, int) {
	return left + right, left - right
}

func Add2(first, second, third int) int {
	return first + second + third
}

func Add(left int, right int) int {
	return left + right
}

func TestRangeSlice(t *testing.T) {
	slice := make([]string, 10)
	slice[0] = "slice0"
	slice[1] = "slice1"
	slice[2] = "slice2"
	for itemIndex, itemValue := range slice {
		fmt.Println(itemIndex, "---->", itemValue)
	}
}

func PrintfString() {
	fmt.Printf("星期 %s 至星期 %s ", "3", "6")
}

func TestRange(t *testing.T) {
	nums := [5]int{4, 5, 6, 8, 3}
	for _, num := range nums {
		fmt.Println(num)
	}
}

func TestMap2(t *testing.T) {
	maps := make(map[string]string)
	maps["Name"] = "MLL"
	maps["Age"] = "22"
	maps["Address"] = "BEIJING"

	_, hasAddress := maps["Address"]
	fmt.Println(maps)
	fmt.Println(hasAddress)
	delete(maps, "Address")
	_, hasAddress = maps["Address"]
	fmt.Println(hasAddress)
}

func TestQuickMaps(t *testing.T) {
	quick := map[string]int{"map1": 2, "map2": 4}
	fmt.Println(quick)
}

func TestMap(t *testing.T) {
	maps := make(map[string]int)
	maps["first"] = 9
	maps["second"] = 11
	maps["third"] = 13
	fmt.Println(maps)
	fmt.Println("map.Length=", len(maps))
	delete(maps, "second")
	fmt.Println(maps)

	_, hasFirst := maps["first"]
	fmt.Println(hasFirst)

	_, hasSecond := maps["Second"]
	fmt.Println(hasSecond)
}

func TestSlice4(t *testing.T) {
	slice := []string{"a1", "a2", "a3"}
	slice = append(slice, "append1", "append2")
	fmt.Println(slice)
}

func TestSlice3(t *testing.T) {
	slice := make([]string, 14)
	for i := 0; i <= 10; i++ {
		slice[i] = strconv.Itoa(i + 1)
	}

	for i := 0; i < 6; i++ {
		fmt.Println(slice[i])
	}
	fmt.Println("Nice Print Finished !")
}

func TestSlice2(t *testing.T) {
	slice := []string{"s1", "s2", "s3"}
	slice2 := slice[0:2]
	fmt.Println(slice2)

	dSlice := make([][]string, 10)
	for i := 0; i < 10; i++ {
		dSlice[i] = make([]string, 5)
		for j := 0; j < 5; j++ {
			//将int转换成string
			dSlice[i][j] = strconv.Itoa(j + 10)
		}
	}
	fmt.Println(dSlice)
}

func TestSlice(t *testing.T) {
	var slice = make([]string, 3)
	slice[0] = "sliceItem1"
	slice[1] = "sliceItem2"
	slice[2] = "sliceItem3"
	fmt.Println(slice)

	//append必须手动接收返回值
	slice = append(slice, "append2")
	fmt.Println(slice)

	slice2 := make([]string, 9)
	copy(slice2, slice)
	fmt.Println(slice2)

	//将索引2、3、4的子项复制到新的Slice中
	simpleSlice := slice[2:5]
	fmt.Println(simpleSlice)
	fmt.Println("simpleSlice的长度为：", len(simpleSlice))
}

func TestArray(t *testing.T) {
	var array [5]string = [5]string{"哈哈", "test", "other", "test2", "test3"}
	fmt.Println(array)
}

func TestIfElse(t *testing.T) {

	var num = 9
	if num < 10 {
		fmt.Println("您输入的数字过小！")
	}
	if num == 9 {
		fmt.Println("您输入的值为9！")
	}

	if num2 := 8; num2 < 5 {
		fmt.Println("您输入的数字小于5！")
	} else {
		fmt.Println("您输入的数字大于5！")
	}

	//num2是局部变量 无法在此处访问 已超出作用域
	// fmt.Println(num2);
}

func TestFor02(t *testing.T) {
	for i := 10; i > 0; i-- {
		fmt.Println(i)
	}

	fmt.Println("...............")
	var time = 15
	for {
		if time <= 0 {
			break
		}
		time--
		fmt.Println(time)
	}
}

func TestConstants2(t *testing.T) {

	const str string = "这是一个常量"
	// str="尝试修改常量  --> cannot assign to str";

	fmt.Println(str)

	const str2 = "常量"
	fmt.Println(str2)

	fmt.Println(math.Sin(0.5))

	const num = 990
	fmt.Println(int64(num))

}

func Test2Variables(t *testing.T) {
	var str string = "测试通过"
	fmt.Println(str)

	var str2 = "这是一个字符串变量"
	fmt.Println(str2)

	var num1, num2, sum = 2, 3, 0
	sum = num1*num2 + num1
	fmt.Println(sum)

	testStr := "快速定义变量不需要使用关键字var"
	fmt.Println(testStr)
}

func TestPrintf(t *testing.T) {
	fmt.Println("这是一个GO0000000")
	fmt.Println("\"...\"")

	fmt.Print("hahha")

	fmt.Println("test", "s4", "s5", "s7")

	fmt.Println(9809 + 7868)

	fmt.Println(true || false)

	var str string = "ceshi"

	fmt.Println(str)

	str2 := "什么"
	fmt.Println(str2)

}
