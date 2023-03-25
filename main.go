package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"go-chat-ai-server/app/di"
	"go-chat-ai-server/app/service"

	auth2 "go-chat-ai-server/ui/handler/auth"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.Logger())
	e.Use(CsrfPostMiddleware)
	e.Use(middleware.CSRF())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("abcdefg1234-session"))))

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.Logger.SetLevel(log.DEBUG)

	sr, _ := di.InitService()
	cs := chatService{chat: sr.ChatService}
	authHandler := auth2.ProvideAuthHandler(sr.AuthService)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/login")
	})
	e.GET("/login", authHandler.LoginEntry)
	e.POST("/login", authHandler.Login)

	e.GET("/signup", authHandler.SignUpEntry)
	e.POST("/signup", authHandler.SignUp)

	authenticated := e.Group("/ai")
	authenticated.Use(SessionAuthMiddleware)
	authenticated.Use(echojwt.WithConfig(auth2.UserJwtConfig()))

	authenticated.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/ai/entry")
	})
	authenticated.GET("/entry", cs.handleCharacterSelectEntry)
	authenticated.GET("/chat", cs.handleChatEntry)
	authenticated.POST("/character", cs.handleSetCharacter)
	authenticated.File("/character/create", "public/character_create.html")
	authenticated.POST("/character/create", cs.handleCharacterCreate)
	authenticated.GET("/logout", authHandler.Logout)

	api := authenticated.Group("/api")
	api.Use(APIKeyAuthMiddleware)
	api.POST("/chat", cs.handleChatSend)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

type chatService struct {
	chat service.Chat
}

func (s chatService) handleSetCharacter(c echo.Context) error {
	text := c.FormValue("name")
	if text == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"msg": "キャラクター名を入力してください",
		})
	}
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["name"] = text
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "character.html", sess.Values["name"])
}

func (s chatService) handleGetCharacter(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "character.html", sess.Values["name"])
}

func (s chatService) handleCharacterSelectEntry(c echo.Context) error {
	names, err := s.chat.Characters(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	params := struct {
		Names []string
	}{
		Names: names,
	}

	return c.Render(http.StatusOK, "character_select.html", params)
}

func (s chatService) handleChatEntry(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	params := struct {
		Name   string
		ApiKey string
	}{
		Name:   sess.Values["name"].(string),
		ApiKey: os.Getenv("MY_API_TOKEN"),
	}

	return c.Render(http.StatusOK, "chat.html", params)

}

func (s chatService) handleChatSend(c echo.Context) error {
	type payloadType struct {
		Text string `json:"text"`
	}

	payload := new(payloadType)
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid JSON")
	}

	c.Logger().Debugf("body %s\n\n", payload)

	if payload.Text == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"msg": "質問を入力してください",
		})
	}
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	name := sess.Values["name"].(string)
	rt, err := s.chat.Talk(c.Request().Context(), payload.Text, name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"msg": rt,
	})
}

func (s chatService) handleCharacterCreate(c echo.Context) error {
	name := c.FormValue("name")
	bio := c.FormValue("bio")

	err := s.chat.CharacterCreate(c.Request().Context(), name, bio)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.Redirect(http.StatusFound, "/")
}

func APIKeyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if authToken != fmt.Sprintf("Bearer %s", os.Getenv("MY_API_TOKEN")) {
			c.Logger().Debugf("invalid api key: ", authToken)
			return fmt.Errorf("invalid api key")
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", sess.Values["token"]))
		c.Logger().Debugf("session: %v", sess.Values)
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func CsrfPostMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		csrfCookie, err := c.Request().Cookie("_csrf")
		if err != nil {
			c.Logger().Debugf("csrf cookie error: %s", err.Error())
		}
		csrf := csrfCookie.Value
		if csrf != "" {
			c.Request().Header.Add("X-CSRF-Token", csrf)
			c.Logger().Debugf("csrf: %s", c.Request().FormValue("_csrf"))
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
