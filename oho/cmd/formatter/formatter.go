package formatter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 格式化器命令
var Cmd = &cobra.Command{
	Use:   "formatter",
	Short: "格式化器状态",
	Long:  "获取格式化器状态",
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "获取格式化器状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, "/formatter")
		if err != nil {
			return err
		}

		var status []types.FormatterStatus
		if err := json.Unmarshal(resp, &status); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(status, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(status) == 0 {
			fmt.Println("没有格式化器")
			return nil
		}

		fmt.Println("格式化器状态:")
		for _, s := range status {
			icon := "❌"
			if s.Status == "running" {
				icon = "✅"
			}
			fmt.Printf("%s %s (状态：%s)\n", icon, s.Name, s.Status)
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(statusCmd)
}
