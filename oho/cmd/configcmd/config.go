package configcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 配置命令
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理命令",
	Long:  "获取和更新 OpenCode 配置",
}

var (
 getCmd = &cobra.Command{
  Use:   "get",
  Short: "获取配置",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Get(ctx, "/config")
   if err != nil {
    return err
   }

   var cfg types.Config
   if err := json.Unmarshal(resp, &cfg); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(cfg, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   fmt.Println("当前配置:")
   fmt.Printf("  默认模型：%s\n", cfg.DefaultModel)
   fmt.Printf("  主题：%s\n", cfg.Theme)
   fmt.Printf("  语言：%s\n", cfg.Language)
   fmt.Printf("  最大 Token：%d\n", cfg.MaxTokens)
   fmt.Printf("  温度：%.2f\n", cfg.Temperature)
   if len(cfg.AutoApprove) > 0 {
    fmt.Printf("  自动批准：%s\n", strings.Join(cfg.AutoApprove, ", "))
   }
   return nil
  },
 }

 setCmd = &cobra.Command{
  Use:   "set",
  Short: "更新配置",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   // 构建更新请求
   updates := make(map[string]interface{})
   
   if theme != "" {
    updates["theme"] = theme
   }
   if language != "" {
    updates["language"] = language
   }
   if defaultModel != "" {
    updates["defaultModel"] = defaultModel
   }
   if maxTokens > 0 {
    updates["maxTokens"] = maxTokens
   }
   if temperature > 0 {
    updates["temperature"] = temperature
   }
   if len(autoApprove) > 0 {
    updates["autoApprove"] = autoApprove
   }

   if len(updates) == 0 {
    return fmt.Errorf("请提供至少一个要更新的配置项")
   }

   resp, err := c.Patch(ctx, "/config", updates)
   if err != nil {
    return err
   }

   var cfg types.Config
   if err := json.Unmarshal(resp, &cfg); err != nil {
    return err
   }

   fmt.Println("配置已更新")
   return nil
  },
 }

 providersCmd = &cobra.Command{
  Use:   "providers",
  Short: "列出提供商和默认模型",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Get(ctx, "/config/providers")
   if err != nil {
    return err
   }

   var result struct {
    Providers []types.Provider        `json:"providers"`
    Default   map[string]string       `json:"default"`
   }
   if err := json.Unmarshal(resp, &result); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   fmt.Println("可用提供商:")
   for _, p := range result.Providers {
    fmt.Printf("  - %s (%s)\n", p.Name, p.ID)
   }

   fmt.Println("\n默认模型:")
   for provider, model := range result.Default {
    fmt.Printf("  %s: %s\n", provider, model)
   }

   return nil
  },
 }

 theme        string
 language     string
 defaultModel string
 maxTokens    int
 temperature  float64
 autoApprove  []string
)

func init() {
 Cmd.AddCommand(getCmd)
 Cmd.AddCommand(setCmd)
 Cmd.AddCommand(providersCmd)

 setCmd.Flags().StringVar(&theme, "theme", "", "主题名称")
 setCmd.Flags().StringVar(&language, "language", "", "语言设置")
 setCmd.Flags().StringVar(&defaultModel, "model", "", "默认模型")
 setCmd.Flags().IntVar(&maxTokens, "max-tokens", 0, "最大 Token 数")
 setCmd.Flags().Float64Var(&temperature, "temperature", 0, "温度参数")
 setCmd.Flags().StringSliceVar(&autoApprove, "auto-approve", nil, "自动批准的工具列表")
}
