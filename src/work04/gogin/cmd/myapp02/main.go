package main

import (
	"fmt"
	"gogin/internal/app/myapp02/conf"
	"gogin/utils"
	"gogin/utils/transport/grpc"
	"gogin/utils/transport/http"
	"io/ioutil"
	"os"
)

func newApp(hs *http.Server, gs *grpc.Server) *App {
	return New(
		Server(
			hs,
			gs,
		),
	)
}
func main() {
	//加载json配置文件
	jsonFile, err := os.Open("../../configs/config.json")
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}

	//读取json文件内容
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//转struct
	var bc conf.Bootstrap
	err = utils.HandleContext(byteValue, &bc)
	if err != nil {
		panic(err)
	}

	fmt.Println(bc)

	app, err := initApp(bc.Server, bc.Data)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
