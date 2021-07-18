package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nityanandagohain/log-kv-store/apigen"
)

// (POST /v1/cache)
func (h *Handler) AddKeyValue(ctx echo.Context) error {
	request := &apigen.AddKeyVal{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if err := h.Store.Put(request.Key, request.Value); err != nil {
		return ThrowError(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

// (GET /v1/cache/{key})
func (h *Handler) GetKey(ctx echo.Context, key string) error {
	val, err := h.Store.Get(key)
	if err != nil {
		return ThrowError(ctx, http.StatusInternalServerError, err.Error())

	}
	return ctx.String(http.StatusOK, val)
}
