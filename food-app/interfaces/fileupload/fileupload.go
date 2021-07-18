package fileupload

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	// "github.com/minio/minio-go/v6"
)

func NewFileUpload() *fileUpload {
	return &fileUpload{}
}

type fileUpload struct{}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

var _ UploadFileInterface = &fileUpload{}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer f.Close()

	size := file.Size
	fmt.Println("the size: ", size)
	if size > int64(512000) {
		return "", errors.New("sorry, please upload an Image of 500KB or less")
	}
	// 仅需要开头的 512 字节用于判断一个文件的类型，因此不需要读取整个文件
	buffer := make([]byte, size)
	f.Read(buffer)
	fileType := http.DetectContentType(buffer)
	// 判断是否为图片
	if !strings.HasPrefix(fileType, "image") {
		return "", errors.New("please upload a valid image")
	}
	filePath := FormatFile(file.Filename)

	// accessKey := os.Getenv("DO_SPACES_KEY")
	// secKey := os.Getenv("DO_SPACES_SECRET")
	// endpoint := os.Getenv("DO_SPACES_ENDPOINT")
	// ssl := true

	// 初始化一个 DigitalOcean Spaces 的客户端
	// client, err := minio.New(endpoint, accessKey, secKey, ssl)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fileBytes := bytes.NewReader(buffer)
	// cacheControl := "max-age=31536000"
	// // 公有化
	// userMetaData := map[string]string{"x-amz-acl": "public-read"}
	// n, err := client.PutObject("chodapi", filePath, fileBytes, size, minio.PutObjectOptions{ContentType: fileType, CacheControl: cacheControl, UserMetadata: userMetaData})
	// if err != nil {
	// 	fmt.Println("the error", err)
	// 	return "", errors.New("something went wrong")
	// }
	// fmt.Println("Successfully uploaded bytes: ", n)
	return filePath, nil
}
