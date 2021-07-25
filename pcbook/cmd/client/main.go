package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pcbook/pb"
	"pcbook/sample"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	// for i := 0; i < 10; i++ {
	// 	createLaptop(laptopClient)
	// }

	// filter := &pb.Filter{
	// 	MaxPriceUsd: 3000,
	// 	MinCpuCores: 4,
	// 	MinCpuGhz:   2.5,
	// 	MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	// }

	// searchLaptop(laptopClient, filter)

	// testUploadImage(laptopClient)

	testRateLaptop(laptopClient)
}

func createLaptop(laptopClient pb.LaptopServiceClient, laptop *pb.Laptop) {

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// 不是大问题
			log.Print("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}

	log.Printf("created laptop with id: %s", res.Id)
}

// func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
// 	log.Print("search filter: ", filter)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	req := &pb.SearchLaptopRequest{Filter: filter}
// 	stream, err := laptopClient.SearchLaptop(ctx, req)
// 	if err != nil {
// 		log.Fatal("cannot search laptop: ", err)
// 	}

// 	for {
// 		res, err := stream.Recv()
// 		if err == io.EOF {
// 			return
// 		}
// 		if err != nil {
// 			log.Fatal("cannot receive response: ", err)
// 		}

// 		laptop := res.GetLaptop()
// 		log.Print("- found: ", laptop.GetId())
// 		log.Print("  + brand: ", laptop.GetBrand())
// 		log.Print("  + name: ", laptop.GetName())
// 		log.Print("  + cpu cores: ", laptop.GetCpu().GetNumberCores())
// 		log.Print("  + cpu min ghz: ", laptop.GetCpu().GetMinGhz())
// 		log.Print("  + ram: ", laptop.GetRam())
// 		log.Print("  + price: ", laptop.GetPriceUsd())
// 	}
// }

// func uploadImage(
// 	laptopClient pb.LaptopServiceClient,
// 	laptopID string,
// 	imagePath string,
// ) {
// 	// 文件路径读取图片
// 	file, err := os.Open(imagePath)
// 	if err != nil {
// 		log.Fatal("cannot open image file: ", err)
// 	}
// 	defer file.Close()

// 	// 创建一个超时 5 秒停止 upload 的 ctx，传递给 laptopClient.UploadImage
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// 生成流
// 	stream, err := laptopClient.UploadImage(ctx)
// 	if err != nil {
// 		log.Fatal("cannot upload image: ", err)
// 	}

// 	// 创建请求
// 	req := &pb.UploadImageRequest{
// 		Data: &pb.UploadImageRequest_Info{
// 			Info: &pb.ImageInfo{
// 				LaptopId:  laptopID,
// 				ImageType: filepath.Ext(imagePath),
// 			},
// 		},
// 	}

// 	// 发送请求
// 	err = stream.Send(req)
// 	if err != nil {
// 		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
// 	}

// 	// chunk 形式读取
// 	reader := bufio.NewReader(file)
// 	buffer := make([]byte, 1024)

// 	// 满 1024 个 byte 由 stream.Send 发送一次，再进入下一个循环直至 EOF
// 	for {
// 		// 读取数据至 buffer
// 		n, err := reader.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("cannot read chunk to buffer: ", err)
// 		}

// 		req := &pb.UploadImageRequest{
// 			Data: &pb.UploadImageRequest_ChunkData{
// 				ChunkData: buffer[:n],
// 			},
// 		}

// 		err = stream.Send(req)
// 		if err != nil {
// 			// stream.RecvMsg(nil) 发送 nil 意为不期望获得任何信息，仅希望获取错误
// 			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
// 		}
// 	}

// 	// 结束发送后，从服务端获取响应
// 	res, err := stream.CloseAndRecv()
// 	if err != nil {
// 		log.Fatal("cannot receive response: ", err)
// 	}

// 	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
// }

// func testUploadImage(laptopClient pb.LaptopServiceClient) {
// 	laptop := sample.NewLaptop()
// 	createLaptop(laptopClient, laptop)
// 	uploadImage(laptopClient, laptop.GetId(), "tmp/laptop.png")
// }

func rateLaptop(
	laptopClient pb.LaptopServiceClient,
	laptopID []string,
	scores []float64,
) error {
	// 创建超时 5 秒则结束的 ctx
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot rate laptop: %v", err)
	}

	// 注意请求和响应都是并发的发送的，因此我们需要创建一个新的 go routine 来获取响应
	// 在此 go routine 中，使用 for 循环来调用 stream.Recv() 来获取服务端的响应
	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more responses")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
				return
			}

			log.Print("received response: ", res)
		}
	}()

	// stream.Send 发送请求至服务端，其中 stream.RecvMsg 只有在发送请求错误时调用，
	// 并以 nil 作为入参，意为不期望收到任何信息
	for i, laptopID := range laptopID {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopID,
			Score:    scores[i],
		}

		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}

		log.Print("sent request: ", req)
	}

	// 注意必须在发送完所有请求之后调用 stream.CloseSend() 来告诉服务端我们不再发送任何数据
	// 最后通过 waitResponse 通道读取错误信息
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close send: %v", err)
	}

	err = <-waitResponse
	return err
}

// 测试 rate laptop
func testRateLaptop(laptopClient pb.LaptopServiceClient) {
	// 创建 3 个 laptop
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		createLaptop(laptopClient, laptop)
	}

	// 创建打分
	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n)? ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := rateLaptop(laptopClient, laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}
