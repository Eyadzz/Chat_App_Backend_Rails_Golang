package persistence

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDB() *gorm.DB {
	// Database connection configuration (without specifying the database name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"))
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database server:", err)
	}

	// Create the database if it doesn't exist
	dbName := os.Getenv("DB_NAME")
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)
	if err := db.Exec(sql).Error; err != nil {
		log.Fatal("Failed to create database:", err)
	}

	// Reconnect to the database with the specified database name
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), dbName)
	db, err = gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func ConnectToRedis() *redis.Client {

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))

	if err != nil {
		log.Fatal("Failed to connect to redis!")
	}

	redisClient := redis.NewClient(opt)

	return redisClient
}
