package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 命令管理
var Cmd = &cobra.Command{
	Use:   "command",
	Short: "命令管理",
	Long:  "列出和管理斜杠命令",
}

var listCmd = &cobra.Command{
 Use:   "list",
 Short: "列出所有命令",
 RunE: func(cmd *cobra.Command, args []string) error {
  c := client.NewClient()
  ctx := context.Background()

  resp, err := c.Get(ctx, "/command")
  if err != nil {
   return err
  }

  var commands []types.Command
  if err := json.Unmarshal(resp, &commands); err != nil {
   return err
  }

  if config.Get().JSON {
   data, _ := json.MarshalIndent(commands, "", "  ")
   fmt.Println(string(data))
   return nil
  }

  if len(commands) == 0 {
   fmt.Println("没有可用命令")
   return nil
  }

  fmt.Printf("共 %d 个命令:\n\n", len(commands))
  for _, c := range commands {
   fmt.Printf("/%s\n", c.Name)
   fmt.Printf("   描述：%s\n", c.Description)
   if c.Usage != "" {
    fmt.Printf("   用法：%s\n", c.Usage)
   }
   fmt.Println()
  }

  return nil
 },
}

func init() {
 Cmd.AddCommand(listCmd)
}
