package main

import (
	pb "github.com/my/repo/internal/helloworld/service/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	port = ":8086"
)

// Streamer 服务端
type Streamer struct{}

// SayHelloStream 实现了 ChatServer 接口中定义的 BidStream 方法
func (s *Streamer) SayHelloStream(stream pb.Greeter_SayHelloStreamServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("收到客户端通过context发出的终止信号")
			return ctx.Err()
		default:
			md, ok := metadata.FromIncomingContext(stream.Context()) // get context from stream
			if !ok {
				log.Println("客户端发送的header信息接收失败")
				return nil
			}
			log.Println("客户端发送的header信息：", md)

			// 接收从客户端发来的消息
			data, err := stream.Recv()
			log.Println("内容：", data)
			if err == io.EOF {
				log.Println("客户端发送的数据流结束")
				return nil
			}
			if err != nil {
				log.Println("接收数据出错:", err)
				return err
			}

			// 如果接收正常，则根据接收到的 字符串 执行相应的指令
			switch data.Input {
			case "结束对话\n":
				log.Println("收到'结束对话'指令")
				if err := stream.Send(&pb.HelloReply{Output: "收到结束指令"}); err != nil {
					return err
				}
				// 收到结束指令时，通过 return nil 终止双向数据流
				return nil

			case "返回数据流\n":
				log.Println("收到'返回数据流'指令")
				// 收到 收到'返回数据流'指令， 连续返回 10 条数据
				for i := 0; i < 10; i++ {
					if err := stream.Send(&pb.HelloReply{Output: "数据流 #" + strconv.Itoa(i)}); err != nil {
						return err
					}
				}

			default:
				// 缺省情况下， 返回 '服务端返回: ' + 输入信息
				log.Printf("[收到消息]: %s", data.Input)

				// create and send header
				header := metadata.Pairs("server-header-key", "val")
				stream.SendHeader(header)
				// create and set trailer
				trailer := metadata.Pairs("server-trailer-key", "val")
				stream.SetTrailer(trailer)

				if err := stream.Send(&pb.HelloReply{Output: "服务端返回: " + data.Input}); err != nil {
					return err
				}
			}
		}
	}
}

func main() {
	log.Println("启动服务端...")
	server := grpc.NewServer()

	// 注册 ChatServer
	pb.RegisterGreeterServer(server, &Streamer{})

	address, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	if err := server.Serve(address); err != nil {
		panic(err)
	}
}