package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// echo FWのコンテキストを
type ITaskController interface {
	GetAllTasks(c echo.Context) error
	GetTaskById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type taskController struct {
	tu usecase.ITaskUsecase
}

// コンストラクタ
// 　外側で
func NewTaskController(tu usecase.ITaskUsecase) ITaskController {
	// ポインタをreturnで返す
	return &taskController{tu}
}

// JWTトークンに組み込まれている、 decodeした内容に、　自動的に格納指定うれる
func (tc *taskController) GetAllTasks(c echo.Context) error {
	// JWTからdecodeした値を読み込んでいく
	user := c.Get("user").(*jwt.Token)
	// 取り出して
	claims := user.Claims.(jwt.MapClaims)
	// くれ＾夢の中に存在する user_idに代入する
	userId := claims["user_id"]

	// any型しているから、float64にアサーションする
	tasksRes, err := tc.tu.GetAllTasks(uint(userId.(float64)))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// contextJSONと言う意味
	return c.JSON(http.StatusOK, tasksRes)
}

func (tc *taskController) GetTaskById(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	// リクエストパラメータから
	id := c.Param("taskId")
	// stringからint型に
	taskId, _ := strconv.Atoi(id)
	//
	taskRes, err := tc.tu.GetTaskById(uint(userId.(float64)), uint(taskId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

// コンテキストの中から、userIdを取得
func (tc *taskController) CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	task := model.Task{}
	// コンテキストバインドを使うことで
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	task.UserId = uint(userId.(float64))

	//
	taskRes, err := tc.tu.CreateTask(task)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// context.JSONで新規作成されたObjectをクライアントに返している
	return c.JSON(http.StatusCreated, taskRes)
}

// コンテキストからuserIDを取得
func (tc *taskController) UpdateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	// int型に変換
	taskId, _ := strconv.Atoi(id)

	task := model.Task{}

	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	//
	taskRes, err := tc.tu.UpdateTask(task, uint(userId.(float64)), uint(taskId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

// コンテキストからUSERIDを取得
func (tc *taskController) DeleteTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	err := tc.tu.DeleteTask(uint(userId.(float64)), uint(taskId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
