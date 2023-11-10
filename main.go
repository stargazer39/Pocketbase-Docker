package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Secret struct {
	Id    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Value string `db:"value" json:"value"`
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

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/secrets/:id",
			Handler: func(c echo.Context) error {
				id := c.PathParam("id")
				var value Secret
				var doc any

				err := app.Dao().DB().Select("value").From("secrets").Where(dbx.Like("name", id)).One(&value)

				if err != nil {
					return err
				}

				if err := json.Unmarshal([]byte(value.Value), &doc); err != nil {
					return err
				}

				return c.JSON(http.StatusOK, doc)
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
