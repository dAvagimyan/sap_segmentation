package main

import (
	"github.com/kelseyhightower/envconfig"
	"os"
	"sap_segmentation/internal/database"
	"sap_segmentation/pkg/logger"
	"sap_segmentation/segmentation"
	"sap_segmentation/segmentation/client"
	"sap_segmentation/segmentation/repo/mysql"
	"time"
)

type Config struct {
	DBHost string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort string `envconfig:"DB_PORT" default:"5432"`
	DBName string `envconfig:"DB_NAME" default:"mesh_group"`
	DBUser string `envconfig:"DB_USER" default:"postgres"`
	DBPass string `envconfig:"DB_PASSWORD" default:"postgres"`

	ConnUri          string `envconfig:"CONN_URI" default:"http://bsm.api.iql.ru/ords/bsm/segmentation/get_segmentation"`
	ConnAuthLoginPWD string `envconfig:"CONN_AUTH_LOGIN_PWD" default:"4Dfddf5:jKlljHGH"`
	ConnUserAgent    string `envconfig:"CONN_USER_AGENT" default:"spacecount-test"`
	ConnTimeout      int    `envconfig:"CONN_TIMEOUT" default:"5"`
	ConnInterval     int    `envconfig:"CONN_INTERVAL" default:"1500"`

	ImportBatchSize  int `envconfig:"IMPORT_BATCH_SIZE" default:"50"`
	LogCleanupMaxAge int `envconfig:"LOG_CLEANUP_MAX_AGE" default:"7"`
}

func main() {
	myLogger, err := logger.NewLogger(os.Stdout, `log`)
	if err != nil {
		panic(err)
	}

	var conf Config
	err = envconfig.Process(``, &conf)
	if err != nil {
		myLogger.GetLogger().Println(err)
	}

	if err := myLogger.RemoveOldFile(conf.LogCleanupMaxAge); err != nil {
		myLogger.GetLogger().Println(err)
	}

	conn, err := database.MysqlOpenConnection(database.Config{
		Host:         conf.DBHost,
		Port:         conf.DBPort,
		Name:         conf.DBName,
		User:         conf.DBUser,
		Pass:         conf.DBPass,
		Debug:        "false",
		MaxIdleConns: 1,
		MaxOpenConns: 10,
		ConnMaxLife:  time.Minute * 20,
	})

	if err != nil {
		myLogger.GetLogger().Println(err)
		os.Exit(0)
	}

	rep := mysql.NewSegmentationRepository(conn)

	cl := client.NewApiClient(client.Config{
		Uri:          conf.ConnUri,
		AuthLoginPwd: conf.ConnAuthLoginPWD,
		UserAgent:    conf.ConnUserAgent,
		Interval:     conf.ConnInterval,
		Timeout:      conf.ConnTimeout,
		Limit:        conf.ImportBatchSize,
	})

	loader := segmentation.NewLoader(cl, rep, myLogger.GetLogger())

	loader.Load(0, conf.ImportBatchSize, conf.ConnInterval)
}
