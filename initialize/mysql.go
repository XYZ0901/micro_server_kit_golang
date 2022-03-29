package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func sqlInit() {
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		Cfg.MysqlConfig.User, Cfg.MysqlConfig.Pwd, Cfg.MysqlConfig.Host,
		Cfg.MysqlConfig.Port, Cfg.MysqlConfig.DbName, Cfg.MysqlConfig.Arguments)

	var err error
	MysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false, // 是否跳过事务
	})
	if err != nil {
		Logger.Sugar().Fatalf("Failed to connect mysql: %s", err.Error())
	}
	sqlDB, err := MysqlDb.DB()
	if err != nil {
		Logger.Sugar().Fatalf("Failed to DB mysql: %s", err.Error())
	}
	// TODO: these parameters can write into config
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetConnMaxIdleTime(0)
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetConnMaxLifetime(0)
}
