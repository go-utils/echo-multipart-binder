package util

import "github.com/labstack/echo/v4"

// BindFunc - custom binder type
type BindFunc func(interface{}, echo.Context) error

// Bind - custom binder func
func (fn BindFunc) Bind(i interface{}, ctx echo.Context) error {
	return fn(i, ctx)
}
