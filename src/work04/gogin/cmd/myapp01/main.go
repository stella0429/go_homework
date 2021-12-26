package main

import (
	"fmt"
	"gogin/api/myapp01/v1/order"
	"gogin/api/myapp01/v1/user"
	"gogin/routers"
)

func main() {
	//使用wire构建依赖
	event := InitializeEvent("hello_world")
	event.Start()

	// 加载api服务的路由配置
	routers.Include(user.Routers, order.Routers)
	// 初始化路由
	r := routers.Init()
	//启动服务并监听8082端口
	if err := r.Run(":8082"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
