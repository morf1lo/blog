package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	blog "github.com/morf1lo/blog"
	"github.com/morf1lo/blog/internal/config"
	"github.com/morf1lo/blog/internal/handler"
	"github.com/morf1lo/blog/internal/repository"
	"github.com/morf1lo/blog/internal/service"
	"github.com/redis/go-redis/v9"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := initEnv(); err != nil {
		logrus.Fatalf("error initializing env: %s", err.Error())
	}

	dbConfig := &config.DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSL"),
	}
	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		logrus.Fatalf("error opening postgres db: %s", err.Error())
	}

	rdbConfig := &redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: viper.GetInt("redis.db"),
		Protocol: viper.GetInt("redis.protocol"),
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	rdb := redis.NewClient(rdbConfig)

	repos := repository.New(db)
	services := service.New(repos, rdb)
	handlers := handler.New(services)

	srv := &blog.Server{}
	go func() {
		if err := srv.Run(viper.GetString("app.port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error while running server: %s", err.Error())
		}
	}()

	log.Print("Server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initEnv() error {
	return godotenv.Load()
}
