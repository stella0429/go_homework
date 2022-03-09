package main

import (
	"bufio"
	"context"
	pb "github.com/my/repo/internal/helloworld/service/helloworld"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8086"
)

func main() {
	// 创建连接
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("连接失败: [%v]\n", err)
		return
	}
	defer conn.Close()

	// 声明客户端
	client := pb.NewGreeterClient(conn)

	// 声明 context
	ctx := context.Background()

	//header 数据的写入
	headerData := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())), "token", "123")
	ctxH := metadata.NewOutgoingContext(ctx, headerData)
	// 后续也可以往后面添加数据
	ctxH = metadata.AppendToOutgoingContext(ctxH, "kay1", "val1", "key2", "val2")

	// 创建双向数据流
	stream, err := client.SayHelloStream(ctxH)
	if err != nil {
		log.Printf("创建数据流失败: [%v]\n", err)
	}

	// 启动一个 goroutine 接收命令行输入的指令
	go func() {
		log.Println("请输入消息...")
		data := bufio.NewReader(os.Stdin)
		for {
			// 获取 命令行输入的字符串， 以回车 \n 作为结束标志
			dataString, _ := data.ReadString('\n')

			// 向服务端发送 指令
			if err := stream.Send(&pb.HelloRequest{Input: dataString}); err != nil {
				return
			}
		}
	}()

	for {
		// 接收 header
		header, _ := stream.Header()
		log.Println("---------从服务端接收到的header：", header)
		// 接收 trailer
		trailer := stream.Trailer()
		log.Println("---------从服务端接收到的trailer：", trailer)

		// 接收从 服务端返回的数据流
		dataRes, err := stream.Recv()
		if err == io.EOF {
			log.Println("⚠️ 收到服务端的结束信号")
			break //如果收到结束信号，则退出“接收循环”，结束客户端程序
		}

		if err != nil {
			// TODO: 处理接收错误
			log.Println("接收数据出错:", err)
		}

		// 没有错误的情况下，打印来自服务端的消息
		log.Printf("[客户端收到]: %s", dataRes.Output)
	}
}
