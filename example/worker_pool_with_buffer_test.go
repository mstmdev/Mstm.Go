package example

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerPoolWithBuffer(t *testing.T) {

	jobs := make(chan int)
	results := make(chan int, 100)
	for work := 1; work <= 3; work++ {
		go worker01(work, jobs, results)
	}

	for job := 1; job <= 9; job++ {
		//如果jobs没有缓冲 因为3个worker在第一次for循环的时候被阻塞了
		//导致没有channel准备接受数据  因而这边的发送也会被阻塞
		jobs <- job
	}

	close(jobs)

	for result := 1; result <= 9; result++ {
		fmt.Println("result:", <-results)
	}
}

func worker01(workerId int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("workId:%d  job:%d\r\n", workerId, job)
		time.Sleep(time.Second)
		//如果results没有缓冲 此处则会阻塞线程
		//即三个worker都会被阻塞
		results <- job * 2
	}
}
