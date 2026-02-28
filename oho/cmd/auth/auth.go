package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
)

// Cmd 认证命令
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "认证管理",
	Long:  "管理认证凭据",
}

var (
	credentials []string

	setCmd = &cobra.Command{
		Use:   "set <provider>",
		Short: "设置认证凭据",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			// 解析凭据
			credsMap := make(map[string]string)
			for _, c := range credentials {
				if idx := indexOf(c, "="); idx > 0 {
					key := c[:idx]
					value := c[idx+1:]
					credsMap[key] = value
				}
			}

			if len(credsMap) == 0 {
				return fmt.Errorf("请提供至少一个凭据 (--credentials key=value)")
			}

			req := map[string]interface{}{
				"providerId":  args[0],
				"credentials": credsMap,
			}

			resp, err := c.Put(ctx, fmt.Sprintf("/auth/%s", args[0]), req)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Printf("提供商 %s 的认证凭据已设置\n", args[0])
			}
			return nil
		},
	}
)

func init() {
	Cmd.AddCommand(setCmd)

	setCmd.Flags().StringArrayVar(&credentials, "credentials", nil, "认证凭据 (key=value 格式)")
	_ = setCmd.MarkFlagRequired("credentials")
}

func indexOf(s string, substr string) int {
	for i := 0; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
