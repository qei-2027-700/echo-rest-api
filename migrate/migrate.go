package main
// mainパッケージに所属させる

import (
	"fmt"
	"go-rest-api/db"
	"go-rest-api/model"
)

func main() {
	dbConn := db.NewDB()

	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	// DBに反映させたいモデル構造 インスタンス化して
	dbConn.AutoMigrate(&model.User{}, &model.Task{})
}
// go run migrate/migrate.go
// main関数が

// GO_ENV=dev go run migrate/migrate.go