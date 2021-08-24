package main

import (
	"fmt"

	"protobuf-json/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	exampleTodo := &pb.Todo{
		TodoId:    uuid.New().String(),
		Status:    pb.Status_STATUS_CREATED,
		Content:   "Do stuff",
		CreatedAt: timestamppb.Now(),
	}

	m := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	//
	jsonBytes, _ := m.Marshal(exampleTodo)
	fmt.Println(string(jsonBytes))
}
