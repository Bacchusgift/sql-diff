package web

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/Bacchusgift/sql-diff/internal/ai"
	"github.com/Bacchusgift/sql-diff/internal/config"
	"github.com/Bacchusgift/sql-diff/internal/differ"
	"github.com/Bacchusgift/sql-diff/internal/parser"
	"gopkg.in/yaml.v3"
)

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// ConfigSaveRequest 配置保存请求
type ConfigSaveRequest struct {
	AI     config.AIConfig `json:"ai"`
	SaveTo string          `json:"save_to"` // "file" or "env"
}

// ConfigExportResponse 配置导出响应
type ConfigExportResponse struct {
	Commands []string `json:"commands"`
}

// CompareRequest 比对请求
type CompareRequest struct {
	SourceSQL string `json:"source_sql"`
	TargetSQL string `json:"target_sql"`
	EnableAI  bool   `json:"enable_ai"`
}

// CompareResponse 比对响应
type CompareResponse struct {
	HasChanges bool               `json:"has_changes"`
	Summary    string             `json:"summary"`
	DDLs       []string           `json:"ddls"`
	AIAnalysis *ai.AnalysisResult `json:"ai_analysis,omitempty"`
}

// AIGenerateCreateRequest AI 生成 CREATE TABLE 请求
type AIGenerateCreateRequest struct {
	Description string `json:"description"`
}

