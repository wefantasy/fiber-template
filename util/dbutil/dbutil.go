package dbutil

import (
	"app/conf"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Quoter 定义了SQL标识符的引用行为接口
type Quoter interface {
	Quote(identifier string) string
}

// DoubleQuoteQuoter 使用双引号 (ANSI SQL, PostgreSQL)
type DoubleQuoteQuoter struct{}

func (q DoubleQuoteQuoter) Quote(s string) string { return `"` + s + `"` }

// BacktickQuoter 使用反引号 (MySQL, MariaDB)
type BacktickQuoter struct{}

func (q BacktickQuoter) Quote(s string) string { return "`" + s + "`" }

// BracketQuoter 使用方括号 (SQL Server)
type BracketQuoter struct{}

func (q BracketQuoter) Quote(s string) string { return "[" + s + "]" }

// NoOpQuoter 不进行任何引用 (用于不需要或自定义处理的场景)
type NoOpQuoter struct{}

func (q NoOpQuoter) Quote(s string) string { return s }

// columnInfo 存储从结构体字段中提取的关键信息
type columnInfo struct {
	Name  string        // db tag 的值
	Value reflect.Value // 字段的 reflect.Value
	IsPK  bool          // tag中是否包含 "pk"
}

// structCache 用于缓存已解析的结构体信息，避免重复反射
var structCache = &sync.Map{}

// parseStruct 解析结构体，提取字段信息并缓存结果
func parseStruct(v reflect.Value) []columnInfo {
	t := v.Type()
	// 检查类型缓存是否存在，如果存在则直接返回缓存的结果
	if cached, ok := structCache.Load(t); ok {
		cachedCols := cached.([]columnInfo)
		cols := make([]columnInfo, len(cachedCols))
		typeInfo := cached.([]columnInfo)
		for i := 0; i < len(typeInfo); i++ {
			cols[i] = columnInfo{
				Name:  typeInfo[i].Name,
				Value: v.Field(i),
				IsPK:  typeInfo[i].IsPK,
			}
		}
		return cols
	}

	// 首次解析该类型，初始化类型信息
	var typeInfo []columnInfo
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}
		// 解析db标签内容（支持"column,pk"格式）
		tagParts := strings.Split(dbTag, ",")
		colName := tagParts[0]
		isPK := false
		// 检查是否包含主键标记
		if len(tagParts) > 1 && tagParts[1] == "pk" {
			isPK = true
		}
		typeInfo = append(typeInfo, columnInfo{
			Name: colName,
			IsPK: isPK,
		})
	}
	structCache.Store(t, typeInfo)

	// 生成包含当前实例值的列信息
	cols := make([]columnInfo, len(typeInfo))
	for i, info := range typeInfo {
		cols[i] = columnInfo{
			Name:  info.Name,
			IsPK:  info.IsPK,
			Value: v.Field(i),
		}
	}
	return cols
}

