/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	blogpb "blogpb"
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blog posts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 创建请求
		req := &blogpb.ListBlogReq{}
		// 调用 ListBlogs 并返回一个 stream
		stream, err := client.ListBlogs(context.Background(), req)
		if err != nil {
			return err
		}
		// 开始遍历
		for {
			// stream.Recv() 返回在当前遍历下的一个 ListBlogRes 的指针
			res, err := stream.Recv()
			// 如果 stream 结束，结束循环
			if err == io.EOF {
				break
			}
			// 如果出错，返回错误
			if err != nil {
				return err
			}
			// 其它情况则打印博客
			fmt.Println(res.GetBlog())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
