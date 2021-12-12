package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	blogpb "blogpb"
)

// Global variables for db connection , collection and context
var db *mongo.Client
var blogdb *mongo.Collection
var mongoCtx context.Context

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

// 实现 `pb/blog.pb.go` 中 BlogServiceServer 接口定义（898行起）：
//
// type BlogServiceServer interface {
//	CreateBlog(context.Context, *CreateBlogReq) (*CreateBlogRes, error)
//	ReadBlog(context.Context, *ReadBlogReq) (*ReadBlogRes, error)
//	UpdateBlog(context.Context, *UpdateBlogReq) (*UpdateBlogRes, error)
//	DeleteBlog(context.Context, *DeleteBlogReq) (*DeleteBlogRes, error)
//	ListBlog(context.Context, *ListBlogReq) (*ListBlogRes, error)
// }
type BlogServiceServer struct {
}

func (s *BlogServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()

	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	result, err := blogdb.InsertOne(ctx, data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid := result.InsertedID.(primitive.ObjectID)
	blog.Id = oid.Hex()
	return &blogpb.CreateBlogRes{Blog: blog}, nil
}

func (s *BlogServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := blogdb.FindOne(ctx, bson.M{"_id": oid})
	data := BlogItem{}
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find blog with ObjectId %s: %v", req.GetId(), err))
	}
	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}
	return response, nil
}

func (s *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	update := bson.M{
		"author_id": blog.GetAuthorId(),
		"title":     blog.GetTitle(),
		"content":   blog.GetContent(),
	}

	filter := bson.M{"_id": oid}

	result := blogdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := BlogItem{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find blog with supplied ObjectId: %v", err))
	}

	return &blogpb.UpdateBlogRes{Blog: &blogpb.Blog{
		Id:       decoded.ID.Hex(),
		AuthorId: decoded.AuthorID,
		Title:    decoded.Title,
		Content:  decoded.Content,
	}}, nil
}

func (s *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	_, err = blogdb.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete blog with id %s: %v", req.GetId(), err))
	}
	return &blogpb.DeleteBlogRes{Success: true}, nil
}

func (s *BlogServiceServer) ListBlogs(req *blogpb.ListBlogReq, stream blogpb.BlogService_ListBlogsServer) error {
	data := &BlogItem{}
	cursor, err := blogdb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error while decoding data: %v", err))
		}
		stream.Send(&blogpb.ListBlogRes{Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		}})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	return nil
}

func main() {
	// # STEP 1

	// 配置 `log` 包
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting server on port :50051")

	// 启动 listener，50051 为默认 gRPC 端口
	listener, err := net.Listen("tcp", ":50051")
	// 异常处理
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	// # STEP 2

	// 设置 options，此处可以配置 TLS 支持等
	opts := []grpc.ServerOption{}
	// 创建新的 gRPC 服务器
	s := grpc.NewServer(opts...)
	// BlogService 实例
	srv := BlogServiceServer{}
	blogpb.RegisterBlogServiceServer(s, &srv)

	// # STEP 3

	// 初始化 MongoDB 客户端
	fmt.Println("Connecting to MongoDB")

	// 非nil的空context
	mongoCtx = context.Background()

	// 使用 context 与 options 进行连接
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// 错误处理
	if err != nil {
		log.Fatal(err)
	}

	// 判断连接是否成功
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to MongoDB")
	}

	// 绑定 collection 至全局变量，以便其他方法进行调用
	blogdb = db.Database("dev-db").Collection("blog")

	// # STEP 4

	// 在子 routine 中启动服务器
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
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}
