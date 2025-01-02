package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/pkg/homeset/http/api"
)

type Router struct {
	engine http.Handler

	// Use cases
	createUser users.CreateUser
	getUser    users.GetUser
}

// RouterConfig used to provide data for NewRouter.
type RouterConfig struct {
	CreateUser users.CreateUser
	GetUser    users.GetUser
}

func NewRouter(config RouterConfig) (http.Handler, error) {
	r := http.NewServeMux()

	router := &Router{
		createUser: config.CreateUser,
		getUser:    config.GetUser,
	}

	api.HandlerFromMux(router, r)
	router.engine = r

	err := router.setupMiddleware()
	if err != nil {
		return nil, fmt.Errorf("setupMiddleware: %w", err)
	}

	return router.engine, nil
}

func (router *Router) setupMiddleware() error {
	// Swagger
	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("api.GetSwagger: %w", err)
	}

	swagger.Servers = nil

	router.engine = oapimiddleware.OapiRequestValidatorWithOptions(
		swagger,
		&oapimiddleware.Options{
			ErrorHandler: ErrorHandler,
		},
	)(router.engine)

	return nil
}

func response(w http.ResponseWriter, code int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}
