package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接（Supabase Postgres）
func InitDB() error {

	// 从环境变量读取
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("============>DATABASE_URL:", dsn)

	if dsn == "" {
		return fmt.Errorf("DATABASE_URL 未设置")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Info),
		PrepareStmt:          false, // 🔥 关键
		DisableAutomaticPing: false,
	})

	if err != nil {
		return err
	}

	log.Println("Supabase PostgreSQL 连接成功")

	log.Println("数据库连接成功（使用现有表结构）")

	if err != nil {
		return err
	}

	log.Println("数据库表迁移成功 (PostgreSQL)")

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
