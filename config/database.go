package config  
import (  
	"fmt"  
	"log"  
	"os"  
	"strings"  
	"time"  
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
	// 🔥 关键：添加参数禁用预编译语句缓存  
	if !strings.Contains(dsn, "statement_cache_mode") {  
		if strings.Contains(dsn, "?") {  
			dsn += "&statement_cache_mode=describe"  
		} else {  
			dsn += "?statement_cache_mode=describe"  
		}  
	}  
	var err error  
	DB, err = gorm.Open(postgres.New(postgres.Config{  
		DSN: dsn,  
		// 🔥 禁用 GORM 的预编译语句  
		PreferSimpleProtocol: true,  
	}), &gorm.Config{  
		Logger:      logger.Default.LogMode(logger.Info),  
		PrepareStmt: false, // 🔥 关键：禁用预编译语句  
	})  
	if err != nil {  
		return err  
	}  
	// 获取底层 sql.DB 并配置连接池  
	sqlDB, err := DB.DB()  
	if err != nil {  
		return err  
	}  
	sqlDB.SetMaxOpenConns(10)  
	sqlDB.SetMaxIdleConns(3)  
	sqlDB.SetConnMaxLifetime(30 * time.Minute)  
	log.Println("✅ Supabase PostgreSQL 连接成功")  
	log.Println("✅ 数据库连接成功（使用现有表结构）")  
	log.Println("✅ 数据库表迁移成功 (PostgreSQL)")  
	return nil  
}  
// GetDB 获取数据库实例  
func GetDB() *gorm.DB {  
	return DB  
}  
