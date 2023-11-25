package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	mock_handler "github.com/ex-rate/auth-service/internal/mocks"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	service "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler_Confirm_JSON_Email(t *testing.T) {
	type args struct {
		body schema.Registration
	}
	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
		dbErr      error
		args       args
	}{
		{
			name:       "valid JSON with email",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusOK,
			dbErr:      nil,
			args: args{
				schema.Registration{
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FullName:       "test",
				},
			},
		},
		{
			name:       "email exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			dbErr:      api_errors.ErrEmailAlreadyExists,
			args: args{
				schema.Registration{
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FullName:       "test",
				},
			},
		},
		{
			name:       "username exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			dbErr:      api_errors.ErrUsernameAlreadyExists,
			args: args{
				schema.Registration{
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FullName:       "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			regRepo := mock_handler.NewMockregistrationRepo(ctrl)
			tokenSrv := token.New("secret")
			regSrv := service.New(regRepo, tokenSrv)

			h := &handler{
				registration: regSrv,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			bodyJSON, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			regRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.dbErr)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON))
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}

func Test_handler_Confirm_JSON_Phone(t *testing.T) {
	type args struct {
		body schema.Registration
	}
	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
		dbErr      error
		args       args
	}{
		{
			name:       "valid JSON with phone",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusOK,
			dbErr:      nil,
			args: args{
				schema.Registration{
					PhoneNumber:    "79551234545",
					HashedPassword: "test1",
					Username:       "test",
					FullName:       "test",
				},
			},
		},
		{
			name:       "phone already exists",
			method:     http.MethodPost,
			url:        "/confirm",
			statusCode: http.StatusBadRequest,
			dbErr:      api_errors.ErrPhoneAlreadyExists,
			args: args{
				schema.Registration{
					PhoneNumber:    "79551234545",
					HashedPassword: "test1",
					Username:       "test",
					FullName:       "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			regRepo := mock_handler.NewMockregistrationRepo(ctrl)
			tokenSrv := token.New("secret")
			regSrv := service.New(regRepo, tokenSrv)

			h := &handler{
				registration: regSrv,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			bodyJSON, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			regRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.dbErr)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON))
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}
