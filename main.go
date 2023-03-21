package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"go-chat-ai-server/app/di"
	"go-chat-ai-server/app/service"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.Logger.SetLevel(log.DEBUG)

	sr, _ := di.InitService()
	cs := chatService{chat: sr.ChatService}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("abcdefg1234-session"))))
	e.GET("/", cs.handleCharacterSelectEntry)
	e.GET("/chat", cs.handleChatEntry)
	e.POST("/chat", cs.handleChatSend)
	e.POST("/character", cs.handleSetCharacter)
	e.GET("/character", cs.handleGetCharacter)
	e.File("/character/create", "public/character_create.html")
	e.POST("/character/create", cs.handleCharacterCreate)
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
	return c.Render(http.StatusOK, "chat.html", sess.Values["name"])

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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
