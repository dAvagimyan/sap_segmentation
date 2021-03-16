package database

import "time"

type Config struct {
	Host  string
	Port  string
	Name  string
	User  string
	Pass  string
	Debug string

	MaxIdleConns int
	MaxOpenConns int
	ConnMaxLife  time.Duration
}
