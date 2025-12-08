package dao

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/models"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func ConnectToMySQL(debug bool, autoMigrate bool) (*gorm.DB, error) {
	//生成连接URL
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		core.Config.Db.Mysql.User,
		core.Config.Db.Mysql.Passwd,
		core.Config.Db.Mysql.Addr,
		core.Config.Db.Mysql.Port,
		core.Config.Db.Mysql.DefaultDB,
	)

	//接入ZapCore日志核心
	gormLogger := zapgorm2.New(core.Logger)
	gormLogger.SetAsDefault()                         // 全局生效，可选
	gormLogger.SlowThreshold = 300 * time.Millisecond // 慢查询阈值

	//连接数据库
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger, //用zap接管日志输出
	})
	if err != nil {
		return nil, err
	}

	//数据库迁移
	if autoMigrate == true {
		err = conn.AutoMigrate(
			&models.Student{},
			&models.Course{},
			&models.Relation{},
		)

		if err != nil {
			return nil, err
		}
	}

	//将连接全局化
	if debug {
		return conn.Debug(), nil
	} else {
		return conn, nil
	}

}
