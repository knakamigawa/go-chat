package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go-chat/app/use_case"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Logger.SetLevel(log.DEBUG)
	e.Use(APIKeyAuthMiddleware)

	e.GET("/", handleChat)
	e.Logger.Fatal(e.Start(":8081"))
}

func handleChat(c echo.Context) error {
	text := c.FormValue("text")
	if text == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"msg": "質問を入力してください",
		})
	}
	rt, err := use_case.SendTalk(text)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"msg": rt,
	})
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
