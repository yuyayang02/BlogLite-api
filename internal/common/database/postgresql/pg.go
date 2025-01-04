package postgresql

import (
	"database/sql"
	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("数据库未初始化或初始化失败")
	}
	return db
}

func PostgresDSN(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func SqliteDNS(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}

func Init(dialector gorm.Dialector) (closeFn func() error) {
	var err error
	db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败，请检查数据库配置和网络连接")
	}

	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal("数据库连接错误，请检查启动配置")
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(50)

	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}

	return sqlDB.Close
}
