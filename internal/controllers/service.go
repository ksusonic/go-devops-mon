package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func (c *Controller) pingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if c.db == nil {
		render.Render(w, r, ErrInternalError(fmt.Errorf("no database provided")))
		return
	}

	if err := c.db.PingContext(ctx); err != nil {
		log.Println("DB ping error:", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}
}
