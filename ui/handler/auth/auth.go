package auth

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go-chat-ai-server/app/service"
	"net/http"
)

func ProvideAuthHandler(authService service.Auth) Handler {
	return Handler{authService: authService}
}

type Handler struct {
	authService service.Auth
}

func (h Handler) LoginEntry(c echo.Context) error {
	param := struct {
		ErrorMessage string
	}{
		ErrorMessage: "",
	}
	return c.Render(http.StatusOK, "login.html", param)
}

func (h Handler) Login(c echo.Context) error {
	emil := c.FormValue("email")
	password := c.FormValue("password")
	param := service.UserAuthenticationParam{
		Email:    emil,
		Password: password,
	}

	token, err := h.authService.Login(c.Request().Context(), param)
	if err != nil {
		return loginError(c, err)
	}

	err = h.setToken(c, token)
	if err != nil {
		return loginError(c, err)
	}

	return c.Redirect(http.StatusFound, "/ai")
}

func loginError(c echo.Context, err error) error {
	c.Logger().Debugf("loginError : %s", err.Error())

	param := struct {
		ErrorMessage string
	}{
		ErrorMessage: "ログインに失敗しました",
	}
	return c.Render(http.StatusOK, "login.html", param)
}

func (h Handler) Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/")
}

func (h Handler) SignUpEntry(c echo.Context) error {
	param := struct {
		Name         string
		Email        string
		ErrorMessage string
	}{
		Name:         "",
		Email:        "",
		ErrorMessage: "",
	}
	return c.Render(http.StatusOK, "signup.html", param)
}

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

func (h Handler) SignUp(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	signUpInput := SignUpInput{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err := signUpInput.Validate()
	if err != nil {
		errs := err.(validation.Errors)

		var errorMessages []string
		for k, err := range errs {
			c.Logger().Debug(k + ": " + err.Error())
			errorMessages = append(errorMessages, k+": "+err.Error())
		}

		return errorSignUp(c, signUpInput, errorMessages)
	}

	param := service.NewUserParam{
		Name:     signUpInput.Name,
		Email:    signUpInput.Email,
		Password: signUpInput.Password,
	}

	token, err := h.authService.Register(c.Request().Context(), param)
	if err != nil {
		return errorSignUp(c, signUpInput, []string{err.Error()})
	}

	err = h.setToken(c, token)
	if err != nil {
		return errorSignUp(c, signUpInput, []string{err.Error()})
	}

	return c.Redirect(http.StatusFound, "/ai")
}

func errorSignUp(c echo.Context, signUpInput SignUpInput, errorMessages []string) error {
	c.Logger().Debugf("errorSignUp : %s", errorMessages)

	errorMessages = append([]string{"サインアップに失敗しました"}, errorMessages...)
	param := struct {
		Name          string
		Email         string
		ErrorMessages []string
	}{
		Name:          signUpInput.Name,
		Email:         signUpInput.Email,
		ErrorMessages: errorMessages,
	}
	return c.Render(http.StatusOK, "signup.html", param)
}

func (h Handler) LoginCheck(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	return sess.Values["token"] != nil
}

func (h Handler) setToken(c echo.Context, token string) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["token"] = token
	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return nil
}

func UserJwtConfig() echojwt.Config {
	return service.JwtConfig()
}

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}

func (s SignUpInput) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(
			&s.Name,
			validation.Required.Error("名前は必須入力です"),
			validation.RuneLength(5, 20).Error("名前は5～20文字です"),
		),
		validation.Field(
			&s.Email,
			validation.Required.Error("メールアドレスは必須入力です"),
			validation.RuneLength(5, 40).Error("メールアドレスは5～40文字です"),
			is.Email.Error("メールアドレスを入力して下さい"),
		),
		validation.Field(
			&s.Password,
			validation.Required.Error("パスワード必須入力です"),
			validation.RuneLength(8, 100).Error("パスワードは8～100文字です"),
		),
	)
}
