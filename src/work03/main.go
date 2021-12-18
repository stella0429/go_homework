package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	serveAddress string = ":8081"
)

var (
	severExit = make(chan struct{})
)

//serve启动处理函数
func register(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello,welcome to register！")
}

//serve关闭处理函数
func shutDown(w http.ResponseWriter, r *http.Request) {
	severExit <- struct{}{}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	//设置多路复用处理函数
	mux := http.NewServeMux()

	//模拟单个serve的启动和关闭
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/exit", shutDown)

	server := http.Server{
		Addr:    serveAddress,
		Handler: mux,
	}

	//监听serve,当监听到server.Shutdown()，返回错误；
	g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			log.Println("Listen serve goroutine exit...! The error is : ", err.Error())
		}
		return errors.Wrapf(err, fmt.Sprintf("Listen serve goroutine return error!"))
	})

	//当接收到shutdown请求，会执行case <-severExit分支，然后延时5s后真正执行Shutdown(),并返回一个err；
	//Listen serve goroutine此时接收到Shutdown退出信号，也直接return；
	//根据源码实现，当errgroup有一个错误返回；就会调用cancel()函数，关闭相关的所有goroutine;
	//因此，signal goroutine会进入case <-ctx.Done():分支，直接return；
	g.Go(func() error {
		select {
		case <-severExit:
			log.Println("Shutdown goroutine will exit...!")
		case <-ctx.Done():
			log.Println("Shutdown goroutine will exit...！ ctx error,", ctx.Err().Error())
		}

		//虽然接受到了退出信号，模拟优雅退出处理方法，延迟5s，等到其他可能未处理完的任务处理完，再退出
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(timeoutCtx)
		return err
	})

	//信号注册和处理，当接收到退出信号，return错误；
	//根据源码实现，当errgroup有一个错误返回；就会调用cancel()函数，关闭相关的所有goroutine;
	//因此，shutdown的goroutine会进入case <-ctx.Done():分支，再调用Shutdown(),将Listen serve goroutine的一起退出
	g.Go(func() error {
		//初始化一个channel，用来接收来来自linux的os.Signal信号值
		exitSignal := make(chan os.Signal, 0)
		//注册此通道用于接收的特定信号（Ctrl+C触发）
		signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

		select {
		case signalContext := <-exitSignal:
			log.Println("signal goroutine exit...!")
			return fmt.Errorf("signal goroutine exit...! os signal is,: %v", signalContext)
		case <-ctx.Done():
			log.Println("signal goroutine exit...! ctx error is : ,", ctx.Err().Error())
			return ctx.Err()
		}
	})

	// g.Wait 等待所有 go执行完毕后执行
	err := g.Wait()
	fmt.Println("g.Wait() return error is :", err)
	fmt.Println(ctx.Err())
}
