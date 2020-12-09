package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func AllSignal() {
	// 接收信号的通道，必须设置其缓冲区，以防止错过信号
	c := make(chan os.Signal, 1)
	// 所有信号都会发送到c上
	signal.Notify(c)
	// 阻塞等待信号
	s := <-c
	fmt.Printf("Got signal: %s\n", s)
}

func SpecificSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	s := <-c
	fmt.Printf("Got signal: %s\n", s)
}

func MultipleChannel() {
	c1, c2 := make(chan os.Signal, 1), make(chan os.Signal, 1)

	signal.Notify(c1, syscall.SIGINT)
	signal.Notify(c2, syscall.SIGINT)

	fmt.Printf("Go signal: %s\n", <-c1)
	fmt.Printf("Go signal: %s\n", <-c2)
}

func SingleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGINT)
	s := <-c
	fmt.Printf("Got signal: %s\n", s)
	signal.Stop(c)
	time.Sleep(2 * time.Second)

	s = <-c
	fmt.Printf("Got signal: %s\n", s)
}

func main() {
	SingleSignal()
	/*
		fmt.Println("signal main")
		time.Sleep(3 * time.Second)
	*/
}
