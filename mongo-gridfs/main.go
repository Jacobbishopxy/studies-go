package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 连接 MongoDB
func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client

	uri := "mongodb://localhost:27017"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(5)

	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		fmt.Println(err.Error())
	}

	return client
}

func UploadFile(file, filename string) {
	// 读取本地数据
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// 连接数据库
	conn := InitiateMongoClient()
	bucket, err := gridfs.NewBucket(conn.Database("myFiles"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// 文件流
	uploadStream, err := bucket.OpenUploadStream(filename) // filename 为存储于数据库中的文件名称
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("Write file to DB was successful. File size: %d \n", fileSize)
}

func DownloadFile(filename string) {
	conn := InitiateMongoClient()

	// GridFs 文件
	db := conn.Database("myFiles")
	fsFiles := db.Collection("fs.files")
	ctx, cfn := context.WithTimeout(context.Background(), 10*time.Second)
	// timeout 后 cancel
	defer cfn()

	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(db)

	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(filename, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v \n", dStream)
	ioutil.WriteFile(filename, buf.Bytes(), 0600)
}

func main() {

	ud := os.Args[1]
	file := os.Args[2]
	filename := path.Base(file)

	switch ud {
	case "upload":
		UploadFile(file, filename)
	case "download":
		DownloadFile(filename)
	default:
		fmt.Println("Please use following:")
		fmt.Println("go run main.go upload <your filename>")
		fmt.Println("go run main.go download <your filename>")
	}
}
