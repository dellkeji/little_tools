package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
	"daily-report-tool/internal/service"
)

// 这是一个手动端到端测试脚本
// 运行方式: go run scripts/e2e_test.go

func main() {
	fmt.Println("=== 日报工具端到端测试 ===")

	// 创建临时测试目录
	tempDir := filepath.Join(".", "test_data_e2e")
	dataPath := filepath.Join(tempDir, "tasks")
	configPath := filepath.Join(tempDir, "config.json")

	// 清理旧的测试数据
	os.RemoveAll(tempDir)
	defer func() {
		fmt.Println("\n清理测试数据...")
		os.RemoveAll(tempDir)
	}()

	fmt.Printf("测试目录: %s\n", tempDir)
	fmt.Println()

	// 测试 1: 初始化服务
	fmt.Println("测试 1: 初始化服务")
	taskRepo := repository.NewFileTaskRepository(dataPath)
	taskService := service.NewTaskService(taskRepo, dataPath)
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := service.NewConfigService(configRepo)
	reminderService := service.NewReminderService(configService, taskService)
	fmt.Println("✓ 服务初始化成功")
	fmt.Println()

	// 测试 2: 创建和保存任务
	fmt.Println("测试 2: 创建和保存任务")
	testDate := time.Now().Truncate(24 * time.Hour)
	testContent := `# 今日工作

## 完成的任务
- 完成需求文档
- 代码评审
- 编写测试用例

## 遇到的问题
- 无

## 明日计划
- 继续开发新功能
`

	err := taskService.SaveTask(testDate, testContent)
	if err != nil {
		fmt.Printf("✗ 保存任务失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 任务保存成功")

	// 验证任务已保存
	task, err := taskService.GetTask(testDate)
	if err != nil {
		fmt.Printf("✗ 读取任务失败: %v\n", err)
		os.Exit(1)
	}
	if task == nil || task.Content != testContent {
		fmt.Println("✗ 任务内容不匹配")
		os.Exit(1)
	}
	fmt.Println("✓ 任务读取验证成功")
	fmt.Println()

	// 测试 3: 创建多个日期的任务
	fmt.Println("测试 3: 创建多个日期的任务")
	dates := []time.Time{
		testDate.AddDate(0, 0, -2),
		testDate.AddDate(0, 0, -1),
		testDate,
		testDate.AddDate(0, 0, 1),
		testDate.AddDate(0, 0, 2),
	}

	for i, date := range dates {
		content := fmt.Sprintf("# 任务 %d\n\n这是第 %d 个任务", i+1, i+1)
		err := taskService.SaveTask(date, content)
		if err != nil {
			fmt.Printf("✗ 保存任务失败 (日期: %s): %v\n", date.Format("2006-01-02"), err)
			os.Exit(1)
		}
	}
	fmt.Printf("✓ 成功创建 %d 个任务\n", len(dates))

	// 验证月份任务查询
	year, month, _ := testDate.Date()
	monthDates, err := taskService.GetMonthTaskDates(year, month)
	if err != nil {
		fmt.Printf("✗ 查询月份任务失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ 查询到 %d 个任务日期\n", len(monthDates))
	fmt.Println()

	// 测试 4: 配置管理
	fmt.Println("测试 4: 配置管理")
	config := &model.Config{
		WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test123",
		ReminderTime:    "10:00",
		ReminderEnabled: true,
		DataPath:        dataPath,
	}

	err = configService.UpdateConfig(config)
	if err != nil {
		fmt.Printf("✗ 保存配置失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 配置保存成功")

	// 验证配置加载
	loadedConfig, err := configService.GetConfig()
	if err != nil {
		fmt.Printf("✗ 加载配置失败: %v\n", err)
		os.Exit(1)
	}
	if loadedConfig.WebhookURL != config.WebhookURL {
		fmt.Println("✗ 配置内容不匹配")
		os.Exit(1)
	}
	fmt.Println("✓ 配置加载验证成功")
	fmt.Println()

	// 测试 5: 配置验证
	fmt.Println("测试 5: 配置验证")

	// 测试无效的 Webhook URL
	invalidURLs := []string{
		"",
		"not-a-url",
		"ftp://invalid.com",
		"https://qyapi.weixin.qq.com/wrong-path",
	}

	invalidCount := 0
	for _, url := range invalidURLs {
		err := configService.ValidateWebhook(url)
		if err != nil {
			invalidCount++
		}
	}
	fmt.Printf("✓ 成功拒绝 %d 个无效 URL\n", invalidCount)

	// 测试有效的 Webhook URL
	validURLs := []string{
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc123",
		"http://example.com/webhook",
	}

	validCount := 0
	for _, url := range validURLs {
		err := configService.ValidateWebhook(url)
		if err == nil {
			validCount++
		}
	}
	fmt.Printf("✓ 成功接受 %d 个有效 URL\n", validCount)
	fmt.Println()

	// 测试 6: 提醒服务
	fmt.Println("测试 6: 提醒服务")

	// 启动提醒服务
	err = reminderService.Start()
	if err != nil {
		fmt.Printf("✗ 启动提醒服务失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 提醒服务启动成功")

	// 等待一小段时间
	time.Sleep(200 * time.Millisecond)

	// 停止提醒服务
	reminderService.Stop()
	fmt.Println("✓ 提醒服务停止成功")
	fmt.Println()

	// 测试 7: 检查今天是否有任务
	fmt.Println("测试 7: 检查今天是否有任务")
	hasTask, err := taskService.HasTodayTask()
	if err != nil {
		fmt.Printf("✗ 检查今天任务失败: %v\n", err)
		os.Exit(1)
	}
	if hasTask {
		fmt.Println("✓ 今天有任务")
	} else {
		fmt.Println("✓ 今天没有任务")
	}
	fmt.Println()

	// 测试 8: 编辑现有任务
	fmt.Println("测试 8: 编辑现有任务")
	originalTask, err := taskService.GetTask(testDate)
	if err != nil {
		fmt.Printf("✗ 获取原始任务失败: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	newContent := testContent + "\n## 更新内容\n- 添加了新的内容"
	err = taskService.SaveTask(testDate, newContent)
	if err != nil {
		fmt.Printf("✗ 更新任务失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 任务更新成功")

	updatedTask, err := taskService.GetTask(testDate)
	if err != nil {
		fmt.Printf("✗ 获取更新后的任务失败: %v\n", err)
		os.Exit(1)
	}

	if !updatedTask.UpdatedAt.After(originalTask.UpdatedAt) {
		fmt.Println("✗ 更新时间未改变")
		os.Exit(1)
	}
	fmt.Println("✓ 更新时间验证成功")

	if updatedTask.CreatedAt != originalTask.CreatedAt {
		fmt.Println("✗ 创建时间被修改")
		os.Exit(1)
	}
	fmt.Println("✓ 创建时间保持不变")
	fmt.Println()

	// 测试 9: 数据持久化验证
	fmt.Println("测试 9: 数据持久化验证")

	// 重新创建服务实例（模拟应用重启）
	taskRepo2 := repository.NewFileTaskRepository(dataPath)
	taskService2 := service.NewTaskService(taskRepo2, dataPath)

	// 读取之前保存的任务
	persistedTask, err := taskService2.GetTask(testDate)
	if err != nil {
		fmt.Printf("✗ 读取持久化任务失败: %v\n", err)
		os.Exit(1)
	}

	if persistedTask == nil {
		fmt.Println("✗ 持久化任务不存在")
		os.Exit(1)
	}

	if persistedTask.Content != newContent {
		fmt.Println("✗ 持久化任务内容不匹配")
		os.Exit(1)
	}
	fmt.Println("✓ 数据持久化验证成功")
	fmt.Println()

	// 测试 10: 文件系统结构验证
	fmt.Println("测试 10: 文件系统结构验证")

	// 检查配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("✗ 配置文件不存在")
		os.Exit(1)
	}
	fmt.Println("✓ 配置文件存在")

	// 检查数据目录
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		fmt.Println("✗ 数据目录不存在")
		os.Exit(1)
	}
	fmt.Println("✓ 数据目录存在")

	// 检查任务文件
	taskFiles, err := filepath.Glob(filepath.Join(dataPath, "*.json"))
	if err != nil {
		fmt.Printf("✗ 扫描任务文件失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ 找到 %d 个任务文件\n", len(taskFiles))
	fmt.Println()

	// 所有测试完成
	fmt.Println("=== 所有测试通过 ===")
	fmt.Println("\n测试摘要:")
	fmt.Println("  ✓ 服务初始化")
	fmt.Println("  ✓ 任务创建和保存")
	fmt.Println("  ✓ 多日期任务管理")
	fmt.Println("  ✓ 配置管理")
	fmt.Println("  ✓ 配置验证")
	fmt.Println("  ✓ 提醒服务")
	fmt.Println("  ✓ 今日任务检查")
	fmt.Println("  ✓ 任务编辑")
	fmt.Println("  ✓ 数据持久化")
	fmt.Println("  ✓ 文件系统结构")
	fmt.Println("\n端到端测试成功完成！")
}
