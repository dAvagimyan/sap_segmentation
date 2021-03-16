package database

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

func MysqlOpenConnection(config Config) (*sqlx.DB, error) {
	conn, errOpen := sqlx.Open(`mysql`, fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True",
		config.User,
		config.Pass,
		config.Host,
		config.Name,
	))

	if errOpen != nil {
		log.Println("Ошибка подключения к бд mysql: ", config, errOpen)
		return nil, errOpen
	}

	conn.SetMaxIdleConns(config.MaxIdleConns)
	conn.SetMaxOpenConns(config.MaxOpenConns)
	conn.SetConnMaxLifetime(config.ConnMaxLife)

	if errPing := conn.Ping(); errPing != nil {
		log.Println("Ошибка пинга к бд mysql", config, errOpen)
		return nil, errPing
	}

	return conn, nil
}
