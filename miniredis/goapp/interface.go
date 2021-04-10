package main

import (
	"net/http"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewController makes new HTTP controller.
func NewController(
	c dapr.Client,
	createInputPort CreateMessageInputPort,
) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/messages", createMessage(c, createInputPort))

	return e
}

func createMessage(cli dapr.Client, port CreateMessageInputPort) func(c echo.Context) error {
	return func(c echo.Context) error {
		type reqBody struct {
			Message string `json:"message"`
		}
		b := &reqBody{}

		if err := c.Bind(b); err != nil {
			return err
		}

		in := &CreateMessageInputData{
			Message: b.Message,
		}

		out, err := port.Do(c.Request().Context(), in)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, out.Message)
	}
}
