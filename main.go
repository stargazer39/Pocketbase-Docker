package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Secret struct {
	Id    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Value any    `db:"value" json:"value"`
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /hello" route to the app router (echo)
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/secrets",
			Handler: func(c echo.Context) error {
				var value map[string]interface{}

				if err := c.Bind(&value); err != nil {
					log.Println("not bound", err)
					return nil
				}

				_, err := app.Dao().DB().
					NewQuery("INSERT INTO secrets(name, value) VALUES ({:name},{:value}) ON CONFLICT(name) DO UPDATE SET value = {:value} WHERE name = {:name}").
					Bind(value).
					Execute()

				return err
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
