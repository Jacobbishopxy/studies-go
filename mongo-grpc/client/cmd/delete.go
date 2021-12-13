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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a blog by its ID",
	Long: `Delete a blog by its MongoDB Unique identifier.

	If no blog post is found for the ID it will return a 'Not Found' error.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		req := &blogpb.DeleteBlogReq{Id: id}

		// 我们只需要错误信息，成功时只需要打印成功信息即可
		_, err = client.DeleteBlog(context.Background(), req)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted blog with ID: %s\n", id)

		return nil
	},
}

func init() {
	deleteCmd.Flags().StringP("id", "i", "", "The id of the blog")
	deleteCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
