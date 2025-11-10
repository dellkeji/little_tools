package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
	"daily-report-tool/internal/service"
)

// TestTaskCreationEditAndSaveFlow 测试完整的任务创建、编辑和保存流程
func TestTaskCreationEditAndSaveFlow(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, "tasks")
	configPath := filepath.Join(tempDir, "config.json")

	// 初始化仓库和服务
	taskRepo := repository.NewFileTaskRepository(dataPath)
	taskService := service.NewTaskService(taskRepo, dataPath)

	// 测试日期
	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)

	t.Run("创建新任务", func(t *testing.T) {
		content := "# 今日工作\n\n- 完成需求文档\n- 代码评审"

		err := taskService.SaveTask(testDate, content)
		if err != nil {
			t.Fatalf("保存任务失败: %v", err)
		}

		// 验证任务已保存
		task, err := taskService.GetTask(testDate)
		if err != nil {
			t.Fatalf("获取任务失败: %v", err)
		}

		if task == nil {
			t.Fatal("任务不应为空")
		}

		if task.Content != content {
			t.Errorf("任务内容不匹配，期望: %s, 实际: %s", content, task.Content)
		}

		if task.Date.Format("2006-01-02") != testDate.Format("2006-01-02") {
			t.Errorf("任务日期不匹配")
		}
	})

	t.Run("编辑现有任务", func(t *testing.T) {
		// 获取原始任务
		originalTask, err := taskService.GetTask(testDate)
		if err != nil {
			t.Fatalf("获取原始任务失败: %v", err)
		}

		// 等待一秒确保更新时间不同
		time.Sleep(1 * time.Second)

		// 更新任务内容
		newContent := "# 今日工作\n\n- 完成需求文档\n- 代码评审\n- 编写测试用例"
		err = taskService.SaveTask(testDate, newContent)
		if err != nil {
			t.Fatalf("更新任务失败: %v", err)
		}

		// 验证任务已更新
		updatedTask, err := taskService.GetTask(testDate)
		if err != nil {
			t.Fatalf("获取更新后的任务失败: %v", err)
		}

		if updatedTask.Content != newContent {
			t.Errorf("任务内容未更新")
		}

		if !updatedTask.UpdatedAt.After(originalTask.UpdatedAt) {
			t.Errorf("更新时间应该晚于原始时间")
		}

		if updatedTask.CreatedAt != originalTask.CreatedAt {
			t.Errorf("创建时间不应改变")
		}
	})

	t.Run("保存多个日期的任务", func(t *testing.T) {
		dates := []time.Time{
			time.Date(2025, 11, 11, 0, 0, 0, 0, time.Local),
			time.Date(2025, 11, 12, 0, 0, 0, 0, time.Local),
			time.Date(2025, 11, 13, 0, 0, 0, 0, time.Local),
		}

		for i, date := range dates {
			content := "任务内容 " + string(rune('A'+i))
			err := taskService.SaveTask(date, content)
			if err != nil {
				t.Fatalf("保存任务失败 (日期: %s): %v", date.Format("2006-01-02"), err)
			}
		}

		// 验证所有任务都已保存
		for i, date := range dates {
			task, err := taskService.GetTask(date)
			if err != nil {
				t.Fatalf("获取任务失败 (日期: %s): %v", date.Format("2006-01-02"), err)
			}

			expectedContent := "任务内容 " + string(rune('A'+i))
			if task.Content != expectedContent {
				t.Errorf("任务内容不匹配 (日期: %s)", date.Format("2006-01-02"))
			}
		}
	})

	// 清理
	_ = configPath
}

