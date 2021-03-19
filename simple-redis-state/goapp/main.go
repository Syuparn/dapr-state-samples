package main

import (
	"github.com/labstack/echo"
)

func main() {
	c, teardown := NewContainer()
	defer teardown()

	c.Invoke(func(e *echo.Echo) {
		e.Logger.Fatal(e.Start(":80"))
	})
}
