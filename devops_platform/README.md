# DevOps 管理平台

一个前后端分离的DevOps管理平台，支持Jenkins节点管理、MySQL操作和Redis管理。

## 技术栈

### 后端
- Golang 1.21+
- Gin Web框架
- GORM (数据库ORM)
- Jenkins Go客户端
- go-redis
- go-sql-driver/mysql

### 前端
- React 18
- TypeScript
- Ant Design
- Axios

## 项目结构

```
├── backend/          # 后端代码
│   ├── cmd/         # 应用入口
│   ├── internal/    # 内部代码
│   ├── pkg/         # 公共包
│   └── config/      # 配置文件
├── frontend/        # 前端代码
└── README.md
```

## 快速开始

### 后端
```bash
cd backend
go mod download
go run cmd/main.go
```

### 前端
```bash
cd frontend
npm install
npm start
```
