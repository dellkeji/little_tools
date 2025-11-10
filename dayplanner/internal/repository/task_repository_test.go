package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"daily-report-tool/internal/model"
)

func TestFileTaskRepository_SaveAndGetByDate(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "task_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewFileTaskRepository(tempDir)

	// 创建测试任务
	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.UTC)
	task := &model.Task{
		Date:      testDate,
		Content:   "# 测试任务\n\n- 任务1\n- 任务2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 测试保存
	err = repo.Save(task)
	if err != nil {
		t.Fatalf("保存任务失败: %v", err)
	}

	// 测试读取
	loadedTask, err := repo.GetByDate(testDate)
	if err != nil {
		t.Fatalf("读取任务失败: %v", err)
	}

	if loadedTask == nil {
		t.Fatal("读取的任务为 nil")
	}

	if loadedTask.Content != task.Content {
		t.Errorf("任务内容不匹配: 期望 %s, 实际 %s", task.Content, loadedTask.Content)
	}
}

func TestFileTaskRepository_GetByDate_NotExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "task_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewFileTaskRepository(tempDir)

	// 读取不存在的任务
	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.UTC)
	task, err := repo.GetByDate(testDate)
	if err != nil {
		t.Fatalf("读取不存在的任务应该返回 nil 而不是错误: %v", err)
	}

	if task != nil {
		t.Error("不存在的任务应该返回 nil")
	}
}

func TestFileTaskRepository_HasTask(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "task_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewFileTaskRepository(tempDir)

	testDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.UTC)

	// 检查不存在的任务
	exists, err := repo.HasTask(testDate)
	if err != nil {
		t.Fatalf("检查任务失败: %v", err)
	}
	if exists {
		t.Error("不存在的任务应该返回 false")
	}

	// 保存任务
	task := &model.Task{
		Date:      testDate,
		Content:   "测试内容",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Save(task)
	if err != nil {
		t.Fatalf("保存任务失败: %v", err)
	}

	// 再次检查
	exists, err = repo.HasTask(testDate)
	if err != nil {
		t.Fatalf("检查任务失败: %v", err)
	}
	if !exists {
		t.Error("存在的任务应该返回 true")
	}
}

func TestFileTaskRepository_GetTaskDates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "task_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewFileTaskRepository(tempDir)

	// 创建多个任务
	dates := []time.Time{
		time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 11, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		task := &model.Task{
			Date:      date,
			Content:   "测试内容",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = repo.Save(task)
		if err != nil {
			t.Fatalf("保存任务失败: %v", err)
		}
	}

	// 查询整个月份的任务
	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 11, 30, 23, 59, 59, 0, time.UTC)

	taskDates, err := repo.GetTaskDates(startDate, endDate)
	if err != nil {
		t.Fatalf("获取任务日期失败: %v", err)
	}

	if len(taskDates) != len(dates) {
		t.Errorf("任务日期数量不匹配: 期望 %d, 实际 %d", len(dates), len(taskDates))
	}
}

func TestFileTaskRepository_GetTaskDates_EmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "task_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 使用不存在的目录
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	repo := NewFileTaskRepository(nonExistentDir)

	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)

	taskDates, err := repo.GetTaskDates(startDate, endDate)
	if err != nil {
		t.Fatalf("空目录应该返回空列表而不是错误: %v", err)
	}

	if len(taskDates) != 0 {
		t.Errorf("空目录应该返回空列表: 实际 %d", len(taskDates))
	}
}
