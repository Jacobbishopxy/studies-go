package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	blogpb "blogpb"
	"server/blog"
)

func main() {
	// 配置 `log` 包
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting server on port :50051")

	// 启动 listener，50051 为默认 gRPC 端口
	listener, err := net.Listen("tcp", ":50051")
	// 异常处理
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	// 设置 options，此处可以配置 TLS 支持等
	opts := []grpc.ServerOption{}
	// 创建新的 gRPC 服务器
	s := grpc.NewServer(opts...)
	// BlogService 实例
	srv, err := blog.NewBlogServiceServer("mongodb://localhost:27017", "blog", "posts")
	if err != nil {
		log.Fatalf("Unable to create new blog service: %v", err)
	}
	// 注册 BlogService 的 gRPC 服务
	blogpb.RegisterBlogServiceServer(s, srv)

	// 在子 routine 中启动 gRPC 服务器
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// 使用 SHUTDOWN HOOK 用于正确的停止服务器
	// 创建一个 channel 用于接收 OS 信号
	c := make(chan os.Signal, 1)

	// CTRL+C 信号，忽略其他信号
	signal.Notify(c, os.Interrupt)

	// 阻塞主 routine 直到信号被接收
	<-c

	// 在接收到 CTRL+C 后关闭服务器
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	srv.Disconnect()
	fmt.Println("Done.")
}
