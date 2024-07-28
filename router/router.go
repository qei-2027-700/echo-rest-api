package router

import (
	"go-rest-api/controller"
	"net/http"

	// "net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	e := echo.New() // echoのインスタンスを作成

	// CORSのミドルウェアを追加している
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 許可するドメインを追加するフロントのドメインを追加
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		// 許可するヘッダーの追加
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		//　許可するHTTPメソッド
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		// クッキーの送受信を可能にする
		AllowCredentials: true,
	}))

	// CSRFのミドルウェアを設定
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteNoneMode, // 自動的にセキュアモードがtrueになる postmanで動作確認をしたい場合は不都合
		CookieSameSite: http.SameSiteDefaultMode,
		//CookieMaxAge:   60, // 有効期限を変更できる。秒単位
	}))

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)

	t := e.Group("/tasks")
	// JWTのミドルウェアを適用されるようにする
	t.Use(echojwt.WithConfig(echojwt.Config{
		// 同じシークレットキーを指定する
		SigningKey: []byte(os.Getenv("SECRET")),
		// クライアントから送られてくる　確認する必要がある
		// 指定している
		TokenLookup: "cookie:token",
	}))

	// タスク関係のエンドポイントを設定
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)

	return e // echoのインスタンスを返す
}
