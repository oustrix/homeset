package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/oustrix/homeset/view/templates"
)

func (router *Router) PageIndex(w http.ResponseWriter, r *http.Request) {
	if err := templates.PageIndex().Render(context.Background(), w); err != nil {
		slog.Error(err.Error())
		return
	}
}