// deReference 辅助函数，获取指针指向的实际 Value
func deReference(o any) reflect.Value {
	v := reflect.ValueOf(o)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// Builder 是用于链式生成 SQL 片段的构建器
type Builder struct {
	cols        []columnInfo
	prefix      string
	filters     []func(c columnInfo) bool
	customWhere []string
	orderBy     string
	limit       string
	quoter      Quoter
}

// NewBuilder 创建一个新的构建器实例
func NewBuilder(o any) *Builder {
	b := &Builder{
		quoter: NoOpQuoter{},
	}
	if o == nil {
		return b
	}
	v := deReference(o)
	if v.Kind() != reflect.Struct {
		panic("dbutil: NewBuilder expects a struct or a pointer to a struct")
	}
	b.cols = parseStruct(v)
	if strings.Contains(conf.DB.Type, "mysql") {
		b = b.WithQuoter(BacktickQuoter{})
	}
	return b
}

// WithQuoter 设置一个自定义的 Quoter
func (b *Builder) WithQuoter(q Quoter) *Builder {
	if q != nil {
		b.quoter = q
	}
	return b
}

// WithPrefix 为所有列名添加前缀 (e.g., "user.")
func (b *Builder) WithPrefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

// OnlyNonZero 添加一个过滤器，只保留值非零的字段
func (b *Builder) OnlyNonZero() *Builder {
	b.filters = append(b.filters, func(c columnInfo) bool {
		return c.Value.IsValid() && !c.Value.IsZero()
	})
	return b
}

// ExcludePK 添加一个过滤器，排除主键字段 (常用于 UPDATE)
func (b *Builder) ExcludePK() *Builder {
	b.filters = append(b.filters, func(c columnInfo) bool {
		return !c.IsPK
	})
	return b
}

// WithCustomWhere 添加用户自定义的WHERE条件
// e.g., WithCustomWhere("age > :min_age", "name LIKE :pattern")
func (b *Builder) WithCustomWhere(clauses ...string) *Builder {
	b.customWhere = append(b.customWhere, clauses...)
	return b
}

// WithOrderBy 添加ORDER BY子句
// e.g., WithOrderBy("created_at DESC")
func (b *Builder) WithOrderBy(orderBy string) *Builder {
	b.orderBy = orderBy
	return b
}

// WithLimit 添加LIMIT子句
func (b *Builder) WithLimit(limit int) *Builder {
	b.limit = fmt.Sprintf("LIMIT %d", limit)
	return b
}

// WithLimitOffset 添加LIMIT和OFFSET子句
func (b *Builder) WithLimitOffset(limit, offset int) *Builder {
	b.limit = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	return b
}

// applyFilters 执行所有已注册的过滤器
func (b *Builder) applyFilters() []columnInfo {
	if b.cols == nil {
		return nil
	}
	var filtered []columnInfo
	for _, c := range b.cols {
		include := true
		for _, f := range b.filters {
			if !f(c) {
				include = false
				break
			}
		}
		if include {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// BuildColumns 生成列名列表
// 用法: builder.BuildColumns(",") -> "id,name,email"
func (b *Builder) BuildColumns(separator string) string {
	cols := b.applyFilters()
	var names []string
	for _, c := range cols {
		names = append(names, b.prefix+b.quoter.Quote(c.Name))
	}
	return strings.Join(names, separator)
}

// BuildColumnsWithAlias 生成带别名的列名列表，用于 SELECT 查询
// 用法: builder.WithPrefix("u.").BuildColumnsWithAlias(",") -> "u.id AS "id", u.name AS "name""
func (b *Builder) BuildColumnsWithAlias(separator string) string {
	if b.cols == nil {
		return ""
	}
	var names []string
	for _, c := range b.cols {
		aliased := fmt.Sprintf(`%s%s AS "%s"`, b.prefix, b.quoter.Quote(c.Name), c.Name)
		names = append(names, aliased)
	}
	return strings.Join(names, separator)
}

// BuildNamedPlaceholders 生成用于 sqlx 的命名占位符
// 用法: builder.BuildNamedPlaceholders(",") -> ":id,:name,:email"
func (b *Builder) BuildNamedPlaceholders(separator string) string {
	cols := b.applyFilters()
	var placeholders []string
	for _, c := range cols {
		placeholders = append(placeholders, ":"+c.Name)
	}
	return strings.Join(placeholders, separator)
}

// BuildSetClauses 生成用于 UPDATE 的 SET 子句
// 用法: builder.BuildSetClauses(",") -> "name=:name,email=:email"
func (b *Builder) BuildSetClauses(separator string) string {
	cols := b.applyFilters()
	var clauses []string
	for _, c := range cols {
		clauses = append(clauses, fmt.Sprintf("%s%s=:%s", b.prefix, b.quoter.Quote(c.Name), c.Name))
	}
	return strings.Join(clauses, separator)
}

// BuildWhereClauses 生成用于 WHERE 的条件子句, 智能合并自动生成和自定义的条件
// 用法: builder.BuildWhereClauses(" AND ") -> "id=:id AND name=:name"
func (b *Builder) BuildWhereClauses(separator string) string {
	autoClauses := []string{}
	filteredCols := b.applyFilters()
	for _, c := range filteredCols {
		autoClauses = append(autoClauses, fmt.Sprintf("%s%s=:%s", b.prefix, b.quoter.Quote(c.Name), c.Name))
	}

	allClauses := append(autoClauses, b.customWhere...)

	if len(allClauses) == 0 {
		return ""
	}

	return strings.Join(allClauses, separator)
}

// BuildSelectQuery 组装一个完整的 SELECT 查询语句
func (b *Builder) BuildSelectQuery(tableName string) string {
	if tableName == "" {
		return ""
	}

	var sb strings.Builder

	// 1. SELECT columns
	sb.WriteString("SELECT ")
	sb.WriteString(b.BuildColumnsWithAlias(", "))

	// 2. FROM table
	sb.WriteString(" FROM ")
	sb.WriteString(b.quoter.Quote(tableName))

	// 3. WHERE clause
	whereClause := b.BuildWhereClauses(" AND ")
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	// 4. ORDER BY clause
	if b.orderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(b.orderBy)
	}

	// 5. LIMIT / OFFSET clause
	if b.limit != "" {
		sb.WriteString(" ") // LIMIT前通常有个空格
		sb.WriteString(b.limit)
	}

	return sb.String()
}
