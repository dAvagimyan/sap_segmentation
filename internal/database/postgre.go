package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func PostgresOpenConnection(config Config) (*sqlx.DB, error) {
	conn, errOpen := sqlx.Open(
		`postgres`,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host,
			config.Port,
			config.User,
			config.Pass,
			config.Name,
		))

	if errOpen != nil {
		return nil, errOpen
	}

	conn.SetMaxIdleConns(config.MaxIdleConns)
	conn.SetMaxOpenConns(config.MaxOpenConns)
	conn.SetConnMaxLifetime(config.ConnMaxLife)

	if errPing := conn.Ping(); errPing != nil {
		return nil, errPing
	}

	return conn, nil
}