// AIGenerateCreateResponse AI 生成 CREATE TABLE 响应
type AIGenerateCreateResponse struct {
	SQL     string `json:"sql"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// AIGenerateAlterRequest AI 生成 ALTER TABLE 请求
type AIGenerateAlterRequest struct {
	CurrentDDL  string `json:"current_ddl"`
	Description string `json:"description"`
}

// AIGenerateAlterResponse AI 生成 ALTER TABLE 响应
type AIGenerateAlterResponse struct {
	SQLs    []string `json:"sqls"`
	Success bool     `json:"success"`
	Error   string   `json:"error"`
}

// ParseCreateRequest SQL解析请求
type ParseCreateRequest struct {
	SQL string `json:"sql"`
}

// ParseCreateResponse SQL解析响应
type ParseCreateResponse struct {
	Success   bool         `json:"success"`
	TableName string       `json:"table_name,omitempty"`
	Columns   []ColumnInfo `json:"columns,omitempty"`
	Error     string       `json:"error,omitempty"`
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Length    string `json:"length,omitempty"`
	Default   string `json:"default,omitempty"`
	Nullable  bool   `json:"nullable"`
	IsPrimary bool   `json:"is_primary"`
	IsUnique  bool   `json:"is_unique"`
	Comment   string `json:"comment,omitempty"`
	AutoInc   bool   `json:"auto_inc"`
	Unsigned  bool   `json:"unsigned"`
}

// handleConfig 处理配置相关请求
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetConfig(w, r)
	case http.MethodPost:
		s.handleSaveConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetConfig 获取配置
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(".sql-diff-config.yaml")
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "加载配置失败", err.Error())
		return
	}

	// 脱敏 API Key
	if cfg.AI.APIKey != "" {
		cfg.AI.APIKey = maskAPIKey(cfg.AI.APIKey)
	}

	s.writeJSON(w, http.StatusOK, cfg)
}

// handleSaveConfig 保存配置
func (s *Server) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var req ConfigSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "解析请求失败", err.Error())
		return
	}

	// 如果 API Key 是脱敏的,则从现有配置中读取
	if strings.HasSuffix(req.AI.APIKey, "***") {
		existingCfg, err := config.LoadConfig(".sql-diff-config.yaml")
		if err == nil {
			req.AI.APIKey = existingCfg.AI.APIKey
		}
	}

	cfg := &config.Config{AI: req.AI}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		s.writeError(w, http.StatusBadRequest, "配置验证失败", err.Error())
		return
	}

	if req.SaveTo == "file" {
		// 保存到文件
		data, err := yaml.Marshal(cfg)
		if err != nil {
			s.writeError(w, http.StatusInternalServerError, "序列化配置失败", err.Error())
			return
		}

		// 尝试保存到当前目录,如果失败则保存到用户目录
		configPath := ".sql-diff-config.yaml"
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			homeDir, _ := os.UserHomeDir()
			configPath = homeDir + "/.sql-diff-config.yaml"
			if err := os.WriteFile(configPath, data, 0644); err != nil {
				s.writeError(w, http.StatusInternalServerError, "保存配置文件失败", err.Error())
				return
			}
		}

		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "配置已保存到 " + configPath,
			"path":    configPath,
		})
	} else {
		// 仅返回环境变量命令,不实际保存
		commands := cfg.SaveToEnv()
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"message":  "请手动执行以下命令设置环境变量",
			"commands": commands,
		})
	}
}

// handleConfigExport 导出环境变量
func (s *Server) handleConfigExport(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(".sql-diff-config.yaml")
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "加载配置失败", err.Error())
		return
	}

	commands := cfg.SaveToEnv()
	s.writeJSON(w, http.StatusOK, ConfigExportResponse{Commands: commands})
}

// handleCompare 处理表结构比对
func (s *Server) handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "解析请求失败", err.Error())
		return
	}

	// 验证输入
	if strings.TrimSpace(req.SourceSQL) == "" || strings.TrimSpace(req.TargetSQL) == "" {
		s.writeError(w, http.StatusBadRequest, "源表和目标表 SQL 不能为空", "")
		return
	}

	// 解析 SQL
	p := parser.NewParser()
	sourceSchema, err := p.Parse(req.SourceSQL)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "解析源表失败", err.Error())
		return
	}

	targetSchema, err := p.Parse(req.TargetSQL)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "解析目标表失败", err.Error())
		return
	}

	// 比对差异
	d := differ.NewDiffer(sourceSchema, targetSchema)
	diff := d.Compare()

	response := CompareResponse{
		HasChanges: diff.HasChanges(),
		Summary:    diff.Summary(),
		DDLs:       diff.GenerateDDL(sourceSchema.Name),
	}

	// AI 分析
	if req.EnableAI {
		cfg, err := config.LoadConfig(".sql-diff-config.yaml")
		if err == nil && cfg.AI.Enabled {
			provider, err := ai.NewProvider(&cfg.AI)
			if err == nil {
				result, err := provider.Analyze(req.SourceSQL, req.TargetSQL, diff.Summary())
				if err == nil {
					response.AIAnalysis = result
				}
			}
		}
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleAIGenerateCreate 处理 AI 生成 CREATE TABLE
func (s *Server) handleAIGenerateCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AIGenerateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "解析请求失败", err.Error())
		return
	}

	if strings.TrimSpace(req.Description) == "" {
		s.writeError(w, http.StatusBadRequest, "描述不能为空", "")
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(".sql-diff-config.yaml")
	if err != nil {
		s.writeJSON(w, http.StatusOK, AIGenerateCreateResponse{
			Success: false,
			Error:   "加载配置失败: " + err.Error(),
		})
		return
	}

	if !cfg.AI.Enabled {
		s.writeJSON(w, http.StatusForbidden, AIGenerateCreateResponse{
			Success: false,
			Error:   "AI 功能未启用，请先配置 AI",
		})
		return
	}

	// 创建 AI Provider
	provider, err := ai.NewProvider(&cfg.AI)
	if err != nil {
		s.writeJSON(w, http.StatusOK, AIGenerateCreateResponse{
			Success: false,
			Error:   "AI 初始化失败: " + err.Error(),
		})
		return
	}

	// 生成 SQL
	sql, err := provider.GenerateCreateTable(req.Description)
	if err != nil {
		s.writeJSON(w, http.StatusOK, AIGenerateCreateResponse{
			Success: false,
			Error:   "生成失败: " + err.Error(),
		})
		return
	}

	s.writeJSON(w, http.StatusOK, AIGenerateCreateResponse{
		Success: true,
		SQL:     sql,
	})
}

// handleAIGenerateAlter 处理 AI 生成 ALTER TABLE
func (s *Server) handleAIGenerateAlter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AIGenerateAlterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "解析请求失败", err.Error())
		return
	}

	if strings.TrimSpace(req.CurrentDDL) == "" || strings.TrimSpace(req.Description) == "" {
		s.writeError(w, http.StatusBadRequest, "表结构和描述不能为空", "")
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(".sql-diff-config.yaml")
	if err != nil {
		s.writeJSON(w, http.StatusOK, AIGenerateAlterResponse{
			Success: false,
			Error:   "加载配置失败: " + err.Error(),
		})
		return
	}

	if !cfg.AI.Enabled {
		s.writeJSON(w, http.StatusForbidden, AIGenerateAlterResponse{
			Success: false,
			Error:   "AI 功能未启用，请先配置 AI",
		})
		return
	}

	// 创建 AI Provider
	provider, err := ai.NewProvider(&cfg.AI)
	if err != nil {
		s.writeJSON(w, http.StatusOK, AIGenerateAlterResponse{
			Success: false,
			Error:   "AI 初始化失败: " + err.Error(),
		})
		return
	}

	// 生成 SQL
	sql, err := provider.GenerateAlterTable(req.CurrentDDL, req.Description)
	if err != nil {
		s.writeJSON(w, http.StatusOK, AIGenerateAlterResponse{
			Success: false,
			Error:   "生成失败: " + err.Error(),
		})
		return
	}

	// 处理多条 SQL 语句
	sqls := strings.Split(sql, "\n")
	var result []string
	for _, stmt := range sqls {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	s.writeJSON(w, http.StatusOK, AIGenerateAlterResponse{
		Success: true,
		SQLs:    result,
	})
}

// handleParseCreate 处理 SQL 解析请求
func (s *Server) handleParseCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ParseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "解析请求失败", err.Error())
		return
	}

	if strings.TrimSpace(req.SQL) == "" {
		s.writeError(w, http.StatusBadRequest, "SQL语句不能为空", "")
		return
	}

	// 解析 SQL
	p := parser.NewParser()
	schema, err := p.Parse(req.SQL)
	if err != nil {
		s.writeJSON(w, http.StatusOK, ParseCreateResponse{
			Success: false,
			Error:   "解析失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	columns := make([]ColumnInfo, 0, len(schema.Columns))
	for _, col := range schema.Columns {
		colInfo := ColumnInfo{
			Name:     col.Name,
			Type:     col.Type,
			Length:   col.Length,
			Default:  col.DefaultValue,
			Nullable: !col.NotNull,
			AutoInc:  col.AutoInc,
			Unsigned: col.Unsigned,
			Comment:  col.Comment,
		}

		// 检查是否为主键
		for _, pk := range schema.PrimaryKeys {
			if pk == col.Name {
				colInfo.IsPrimary = true
				break
			}
		}

		// 检查是否为唯一索引
		for _, idx := range schema.Indexes {
			if idx.Type == "UNIQUE" {
				for _, idxCol := range idx.Columns {
					if idxCol == col.Name {
						colInfo.IsUnique = true
						break
					}
				}
			}
		}

		columns = append(columns, colInfo)
	}

	s.writeJSON(w, http.StatusOK, ParseCreateResponse{
		Success:   true,
		TableName: schema.Name,
		Columns:   columns,
	})
}

// writeJSON 写入 JSON 响应
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError 写入错误响应
func (s *Server) writeError(w http.ResponseWriter, status int, message, details string) {
	s.writeJSON(w, status, ErrorResponse{
		Success: false,
		Error:   message,
		Details: details,
	})
}

// maskAPIKey 脱敏 API Key
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:3] + "-***"
}
