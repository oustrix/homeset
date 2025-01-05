package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/pkg/homeset/http/api"
)

type Middleware = func(http.Handler) http.Handler

type Router struct {
	engine http.Handler

	// Use cases
	createUser users.CreateUser
	getUser    users.GetUser
}

// RouterConfig used to provide data for NewRouter.
type RouterConfig struct {
	CreateUser  users.CreateUser
	GetUser     users.GetUser
	Middlewares []Middleware
}

func NewRouter(config RouterConfig) (http.Handler, error) {
	r := http.NewServeMux()

	router := &Router{
		createUser: config.CreateUser,
		getUser:    config.GetUser,
	}

	api.HandlerFromMux(router, r)
	router.engine = r

	for _, middleware := range config.Middlewares {
		router.engine = middleware(router.engine)
	}

	return router.engine, nil
}

func response(w http.ResponseWriter, code int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}
