package session

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

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
	sessionID      string
	parentID       string
	title          string
	messageID      string
	providerID     string
	modelID        string
	permissionID   string
	permissionResp string
	rememberPerm   bool
	runningOnly    bool
	limit          int
	offset         int
	sortBy         string
	sortOrder      string
	statusFilter   string
	// list 过滤参数
	filterID        string
	filterTitle     string
	filterCreated   int64
	filterUpdated   int64
	filterProjectID string
	filterDirectory string
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
	Cmd.AddCommand(submitCmd)

	// 全局会话标志
	Cmd.PersistentFlags().StringVarP(&sessionID, "session", "s", "", "会话 ID")

	// listCmd 标志
	listCmd.Flags().BoolVar(&runningOnly, "running", false, "只显示正在运行的会话")
	listCmd.Flags().IntVar(&limit, "limit", 0, "限制结果数量")
	listCmd.Flags().IntVar(&offset, "offset", 0, "分页偏移量")
	listCmd.Flags().StringVar(&sortBy, "sort", "updated", "排序字段 (created/updated)")
	listCmd.Flags().StringVar(&sortOrder, "order", "desc", "排序顺序 (asc/desc)")
	listCmd.Flags().StringVar(&statusFilter, "status", "", "按状态过滤 (running/completed/error/aborted/idle)")
	// 过滤参数
	listCmd.Flags().StringVar(&filterID, "id", "", "按 ID 过滤（支持模糊查询）")
	listCmd.Flags().StringVar(&filterTitle, "title", "", "按标题过滤（支持模糊查询）")
	listCmd.Flags().Int64Var(&filterCreated, "created", 0, "按创建时间过滤（时间戳，精确匹配）")
	listCmd.Flags().Int64Var(&filterUpdated, "updated", 0, "按更新时间过滤（时间戳，精确匹配）")
	listCmd.Flags().StringVar(&filterProjectID, "project-id", "", "按项目 ID 过滤（支持模糊查询）")
	listCmd.Flags().StringVar(&filterDirectory, "directory", "", "按目录过滤（支持模糊查询）")

	// createCmd 标志
	createCmd.Flags().StringVar(&parentID, "parent", "", "父会话 ID（用于创建子会话）")
	createCmd.Flags().StringVar(&title, "title", "", "会话标题")
}

// listCmd 列出所有会话
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有会话",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		// 获取会话列表
		resp, err := c.Get(ctx, "/session")
		if err != nil {
			return err
		}

		var sessions []types.Session
		if err := json.Unmarshal(resp, &sessions); err != nil {
			return err
		}

		// 获取会话状态（用于状态过滤）
		var statusMap map[string]types.SessionStatus
		if statusFilter != "" || runningOnly {
			statusResp, err := c.Get(ctx, "/session/status")
			if err != nil {
				return err
			}
			if err := json.Unmarshal(statusResp, &statusMap); err != nil {
				return err
			}
		}

		// 应用状态过滤
		if runningOnly {
			var filteredSessions []types.Session
			for _, session := range sessions {
				if status, exists := statusMap[session.ID]; exists && status.IsWorking {
					filteredSessions = append(filteredSessions, session)
				}
			}
			sessions = filteredSessions
		} else if statusFilter != "" {
			var filteredSessions []types.Session
			for _, session := range sessions {
				status, exists := statusMap[session.ID]
				match := false
				switch statusFilter {
				case "running":
					// running 状态必须在 statusMap 中存在且 IsWorking=true
					match = exists && status.IsWorking
				case "completed", "idle":
					// completed/idle: 如果 statusMap 中没有该会话，说明不在工作中，视为 completed
					// 如果存在且不是 working 状态，也是 completed
					if !exists {
						match = true
					} else {
						match = !status.IsWorking
					}
				case "error":
					match = exists && status.Status == "error"
				case "aborted":
					match = exists && status.Status == "aborted"
				default:
					match = true
				}
				if match {
					filteredSessions = append(filteredSessions, session)
				}
			}
			sessions = filteredSessions
		}

		// 应用字段过滤
		sessions = applyFieldFilter(sessions)

		// 应用排序
		sort.Slice(sessions, func(i, j int) bool {
			var less bool
			switch sortBy {
			case "created":
				less = sessions[i].Time.Created < sessions[j].Time.Created
			case "updated", "":
				less = sessions[i].Time.Updated < sessions[j].Time.Updated
			default:
				less = sessions[i].Time.Updated < sessions[j].Time.Updated
			}
			if sortOrder == "desc" {
				return !less
			}
			return less
		})

		// 应用分页
		if limit > 0 {
			start := offset
			if start > len(sessions) {
				start = len(sessions)
			}
			end := start + limit
			if end > len(sessions) {
				end = len(sessions)
			}
			sessions = sessions[start:end]
		}

		return outputSessions(sessions)
	},
}

// createCmd 创建新会话
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新会话",
	Long:  "创建一个新的 OpenCode 会话，可选择指定父会话和标题",
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

// applyFieldFilter 应用字段过滤
func applyFieldFilter(sessions []types.Session) []types.Session {
	// 如果没有设置任何过滤条件，返回原列表
	if filterID == "" && filterTitle == "" && filterCreated == 0 &&
		filterUpdated == 0 && filterProjectID == "" && filterDirectory == "" {
		return sessions
	}

	var filtered []types.Session
	for _, s := range sessions {
		match := true

		// ID 过滤（支持模糊查询）
		if filterID != "" && !contains(s.ID, filterID) {
			match = false
		}

		// 标题过滤（支持模糊查询）
		if filterTitle != "" && !contains(s.Title, filterTitle) {
			match = false
		}

		// 创建时间过滤（精确匹配）
		if filterCreated != 0 && s.Time.Created != filterCreated {
			match = false
		}

		// 更新时间过滤（精确匹配）
		if filterUpdated != 0 && s.Time.Updated != filterUpdated {
			match = false
		}

		// 项目 ID 过滤（支持模糊查询）
		if filterProjectID != "" && !contains(s.ProjectID, filterProjectID) {
			match = false
		}

		// 目录过滤（支持模糊查询）
		if filterDirectory != "" && !contains(s.Directory, filterDirectory) {
			match = false
		}

		if match {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// contains 检查是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	if s == "" || substr == "" {
		return true
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
