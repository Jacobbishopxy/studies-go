package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func FileDownload(c *gin.Context) {
	filename := c.Query("filename")
	conn := InitiateMongoClient()

	// GridFs 文件
	db := conn.Database("files")
	fsFiles := db.Collection("fs.files")
	ctx, cfn := context.WithTimeout(context.Background(), 10*time.Second)
	// timeout 后 cancel
	defer cfn()

	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{"filename": filename}).Decode(&results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(filename, &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileContentDisposition := fmt.Sprintf("attachment; filename=%s", filename)
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("File size to download: %v \n", dStream),
	})
}
