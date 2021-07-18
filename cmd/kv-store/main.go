package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nityanandagohain/log-kv-store/apigen"
	handler "github.com/nityanandagohain/log-kv-store/pkg/web"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())

	// Panic Handler
	e.Use(middleware.Recover())

	// RequestID
	e.Use(middleware.RequestID())

	swagger, err := apigen.GetSwagger()
	if err != nil {
		log.Fatalf(err.Error())
	}
	e.Use(handler.ValidateRequests(swagger, &handler.Options{}))

	cacheHandler := handler.NewHandler()
	apigen.RegisterHandlersWithBaseURL(e, cacheHandler, "")
	e.Start(":3000")
}
