package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Hello, World!")
	//})

	e.Static("/", "public")
	e.Logger.Fatal(e.Start(":1323"))
}