// TestMonthTaskDatesQuery 测试日历导航和日期查询功能
func TestMonthTaskDatesQuery(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, "tasks")

	// 初始化仓库和服务
	taskRepo := repository.NewFileTaskRepository(dataPath)
	taskService := service.NewTaskService(taskRepo, dataPath)

	t.Run("查询空月份", func(t *testing.T) {
		dates, err := taskService.GetMonthTaskDates(2025, time.November)
		if err != nil {
			t.Fatalf("查询月份任务日期失败: %v", err)
		}

		if len(dates) != 0 {
			t.Errorf("空月份应该返回0个日期，实际: %d", len(dates))
		}
	})

	t.Run("查询有任务的月份", func(t *testing.T) {
		// 创建多个任务
		testDates := []time.Time{
			time.Date(2025, 11, 5, 0, 0, 0, 0, time.Local),
			time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local),
			time.Date(2025, 11, 15, 0, 0, 0, 0, time.Local),
			time.Date(2025, 11, 20, 0, 0, 0, 0, time.Local),
		}

		for _, date := range testDates {
			err := taskService.SaveTask(date, "测试任务")
			if err != nil {
				t.Fatalf("保存任务失败: %v", err)
			}
		}

		// 查询11月的任务日期
		dates, err := taskService.GetMonthTaskDates(2025, time.November)
		if err != nil {
			t.Fatalf("查询月份任务日期失败: %v", err)
		}

		if len(dates) != len(testDates) {
			t.Errorf("期望 %d 个日期，实际: %d", len(testDates), len(dates))
		}

		// 验证日期是否正确
		dateMap := make(map[string]bool)
		for _, date := range dates {
			dateMap[date.Format("2006-01-02")] = true
		}

		for _, testDate := range testDates {
			dateStr := testDate.Format("2006-01-02")
			if !dateMap[dateStr] {
				t.Errorf("缺少日期: %s", dateStr)
			}
		}
	})

	t.Run("跨月份查询", func(t *testing.T) {
		// 在10月创建任务
		octDate := time.Date(2025, 10, 25, 0, 0, 0, 0, time.Local)
		err := taskService.SaveTask(octDate, "10月任务")
		if err != nil {
			t.Fatalf("保存10月任务失败: %v", err)
		}

		// 查询11月，不应包含10月的任务
		novDates, err := taskService.GetMonthTaskDates(2025, time.November)
		if err != nil {
			t.Fatalf("查询11月任务日期失败: %v", err)
		}

		for _, date := range novDates {
			if date.Month() != time.November {
				t.Errorf("11月查询结果包含其他月份的日期: %s", date.Format("2006-01-02"))
			}
		}

		// 查询10月，应包含10月的任务
		octDates, err := taskService.GetMonthTaskDates(2025, time.October)
		if err != nil {
			t.Fatalf("查询10月任务日期失败: %v", err)
		}

		found := false
		for _, date := range octDates {
			if date.Format("2006-01-02") == octDate.Format("2006-01-02") {
				found = true
				break
			}
		}

		if !found {
			t.Error("10月查询结果应包含10月25日的任务")
		}
	})
}

// TestConfigurationManagement 测试配置更新和应用
func TestConfigurationManagement(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 初始化仓库和服务
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := service.NewConfigService(configRepo)

	t.Run("加载默认配置", func(t *testing.T) {
		config, err := configService.GetConfig()
		if err != nil {
			t.Fatalf("加载配置失败: %v", err)
		}

		if config == nil {
			t.Fatal("配置不应为空")
		}

		// 验证默认值
		if config.ReminderTime != "10:00" {
			t.Errorf("默认提醒时间应为 10:00，实际: %s", config.ReminderTime)
		}

		if config.ReminderEnabled {
			t.Error("默认提醒应为禁用状态")
		}
	})

	t.Run("更新配置", func(t *testing.T) {
		newConfig := &model.Config{
			WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test123",
			ReminderTime:    "09:30",
			ReminderEnabled: true,
			DataPath:        "./data/tasks",
		}

		err := configService.UpdateConfig(newConfig)
		if err != nil {
			t.Fatalf("更新配置失败: %v", err)
		}

		// 重新加载配置验证
		loadedConfig, err := configService.GetConfig()
		if err != nil {
			t.Fatalf("重新加载配置失败: %v", err)
		}

		if loadedConfig.WebhookURL != newConfig.WebhookURL {
			t.Errorf("Webhook URL 不匹配")
		}

		if loadedConfig.ReminderTime != newConfig.ReminderTime {
			t.Errorf("提醒时间不匹配")
		}

		if loadedConfig.ReminderEnabled != newConfig.ReminderEnabled {
			t.Errorf("提醒启用状态不匹配")
		}
	})

	t.Run("验证无效的Webhook URL", func(t *testing.T) {
		invalidURLs := []string{
			"",
			"not-a-url",
			"ftp://invalid.com",
			"http://",
			"https://qyapi.weixin.qq.com/wrong-path",
			"https://qyapi.weixin.qq.com/cgi-bin/webhook/send", // 缺少key参数
		}

		for _, url := range invalidURLs {
			err := configService.ValidateWebhook(url)
			if err == nil {
				t.Errorf("应该拒绝无效的 URL: %s", url)
			}
		}
	})

	t.Run("验证有效的Webhook URL", func(t *testing.T) {
		validURLs := []string{
			"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc123",
			"http://example.com/webhook",
			"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX",
		}

		for _, url := range validURLs {
			err := configService.ValidateWebhook(url)
			if err != nil {
				t.Errorf("应该接受有效的 URL: %s, 错误: %v", url, err)
			}
		}
	})

	t.Run("验证无效的时间格式", func(t *testing.T) {
		invalidConfig := &model.Config{
			WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
			ReminderTime:    "25:00", // 无效时间
			ReminderEnabled: true,
			DataPath:        "./data/tasks",
		}

		err := configService.UpdateConfig(invalidConfig)
		if err == nil {
			t.Error("应该拒绝无效的时间格式")
		}
	})
}

