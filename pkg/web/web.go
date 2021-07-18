package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/labstack/echo/v4"
	"github.com/nityanandagohain/log-kv-store/pkg/store"
)

type Handler struct {
	Store store.Store
}

func NewHandler() *Handler {
	store, err := store.NewFileStore("../..")
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	return &Handler{
		Store: store,
	}
}

func ThrowError(ctx echo.Context, statuscode int, message string) error {
	ctx.JSON(statuscode, map[string]string{
		"message": message,
	})
	return errors.New(message)
}

type Options struct {
	Options      openapi3filter.Options
	ParamDecoder openapi3filter.ContentParameterDecoder
	UserData     interface{}
}

func ValidateRequests(swagger *openapi3.T, options *Options) echo.MiddlewareFunc {
	loader := openapi3.NewLoader()
	err := swagger.Validate(loader.Context)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// router, err := legacyrouter.NewRouter(swagger)
	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// router := openapi3filter.NewRouter().WithSwagger(swagger)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			req := ctx.Request()

			// find route
			route, pathParams, err := router.FindRoute(req)
			if err != nil {
				return err
			}

			// Validate request
			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    req,
				PathParams: pathParams,
				Route:      route,
			}

			if err := openapi3filter.ValidateRequest(ctx.Request().Context(), requestValidationInput); err != nil {
				return ThrowError(ctx, http.StatusBadRequest, err.Error())
			}

			return next(ctx)
		}
	}
}
