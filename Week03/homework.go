package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	xerrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {

	//使用默认路由创建 HTTP Server
	server := http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: http.DefaultServeMux,
	}

	// 注册处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome Week03 Go! Go! Go!")
	})

	g := new(errgroup.Group)
	sig := make(chan os.Signal, 1)

	// 注册linux信号
	g.Go(func() (err error) {
		// errgroup不处理panic，因此要在此方法内处理
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err = xerrors.New(fmt.Sprintf("signal notify: panic recovered: %s\n%s", r, buf))
			}
		}()
		// 注册监听的信号
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		return nil
	})

	// Server Setup
	g.Go(func() (err error) {
		// errgroup不处理panic，因此要在此方法内处理
		defer func() {

			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				errStr := fmt.Sprintf("server setup: panic recovered: %s\n%s", r, buf)
				if err != nil {
					err = xerrors.Wrap(err, errStr)
				}
				err = xerrors.New(errStr)
			}
		}()
		err = server.ListenAndServe()
		return err
	})

	// Server Shutdown
	g.Go(func() (err error) {
		// errgroup不处理panic，因此要在此方法内处理
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				errStr := fmt.Sprintf("server shutdown: panic recovered: %s\n%s", r, buf)
				if err != nil {
					err = xerrors.Wrap(err, errStr)
				}
				err = xerrors.New(errStr)
			}
		}()

		// 阻塞等待结束信号
		<-sig
		// 优雅关闭
		err = server.Shutdown(context.Background())
		return err
	})

	if err := g.Wait(); err != http.ErrServerClosed {
		fmt.Println(err)
	} else {
		fmt.Println("HTTP Server Shutdown Gracefully...")
	}
}
