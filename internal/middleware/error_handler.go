package middleware

import (
	"blog_platform/internal/errs"
	"blog_platform/internal/sqlerr"
	"blog_platform/internal/utils"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GlobalErrorHandler(err error, c echo.Context) {
	log.Printf("[GLOBAL ERROR HANDLER] Path: %s | Error: %v", c.Request().URL.Path, err)

	err = sqlerr.HandleError(err)

	var statusCode int
	var message string

	var customErr *errs.HTTPError
	var echoErr *echo.HTTPError

	switch {
	case errors.As(err, &customErr): // It's custom HTTP errors
		statusCode = customErr.Code
		message = customErr.Message
	case errors.As(err, &echoErr): // It's built-in echo error
		statusCode = echoErr.Code
		if msg, ok := echoErr.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(echoErr.Code)
		}

	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	if !c.Response().Committed {
		c.JSON(
			statusCode,
			utils.APIResponse{
				Success: false,
				Message: message,
			})
	}
}
