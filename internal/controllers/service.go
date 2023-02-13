package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func (c *Controller) pingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := c.Storage.Ping(ctx); err != nil {
		c.Logger.Error("DB ping error", zap.Error(err))
		render.Render(w, r, ErrInternalError(err, c.Logger))
		return
	}
}
