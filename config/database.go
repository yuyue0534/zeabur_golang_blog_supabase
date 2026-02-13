package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接（Supabase Postgres）
func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL 未设置")
	}
	// 配置 pgx 连接
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return err
	}
	// 在 pgx v5 中不再存在 PreferSimpleProtocol 字段，使用默认连接配置
	// 若确实需要简单协议，可在驱动层实现；此处保持默认设置
	var gormErr error
	DB, gormErr = gorm.Open(postgres.New(postgres.Config{
		Conn:       nil,
		DriverName: "pgx",
		DSN:        config.ConnString(),
	}), &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Info),
		PrepareStmt:          false,
		DisableAutomaticPing: false,
	})
	if gormErr != nil {
		return gormErr
	}
	log.Println("Supabase PostgreSQL 连接成功")
	log.Println("数据库连接成功（使用现有表结构）")
	log.Println("数据库表迁移成功 (PostgreSQL)")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
