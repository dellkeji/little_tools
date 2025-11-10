package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"daily-report-tool/internal/repository"
)

func TestTaskService_GetTask(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建仓库和服务
	taskRepo := repository.NewFileTaskRepository(tempDir)
	taskService := NewTaskService(taskRepo, tempDir)

	// 测试日期
	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)

	// 测试获取不存在的任务
	task, err := taskService.GetTask(testDate)
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}
	if task != nil {
		t.Errorf("期望任务为 nil，实际得到: %v", task)
	}

	// 保存一个任务
	content := "# 测试任务\n\n- 任务1\n- 任务2"
	if err := taskService.SaveTask(testDate, content); err != nil {
		t.Fatalf("保存任务失败: %v", err)
	}

	// 再次获取任务
	task, err = taskService.GetTask(testDate)
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}
	if task == nil {
		t.Fatal("期望获取到任务，但得到 nil")
	}
	if task.Content != content {
		t.Errorf("任务内容不匹配，期望: %s，实际: %s", content, task.Content)
	}
}

func TestTaskService_SaveTask(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建仓库和服务
	taskRepo := repository.NewFileTaskRepository(tempDir)
	taskService := NewTaskService(taskRepo, tempDir)

	// 测试日期
	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)
	content := "# 测试任务"

	// 保存新任务
	if err := taskService.SaveTask(testDate, content); err != nil {
		t.Fatalf("保存任务失败: %v", err)
	}

	// 验证文件是否创建
	fileName := testDate.Format("2006-01-02") + ".json"
	filePath := filepath.Join(tempDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("任务文件未创建: %s", filePath)
	}

	// 更新任务
	updatedContent := "# 更新的任务"
	if err := taskService.SaveTask(testDate, updatedContent); err != nil {
		t.Fatalf("更新任务失败: %v", err)
	}

	// 验证任务内容已更新
	task, err := taskService.GetTask(testDate)
	if err != nil {
		t.Fatalf("获取任务失败: %v", err)
	}
	if task.Content != updatedContent {
		t.Errorf("任务内容未更新，期望: %s，实际: %s", updatedContent, task.Content)
	}
}

func TestTaskService_GetMonthTaskDates(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建仓库和服务
	taskRepo := repository.NewFileTaskRepository(tempDir)
	taskService := NewTaskService(taskRepo, tempDir)

	// 保存多个任务
	dates := []time.Time{
		time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local),
		time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local),
		time.Date(2025, 11, 20, 0, 0, 0, 0, time.Local),
		time.Date(2025, 12, 5, 0, 0, 0, 0, time.Local), // 不同月份
	}

	for _, date := range dates {
		if err := taskService.SaveTask(date, "测试内容"); err != nil {
			t.Fatalf("保存任务失败: %v", err)
		}
	}

	// 获取 11 月的任务日期
	taskDates, err := taskService.GetMonthTaskDates(2025, time.November)
	if err != nil {
		t.Fatalf("获取月份任务日期失败: %v", err)
	}

	// 验证结果
	if len(taskDates) != 3 {
		t.Errorf("期望 3 个任务日期，实际得到: %d", len(taskDates))
	}

	// 验证日期是否正确
	expectedDates := map[string]bool{
		"2025-11-01": true,
		"2025-11-10": true,
		"2025-11-20": true,
	}

	for _, date := range taskDates {
		dateStr := date.Format("2006-01-02")
		if !expectedDates[dateStr] {
			t.Errorf("意外的任务日期: %s", dateStr)
		}
	}
}

func TestTaskService_HasTodayTask(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建仓库和服务
	taskRepo := repository.NewFileTaskRepository(tempDir)
	taskService := NewTaskService(taskRepo, tempDir)

	// 测试今天没有任务的情况
	hasTask, err := taskService.HasTodayTask()
	if err != nil {
		t.Fatalf("检查今天任务失败: %v", err)
	}
	if hasTask {
		t.Errorf("期望今天没有任务，但返回有任务")
	}

	// 保存今天的任务
	today := time.Now().Truncate(24 * time.Hour)
	if err := taskService.SaveTask(today, "今天的任务"); err != nil {
		t.Fatalf("保存今天任务失败: %v", err)
	}

	// 再次检查
	hasTask, err = taskService.HasTodayTask()
	if err != nil {
		t.Fatalf("检查今天任务失败: %v", err)
	}
	if !hasTask {
		t.Errorf("期望今天有任务，但返回没有任务")
	}
}

func TestTaskService_EnsureDataDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, "data", "tasks")

	// 创建仓库和服务
	taskRepo := repository.NewFileTaskRepository(dataPath)
	taskService := NewTaskService(taskRepo, dataPath)

	// 验证数据目录不存在
	if _, err := os.Stat(dataPath); !os.IsNotExist(err) {
		t.Fatalf("数据目录不应该存在")
	}

	// 调用 SaveTask，应该自动创建目录
	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)
	if err := taskService.SaveTask(testDate, "测试"); err != nil {
		t.Fatalf("保存任务失败: %v", err)
	}

	// 验证数据目录已创建
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("数据目录未创建: %s", dataPath)
	}
}
