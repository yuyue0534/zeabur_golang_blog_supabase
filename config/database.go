package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	// 优先使用 SUPABASE_DB_URL，其次 DATABASE_URL（兼容你现在的配置）
	dsn := strings.TrimSpace(os.Getenv("SUPABASE_DB_URL"))
	if dsn == "" {
		dsn = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	}
	if dsn == "" {
		return fmt.Errorf("SUPABASE_DB_URL / DATABASE_URL 未设置")
	}

	// 如果是 URL 形式，确保带 sslmode=require（Supabase 常见要求）
	dsn = ensureSSLModeRequire(dsn)

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,

		// ✅ 关键：避免 pgx 隐式 prepared statements（对 Supabase pooler / pgbouncer 更友好）
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 可选：设置连接池参数（避免部署环境连接不稳定）
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Supabase PostgreSQL 连接成功")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func ensureSSLModeRequire(dsn string) string {
	// gorm postgres DSN 既支持 key=val，也支持 postgres:// URL
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return dsn // 解析失败就原样返回，让上层报错更直观
		}
		q := u.Query()
		if q.Get("sslmode") == "" {
			q.Set("sslmode", "require")
			u.RawQuery = q.Encode()
		}
		return u.String()
	}

	// key=val 形式：如果没写 sslmode，补上
	if !strings.Contains(dsn, "sslmode=") {
		return dsn + " sslmode=require"
	}
	return dsn
}
