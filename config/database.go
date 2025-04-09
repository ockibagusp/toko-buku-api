package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"toko-buku-api/pkg/logger"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func NewDatabase(viper *viper.Viper, log *logger.Logger) *sql.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	maxConnection := viper.GetInt("database.pool.max")
	idleConnection := viper.GetInt("database.pool.idle")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")
	maxIdleTimeConnection := viper.GetInt("database.pool.idletime")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", username, password, host, port, database)

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
	// 		SlowThreshold:             time.Second * 5,
	// 		Colorful:                  false,
	// 		IgnoreRecordNotFoundError: true,
	// 		ParameterizedQueries:      true,
	// 		LogLevel:                  logger.Info,
	// 	}),
	// })
	cfg := mysql.Config{
		User:   username,
		Passwd: password,
		Addr:   fmt.Sprintf("%s:%d", host, port),
		DBName: database,
		Loc: func() *time.Location {
			loc, err := time.LoadLocation("Asia/Jakarta")
			if err != nil {
				log.Fatal(context.Background(), nil, "failed to load location: %v", err)
			}
			return loc
		}(),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(context.Background(), nil, "failed to connect database: %v", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(context.Background(), nil, "failed to ping is still alive: %v", pingErr)
	}

	db.SetMaxOpenConns(maxConnection)
	db.SetMaxIdleConns(idleConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))
	db.SetConnMaxIdleTime(time.Second * time.Duration(maxIdleTimeConnection))

	return db
}
