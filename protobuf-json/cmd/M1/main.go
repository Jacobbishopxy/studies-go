package main

import (
	"encoding/json"
	"fmt"

	"protobuf-json/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	exampleTodo := &pb.Todo{
		TodoId:    uuid.New().String(),
		Status:    pb.Status_STATUS_CREATED,
		Content:   "Do stuff",
		CreatedAt: timestamppb.Now(),
	}

	jsonBytes, _ := json.MarshalIndent(exampleTodo, "", "    ")
	fmt.Println(string(jsonBytes))
}
