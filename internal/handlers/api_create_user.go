package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/oustrix/homeset/internal/domain"
	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/pkg/homeset/http/api"
)

func (router *Router) APICreateUser(w http.ResponseWriter, r *http.Request) {
	var request api.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := router.createUser(r.Context(), unmarshalCreateUserRequest(request))
	if err != nil {
		handleCreateUserError(w, err)
		return
	}

	response(w, http.StatusCreated, marshalCreateUserResponse(result))
}

func unmarshalCreateUserRequest(req api.CreateUserRequest) users.CreateUserParams {
	return users.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	}
}

func marshalCreateUserResponse(result users.CreateUserResult) api.CreateUserResponse {
	return api.CreateUserResponse{
		User: marshalUser(result.User),
	}
}

func marshalUser(user models.User) api.User {
	return api.User{
		Username: user.Username,
	}
}

func handleCreateUserError(w http.ResponseWriter, err error) {
	var businessErr domain.Error
	if errors.As(err, &businessErr) {
		// Ok, process later
	} else {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch {
	case errors.Is(businessErr, users.ErrUserAlreadyExists):
		responseError(w, http.StatusConflict, err.Error())
	default:
		responseError(w, http.StatusInternalServerError, err.Error())
	}
}
