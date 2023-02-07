package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type ErrResponse struct {
	Logger         *zap.Logger `json:"-"`
	Err            error       `json:"-"`
	HTTPStatusCode int         `json:"-"`
	StatusText     string      `json:"status"`
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	e.Logger.Error("rendering err reponse",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Int("status", e.HTTPStatusCode),
		zap.Error(e.Err),
	)
	return nil
}

func ErrBadRequest(err error, logger *zap.Logger) render.Renderer {
	return &ErrResponse{
		Logger:         logger,
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad request.",
	}
}

func ErrInternalError(err error, logger *zap.Logger) render.Renderer {
	return &ErrResponse{
		Logger:         logger,
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal error.",
	}
}
