package service

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xwb1989/sqlparser"
)

type MySQLService struct {
	db *sql.DB
	mu sync.RWMutex
}

var mysqlService *MySQLService

func GetMySQLService() *MySQLService {
	if mysqlService == nil {
		mysqlService = &MySQLService{}
	}
	return mysqlService
}

func (s *MySQLService) Connect(host, port, username, password, database string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("MySQL连接测试失败: %w", err)
	}

	s.db = db
	return nil
}

func (s *MySQLService) ValidateSQL(query string) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return fmt.Errorf("SQL语句不能为空")
	}

	// 使用sqlparser验证SQL语法
	_, err := sqlparser.Parse(query)
	if err != nil {
		return fmt.Errorf("SQL语法错误: %w", err)
	}

	// 检查危险操作
	upperQuery := strings.ToUpper(query)
	dangerousKeywords := []string{"DROP DATABASE", "DROP TABLE", "TRUNCATE"}
	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperQuery, keyword) {
			return fmt.Errorf("禁止执行危险操作: %s", keyword)
		}
	}

	return nil
}

func (s *MySQLService) ExecuteSQL(query string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return nil, fmt.Errorf("MySQL未连接")
	}

	// 验证SQL
	if err := s.ValidateSQL(query); err != nil {
		return nil, err
	}

	// 判断是查询还是执行
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "SHOW") || strings.HasPrefix(upperQuery, "DESCRIBE") {
		return s.executeQuery(query)
	}

	return s.executeExec(query)
}

func (s *MySQLService) executeQuery(query string) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		result = append(result, row)
	}

	return result, nil
}

func (s *MySQLService) executeExec(query string) (map[string]interface{}, error) {
	result, err := s.db.Exec(query)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId()

	return map[string]interface{}{
		"rowsAffected": rowsAffected,
		"lastInsertId": lastInsertId,
	}, nil
}

func (s *MySQLService) GetDatabases() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return nil, fmt.Errorf("MySQL未连接")
	}

	rows, err := s.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	databases := make([]string, 0)
	for rows.Next() {
		var db string
		if err := rows.Scan(&db); err != nil {
			return nil, err
		}
		databases = append(databases, db)
	}

	return databases, nil
}

func (s *MySQLService) GetTables(database string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return nil, fmt.Errorf("MySQL未连接")
	}

	query := fmt.Sprintf("SHOW TABLES FROM `%s`", database)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}
