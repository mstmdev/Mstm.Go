package example

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)


func httpGet() {
	resp, err := http.Get("http://localhost:26039/HttpTest/HttpGetTest?arg1=hello&arg2=world")
	printlnHttpInfo(resp, err)

}

func httpPost() {
	resp, err := http.PostForm("http://localhost:26039/HttpTest/HttpPostTest", url.Values{"arg1": {"postArg1"}, "arg2": {"postArg2"}})
	printlnHttpInfo(resp, err)
}

func httpsGet() {
	resp, err := http.Get("http://localhost:26039/HttpTest/RequireHttpsTest?arg1=hello&arg2=world")
	printlnHttpInfo(resp, err)
}

//只会返回Head信息没有Body信息
func httpHead() {
	resp, err := http.Head("http://localhost:26039/HttpTest/HttpHeadTest?arg1=head1&arg2=head2")
	printlnHttpInfo(resp, err)
}

//一个简单的http服务器
func httpServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	s := &http.Server{
		Addr: ":80",
		//Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func httpFileServer() {
	log.Fatal(http.ListenAndServe(":8181", http.FileServer(http.Dir("E:\\"))))

}

func printlnHttpInfo(resp *http.Response, err error) {
	if err != nil {
		fmt.Println("Get失败")
		return
	}
	defer resp.Body.Close()
	fmt.Println("GET结束")
	fmt.Println("错误信息：", err)
	fmt.Println("响应内容：")
	fmt.Println("Cookie:", resp.Cookies())
	_, location := resp.Location()
	fmt.Println("Location:", location)
	fmt.Println("CotentLength", resp.ContentLength)
	fmt.Println("Header:", resp.Header)
	fmt.Println("Server:", resp.Header.Get("Server"))
	fmt.Println("StatuCodeL:", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll失败")
		return
	}
	fmt.Println("Body:", string(body))
}

func httpClientDemo() {
	for {
		client := &http.Client{}
		req, errreq := http.NewRequest("GET", "http://www.baidu.com", nil)
		req.Header.Add("User-Agent", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/7.0; SLCC2;)")
		req.Header.Add("Referer", "http://www.google.com")
		resp, errResp := client.Do(req)
		resp.Body.Close()
		fmt.Println(resp)
		fmt.Println("--------------------------------------------------------")
		fmt.Println(errreq)
		fmt.Println(errResp)
	}
}
