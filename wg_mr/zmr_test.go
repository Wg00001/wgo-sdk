package wg_mr

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMapReduce(t *testing.T) {
	testMrStruct()
	testNoMr()
}

func exampleMr() *MapReduce[int, int, int] {
	return New[int, int, int]().
		Generate(func(source chan<- int) {
			for i := 0; i <= 100; i++ {
				source <- i
			}
		}).
		Mapper(func(item int, writer Writer[int], cancel func(error)) {
			fmt.Println("mapper:", item)
			time.Sleep(time.Second)
			writer.Write(item)
		}).
		Reducer(func(pipe <-chan int, writer Writer[int], cancel func(error)) {
			res := 0
			for v := range pipe {
				res += v
				fmt.Println("reducer:", v)
			}
			writer.Write(res)
		}).
		WithWorkers(101)
}

func testMrStruct() {
	mr := exampleMr()
	res, err := mr.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func testNoMr() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}()
	pipe := make(chan int)
	wg.Add(1)
	go func() {
		wg2 := sync.WaitGroup{}
		for i := range ch {
			wg2.Add(1)
			go func(i int) {
				fmt.Println("mapper:", i)
				time.Sleep(time.Second)
				pipe <- i
				wg2.Done()
			}(i)
		}
		go func() {
			wg2.Wait()
			close(pipe)
		}()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := range pipe {
			fmt.Println("reducer:", i)
		}
		wg.Done()
	}()
	wg.Wait()
}
