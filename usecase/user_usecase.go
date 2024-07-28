package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	// 値渡しで受け取って 返り値の型
	SignUp(user model.User) (model.UserResponse, error)

	// JWTトークンを受け取るためstring型にしている
	Login(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository // 型を指定している
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	// func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	// interfaceだけに依存する
	return &userUsecase{ur, uv} // 構造体の実態を
	// return &userUsecase{ur} // 構造体の実態を
}

// PWをハッシュ化する
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}

	// ハッシュかする処理
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		// インスタンスを作成して返す
		return model.UserResponse{}, err
	}

	newUser := model.User{Email: user.Email, Password: string(hash)}

	//
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}

	// 実態を作成して、変数に代入
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil // エラーの返り血はnilとなる
}

// DB内に存在するかを判定する
func (uu *userUsecase) Login(user model.User) (string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}

	storedUser := model.User{}

	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}

	// 平文のPWが一致するかどうか
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	// 一致する場合は、
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		// 有効期限 12h
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	// JWTトークンを生成している
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	// JWTトークンを返している
	return tokenString, nil
}
