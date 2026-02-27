package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 工具命令
var Cmd = &cobra.Command{
	Use:   "tool",
	Short: "工具命令",
	Long:  "列出和管理实验性工具",
}

var (
	providerID string
	modelID    string

	idsCmd = &cobra.Command{
		Use:   "ids",
		Short: "列出所有工具 ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Get(ctx, "/experimental/tool/ids")
			if err != nil {
				return err
			}

			var toolIDs types.ToolIDs
			if err := json.Unmarshal(resp, &toolIDs); err != nil {
				return err
			}

			if config.Get().JSON {
				data, _ := json.MarshalIndent(toolIDs, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			if len(toolIDs.IDs) == 0 {
				fmt.Println("没有可用工具")
				return nil
			}

			fmt.Printf("共 %d 个工具:\n\n", len(toolIDs.IDs))
			for _, id := range toolIDs.IDs {
				fmt.Printf("🔧 %s\n", id)
			}

			return nil
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "列出指定模型的工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			if providerID == "" || modelID == "" {
				return fmt.Errorf("请提供 --provider 和 --model 参数")
			}

			c := client.NewClient()
			ctx := context.Background()

			queryParams := map[string]string{
				"provider": providerID,
				"model":    modelID,
			}

			resp, err := c.GetWithQuery(ctx, "/experimental/tool", queryParams)
			if err != nil {
				return err
			}

			var toolList types.ToolList
			if err := json.Unmarshal(resp, &toolList); err != nil {
				return err
			}

			if config.Get().JSON {
				data, _ := json.MarshalIndent(toolList, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			if len(toolList.Tools) == 0 {
				fmt.Println("没有可用工具")
				return nil
			}

			fmt.Printf("共 %d 个工具:\n\n", len(toolList.Tools))
			for _, t := range toolList.Tools {
				fmt.Printf("🔧 %s\n", t.Name)
				fmt.Printf("   描述：%s\n", t.Description)
				fmt.Println()
			}

			return nil
		},
	}
)

func init() {
	Cmd.AddCommand(idsCmd)
	Cmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&providerID, "provider", "", "提供商 ID")
	listCmd.Flags().StringVar(&modelID, "model", "", "模型 ID")
	listCmd.MarkFlagRequired("provider")
	listCmd.MarkFlagRequired("model")
}
