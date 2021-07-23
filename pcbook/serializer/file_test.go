package serializer_test

import (
	"testing"

	"pcbook/pb"
	"pcbook/sample"
	"pcbook/serializer"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	// 所有的单元测试可以并发执行，任何数据竞争可被轻易发现
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()

	// 测试写入 binary 文件
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	// 测试写入 JSON 文件
	err = serializer.WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}

	// 测试读取 binary 文件
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)

	// 测试读取 JSON 文件
	// err = se

	require.True(t, proto.Equal(laptop1, laptop2))

}
