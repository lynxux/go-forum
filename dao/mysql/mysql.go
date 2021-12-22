package mysql

import (
	"fmt"
	"web_app/settings"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var dbSqlx *sqlx.DB

func Init(mySqlConfig *settings.MySqlConfig) (err error) {
	//dsn := "root:123456@tcp(127.0.0.1:3306)/sql_demo?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mySqlConfig.User,
		mySqlConfig.Password,
		mySqlConfig.Host,
		mySqlConfig.Port,
		mySqlConfig.DbName,
	)
	dbSqlx, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		//fmt.Println("Connect DB failed!", err)
		zap.L().Error("connect DB failed", zap.Error(err))
		return err
	}
	dbSqlx.SetMaxOpenConns(mySqlConfig.MaxOpenConns)
	dbSqlx.SetMaxIdleConns(mySqlConfig.MaxIdleConns)
	return
}

func Close() {
	_ = dbSqlx.Close()
}
