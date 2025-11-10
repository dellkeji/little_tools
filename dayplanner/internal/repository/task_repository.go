package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/util"
)

// TaskRepository 定义任务数据访问接口
type TaskRepository interface {
	// GetByDate 获取指定日期的任务
	GetByDate(date time.Time) (*model.Task, error)

	// Save 保存任务
	Save(task *model.Task) error

	// GetTaskDates 获取日期范围内有任务的日期列表
	GetTaskDates(startDate, endDate time.Time) ([]time.Time, error)

	// HasTask 检查指定日期是否有任务
	HasTask(date time.Time) (bool, error)
}

// FileTaskRepository 基于文件系统的任务仓库实现
type FileTaskRepository struct {
	dataPath string
}

// NewFileTaskRepository 创建新的文件任务仓库
func NewFileTaskRepository(dataPath string) *FileTaskRepository {
	return &FileTaskRepository{
		dataPath: dataPath,
	}
}

// GetByDate 获取指定日期的任务
func (r *FileTaskRepository) GetByDate(date time.Time) (*model.Task, error) {
	filePath := r.getTaskFilePath(date)
	util.Debug("读取任务文件: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			util.Debug("任务文件不存在: %s", filePath)
			return nil, nil // 文件不存在返回 nil，不是错误
		}
		util.Error("读取任务文件失败: %s, 错误: %v", filePath, err)
		return nil, fmt.Errorf("读取任务文件失败: %w", err)
	}

	var task model.Task
	if err := json.Unmarshal(data, &task); err != nil {
		util.Error("解析任务数据失败: %s, 错误: %v", filePath, err)
		return nil, fmt.Errorf("解析任务数据失败: %w", err)
	}

	util.Info("成功读取任务: %s", date.Format("2006-01-02"))
	return &task, nil
}

// Save 保存任务
func (r *FileTaskRepository) Save(task *model.Task) error {
	// 确保数据目录存在
	if err := os.MkdirAll(r.dataPath, 0755); err != nil {
		util.Error("创建数据目录失败: %s, 错误: %v", r.dataPath, err)
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	filePath := r.getTaskFilePath(task.Date)
	util.Debug("保存任务文件: %s", filePath)

	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		util.Error("序列化任务数据失败: 错误: %v", err)
		return fmt.Errorf("序列化任务数据失败: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		util.Error("写入任务文件失败: %s, 错误: %v", filePath, err)
		return fmt.Errorf("写入任务文件失败: %w", err)
	}

	util.Info("成功保存任务: %s", task.Date.Format("2006-01-02"))
	return nil
}

// GetTaskDates 获取日期范围内有任务的日期列表
func (r *FileTaskRepository) GetTaskDates(startDate, endDate time.Time) ([]time.Time, error) {
	var dates []time.Time
	util.Debug("查询任务日期范围: %s 到 %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// 确保数据目录存在
	if _, err := os.Stat(r.dataPath); os.IsNotExist(err) {
		util.Debug("数据目录不存在: %s", r.dataPath)
		return dates, nil // 目录不存在返回空列表
	}

	entries, err := os.ReadDir(r.dataPath)
	if err != nil {
		util.Error("读取数据目录失败: %s, 错误: %v", r.dataPath, err)
		return nil, fmt.Errorf("读取数据目录失败: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 解析文件名格式: YYYY-MM-DD.json
		name := entry.Name()
		if len(name) < 15 || filepath.Ext(name) != ".json" {
			continue
		}

		dateStr := name[:10] // 提取 YYYY-MM-DD 部分
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // 跳过无效的文件名
		}

		// 检查日期是否在范围内
		if (date.Equal(startDate) || date.After(startDate)) &&
			(date.Equal(endDate) || date.Before(endDate)) {
			dates = append(dates, date)
		}
	}

	return dates, nil
}

// HasTask 检查指定日期是否有任务
func (r *FileTaskRepository) HasTask(date time.Time) (bool, error) {
	filePath := r.getTaskFilePath(date)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("检查任务文件失败: %w", err)
	}
	return true, nil
}

// getTaskFilePath 获取任务文件路径
func (r *FileTaskRepository) getTaskFilePath(date time.Time) string {
	fileName := date.Format("2006-01-02") + ".json"
	return filepath.Join(r.dataPath, fileName)
}
