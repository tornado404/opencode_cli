package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 会话命令
var Cmd = &cobra.Command{
	Use:   "session",
	Short: "会话管理命令",
	Long:  "管理 OpenCode 会话，包括创建、删除、更新等操作",
}

var (
	sessionID       string
	parentID        string
	title           string
	messageID       string
	providerID      string
	modelID         string
	permissionID    string
	permissionResp  string
	rememberPerm    bool
)

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(childrenCmd)
	Cmd.AddCommand(todoCmd)
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(forkCmd)
	Cmd.AddCommand(abortCmd)
	Cmd.AddCommand(shareCmd)
	Cmd.AddCommand(unshareCmd)
	Cmd.AddCommand(diffCmd)
	Cmd.AddCommand(summarizeCmd)
	Cmd.AddCommand(revertCmd)
	Cmd.AddCommand(unrevertCmd)
	Cmd.AddCommand(permissionsCmd)

	// 全局会话标志
	Cmd.PersistentFlags().StringVarP(&sessionID, "session", "s", "", "会话 ID")
}

// listCmd 列出所有会话
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有会话",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, "/session")
		if err != nil {
			return err
		}

		var sessions []types.Session
		if err := json.Unmarshal(resp, &sessions); err != nil {
			return err
		}

		return outputSessions(sessions)
	},
}

// createCmd 创建新会话
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新会话",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{}
		if parentID != "" {
			req["parentID"] = parentID
		}
		if title != "" {
			req["title"] = title
		}

		resp, err := c.Post(ctx, "/session", req)
		if err != nil {
			return err
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return err
		}

		fmt.Printf("会话创建成功:\n")
		fmt.Printf("  ID: %s\n", session.ID)
		if session.Title != "" {
			fmt.Printf("  标题：%s\n", session.Title)
		}
		fmt.Printf("  模型：%s\n", session.Model)

		return nil
	},
}

// statusCmd 获取所有会话状态
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "获取所有会话状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, "/session/status")
		if err != nil {
			return err
		}

		var status map[string]types.SessionStatus
		if err := json.Unmarshal(resp, &status); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(status, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for id, s := range status {
			fmt.Printf("%s: %s (就绪：%v, 工作中：%v)\n", id, s.Status, s.IsReady, s.IsWorking)
		}

		return nil
	},
}

// getCmd 获取会话详情
var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "获取会话详情",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, fmt.Sprintf("/session/%s", id))
		if err != nil {
			return err
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return err
		}

		return outputSessions([]types.Session{session})
	},
}

// deleteCmd 删除会话
var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "删除会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Delete(ctx, fmt.Sprintf("/session/%s", id))
		if err != nil {
			return err
		}

		var deleted bool
		if err := json.Unmarshal(resp, &deleted); err != nil {
			return err
		}

		if deleted {
			fmt.Printf("会话 %s 已删除\n", id)
		}

		return nil
	},
}

// updateCmd 更新会话
var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "更新会话属性",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}
		if title == "" {
			return fmt.Errorf("请使用 --title 指定新标题")
		}

		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{"title": title}
		resp, err := c.Patch(ctx, fmt.Sprintf("/session/%s", id), req)
		if err != nil {
			return err
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return err
		}

		fmt.Printf("会话标题已更新为：%s\n", session.Title)
		return nil
	},
}

// childrenCmd 获取子会话
var childrenCmd = &cobra.Command{
	Use:   "children [id]",
	Short: "获取子会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, fmt.Sprintf("/session/%s/children", id))
		if err != nil {
			return err
		}

		var sessions []types.Session
		if err := json.Unmarshal(resp, &sessions); err != nil {
			return err
		}

		return outputSessions(sessions)
	},
}

// todoCmd 获取待办事项
var todoCmd = &cobra.Command{
	Use:   "todo [id]",
	Short: "获取会话待办事项",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, fmt.Sprintf("/session/%s/todo", id))
		if err != nil {
			return err
		}

		var todos []types.Todo
		if err := json.Unmarshal(resp, &todos); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(todos, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, todo := range todos {
			status := "☐"
			if todo.Status == "completed" {
				status = "☑"
			}
			fmt.Printf("%s %s\n", status, todo.Content)
		}

		return nil
	},
}

// initCmd 初始化 AGENTS.md
var initCmd = &cobra.Command{
	Use:   "init [id]",
	Short: "分析应用并创建 AGENTS.md",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}
		if providerID == "" || modelID == "" {
			return fmt.Errorf("请提供 --provider 和 --model 参数")
		}

		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{
			"messageID":  messageID,
			"providerID": providerID,
			"modelID":    modelID,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/init", id), req)
		if err != nil {
			return err
		}

		var success bool
		if err := json.Unmarshal(resp, &success); err != nil {
			return err
		}

		if success {
			fmt.Println("AGENTS.md 创建成功")
		}

		return nil
	},
}

// forkCmd 分叉会话
var forkCmd = &cobra.Command{
	Use:   "fork [id]",
	Short: "在某条消息处分叉会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{}
		if messageID != "" {
			req["messageID"] = messageID
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/fork", id), req)
		if err != nil {
			return err
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return err
		}

		fmt.Printf("会话分叉成功:\n")
		fmt.Printf("  新 ID: %s\n", session.ID)
		return nil
	},
}

