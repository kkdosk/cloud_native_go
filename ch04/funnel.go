package ch04

import (
	"fmt"
	"sync"
	"time"
)

func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int)

	var wg sync.WaitGroup

	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan int) {
			defer wg.Done()

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()
	return dest
}

func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0)

	for i := 0; i < n; i++ {
		ch := make(chan int)
		dests = append(dests, ch)

		go func() {
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}
	return dests
}

func main() {
	sources := make([]<-chan int, 0)

	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources = append(sources, ch)

		go func() {
			defer close(ch)

			for i := 1; i <= 5; i++ {
				ch <- i
				time.Sleep(time.Second)
			}
		}()
	}

	dest := Funnel(sources...)
	for d := range dest {
		fmt.Println(d)
	}

	{
		source := make(chan int)
		dests := Split(source, 5)

		go func() {
			for i := 1; i <= 10; i++ {
				source <- i
			}

			close(source)
		}()

		var wg sync.WaitGroup
		wg.Add(len(dests))

		for i, ch := range dests {
			go func(i int, d <-chan int) {
				defer wg.Done()

				for val := range d {
					fmt.Printf("#%d got %d\n", i, val)
				}
			}(i, ch)
		}

		wg.Wait()
	}
}
