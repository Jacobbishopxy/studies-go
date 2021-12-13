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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	blogpb "blogpb"
)

var cfgFile string

// Client 与 context 的全局变量，可以被用于所有的子命令
var client blogpb.BlogServiceClient
var requestOpts grpc.DialOption

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blogclient",
	Short: "A gRPC client to communicate with the BlogService server",
	Long: `A gRPC client to communicate with the BlogService server.
	You can use this client to create and read blogs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// 初始化配置，读取配置文件和环境变量
	cobra.OnInitialize(initConfig)
	// 配置初始化完成后，初始化 gRPC 客户端
	fmt.Println("Starting Blog Service Client...")
	// 建立带有超时 10 秒的 context，当服务端不响应时
	// requestCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancelFn()
	// 建立非安全的 gRPC options（不带 TLS）
	requestOpts = grpc.WithInsecure()
	// 拨号到服务端，返回一个 gRPC 客户端连接
	conn, err := grpc.Dial("localhost:50051", requestOpts)
	if err != nil {
		log.Fatalf("Unable to establish client connection to localhost:50051: %v", err)
	}
	// 初始化 BlogServiceClient，导入连接
	client = blogpb.NewBlogServiceClient(conn)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".client" (without extension).
		viper.AddConfigPath(home)
		// viper.SetConfigType("yaml")
		viper.SetConfigName(".blogclient")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
