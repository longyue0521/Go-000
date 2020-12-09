package main

import (
	"fmt"
	"time"
)

func do() {
	/*
		defer func() {
			var err error
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
			}
			fmt.Println(err)
		}()
	*/
	panic("fatal error in do")
}

func other() {
	// other 无法捕获do抛出的panic， 因为其属于两个不同的go协程

	go func() {
		// 只能在此处捕获panic
		do()
	}()
}

/*
	do函数抛出的panic，可以在do内捕获，也可以在直接或间接调用do的外层函数捕获。
	但求其与do在同一个go协程内，
	下方案例中do抛出的panic无法在main中捕获，并且会导致main协程异常退出。
*/
func main() {
	defer func() {
		// 无法捕获do函数抛出的panic
		err := recover()
		fmt.Println(err, "in main")
	}()
	// do内的panic，可以在do函数内用捕获，也可以在调用的do的外层recover
	go func() {
		/*
			defer func() {
				err := recover()
				fmt.Println(err, "outer func")
			}()
		*/
		func() {
			// 可在此处捕获panic
			do()
		}()

	}()

	time.Sleep(time.Second)
}
