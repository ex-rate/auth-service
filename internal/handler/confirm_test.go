package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	mock_service "github.com/ex-rate/auth-service/internal/mocks"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/ex-rate/auth-service/internal/service"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// POST /confirm, status code: StatusOK
func TestHandler_Confirm_StatusOK(t *testing.T) {
	type args struct {
		body      schema.Registration
		user      string
		secretKey string
	}
	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
		args       args
	}{
		{
			name:       "valid JSON with email",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusOK,
			args: args{
				body: schema.Registration{
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FirstName:      "test",
					LastName:       "test",
					Patronymic:     "test",
				},
				user:      "test",
				secretKey: "secret",
			},
		},
		{
			name:       "valid JSON with phone",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusOK,
			args: args{
				body: schema.Registration{
					PhoneNumber:    "79999999999",
					HashedPassword: "test1",
					Username:       "test",
					FirstName:      "test",
					LastName:       "test",
					Patronymic:     "test",
				},
				user:      "test",
				secretKey: "secret",
			},
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

			bodyJSON, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			id := uuid.New()

			registrationRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(id, nil)
			tokenRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(nil)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			var actualBody gin.H
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&actualBody)
			require.NoError(t, err)

			token := actualBody["refresh-token"].(string)
			checkUsername(t, token, tt.args.user, tt.args.secretKey)
		})
	}
}

func checkUsername(t *testing.T, token, username, secretKey string) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)

	mapClaims := jwtToken.Claims.(jwt.MapClaims)
	actualUser := mapClaims["user"].(string)

	assert.Equal(t, username, actualUser, fmt.Sprintf("username not equal: expected: %v actual: %v", username, actualUser))
}

// POST /confirm, status code: StatusBadRequest
func TestHandler_Confirm_StatusBadRequest(t *testing.T) {
	type args struct {
		body schema.Registration
	}
	tests := []struct {
		name         string
		method       string
		url          string
		statusCode   int
		args         args
		dbErr        error
		expectedBody gin.H
	}{
		{
			name:       "valid JSON with email: username already exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			args: args{
				schema.Registration{
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FirstName:      "test",
					LastName:       "test",
					Patronymic:     "test",
				},
			},
			dbErr: api_errors.ErrUsernameAlreadyExists,
		},
		{
			name:       "valid JSON with phone: username already exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			args: args{
				schema.Registration{
					PhoneNumber:    "79999999999",
					HashedPassword: "test1",
					Username:       "test",
					FirstName:      "test",
					LastName:       "test",
					Patronymic:     "test",
				},
			},
			dbErr: api_errors.ErrUsernameAlreadyExists,
		},
		{
			name:       "valid JSON with email: email already exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			args: args{
				schema.Registration{
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FirstName:      "test",
					LastName:       "test",
					Patronymic:     "test",
				},
			},
			dbErr: api_errors.ErrEmailAlreadyExists,
		},
		{
			name:       "valid JSON with phone: phone already exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			args: args{
				schema.Registration{
					PhoneNumber:    "79999999999",
					HashedPassword: "test1",
					Username:       "test",
					FirstName:      "test",
					LastName:       "test",
					Patronymic:     "test",
				},
			},
			dbErr: api_errors.ErrPhoneAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			registrationRepo := mock_service.NewMockregistrationRepo(ctrl)

			// services
			registrationSrv := registration.New(registrationRepo, nil)

			service := service.New(registrationSrv, nil)

			h := &handler{
				service: service,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			bodyJSON, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			id := uuid.New()

			registrationRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(id, tt.dbErr)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			tt.expectedBody = gin.H{"message": tt.dbErr.Error()}

			var actualBody gin.H
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&actualBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}