// abortCmd 中止会话
var abortCmd = &cobra.Command{
	Use:   "abort [id]",
	Short: "中止正在运行的会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/abort", id), nil)
		if err != nil {
			return err
		}

		var success bool
		if err := json.Unmarshal(resp, &success); err != nil {
			return err
		}

		if success {
			fmt.Printf("会话 %s 已中止\n", id)
		}

		return nil
	},
}

// shareCmd 分享会话
var shareCmd = &cobra.Command{
	Use:   "share [id]",
	Short: "分享会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/share", id), nil)
		if err != nil {
			return err
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return err
		}

		fmt.Printf("会话已分享\n")
		return nil
	},
}

// unshareCmd 取消分享会话
var unshareCmd = &cobra.Command{
	Use:   "unshare [id]",
	Short: "取消分享会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Delete(ctx, fmt.Sprintf("/session/%s/share", id))
		if err != nil {
			return err
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return err
		}

		fmt.Printf("会话已取消分享\n")
		return nil
	},
}

// diffCmd 获取会话差异
var diffCmd = &cobra.Command{
	Use:   "diff [id]",
	Short: "获取会话差异",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		queryParams := map[string]string{}
		if messageID != "" {
			queryParams["messageID"] = messageID
		}

		resp, err := c.GetWithQuery(ctx, fmt.Sprintf("/session/%s/diff", id), queryParams)
		if err != nil {
			return err
		}

		var diffs []types.FileDiff
		if err := json.Unmarshal(resp, &diffs); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(diffs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, diff := range diffs {
			fmt.Printf("文件：%s (状态：%s)\n", diff.Path, diff.Status)
		}

		return nil
	},
}

// summarizeCmd 总结会话
var summarizeCmd = &cobra.Command{
	Use:   "summarize [id]",
	Short: "总结会话",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}
		if providerID == "" || modelID == "" {
			return fmt.Errorf("请提供 --provider 和 --model 参数")
		}

		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{
			"providerID": providerID,
			"modelID":    modelID,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/summarize", id), req)
		if err != nil {
			return err
		}

		var success bool
		if err := json.Unmarshal(resp, &success); err != nil {
			return err
		}

		if success {
			fmt.Println("会话总结完成")
		}

		return nil
	},
}

// revertCmd 回退消息
var revertCmd = &cobra.Command{
	Use:   "revert [id]",
	Short: "回退消息",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}
		if messageID == "" {
			return fmt.Errorf("请提供 --message 参数")
		}

		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{
			"messageID": messageID,
		}
		if permissionID != "" {
			req["partID"] = permissionID
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/revert", id), req)
		if err != nil {
			return err
		}

		var success bool
		if err := json.Unmarshal(resp, &success); err != nil {
			return err
		}

		if success {
			fmt.Println("消息已回退")
		}

		return nil
	},
}

// unrevertCmd 恢复已回退的消息
var unrevertCmd = &cobra.Command{
	Use:   "unrevert [id]",
	Short: "恢复所有已回退的消息",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/unrevert", id), nil)
		if err != nil {
			return err
		}

		var success bool
		if err := json.Unmarshal(resp, &success); err != nil {
			return err
		}

		if success {
			fmt.Println("已恢复所有回退的消息")
		}

		return nil
	},
}

// permissionsCmd 响应权限请求
var permissionsCmd = &cobra.Command{
	Use:   "permissions [id] [permissionID]",
	Short: "响应权限请求",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := sessionID
		if len(args) > 0 {
			id = args[0]
		}
		if id == "" {
			return fmt.Errorf("请提供会话 ID 或使用 -s 标志")
		}

		permID := permissionID
		if len(args) > 1 {
			permID = args[1]
		}
		if permID == "" {
			return fmt.Errorf("请提供权限 ID")
		}
		if permissionResp == "" {
			return fmt.Errorf("请提供 --response 参数 (allow/deny)")
		}

		c := client.NewClient()
		ctx := context.Background()

		req := map[string]interface{}{
			"response": permissionResp,
			"remember": rememberPerm,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/permissions/%s", id, permID), req)
		if err != nil {
			return err
		}

		var success bool
		if err := json.Unmarshal(resp, &success); err != nil {
			return err
		}

		if success {
			fmt.Printf("权限请求 %s 已响应：%s\n", permID, permissionResp)
		}

		return nil
	},
}

func outputSessions(sessions []types.Session) error {
	if config.Get().JSON {
		data, _ := json.MarshalIndent(sessions, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(sessions) == 0 {
		fmt.Println("没有会话")
		return nil
	}

	fmt.Printf("共 %d 个会话:\n\n", len(sessions))
	for _, s := range sessions {
		fmt.Printf("ID:     %s\n", s.ID)
		if s.Title != "" {
			fmt.Printf("标题：   %s\n", s.Title)
		}
		fmt.Printf("模型：   %s\n", s.Model)
		if s.Agent != "" {
			fmt.Printf("代理：   %s\n", s.Agent)
		}
		fmt.Println("---")
	}

	return nil
}
