package main

import (
	"fmt"
	"sync"
)

func main() {
	n := 10
	ch := make(chan int, n)
	out := make(chan string, 100)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			/*
				for num := range ch {
					out <- fmt.Sprintf("id=%d,num=%d\n", id, num)
				}
			*/
			for {
				// 当close(ch)，ok为false，可退出
				num, ok := <-ch
				if !ok {
					break
				}
				out <- fmt.Sprintf("id=%d,num=%d\n", id, num)
			}
			wg.Done()
		}(i)
	}

	for j := 0; j < 100; j++ {
		ch <- j
	}
	// range ch以阻塞方式获取ch中数据，当close(ch)时，range ch能够知道ch关闭，会正常退出
	// Only the sender should close a channel, never the receiver.
	// Sending on a closed channel will cause a panic.
	/*
		Channels aren't like files;
		you don't usually need to close them.
		Closing is only necessary when the receiver must be told there are no more values coming, such as to terminate a range loop
	*/
	close(ch)
	wg.Wait()

	/*
		// 没有数据时，range out操作会阻塞，可能导致所有协程阻塞
		// range ch以阻塞方式获取数据，当close(ch)时，range ch能感知ch关闭，会正常退出
		// The loop for i := range c receives values from the channel repeatedly until it is closed.
		for s := range out {
			fmt.Println(s)
		}

		for {
			// 没有数据时<-out操作会阻塞，可能导致所有协程阻塞
			// 当close(ch)，ok为false，可退出
			// ok is false if there are no more values to receive and the channel is closed.
			s, ok := <-out
			if !ok {
				break
			}
			fmt.Println(s)
		}
	*/
	i := 0
	for {
		// 利用channel的len属性，找到退出时机。
		if len(out) == 0 {
			break
		}
		s := <-out
		fmt.Println(s)
		i++
	}
	fmt.Println(i)
}
