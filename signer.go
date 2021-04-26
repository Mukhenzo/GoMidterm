package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}){
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := range in {
		wg.Add(1)
		go func(i interface{}) {
			defer wg.Done()

			data := strconv.Itoa((i).(int))

			mu.Lock()
			md5 := DataSignerMd5(data)
			fmt.Println(md5)
			mu.Unlock()

			crcWg := &sync.WaitGroup{}
			crcWg.Add(2)
			var first string
			go func() {
				defer crcWg.Done()
				first = DataSignerCrc32(data)
			}()
			var second string
			go func() {
				defer crcWg.Done()
				second = DataSignerCrc32(md5)
			}()
			crcWg.Wait()
			fmt.Println(first)
			fmt.Println(second)

			out <- first + "~" + second
		}(i)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}){
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go func(i interface{}) {
			defer wg.Done()

			wgJob := &sync.WaitGroup{}
			mu := &sync.Mutex{}
			buffer := make([]string, 6)

			for j := 0; j < 6; j++ {
				wgJob.Add(1)
				data := strconv.Itoa(j) + i.(string)

				go func(j int) {
					defer wgJob.Done()

					data = DataSignerCrc32(data)
					fmt.Println(data)

					mu.Lock()
					buffer[j] = data
					mu.Unlock()
				}(j)
			}

			wgJob.Wait()
			result := ""

			for j := 0; j < 6; j++ {
				result += buffer[j]
			}
			out <- result
			fmt.Println(" ")
		}(i)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}){
	var result []string

	for i := range in {
		result = append(result, i.(string))
	}

	sort.Strings(result)
	out <- strings.Join(result, "_")
}

func ExecutePipeline(hashSignJobs ...job){
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, job := range hashSignJobs {
		wg.Add(1)

		out := make(chan interface{})
		go func(job func(in, out chan interface{}), in, out chan interface{}) {
			defer wg.Done()
			job(in, out)
			close(out)
		}(job, in, out)
		in = out
	}

	wg.Wait()
}

func main() {
	fmt.Println("pew")
}