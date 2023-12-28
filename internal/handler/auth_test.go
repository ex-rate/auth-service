package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	mock_service "github.com/ex-rate/auth-service/internal/mocks"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/ex-rate/auth-service/internal/service"
	"github.com/ex-rate/auth-service/internal/service/auth"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/ex-rate/auth-service/pkg/random"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// POST /password, status code: StatusOK
func TestHandler_AuthWithPassword_StatusOK(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		secretKey    string
		requestBody  schema.AuthWithPassword
		hashPassword string
		statusCode   int
	}{
		{
			name:      "positive case",
			method:    http.MethodPost,
			url:       "/password",
			secretKey: random.String(20),
			requestBody: schema.AuthWithPassword{
				Username: random.String(10),
				Password: random.String(12),
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenRepo := mock_service.NewMocktokenRepo(ctrl)
			authRepo := mock_service.NewMockauthRepo(ctrl)

			// services
			tokenSrv := token.New(tt.secretKey, tokenRepo)
			authSrv := auth.New(authRepo, tokenSrv)

			service := service.New(nil, tokenSrv, authSrv)

			h := &Handler{
				service: service,
			}

			hashPassword, err := token.HashPassword(tt.requestBody.Password)
			require.NoError(t, err)

			tt.hashPassword = hashPassword

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			id := uuid.New()

			authRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(id, nil)
			authRepo.EXPECT().GetHashPassword(gomock.Any(), gomock.Any()).Return(tt.hashPassword, nil)
			tokenRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(id, nil)
			tokenRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(nil)

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}

// POST /password, status code: StatusBadRequest
func TestHandler_AuthWithPassword_IncorrectPassword(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		secretKey    string
		requestBody  schema.AuthWithPassword
		hashPassword string
		statusCode   int
		expectedBody gin.H
	}{
		{
			name:      "positive case",
			method:    http.MethodPost,
			url:       "/password",
			secretKey: random.String(20),
			requestBody: schema.AuthWithPassword{
				Username: random.String(10),
				Password: random.String(12),
			},
			statusCode:   http.StatusBadRequest,
			hashPassword: random.String(15),
			expectedBody: gin.H{
				"message": api_errors.ErrIncorrectPassword.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenRepo := mock_service.NewMocktokenRepo(ctrl)
			authRepo := mock_service.NewMockauthRepo(ctrl)

			// services
			tokenSrv := token.New(tt.secretKey, tokenRepo)
			authSrv := auth.New(authRepo, tokenSrv)

			service := service.New(nil, tokenSrv, authSrv)

			h := &Handler{
				service: service,
			}

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			id := uuid.New()

			authRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(id, nil)
			authRepo.EXPECT().GetHashPassword(gomock.Any(), gomock.Any()).Return(tt.hashPassword, nil)

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			checkBody(t, tt.expectedBody, resp.Body)
		})
	}
}

// POST /password, status code: StatusBadRequest
func TestHandler_AuthWithPassword_UserNotExists(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		secretKey    string
		requestBody  schema.AuthWithPassword
		hashPassword string
		statusCode   int
		expectedBody gin.H
	}{
		{
			name:      "positive case",
			method:    http.MethodPost,
			url:       "/password",
			secretKey: random.String(20),
			requestBody: schema.AuthWithPassword{
				Username: random.String(10),
				Password: random.String(12),
			},
			statusCode:   http.StatusBadRequest,
			hashPassword: random.String(15),
			expectedBody: gin.H{
				"message": api_errors.ErrUserNotExists.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenRepo := mock_service.NewMocktokenRepo(ctrl)
			authRepo := mock_service.NewMockauthRepo(ctrl)

			// services
			tokenSrv := token.New(tt.secretKey, tokenRepo)
			authSrv := auth.New(authRepo, tokenSrv)

			service := service.New(nil, tokenSrv, authSrv)

			h := &Handler{
				service: service,
			}

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			authRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(uuid.UUID{}, api_errors.ErrUserNotExists)

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			checkBody(t, tt.expectedBody, resp.Body)
		})
	}
}
