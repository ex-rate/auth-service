package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	mock_service "github.com/ex-rate/auth-service/internal/mocks"
	"github.com/ex-rate/auth-service/internal/service"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// url: /restore_token, status code: StatusOK
func TestHandler_RestoreToken_StatusOK(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		token        string
		expectedBody gin.H
		statusCode   int
	}{
		{
			name:         "valid token",
			method:       http.MethodPost,
			url:          "/restore_token",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlciI6IkpvaG4gRG9lIiwiZXhwIjoxODE2MjM5MDIyfQ.0F1OedaBS-s7BZCH8jJSfHAaiA8zbC1wmXoRuQ_NHk0",
			expectedBody: gin.H{"message": "successfully created token", "access-token": gomock.Any(), "refresh-token": gomock.Any()},
			statusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenRepo := mock_service.NewMocktokenRepo(ctrl)
			registrationRepo := mock_service.NewMockregistrationRepo(ctrl)

			// services
			tokenSrv := token.New("secret", tokenRepo)
			registrationSrv := registration.New(registrationRepo, tokenSrv)

			service := service.New(registrationSrv, tokenSrv)

			h := &handler{
				service: service,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			tokenRepo.EXPECT().CheckToken(gomock.Any(), gomock.Any()).Return(nil)
			registrationRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
			tokenRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(nil)

			resp := testRequest(t, ts, tt.method, tt.url, nil, map[string]string{"refresh-token": tt.token})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			var actualBody gin.H
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&actualBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}

// url: /restore_token, status code: StatusBadRequest
func TestHandler_RestoreToken_InvalidToken(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		token        string
		expectedBody gin.H
		statusCode   int
	}{
		{
			name:         "invalid token",
			method:       http.MethodPost,
			url:          "/restore_token",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlciI6IkpvaG4gRG9lIiwiZXhwIjoxODE2MjM5MDIyfQ.0F1OedaBS-s7BZCH8jJSfHAaiA8zbC1wmXoRuQ_NHk0",
			expectedBody: gin.H{"message": api_errors.ErrInvalidToken.Error()},
			statusCode:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenRepo := mock_service.NewMocktokenRepo(ctrl)
			registrationRepo := mock_service.NewMockregistrationRepo(ctrl)

			// services
			tokenSrv := token.New("secret", tokenRepo)
			registrationSrv := registration.New(registrationRepo, tokenSrv)

			service := service.New(registrationSrv, tokenSrv)

			h := &handler{
				service: service,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			tokenRepo.EXPECT().CheckToken(gomock.Any(), gomock.Any()).Return(api_errors.ErrInvalidToken)

			resp := testRequest(t, ts, tt.method, tt.url, nil, map[string]string{"refresh-token": tt.token})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			var actualBody gin.H
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&actualBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}
