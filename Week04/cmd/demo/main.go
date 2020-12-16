package main

import (
	"Go-000/Week04/internal/pkg/http"
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	//用于监听signal信号
	signalChan := make(chan os.Signal, 1)

	//初始化资源,获取handler
	handler := InitProjectDemo()

	//注册服务
	s := http.NewHttpServer(":8888")
	s.RegisterGirlsHandler(handler)


	// 创建带有cancel的父context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建errgroup
	group, _ := errgroup.WithContext(ctx)

	//使用errgroup启动rpc服务
	group.Go(func() error {
		return  s.Run()
	})

	// 监听signal信号
	go func() {
		signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	}()

	//监听signal信号,收到signal信号通知其他http服务退出
	go func() {
		for {
			select {
			case <-signalChan:
				log.Println("received os signal, ready cancel other running server")
				cancel()
			case <-ctx.Done():
				//优雅关闭
				s.ShutDown()
				log.Printf("http server gracefull stop")
			}
		}
	}()

	if err := group.Wait(); err != nil {
		// 收到第一个错误后，开始关闭全部server流程
		cancel()
		log.Println(err)
	}
}
