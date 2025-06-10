package dbutil

import (
	"app/model"
	"testing"
)

func TestDbutil(t *testing.T) {
	user := model.User{}

	// 生成 UPDATE SET 子句，只包含非空字段，且排除主键
	// "name=:name"
	updateSet := NewBuilder(&user).OnlyNonZero().ExcludePK().BuildSetClauses(", ")
	t.Log(updateSet)

	// 生成 SELECT 查询列，带前缀和别名
	// "u.id AS "id", u.name AS "name", u.email AS "email""
	selectCols := NewBuilder(&model.User{}).WithPrefix("u.").BuildColumnsWithAlias(", ")
	t.Log(selectCols)

	// 生成 WHERE 条件，所有非空字段用 AND 连接
	// "id=:id AND name=:name"
	whereConditions := NewBuilder(&user).OnlyNonZero().BuildWhereClauses(" AND ")
	t.Log(whereConditions)

	// 1. 复杂查询：结合了结构体字段、自定义WHERE、排序和分页
	query := NewBuilder(&user).
		OnlyNonZero(). // status 不为空时，会自动加入 WHERE status=:status
		WithPrefix("u.").
		WithCustomWhere("u.registration_date > :start_date"). // 自定义WHERE条件
		WithCustomWhere("u.is_active = :is_active").
		WithOrderBy("u.created_at DESC"). // 排序
		WithLimitOffset(10, 20).          // 分页 (取10条，偏移20)
		BuildSelectQuery("users u")       // 从 users 表 (别名u) 构建查询
	t.Log(query)

	// 2. 仅使用自定义条件
	// 有时我们可能不需要根据结构体生成条件，只用自定义的
	query2 := NewBuilder(nil). // 传入 nil 来跳过结构体解析
					WithCustomWhere("is_deleted = :is_deleted").
					WithOrderBy("id ASC").
					BuildSelectQuery("products")
	t.Log(query2)
}
