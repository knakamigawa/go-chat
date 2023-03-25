package service

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go-chat-ai-server/app/repository"
	"go-chat-ai-server/domain/model"
	"os"
	"time"
)

type UserAuthenticationParam struct {
	Email    string
	Password string
}

type NewUserParam struct {
	Name     string
	Email    string
	Password string
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func getJwtSecret() []byte {
	// FIXME レシーバーの構造体に持たせて起動時の処理でセットする
	jwtSecret := os.Getenv("JWT_SECRET")
	return []byte(jwtSecret)
}

func ProvideAuth(userRepository repository.UserMailPassword) Auth {
	return Auth{userRepository: userRepository}
}

type Auth struct {
	userRepository repository.UserMailPassword
}

func JwtConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: getJwtSecret(),
	}
}

func (r Auth) Login(c context.Context, loginRequest UserAuthenticationParam) (string, error) {
	email := model.MakeEmail(loginRequest.Email)
	password, err := model.MakePassword(loginRequest.Password)
	if err != nil {
		return "", err
	}
	user, err := r.userRepository.Of(c, email, password)
	if err != nil {
		return "", err
	}

	claims := &jwtCustomClaims{
		user.ID().String(),
		true,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * 72),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(getJwtSecret())
	if err != nil {
		return "", err
	}

	return t, nil
}

func (r Auth) Register(c context.Context, newUser NewUserParam) (string, error) {
	password, err := model.MakePassword(newUser.Password)
	if err != nil {
		return "", err
	}
	user, err := model.NewUserWithEmailPassword(newUser.Name, newUser.Email, password)
	if err != nil {
		return "", err
	}

	err = r.userRepository.Save(c, user)
	if err != nil {
		return "", err
	}

	return r.Login(c, UserAuthenticationParam{Email: newUser.Email, Password: newUser.Password})
}