// TestReminderServiceIntegration 测试提醒服务集成
func TestReminderServiceIntegration(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, "tasks")
	configPath := filepath.Join(tempDir, "config.json")

	// 初始化仓库和服务
	taskRepo := repository.NewFileTaskRepository(dataPath)
	taskService := service.NewTaskService(taskRepo, dataPath)
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := service.NewConfigService(configRepo)
	reminderService := service.NewReminderService(configService, taskService)

	t.Run("启动和停止提醒服务", func(t *testing.T) {
		// 配置提醒服务
		config := &model.Config{
			WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
			ReminderTime:    "10:00",
			ReminderEnabled: true,
			DataPath:        dataPath,
		}

		err := configService.UpdateConfig(config)
		if err != nil {
			t.Fatalf("更新配置失败: %v", err)
		}

		// 启动提醒服务
		err = reminderService.Start()
		if err != nil {
			t.Fatalf("启动提醒服务失败: %v", err)
		}

		// 等待一小段时间确保服务运行
		time.Sleep(100 * time.Millisecond)

		// 停止提醒服务
		reminderService.Stop()

		// 再次停止应该不会出错
		reminderService.Stop()
	})

	t.Run("禁用提醒时不启动服务", func(t *testing.T) {
		// 配置禁用提醒
		config := &model.Config{
			WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
			ReminderTime:    "10:00",
			ReminderEnabled: false,
			DataPath:        dataPath,
		}

		err := configService.UpdateConfig(config)
		if err != nil {
			t.Fatalf("更新配置失败: %v", err)
		}

		// 启动提醒服务（应该不会真正启动）
		err = reminderService.Start()
		if err != nil {
			t.Fatalf("启动提醒服务失败: %v", err)
		}

		// 停止服务
		reminderService.Stop()
	})

	t.Run("检查今天是否有任务", func(t *testing.T) {
		// 初始状态：今天没有任务
		hasTask, err := taskService.HasTodayTask()
		if err != nil {
			t.Fatalf("检查今天任务失败: %v", err)
		}

		if hasTask {
			t.Error("今天不应该有任务")
		}

		// 创建今天的任务
		today := time.Now().Truncate(24 * time.Hour)
		err = taskService.SaveTask(today, "今天的任务")
		if err != nil {
			t.Fatalf("保存今天的任务失败: %v", err)
		}

		// 再次检查
		hasTask, err = taskService.HasTodayTask()
		if err != nil {
			t.Fatalf("检查今天任务失败: %v", err)
		}

		if !hasTask {
			t.Error("今天应该有任务")
		}
	})
}

