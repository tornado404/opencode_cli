package file

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 文件命令
var Cmd = &cobra.Command{
	Use:   "file",
	Short: "文件管理命令",
	Long:  "管理文件，包括列出、读取内容和状态",
}

var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "列出文件和目录",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		filePath := ""
		if len(args) > 0 {
			filePath = args[0]
		}

		queryParams := map[string]string{}
		if filePath != "" {
			queryParams["path"] = filePath
		}

		resp, err := c.GetWithQuery(ctx, "/file", queryParams)
		if err != nil {
			return err
		}

		var nodes []types.FileNode
		if err := json.Unmarshal(resp, &nodes); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(nodes, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(nodes) == 0 {
			fmt.Println("空目录")
			return nil
		}

		for _, node := range nodes {
			icon := "📄"
			if node.Type == "directory" {
				icon = "📁"
			}
			fmt.Printf("%s %s\n", icon, node.Path)
		}

		return nil
	},
}

var contentCmd = &cobra.Command{
	Use:   "content <path>",
	Short: "读取文件内容",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		queryParams := map[string]string{
			"path": args[0],
		}

		resp, err := c.GetWithQuery(ctx, "/file/content", queryParams)
		if err != nil {
			return err
		}

		var content types.FileContent
		if err := json.Unmarshal(resp, &content); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("文件：%s\n", content.Path)
		fmt.Printf("编码：%s\n\n", content.Encoding)
		fmt.Println(content.Content)

		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "获取已跟踪文件的状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, "/file/status")
		if err != nil {
			return err
		}

		var files []types.File
		if err := json.Unmarshal(resp, &files); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(files, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(files) == 0 {
			fmt.Println("没有已跟踪的文件")
			return nil
		}

		fmt.Printf("共 %d 个已跟踪文件:\n\n", len(files))
		for _, f := range files {
			fmt.Printf("- %s\n", f.Path)
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(contentCmd)
	Cmd.AddCommand(statusCmd)
}
