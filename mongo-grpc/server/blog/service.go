package blog

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	blogpb "blogpb"
)

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
	db       *mongo.Client
	blogdb   *mongo.Collection
	mongoCtx context.Context
}

// 构造 BlogServiceServer 实例
func NewBlogServiceServer(uri, database, collection string) (*BlogServiceServer, error) {
	fmt.Println("Connecting to MongoDB")

	// 非 nil 的空 context
	mongoCtx := context.Background()
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 检查连接
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to MongoDB")
	}
	// 获取 blog 的 collection
	blogdb := db.Database(database).Collection(collection)

	return &BlogServiceServer{
		db:       db,
		blogdb:   blogdb,
		mongoCtx: mongoCtx,
	}, nil
}

// 断开 MongoDB 连接
func (s *BlogServiceServer) Disconnect() {
	s.db.Disconnect(s.mongoCtx)
}

// 创建博客
func (s *BlogServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	// 从 protobuf 请求类型中获取 blog 类型
	// 即 req.Blog，并附带 nil 检查
	blog := req.GetBlog()
	// 转换 blog 成为 BSON 类型
	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	// 插入数据至 MongoDB
	result, err := s.blogdb.InsertOne(ctx, data)
	// 错误处理
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	// 从 result 中获取由 MongoDB 生成的 ObjectId，类型声明为 primitive.ObjectID
	oid := result.InsertedID.(primitive.ObjectID)
	// ObjectId 转换为 string，并重新写入至 proto 类型
	blog.Id = oid.Hex()
	return &blogpb.CreateBlogRes{Blog: blog}, nil
}

// 根据ID查询博客
func (s *BlogServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	// 转换 string id （由 proto 而来） 成为 MongoDB 的 ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := s.blogdb.FindOne(ctx, bson.M{"_id": oid})
	// 创建一个空的 BlogItem
	decoded := BlogItem{}
	// 将结果解码到 data 中
	if err := result.Decode(&decoded); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find blog with ObjectId %s: %v", req.GetId(), err))
	}

	// 将 data 转换为 proto 类型
	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id:       decoded.ID.Hex(),
			AuthorId: decoded.AuthorID,
			Title:    decoded.Title,
			Content:  decoded.Content,
		},
	}

	return response, nil
}

// 更新博客
func (s *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	// 从 protobuf 请求类型中获取 blog 类型
	blog := req.GetBlog()
	// 转换 string 的 Id 至 MongoDB 的 ObjectId
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	// 转换 blog 成为 BSON 类型
	update := bson.M{
		"author_id": blog.GetAuthorId(),
		"title":     blog.GetTitle(),
		"content":   blog.GetContent(),
	}

	// 转换 ObjectId 成为 BSON 的 document 用于 id 查询
	filter := bson.M{"_id": oid}

	// BSON 类型的返回
	// 为了使返回为更新后的 document 而不是原始的 document,需要添加 options
	result := s.blogdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	// result 解码后的 BlogItem 类型
	decoded := BlogItem{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find blog with supplied ObjectId: %v", err))
	}

	// 将 decoded 转换为 proto 类型
	response := &blogpb.UpdateBlogRes{Blog: &blogpb.Blog{
		Id:       decoded.ID.Hex(),
		AuthorId: decoded.AuthorID,
		Title:    decoded.Title,
		Content:  decoded.Content,
	}}

	return response, nil
}

// 删除博客
func (s *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	// 将请求中的 string Id 转换至 MongoDB 的 ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	// DeleteOne 返回的结果为一个包含了被删除的数据条数的结构体（本例总是为一），因此只需要返回布尔值
	_, err = s.blogdb.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete blog with id %s: %v", req.GetId(), err))
	}

	return &blogpb.DeleteBlogRes{Success: true}, nil
}

// 查询博客列表，服务端流式 gRPC
func (s *BlogServiceServer) ListBlogs(req *blogpb.ListBlogReq, stream blogpb.BlogService_ListBlogsServer) error {
	// 初始化 BlogItem 用于存储解码后的数据
	data := &BlogItem{}
	// collection.Find() 返回一个 Cursor 类型的游标，用于遍历数据库中的数据
	cursor, err := s.blogdb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// 函数结束后调用 Close() 关闭游标
	defer cursor.Close(context.Background())

	// cursor.Next() 返回一个 bool 值，用于判断是否还有下一个数据
	for cursor.Next(context.Background()) {
		// data 解码
		err := cursor.Decode(data)
		// 错误处理
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error while decoding data: %v", err))
		}
		// 如果没有错误则将解码后的数据转换为 proto 类型
		response := &blogpb.ListBlogRes{Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		}}
		// 流式发送给客户端
		stream.Send(response)
	}

	// 检查 cursor 是否存在任何错误
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}

	return nil
}
