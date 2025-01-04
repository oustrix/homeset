package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/handlers"
	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/pkg/homeset/http/api"
	"github.com/stretchr/testify/suite"
)

type apiCreateUserSuite struct {
	suite.Suite
}

func TestAPICreateUser(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(apiCreateUserSuite))
}

func (s *apiCreateUserSuite) TestHTTP_OK() {
	var (
		username = uuid.NewString()
		password = uuid.NewString()

		user = models.User{
			Username:     username,
			PasswordHash: uuid.NewString(),
		}
	)

	createUser := func(_ context.Context, params users.CreateUserParams) (users.CreateUserResult, error) {
		s.Require().Equal(username, params.Username)
		s.Require().Equal(password, params.Password)

		return users.CreateUserResult{User: user}, nil
	}

	router, err := handlers.NewRouter(handlers.RouterConfig{
		CreateUser: createUser,
	})
	s.Require().NoError(err)

	requestBody := api.CreateUserRequest{
		Username: username,
		Password: password,
	}
	requestJSONBody, err := json.Marshal(requestBody)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(requestJSONBody)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	s.Assert().Equal(http.StatusCreated, w.Code)

	expectedResponse := api.CreateUserResponse{
		User: api.User{
			Username: user.Username,
		},
	}
	responseJSONBody, err := json.Marshal(expectedResponse)
	s.Require().NoError(err)
	s.Require().Equal(string(responseJSONBody)+"\n", w.Body.String())
}

func (s *apiCreateUserSuite) TestHTTP_Error_EmptyPassword() {
	router, err := handlers.NewRouter(handlers.RouterConfig{})
	s.Require().NoError(err)

	requestBody := api.CreateUserRequest{
		Username: uuid.NewString(),
	}
	requestJSONBody, err := json.Marshal(requestBody)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(requestJSONBody)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	s.Assert().Equal(http.StatusBadRequest, w.Code)

	expectedResponse := api.Error{
		Error:      "Error at \"password\": minimum string length is 4",
		StatusCode: http.StatusBadRequest,
	}
	responseJSONBody, err := json.Marshal(expectedResponse)
	s.Require().NoError(err)
	s.Require().Equal(string(responseJSONBody)+"\n", w.Body.String())
}

func (s *apiCreateUserSuite) TestHTTP_Error_EmptyUsername() {
	router, err := handlers.NewRouter(handlers.RouterConfig{})
	s.Require().NoError(err)

	requestBody := api.CreateUserRequest{
		Password: uuid.NewString(),
	}
	requestJSONBody, err := json.Marshal(requestBody)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(requestJSONBody)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	s.Assert().Equal(http.StatusBadRequest, w.Code)

	expectedResponse := api.Error{
		Error:      "Error at \"username\": minimum string length is 4",
		StatusCode: http.StatusBadRequest,
	}
	responseJSONBody, err := json.Marshal(expectedResponse)
	s.Require().NoError(err)
	s.Require().Equal(string(responseJSONBody)+"\n", w.Body.String())
}

func (s *apiCreateUserSuite) TestHTTP_Error_Unexpected() {
	createUser := func(_ context.Context, _ users.CreateUserParams) (users.CreateUserResult, error) {
		return users.CreateUserResult{}, errors.New("some unexpected error")
	}

	router, err := handlers.NewRouter(handlers.RouterConfig{
		CreateUser: createUser,
	})
	s.Require().NoError(err)

	requestBody := api.CreateUserRequest{
		Username: uuid.NewString(),
		Password: uuid.NewString(),
	}
	requestJSONBody, err := json.Marshal(requestBody)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(string(requestJSONBody)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	s.Assert().Equal(http.StatusInternalServerError, w.Code)

	expectedResponse := api.Error{
		Error:      "some unexpected error",
		StatusCode: http.StatusInternalServerError,
	}
	responseJSONBody, err := json.Marshal(expectedResponse)
	s.Require().NoError(err)
	s.Require().Equal(string(responseJSONBody)+"\n", w.Body.String())
}
