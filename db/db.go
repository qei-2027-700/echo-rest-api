package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// 環境変数をデバッグ出力
	fmt.Println("POSTGRES_USER:", os.Getenv("POSTGRES_USER"))
	fmt.Println("POSTGRES_PW:", os.Getenv("POSTGRES_PW"))
	fmt.Println("POSTGRES_DB:", os.Getenv("POSTGRES_DB"))
	fmt.Println("POSTGRES_HOST:", os.Getenv("POSTGRES_HOST"))
	fmt.Println("POSTGRES_PORT:", os.Getenv("POSTGRES_PORT"))

	// 文字列を生成
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connceted")

	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalln("Failed to get sql.DB from gorm.DB:", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalln("Failed to close database connection:", err)
	}
}

// GO_ENV=dev go run migrate/migrate.go
