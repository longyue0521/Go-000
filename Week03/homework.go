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

	g, ctx := errgroup.WithContext(context.Background())

	// 注册监听的信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	// Start Server
	g.Go(func() (err error) {
		// errgroup不处理panic，因此要在此方法内处理
		// 小技巧：使用命名返回值，在defer内可以修改最终的返回值err
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				errStr := fmt.Sprintf("Server ListenAndServe: panic recovered: %s\n%s", r, buf)
				if err != nil {
					err = xerrors.Wrap(err, errStr)
					return
				}
				err = xerrors.New(errStr)
			}
		}()

		err = server.ListenAndServe()
		return
	})

	// Server Shutdown
	g.Go(func() (err error) {
		// errgroup不处理panic，因此要在此方法内处理
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				errStr := fmt.Sprintf("Server Shutdown/Close: panic recovered: %s\n%s", r, buf)
				if err != nil {
					err = xerrors.Wrap(err, errStr)
					return
				}
				err = xerrors.New(errStr)
			}
		}()

		// 阻塞等待
		select {
		case <-quit:
			// 退出信号，优雅关闭
			// 此后server.ListenAndServe直接返回ErrServerClosed
			err = server.Shutdown(context.Background())
		case <-ctx.Done():
			// 异常情况，优雅关闭
			err = server.Shutdown(context.Background())
		}
		return
	})

	if err := g.Wait(); xerrors.Is(err, http.ErrServerClosed) {
		fmt.Println("HTTP Server Shutdown Gracefully...")
	} else {
		fmt.Println(err)
	}
}
