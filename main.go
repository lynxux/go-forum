package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"

	"go.uber.org/zap"

	"github.com/spf13/viper"
)

//Go web开发通用脚手架
func main() {
	/*1.加载配置
	2.初始化日志
	3.初始化Mysql
	4.初始化redis
	5.注册路由
	6.启动服务（优雅关机）
	*/
	//1 加载配置
	if err := settings.Init(); err != nil {
		fmt.Println("init settings failed! err: ", err)
		return
	}

	//2 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("init logger failed! err: ", err)
		return
	}
	//将缓存同步到文件中
	defer zap.L().Sync()
	//通过zap.L() 访问zap的全局logger
	zap.L().Debug("logger init success")

	//3 初始化Mysql
	if err := mysql.Init(settings.Conf.MySqlConfig); err != nil {
		fmt.Println("init mysql failed! err: ", err)
		return
	}
	defer mysql.Close()

	//4 初始化redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis failed! err: ", err)
		return
	}
	defer redis.Close()

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Println("init snowflake failed! err: ", err)
		return
	}

	//5 注册路由
	r := routes.Setup(settings.Conf.Mode)
	//fmt.Println("main.settings.Conf.Port: ", settings.Conf.Port)
	//fmt.Println("main.settings.Conf.Version: ", settings.Conf.Version)

	//6 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