// TestErrorHandlingAndLogging 测试错误处理和日志记录
func TestErrorHandlingAndLogging(t *testing.T) {
	t.Run("处理不存在的任务文件", func(t *testing.T) {
		tempDir := t.TempDir()
		dataPath := filepath.Join(tempDir, "tasks")

		taskRepo := repository.NewFileTaskRepository(dataPath)
		taskService := service.NewTaskService(taskRepo, dataPath)

		// 尝试获取不存在的任务
		date := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)
		task, err := taskService.GetTask(date)

		// 不应该返回错误，而是返回 nil
		if err != nil {
			t.Fatalf("获取不存在的任务不应返回错误: %v", err)
		}

		if task != nil {
			t.Error("不存在的任务应返回 nil")
		}
	})

	t.Run("处理无效的JSON文件", func(t *testing.T) {
		tempDir := t.TempDir()
		dataPath := filepath.Join(tempDir, "tasks")

		// 创建数据目录
		err := os.MkdirAll(dataPath, 0755)
		if err != nil {
			t.Fatalf("创建数据目录失败: %v", err)
		}

		// 写入无效的JSON文件
		invalidFile := filepath.Join(dataPath, "2025-11-10.json")
		err = os.WriteFile(invalidFile, []byte("invalid json content"), 0644)
		if err != nil {
			t.Fatalf("写入无效JSON文件失败: %v", err)
		}

		taskRepo := repository.NewFileTaskRepository(dataPath)
		taskService := service.NewTaskService(taskRepo, dataPath)

		// 尝试获取任务
		date := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)
		_, err = taskService.GetTask(date)

		// 应该返回错误
		if err == nil {
			t.Error("读取无效JSON文件应返回错误")
		}
	})

	t.Run("处理只读文件系统", func(t *testing.T) {
		// 注意：这个测试在某些系统上可能无法正常工作
		// 因为很难模拟只读文件系统
		t.Skip("跳过只读文件系统测试")
	})

	t.Run("处理配置文件损坏", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")

		// 写入无效的配置文件
		err := os.WriteFile(configPath, []byte("{invalid json}"), 0644)
		if err != nil {
			t.Fatalf("写入无效配置文件失败: %v", err)
		}

		configRepo := repository.NewFileConfigRepository(configPath)
		configService := service.NewConfigService(configRepo)

		// 尝试加载配置
		_, err = configService.GetConfig()

		// 应该返回错误
		if err == nil {
			t.Error("读取无效配置文件应返回错误")
		}
	})
}

// TestDataPersistence 测试数据持久化
func TestDataPersistence(t *testing.T) {
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, "tasks")

	t.Run("任务数据持久化", func(t *testing.T) {
		// 第一次：创建并保存任务
		{
			taskRepo := repository.NewFileTaskRepository(dataPath)
			taskService := service.NewTaskService(taskRepo, dataPath)

			date := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)
			content := "# 持久化测试\n\n这是一个测试任务"

			err := taskService.SaveTask(date, content)
			if err != nil {
				t.Fatalf("保存任务失败: %v", err)
			}
		}

		// 第二次：重新创建服务并读取任务
		{
			taskRepo := repository.NewFileTaskRepository(dataPath)
			taskService := service.NewTaskService(taskRepo, dataPath)

			date := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)
			task, err := taskService.GetTask(date)
			if err != nil {
				t.Fatalf("读取任务失败: %v", err)
			}

			if task == nil {
				t.Fatal("任务不应为空")
			}

			expectedContent := "# 持久化测试\n\n这是一个测试任务"
			if task.Content != expectedContent {
				t.Errorf("任务内容不匹配，期望: %s, 实际: %s", expectedContent, task.Content)
			}
		}
	})

	t.Run("配置数据持久化", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "config.json")

		// 第一次：创建并保存配置
		{
			configRepo := repository.NewFileConfigRepository(configPath)
			configService := service.NewConfigService(configRepo)

			config := &model.Config{
				WebhookURL:      "https://example.com/webhook",
				ReminderTime:    "14:30",
				ReminderEnabled: true,
				DataPath:        dataPath,
			}

			err := configService.UpdateConfig(config)
			if err != nil {
				t.Fatalf("保存配置失败: %v", err)
			}
		}

		// 第二次：重新创建服务并读取配置
		{
			configRepo := repository.NewFileConfigRepository(configPath)
			configService := service.NewConfigService(configRepo)

			config, err := configService.GetConfig()
			if err != nil {
				t.Fatalf("读取配置失败: %v", err)
			}

			if config.WebhookURL != "https://example.com/webhook" {
				t.Error("Webhook URL 未正确持久化")
			}

			if config.ReminderTime != "14:30" {
				t.Error("提醒时间未正确持久化")
			}

			if !config.ReminderEnabled {
				t.Error("提醒启用状态未正确持久化")
			}
		}
	})
}
