package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	// これはinterfaceなので、必ずこのメソッドを実装しなければならない
	SignUp(c echo.Context) error // 返り血はerror interface型
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type userController struct {
	// uuと言う名前で定義
	uu usecase.IUserUsecase
}

// 引数として注入できるように指定おく
func NewUserController(uu usecase.IUserUsecase) IUserController {
	//
	return &userController{uu}
}

// 引数の方と、戻り値の型は、interfaceで書かれているものと全く一緒
func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}

	// userオブジェクトのおポインタを
	if err := c.Bind(&user); err != nil {
		// bindを失敗した場合は、
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		// サインアップに失敗したら
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	// Bindメソッド
	if err := c.Bind(&user); err != nil {
		// 変換に失敗したら、、、
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	tokenString, err := uc.uu.Login(user) // JWTを生成

	if err != nil {
		// エラーメッセージをクライアントに返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// サーバーサイドでクッキーに
	cookie := new(http.Cookie) // cookie構造体を作成
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true // postmanで検証したい場合は、コメアウトしておく
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode // CORS
	c.SetCookie(cookie)                     // httpレスポンスに含めるようにする
	return c.NoContent(http.StatusOK)
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now() // 有効期限がすぐにきれるようにしておく
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
