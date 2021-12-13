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
	"context"
	"fmt"

	"github.com/spf13/cobra"

	blogpb "blogpb"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog post",
	Long: `Create a new blogpost on the server through gRPC.

	A blog post requires an AuthorId, Title and Content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 从 flags 中获取数据
		author, err := cmd.Flags().GetString("author")
		if err != nil {
			return err
		}
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			return err
		}
		content, err := cmd.Flags().GetString("content")
		if err != nil {
			return err
		}

		// 创建 blog 的 protobuffer 信息
		blog := &blogpb.Blog{
			AuthorId: author,
			Title:    title,
			Content:  content,
		}

		// 包裹 blog 信息到 protobuf 的 CreateBlogReq 结构体中
		encoded := blogpb.CreateBlogReq{
			Blog: blog,
		}
		// RPC 调用
		res, err := client.CreateBlog(context.TODO(), &encoded)
		if err != nil {
			return err
		}
		fmt.Printf("Blog created: %s\n", res.Blog.Id)
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("author", "a", "", "Add an author")
	createCmd.Flags().StringP("title", "t", "", "A title for the blog")
	createCmd.Flags().StringP("content", "c", "", "The content for the blog")
	createCmd.MarkFlagRequired("author")
	createCmd.MarkFlagRequired("title")
	createCmd.MarkFlagRequired("content")
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
