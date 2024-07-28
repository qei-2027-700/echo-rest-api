package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
)

func main() {
	// dbをインスタンス化
	db := db.NewDB()

	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()

	// 引数として注入している コンストラクタを起動する
	userRepository := repository.NewUserRepository(db)

	taskRepository := repository.NewTaskRepository(db)

	// 外側でインスタンス化しておいたものを
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	// userUsecase := usecase.NewUserUsecase(userRepository)

	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	// taskUsecase := usecase.NewTaskUsecase(taskRepository)

	// 外側でインスタンス化したユースケースを引数に
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)

	// 外側で 引数として注入している
	e := router.NewRouter(userController, taskController)
	// e := router.NewRouter(userController)

	// e.Startでサーバーを起動する。エラーが発生したら、ログを出す
	e.Logger.Fatal(e.Start(":8080"))
}

// go run main.go
