package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	// 启动http服务
	g.Go(func() error {
		return startServer(ctx, ":8080")
	})

	// 监听 ctrl+c
	g.Go(func() error {
		return cancelSignal(ctx)
	})

	if err := g.Wait(); err != nil {
		fmt.Println("return err:", err.Error())
	}
}

// 启动 http server
func startServer(ctx context.Context, addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello world!")
	})
	s := http.Server{Addr: addr}

	// 监听上下文的状态
	go func() {
		<-ctx.Done()
		fmt.Println("server is shutting down")
		_ = s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

// 强制关闭 http server
func cancelSignal(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-c:
		fmt.Printf("get %v signal\n", s)
		return fmt.Errorf("get %v signal\n", s)
	case <-ctx.Done():
		return nil
	}
}
