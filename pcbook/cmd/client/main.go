package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pcbook/pb"
	"pcbook/sample"
)

func testUploadImage(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	createLaptop(laptopClient, laptop)
	uploadImage(laptopClient, laptop.GetId(), "tmp/laptop.png")
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	testUploadImage(laptopClient)

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

func uploadImage(laptopClient pb.LaptopServiceClient, laptopID string, imagePath string) {
	// 文件路径读取图片
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	// 创建一个超时 5 秒停止 upload 的 ctx，传递给 laptopClient.UploadImage
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 生成流
	stream, err := laptopClient.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	// 创建请求
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	// 发送请求
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	// chunk 形式读取
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	// 满 1024 个 byte 由 stream.Send 发送一次，再进入下一个循环直至 EOF
	for {
		// 读取数据至 buffer
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			// stream.RecvMsg(nil) 发送 nil 意为不期望获得任何信息，仅希望获取错误
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	// 结束发送后，从服务端获取响应
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}
