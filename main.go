/*
 * @Author: tinson.liu
 * @Date: 2021-03-03 12:00:21
 * @LastEditors: tinson.liu
 * @LastEditTime: 2021-04-09 16:34:57
 * @Description: In User Settings Edit
 * @FilePath: /ts-go-server/main.go
 */
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ts-go-server/initialize"
	"ts-go-server/pkg/common"
)

func main() {
	// 初始化配置
	initialize.InitConfig()
	// 初始化日志
	initialize.Logger()
	// 初始化路由
	r := initialize.Routers()
	//初始化数据库
	initialize.Mysql()
	//初始化Redis
	initialize.Redis()
	// 初始校验器
	initialize.Validate()
	// 初始化Casbin
	initialize.Casbin()

	//是否初始化数据(慎用) $ts-go-server init
	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			initialize.InitData()
			fmt.Println("数据初始化成功!")
			os.Exit(1)
		}
	}

	// 关闭cache连接池
	defer common.Redis.Close()

	// 启动服务器
	host := "0.0.0.0"
	port := common.Conf.System.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}
	go func() {
		// 加入pprof性能分析
		if err := http.ListenAndServe(":8005", nil); err != nil {
			common.Log.Error("listen pprof error: ", err)
		}
	}()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Error("listen error: ", err)
		}
	}()

	common.Log.Info(fmt.Sprintf("Server is running at %s:%d", host, port))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Error("Server forced to shutdown: ", err)
	}
	common.Log.Info("Server exiting")
}
