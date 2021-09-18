package example

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerPoolWithoutBuffer(t *testing.T) {

	jobs := make(chan int)
	results := make(chan int)
	for work := 1; work <= 3; work++ {
		go worker02(work, jobs, results)
	}

	go func() {
		for result := 1; result <= 9; result++ {
			fmt.Println("result:", <-results)
		}
	}()

	for job := 1; job <= 9; job++ {
		jobs <- job
	}

	close(jobs)
	close(results)
}

func worker02(workerId int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("workId:%d  job:%d\r\n", workerId, job)
		time.Sleep(time.Second)
		results <- job * 2
	}
}
