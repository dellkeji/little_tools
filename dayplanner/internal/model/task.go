package model

import "time"

// Task 表示用户为特定日期创建的工作任务记录
type Task struct {
	Date      time.Time `json:"date"`       // 任务日期
	Content   string    `json:"content"`    // Markdown 内容
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
