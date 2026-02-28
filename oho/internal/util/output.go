package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anomalyco/oho/internal/config"
)

// Output 输出结果
func Output(data interface{}) error {
	if config.Get().JSON {
		return OutputJSON(data)
	}
	return nil
}

// OutputJSON 以 JSON 格式输出
func OutputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// OutputText 以文本格式输出
func OutputText(format string, args ...interface{}) {
	if !config.Get().JSON {
		fmt.Printf(format, args...)
	}
}

// OutputLine 输出一行文本
func OutputLine(line string) {
	if !config.Get().JSON {
		fmt.Println(line)
	}
}

// OutputTable 输出表格
func OutputTable(headers []string, rows [][]string) {
	if config.Get().JSON {
		// 转换为 map 数组输出
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			rowMap := make(map[string]string)
			for j, cell := range row {
				if j < len(headers) {
					rowMap[headers[j]] = cell
				}
			}
			data[i] = rowMap
		}
		_ = OutputJSON(data)
		return
	}

	// 计算列宽
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 输出表头
	headerLine := ""
	for i, h := range headers {
		headerLine += fmt.Sprintf("%-*s ", colWidths[i], h)
	}
	fmt.Println(headerLine)
	fmt.Println(strings.Repeat("-", len(headerLine)))

	// 输出数据行
	for _, row := range rows {
		rowLine := ""
		for i, cell := range row {
			if i < len(colWidths) {
				rowLine += fmt.Sprintf("%-*s ", colWidths[i], cell)
			}
		}
		fmt.Println(rowLine)
	}
}

// Confirm 确认操作
func Confirm(prompt string) bool {
	if config.Get().JSON {
		return true
	}

	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	_, _ = fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

// ReadStdin 从 stdin 读取内容
func ReadStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}

	data, err := os.ReadFile("/dev/stdin")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Truncate 截断字符串
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Pluralize 复数形式
func Pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
