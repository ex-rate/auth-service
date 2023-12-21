package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ex-rate/auth-service/internal/entities"
	api_errors "github.com/ex-rate/auth-service/internal/errors"
	mock_service "github.com/ex-rate/auth-service/internal/mocks"
	"github.com/ex-rate/auth-service/internal/service"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/ex-rate/auth-service/pkg/random"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PUT /restore_token, status code: StatusOK
func TestHandler_RestoreToken_StatusOK(t *testing.T) {
	type args struct {
		user      string
		secretKey string
	}
	type headers struct {
		accessToken string
	}

	hour := time.Now().Add(1 * time.Hour)
	month := time.Now().Add(30 * 24 * time.Hour)

	tests := []struct {
		name         string
		method       string
		args         args
		url          string
		headers      headers
		refreshToken string
		requestBody  entities.RestoreToken
		statusCode   int
	}{
		{
			name:   "valid token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10),
				secretKey: "secret",
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "valid token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: "secret",
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "valid token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: "secret",
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "valid token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: "secret",
			},
			statusCode: http.StatusOK,
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

			// generating test tokens
			accessToken, accessTokenExpectedExp := generateToken(t, tt.args.secretKey, tt.args.user, hour)
			refreshToken, refreshTokenExpectedExp := generateToken(t, tt.args.secretKey, tt.args.user, month)

			tt.headers.accessToken = fmt.Sprintf("Bearer %s", accessToken)
			tt.requestBody.RefreshToken = refreshToken

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{AuthorizationHeader: tt.headers.accessToken})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			var actualBody gin.H
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&actualBody)
			require.NoError(t, err)

			actualRefreshToken := actualBody["refresh-token"].(string)
			actualAccessToken := actualBody["access-token"].(string)

			checkToken(t, actualRefreshToken, tt.args.user, tt.args.secretKey, float64(refreshTokenExpectedExp))
			checkToken(t, actualAccessToken, tt.args.user, tt.args.secretKey, float64(accessTokenExpectedExp))
		})
	}
}

// PUT /restore_token, status code: StatusBadRequest
func TestHandler_RestoreToken_InvalidToken(t *testing.T) {
	type args struct {
		user      string
		secretKey string
	}
	type headers struct {
		accessToken string
	}

	day := time.Hour * 24

	tests := []struct {
		name          string
		method        string
		args          args
		headers       headers
		requestBody   entities.RestoreToken
		tokenDuration map[string]time.Time
		url           string
		token         string
		expectedBody  gin.H
		statusCode    int
	}{
		{
			name:   "expired access token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(30 * day),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired refresh token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(11),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(30 * day),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(30 * day),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(30 * day),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired refresh token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired refresh token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired refresh token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access & refresh token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access & refresh token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access & refresh token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
		{
			name:   "expired access & refresh token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: "secret",
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-30 * 24 * time.Hour),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidToken.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSrv := token.New("secret", nil)
			registrationSrv := registration.New(nil, tokenSrv)

			service := service.New(registrationSrv, tokenSrv)

			h := &handler{
				service: service,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			// generating test tokens
			accessToken, _ := generateToken(t, tt.args.secretKey, tt.args.user, tt.tokenDuration["access-token"])
			refreshToken, _ := generateToken(t, tt.args.secretKey, tt.args.user, tt.tokenDuration["refresh-token"])

			tt.headers.accessToken = fmt.Sprintf("Bearer %s", accessToken)
			tt.requestBody.RefreshToken = refreshToken

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{AuthorizationHeader: tt.headers.accessToken})
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

// generateToken генерирует токен с заданными данными. Возвращает токен и дату истечения
func generateToken(t *testing.T, secretKey, username string, exp time.Time) (string, int64) {
	token := jwt.New(jwt.SigningMethodHS256)
	key := []byte(secretKey)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp.Unix()
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(key)
	require.NoError(t, err)

	return tokenString, exp.Unix()
}

func checkToken(t *testing.T, token, expectedUser, secretKey string, expectedExp float64) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)

	mapClaims := jwtToken.Claims.(jwt.MapClaims)
	actualExpr := mapClaims["exp"].(float64)
	actualUser := mapClaims["user"].(string)

	assert.Equal(t, expectedExp, actualExpr, fmt.Sprintf("expiration time not equal: expected: %v actual: %v", expectedExp, actualExpr))
	assert.Equal(t, expectedUser, actualUser, fmt.Sprintf("username not equal: expected: %v actual: %v", expectedUser, actualUser))
}
