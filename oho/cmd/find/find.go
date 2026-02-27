package find

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 查找命令
var Cmd = &cobra.Command{
	Use:   "find",
	Short: "查找命令",
	Long:  "在项目中查找文件、符号和文本内容",
}

var (
 textCmd = &cobra.Command{
  Use:   "text <pattern>",
  Short: "在文件中搜索文本",
  Args:  cobra.ExactArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   queryParams := map[string]string{
    "pattern": args[0],
   }

   resp, err := c.GetWithQuery(ctx, "/find", queryParams)
   if err != nil {
    return err
   }

   var matches []types.FindMatch
   if err := json.Unmarshal(resp, &matches); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(matches, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   if len(matches) == 0 {
    fmt.Println("未找到匹配")
    return nil
   }

   fmt.Printf("找到 %d 个匹配:\n\n", len(matches))
   for _, m := range matches {
    fmt.Printf("📄 %s (行 %d)\n", m.Path, m.LineNumber)
    fmt.Printf("   %s\n", m.Lines)
    if len(m.Submatches) > 0 {
     for _, s := range m.Submatches {
      fmt.Printf("   └─ 匹配位置：%d-%d\n", s.Start, s.End)
     }
    }
    fmt.Println()
   }

   return nil
  },
 }

 fileCmd = &cobra.Command{
  Use:   "file <query>",
  Short: "按名称查找文件",
  Args:  cobra.ExactArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   queryParams := map[string]string{
    "query": args[0],
   }

   if fileType := cmd.Flag("type"); fileType != nil && fileType.Value.String() != "" {
    queryParams["type"] = fileType.Value.String()
   }
   if directory := cmd.Flag("directory"); directory != nil && directory.Value.String() != "" {
    queryParams["directory"] = directory.Value.String()
   }
   if limit := cmd.Flag("limit"); limit != nil && limit.Value.String() != "" {
    queryParams["limit"] = limit.Value.String()
   }

   resp, err := c.GetWithQuery(ctx, "/find/file", queryParams)
   if err != nil {
    return err
   }

   var paths []string
   if err := json.Unmarshal(resp, &paths); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(paths, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   if len(paths) == 0 {
    fmt.Println("未找到文件")
    return nil
   }

   fmt.Printf("找到 %d 个文件:\n\n", len(paths))
   for _, p := range paths {
    fmt.Printf("📄 %s\n", p)
   }

   return nil
  },
 }

 symbolCmd = &cobra.Command{
  Use:   "symbol <query>",
  Short: "查找工作区符号",
  Args:  cobra.ExactArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   queryParams := map[string]string{
    "query": args[0],
   }

   resp, err := c.GetWithQuery(ctx, "/find/symbol", queryParams)
   if err != nil {
    return err
   }

   var symbols []types.Symbol
   if err := json.Unmarshal(resp, &symbols); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(symbols, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   if len(symbols) == 0 {
    fmt.Println("未找到符号")
    return nil
   }

   fmt.Printf("找到 %d 个符号:\n\n", len(symbols))
   for _, s := range symbols {
    icon := "🔖"
    switch s.Kind {
    case "function":
     icon = "⚙️"
    case "class":
     icon = "🏛️"
    case "variable":
     icon = "📦"
    case "type":
     icon = "🔤"
    }
    fmt.Printf("%s %s (%s)\n", icon, s.Name, s.Kind)
    fmt.Printf("   位置：%s:%d:%d\n", s.Path, s.Line, s.Column)
    if s.Container != "" {
     fmt.Printf("   容器：%s\n", s.Container)
    }
   }

   return nil
  },
 }
)

func init() {
 Cmd.AddCommand(textCmd)
 Cmd.AddCommand(fileCmd)
 Cmd.AddCommand(symbolCmd)

 fileCmd.Flags().String("type", "", "文件类型限制 (file/directory)")
 fileCmd.Flags().String("directory", "", "搜索目录")
 fileCmd.Flags().Int("limit", 100, "最大结果数")
}
