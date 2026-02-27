package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 代理命令
var Cmd = &cobra.Command{
	Use:   "agent",
	Short: "代理命令",
	Long:  "列出和管理 AI 代理",
}

var listCmd = &cobra.Command{
 Use:   "list",
 Short: "列出所有代理",
 RunE: func(cmd *cobra.Command, args []string) error {
  c := client.NewClient()
  ctx := context.Background()

  resp, err := c.Get(ctx, "/agent")
  if err != nil {
   return err
  }

  var agents []types.Agent
  if err := json.Unmarshal(resp, &agents); err != nil {
   return err
  }

  if config.Get().JSON {
   data, _ := json.MarshalIndent(agents, "", "  ")
   fmt.Println(string(data))
   return nil
  }

  if len(agents) == 0 {
   fmt.Println("没有可用代理")
   return nil
  }

  fmt.Printf("共 %d 个代理:\n\n", len(agents))
  for _, a := range agents {
   fmt.Printf("🤖 %s\n", a.Name)
   fmt.Printf("   ID: %s\n", a.ID)
   fmt.Printf("   描述：%s\n", a.Description)
   if len(a.Tools) > 0 {
    fmt.Printf("   工具：%s\n", a.Tools)
   }
   fmt.Println()
  }

  return nil
 },
}

func init() {
 Cmd.AddCommand(listCmd)
}
