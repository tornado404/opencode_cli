package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 提供商命令
var Cmd = &cobra.Command{
	Use:   "provider",
	Short: "提供商管理命令",
	Long:  "管理 AI 提供商，包括列表、认证和 OAuth",
}

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "列出所有提供商",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Get(ctx, "/provider")
			if err != nil {
				return err
			}

			// 尝试多种 JSON 结构
			var all []types.Provider
			var defaultMap map[string]string
			var connected []string

			// 尝试解析为完整结构
			var result1 struct {
				All       []types.Provider  `json:"all"`
				Default   map[string]string `json:"default"`
				Connected []string          `json:"connected"`
			}
			if err := json.Unmarshal(resp, &result1); err == nil {
				all = result1.All
				defaultMap = result1.Default
				connected = result1.Connected
			} else {
				// 尝试解析为简单数组
				if err := json.Unmarshal(resp, &all); err != nil {
					// 尝试解析为 map
					var result2 map[string]interface{}
					if err := json.Unmarshal(resp, &result2); err == nil {
						if v, ok := result2["all"].([]interface{}); ok {
							for _, item := range v {
								if data, err := json.Marshal(item); err == nil {
									var p types.Provider
									if err := json.Unmarshal(data, &p); err == nil {
										all = append(all, p)
									}
								}
							}
						}
						if v, ok := result2["default"].(map[string]interface{}); ok {
							defaultMap = make(map[string]string)
							for k, val := range v {
								if s, ok := val.(string); ok {
									defaultMap[k] = s
								}
							}
						}
						if v, ok := result2["connected"].([]interface{}); ok {
							for _, item := range v {
								if s, ok := item.(string); ok {
									connected = append(connected, s)
								}
							}
						}
					} else {
						return fmt.Errorf("解析提供商列表失败：%w", err)
					}
				}
			}

			if config.Get().JSON {
				data, _ := json.MarshalIndent(map[string]interface{}{
					"all":       all,
					"default":   defaultMap,
					"connected": connected,
				}, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			fmt.Println("所有提供商:")
			for _, p := range all {
				status := " "
				for _, c := range connected {
					if c == p.ID {
						status = "✓"
						break
					}
				}
				fmt.Printf("  %s %s (%s)\n", status, p.Name, p.ID)
			}

			if len(defaultMap) > 0 {
				fmt.Println("\n默认模型:")
				for provider, model := range defaultMap {
					fmt.Printf("  %s: %s\n", provider, model)
				}
			}

			return nil
		},
	}

	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "获取提供商认证方式",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Get(ctx, "/provider/auth")
			if err != nil {
				return err
			}

			var methods map[string][]types.ProviderAuthMethod
			if err := json.Unmarshal(resp, &methods); err != nil {
				return fmt.Errorf("解析认证方式失败：%w", err)
			}

			if config.Get().JSON {
				data, _ := json.MarshalIndent(methods, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			for provider, authMethods := range methods {
				fmt.Printf("%s:\n", provider)
				for _, m := range authMethods {
					required := ""
					if m.Required {
						required = " (必需)"
					}
					fmt.Printf("  - %s%s: %s\n", m.Type, required, m.Description)
					if m.URL != "" {
						fmt.Printf("    URL: %s\n", m.URL)
					}
				}
			}

			return nil
		},
	}

	oauthAuthorizeCmd = &cobra.Command{
		Use:   "oauth authorize <provider>",
		Short: "使用 OAuth 授权提供商",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			providerID := args[0]

			resp, err := c.Post(ctx, fmt.Sprintf("/provider/%s/oauth/authorize", providerID), nil)
			if err != nil {
				return err
			}

			var auth types.ProviderAuthAuthorization
			if err := json.Unmarshal(resp, &auth); err != nil {
				return fmt.Errorf("解析 OAuth 响应失败：%w", err)
			}

			if config.Get().JSON {
				data, _ := json.MarshalIndent(auth, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			fmt.Println("OAuth 授权信息:")
			fmt.Printf("  授权 URL: %s\n", auth.URL)
			fmt.Printf("  State: %s\n", auth.State)
			fmt.Printf("  Code Challenge: %s\n", auth.CodeChallenge)
			fmt.Println("\n请在浏览器中打开授权 URL 完成认证")

			return nil
		},
	}

	oauthCallbackCmd = &cobra.Command{
		Use:   "oauth callback <provider>",
		Short: "处理 OAuth 回调",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			providerID := args[0]

			resp, err := c.Post(ctx, fmt.Sprintf("/provider/%s/oauth/callback", providerID), nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return fmt.Errorf("解析回调响应失败：%w", err)
			}

			if success {
				fmt.Printf("提供商 %s 的 OAuth 回调处理成功\n", providerID)
			}

			return nil
		},
	}
)

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(authCmd)
	Cmd.AddCommand(oauthAuthorizeCmd)
	Cmd.AddCommand(oauthCallbackCmd)
}
