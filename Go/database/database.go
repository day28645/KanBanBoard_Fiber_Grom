package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n======================================\n", sql)
}

func ConnectDB() {
	dsn := "root:root1@tcp(127.0.0.1:3306)/fibergorm?parseTime=true"
	dial := mysql.Open(dsn)

	var err error
	DB, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	//DB.AutoMigrate(&models.User{}, &models.Login{}, &models.Board{}, &models.BoardMember{}, &models.ColumnBoard{}, &models.Task{}, &models.TaskAssignee{})
}
