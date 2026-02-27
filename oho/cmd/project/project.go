package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 项目命令
var Cmd = &cobra.Command{
	Use:   "project",
	Short: "项目管理命令",
	Long:  "管理 OpenCode 项目",
}

var (
 listCmd = &cobra.Command{
  Use:   "list",
  Short: "列出所有项目",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Get(ctx, "/project")
   if err != nil {
    return err
   }

   var projects []types.Project
   if err := json.Unmarshal(resp, &projects); err != nil {
    return err
   }

   return outputProjects(projects)
  },
 }

 currentCmd = &cobra.Command{
  Use:   "current",
  Short: "获取当前项目",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Get(ctx, "/project/current")
   if err != nil {
    return err
   }

   var project types.Project
   if err := json.Unmarshal(resp, &project); err != nil {
    return err
   }

   return outputProjects([]types.Project{project})
  },
 }

 pathCmd = &cobra.Command{
  Use:   "path",
  Short: "获取当前路径",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Get(ctx, "/path")
   if err != nil {
    return err
   }

   var path types.Path
   if err := json.Unmarshal(resp, &path); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(path, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   fmt.Printf("当前路径：%s\n", path.Current)
   fmt.Printf("主目录：%s\n", path.Home)
   fmt.Printf("Git 仓库：%v\n", path.IsGit)
   return nil
  },
 }

 vcsCmd = &cobra.Command{
  Use:   "vcs",
  Short: "获取 VCS 信息",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Get(ctx, "/vcs")
   if err != nil {
    return err
   }

   var vcs types.VcsInfo
   if err := json.Unmarshal(resp, &vcs); err != nil {
    return err
   }

   if config.Get().JSON {
    data, _ := json.MarshalIndent(vcs, "", "  ")
    fmt.Println(string(data))
    return nil
   }

   fmt.Printf("VCS 类型：%s\n", vcs.Type)
   fmt.Printf("分支：%s\n", vcs.Branch)
   fmt.Printf("提交：%s\n", vcs.Commit)
   fmt.Printf("远程：%s\n", vcs.Remote)
   fmt.Printf("有未提交更改：%v\n", vcs.IsDirty)
   return nil
  },
 }

 instanceDisposeCmd = &cobra.Command{
  Use:   "dispose",
  Short: "销毁当前实例",
  RunE: func(cmd *cobra.Command, args []string) error {
   c := client.NewClient()
   ctx := context.Background()

   resp, err := c.Post(ctx, "/instance/dispose", nil)
   if err != nil {
    return err
   }

   var success bool
   if err := json.Unmarshal(resp, &success); err != nil {
    return err
   }

   if success {
    fmt.Println("实例已销毁")
   }
   return nil
  },
 }
)

func init() {
 Cmd.AddCommand(listCmd)
 Cmd.AddCommand(currentCmd)
 Cmd.AddCommand(pathCmd)
 Cmd.AddCommand(vcsCmd)
 Cmd.AddCommand(instanceDisposeCmd)
}

func outputProjects(projects []types.Project) error {
 if config.Get().JSON {
  data, _ := json.MarshalIndent(projects, "", "  ")
  fmt.Println(string(data))
  return nil
 }

 if len(projects) == 0 {
  fmt.Println("没有项目")
  return nil
 }

 fmt.Printf("共 %d 个项目:\n\n", len(projects))
 for _, p := range projects {
  fmt.Printf("ID:   %s\n", p.ID)
  fmt.Printf("名称：%s\n", p.Name)
  fmt.Printf("路径：%s\n", p.Path)
  fmt.Printf("VCS:  %s\n", p.Vcs)
  fmt.Println("---")
 }

 return nil
}
