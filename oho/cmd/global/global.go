package global

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 全局命令
var Cmd = &cobra.Command{
	Use:   "global",
	Short: "全局命令",
	Long:  "全局操作，包括健康检查和事件流",
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "检查服务器健康状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, "/global/health")
		if err != nil {
			return err
		}

		var health types.HealthResponse
		if err := json.Unmarshal(resp, &health); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(health, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		status := "❌ 不健康"
		if health.Healthy {
			status = "✅ 健康"
		}
		fmt.Printf("服务器状态：%s\n", status)
		fmt.Printf("版本：%s\n", health.Version)
		return nil
	},
}

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "监听全局事件流 (SSE)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		eventChan, errChan, err := c.SSEStream(ctx, "/global/event")
		if err != nil {
			return err
		}

		fmt.Println("正在监听全局事件... (Ctrl+C 停止)")

		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					return nil
				}
				fmt.Printf("%s", event)
			case err, ok := <-errChan:
				if ok && err != nil {
					return err
				}
			case <-ctx.Done():
				return nil
			}
		}
	},
}

func init() {
	Cmd.AddCommand(healthCmd)
	Cmd.AddCommand(eventCmd)
}
