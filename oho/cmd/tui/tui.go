package tui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd TUI 控制命令
var Cmd = &cobra.Command{
	Use:   "tui",
	Short: "TUI 控制命令",
	Long:  "控制 TUI 界面行为",
}

var (
	title     string
	message   string
	variant   string
	command   string
	controlID string
	body      string

	appendPromptCmd = &cobra.Command{
		Use:   "append-prompt <text>",
		Short: "向提示词追加文本",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/append-prompt", map[string]string{"text": args[0]})
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("提示词已追加")
			}
			return nil
		},
	}

	openHelpCmd = &cobra.Command{
		Use:   "open-help",
		Short: "打开帮助对话框",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/open-help", nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("帮助对话框已打开")
			}
			return nil
		},
	}

	openSessionsCmd = &cobra.Command{
		Use:   "open-sessions",
		Short: "打开会话选择器",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/open-sessions", nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("会话选择器已打开")
			}
			return nil
		},
	}

	openThemesCmd = &cobra.Command{
		Use:   "open-themes",
		Short: "打开主题选择器",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/open-themes", nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("主题选择器已打开")
			}
			return nil
		},
	}

	openModelsCmd = &cobra.Command{
		Use:   "open-models",
		Short: "打开模型选择器",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/open-models", nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("模型选择器已打开")
			}
			return nil
		},
	}

	submitPromptCmd = &cobra.Command{
		Use:   "submit-prompt",
		Short: "提交当前提示词",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/submit-prompt", nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("提示词已提交")
			}
			return nil
		},
	}

	clearPromptCmd = &cobra.Command{
		Use:   "clear-prompt",
		Short: "清除提示词",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Post(ctx, "/tui/clear-prompt", nil)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("提示词已清除")
			}
			return nil
		},
	}

	executeCommandCmd = &cobra.Command{
		Use:   "execute-command",
		Short: "执行命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			if command == "" {
				return fmt.Errorf("请提供 --command 参数")
			}

			c := client.NewClient()
			ctx := context.Background()

			req := types.TUICommandRequest{Command: command}
			resp, err := c.Post(ctx, "/tui/execute-command", req)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Printf("命令 %s 已执行\n", command)
			}
			return nil
		},
	}

	showToastCmd = &cobra.Command{
		Use:   "show-toast",
		Short: "显示提示消息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if message == "" {
				return fmt.Errorf("请提供 --message 参数")
			}

			c := client.NewClient()
			ctx := context.Background()

			req := types.TUIToastRequest{
				Title:   title,
				Message: message,
				Variant: variant,
			}

			resp, err := c.Post(ctx, "/tui/show-toast", req)
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("提示消息已显示")
			}
			return nil
		},
	}

	controlNextCmd = &cobra.Command{
		Use:   "control-next",
		Short: "等待下一个控制请求",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Get(ctx, "/tui/control/next")
			if err != nil {
				return err
			}

			fmt.Printf("控制请求：%s\n", string(resp))
			return nil
		},
	}

	controlResponseCmd = &cobra.Command{
		Use:   "control-response",
		Short: "响应控制请求",
		RunE: func(cmd *cobra.Command, args []string) error {
			if body == "" {
				return fmt.Errorf("请提供 --body 参数")
			}

			c := client.NewClient()
			ctx := context.Background()

			var bodyData interface{}
			if err := json.Unmarshal([]byte(body), &bodyData); err != nil {
				return fmt.Errorf("解析 body 失败：%w", err)
			}

			resp, err := c.Post(ctx, "/tui/control/response", map[string]interface{}{"body": bodyData})
			if err != nil {
				return err
			}

			var success bool
			if err := json.Unmarshal(resp, &success); err != nil {
				return err
			}

			if success {
				fmt.Println("控制请求已响应")
			}
			return nil
		},
	}
)

func init() {
	Cmd.AddCommand(appendPromptCmd)
	Cmd.AddCommand(openHelpCmd)
	Cmd.AddCommand(openSessionsCmd)
	Cmd.AddCommand(openThemesCmd)
	Cmd.AddCommand(openModelsCmd)
	Cmd.AddCommand(submitPromptCmd)
	Cmd.AddCommand(clearPromptCmd)
	Cmd.AddCommand(executeCommandCmd)
	Cmd.AddCommand(showToastCmd)
	Cmd.AddCommand(controlNextCmd)
	Cmd.AddCommand(controlResponseCmd)

	showToastCmd.Flags().StringVar(&title, "title", "", "消息标题")
	showToastCmd.Flags().StringVar(&message, "message", "", "消息内容")
	showToastCmd.Flags().StringVar(&variant, "variant", "info", "消息类型 (info/warning/error/success)")
	showToastCmd.MarkFlagRequired("message")

	executeCommandCmd.Flags().StringVar(&command, "command", "", "要执行的命令")
	executeCommandCmd.MarkFlagRequired("command")

	controlResponseCmd.Flags().StringVar(&body, "body", "", "响应体 (JSON 格式)")
	controlResponseCmd.MarkFlagRequired("body")
}
