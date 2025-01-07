package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/oustrix/homeset/view/templates"
)

func (router *Router) PageLogin(w http.ResponseWriter, r *http.Request) {
	if err := templates.PageLogin().Render(context.Background(), w); err != nil {
		slog.Error(err.Error())
		return
	}
}
