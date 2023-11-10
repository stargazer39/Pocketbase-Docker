package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mattn/go-sqlite3"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
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
				return app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
					// txDao.NewExp("total > {:min} AND total < {:max}", dbx.Params{"min": 10, "max": 30})
					// update a record
					var value map[string]interface{}

					if err := c.Bind(&value); err != nil {
						log.Println("not bound", err)
						return nil
					}

					_, err := txDao.DB().
						NewQuery("INSERT INTO secrets(name, value) VALUES ({:name},{:value})").
						Bind(value).
						Execute()

					var sqliteErr sqlite3.Error

					if errors.As(err, &sqliteErr) {
						if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
							_, err := txDao.DB().
								NewQuery("UPDATE secrets SET value = {:value} WHERE name = {:name}").
								Bind(value).
								Execute()

							return err
						}
					}

					if err != nil {
						return err
					}

					// update doc
					return c.JSON(http.StatusOK, value)
				})
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
