package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	mock_service "github.com/ex-rate/auth-service/internal/mocks"
	schema "github.com/ex-rate/auth-service/internal/schemas"
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
		requestBody  schema.RestoreToken
		statusCode   int
	}{
		{
			name:   "valid token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10),
				secretKey: random.String(5),
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "valid token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: random.String(5),
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "valid token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: random.String(5),
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "valid token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: random.String(5),
			},
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenRepo := mock_service.NewMocktokenRepo(ctrl)

			// services
			tokenSrv := token.New(tt.args.secretKey, tokenRepo)
			registrationSrv := registration.New(nil, tokenSrv)

			service := service.New(registrationSrv, tokenSrv, nil)

			h := &Handler{
				service: service,
			}

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			tokenRepo.EXPECT().CheckToken(gomock.Any(), gomock.Any()).Return(nil)
			tokenRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(nil)
			tokenRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)

			// generating test tokens
			accessToken, accessTokenExpectedExp := generateToken(t, tt.args.secretKey, tt.args.user, hour, true)
			refreshToken, refreshTokenExpectedExp := generateToken(t, tt.args.secretKey, tt.args.user, month, true)

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

	day := 24 * time.Hour
	month := 30 * day

	tests := []struct {
		name          string
		method        string
		args          args
		headers       headers
		requestBody   schema.RestoreToken
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
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired refresh token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(11),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking refresh token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired refresh token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking refresh token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired refresh token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking refresh token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired refresh token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking refresh token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access & refresh token",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access & refresh token: long data #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(100),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access & refresh token: long data #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(1000),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
		{
			name:   "expired access & refresh token: long data #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				user:      random.String(10000),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(-1 * time.Hour),
				"refresh-token": time.Now().Add(-month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrInvalidToken),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSrv := token.New(tt.args.secretKey, nil)
			registrationSrv := registration.New(nil, tokenSrv)

			service := service.New(registrationSrv, tokenSrv, nil)

			h := &Handler{
				service: service,
			}

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			// generating test tokens
			accessToken, _ := generateToken(t, tt.args.secretKey, tt.args.user, tt.tokenDuration["access-token"], true)
			refreshToken, _ := generateToken(t, tt.args.secretKey, tt.args.user, tt.tokenDuration["refresh-token"], true)

			tt.headers.accessToken = fmt.Sprintf("Bearer %s", accessToken)
			tt.requestBody.RefreshToken = refreshToken

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{AuthorizationHeader: tt.headers.accessToken})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			checkBody(t, tt.expectedBody, resp.Body)
		})
	}
}

// PUT /restore_token, status code: StatusBadRequest
func TestHandler_RestoreToken_InvalidUsername(t *testing.T) {
	type args struct {
		refreshUsername string
		accessUsername  string
		secretKey       string
	}
	type headers struct {
		accessToken string
	}

	day := 24 * time.Hour
	month := 30 * day

	tests := []struct {
		name          string
		method        string
		args          args
		headers       headers
		requestBody   schema.RestoreToken
		tokenDuration map[string]time.Time
		url           string
		token         string
		expectedBody  gin.H
		statusCode    int
	}{
		{
			name:   "usernames does not match #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				refreshUsername: random.String(10),
				accessUsername:  random.String(5),
				secretKey:       random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidUsername.Error(),
			},
		},
		{
			name:   "usernames does not match #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				refreshUsername: random.String(100),
				accessUsername:  random.String(26),
				secretKey:       random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidUsername.Error(),
			},
		},
		{
			name:   "usernames does not match #3",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				refreshUsername: random.String(50),
				accessUsername:  random.String(70),
				secretKey:       random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"message": api_errors.ErrInvalidUsername.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSrv := token.New(tt.args.secretKey, nil)
			registrationSrv := registration.New(nil, tokenSrv)

			service := service.New(registrationSrv, tokenSrv, nil)

			h := &Handler{
				service: service,
			}

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			// generating test tokens
			accessToken, _ := generateToken(t, tt.args.secretKey, tt.args.accessUsername, tt.tokenDuration["access-token"], true)
			refreshToken, _ := generateToken(t, tt.args.secretKey, tt.args.refreshUsername, tt.tokenDuration["refresh-token"], true)

			tt.headers.accessToken = fmt.Sprintf("Bearer %s", accessToken)
			tt.requestBody.RefreshToken = refreshToken

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{AuthorizationHeader: tt.headers.accessToken})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			checkBody(t, tt.expectedBody, resp.Body)
		})
	}
}

// PUT /restore_token, status code: StatusUnauthorized
func TestHandler_RestoreToken_Unauthorized(t *testing.T) {
	type args struct {
		username  string
		secretKey string
	}
	type headers struct {
		accessToken string
	}

	day := 24 * time.Hour
	month := 30 * day

	tests := []struct {
		name          string
		method        string
		args          args
		headers       headers
		requestBody   schema.RestoreToken
		tokenDuration map[string]time.Time
		url           string
		token         string
		expectedBody  gin.H
		statusCode    int
	}{
		{
			name:   "unauthorized #1",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				username:  random.String(7),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusUnauthorized,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrNotAuthorized.Error()),
			},
		},
		{
			name:   "unauthorized #2",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				username:  random.String(15),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusUnauthorized,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrNotAuthorized.Error()),
			},
		},
		{
			name:   "unauthorized: long data",
			method: http.MethodPut,
			url:    "/restore_token",
			args: args{
				username:  random.String(150),
				secretKey: random.String(5),
			},
			tokenDuration: map[string]time.Time{
				"access-token":  time.Now().Add(1 * time.Hour),
				"refresh-token": time.Now().Add(month),
			},
			statusCode: http.StatusUnauthorized,
			expectedBody: gin.H{
				"message": fmt.Sprintf("error while checking access token: %v", api_errors.ErrNotAuthorized.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSrv := token.New(tt.args.secretKey, nil)
			registrationSrv := registration.New(nil, tokenSrv)

			service := service.New(registrationSrv, tokenSrv, nil)

			h := &Handler{
				service: service,
			}

			r := runTestServer(*h)

			ts := httptest.NewServer(r)
			defer ts.Close()

			// generating test tokens
			accessToken, _ := generateToken(t, tt.args.secretKey, tt.args.username, tt.tokenDuration["access-token"], false)
			refreshToken, _ := generateToken(t, tt.args.secretKey, tt.args.username, tt.tokenDuration["refresh-token"], false)

			tt.headers.accessToken = fmt.Sprintf("Bearer %s", accessToken)
			tt.requestBody.RefreshToken = refreshToken

			bodyJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{AuthorizationHeader: tt.headers.accessToken})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			checkBody(t, tt.expectedBody, resp.Body)
		})
	}
}

func checkBody(t *testing.T, expectedBody gin.H, body io.ReadCloser) {
	var actualBody gin.H
	dec := json.NewDecoder(body)
	err := dec.Decode(&actualBody)
	require.NoError(t, err)

	assert.Equal(t, expectedBody, actualBody)
}

// generateToken генерирует токен с заданными данными. Возвращает токен и дату истечения
func generateToken(t *testing.T, secretKey, username string, exp time.Time, auth bool) (string, int64) {
	token := jwt.New(jwt.SigningMethodHS256)
	key := []byte(secretKey)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp.Unix()
	claims["authorized"] = auth
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
