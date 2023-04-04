package main

import (
	"log"
	"os"

	"github.com/johnnyzhao/auth-server-api/internal/app"
	"github.com/johnnyzhao/auth-server-api/internal/pkg/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

func main() {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		log.Fatal("database dsn not found in env DB_DSN")
	}

	newLogger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger2.Config{
			LogLevel: logger2.Info, // Log level
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}
	r := repo.NewUserRepo(db)
	if err := r.Migrate(); err != nil {
		log.Fatal("db migration error", err)
	}
	server := app.NewServer(app.NewApi(r))
	server.Setup()
	server.Run(":8080")
}
