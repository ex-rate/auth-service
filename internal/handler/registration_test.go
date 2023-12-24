package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/ex-rate/auth-service/internal/service"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/ex-rate/auth-service/pkg/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// POST /signup, status code: StatusBadRequest
func TestHandler_SignUp_StatusBadRequest(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		body       string
		statusCode int
	}{
		{
			name:   "invalid JSON",
			method: http.MethodPost,
			url:    "/signup",
			body: `{
				"email": "sss",
				"hash_password": "tasfasest",`,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// в этом тесте БД не нужна

			// services
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

			body := strings.NewReader(tt.body)

			resp := testRequest(t, ts, tt.method, tt.url, body, map[string]string{})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}

// POST /signup, status code: StatusPermanentRedirect
func TestHandler_SignUp_StatusPermanentRedirect(t *testing.T) {
	type args struct {
		body schema.Registration
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
			url:        "/signup",
			statusCode: http.StatusPermanentRedirect,
			args: args{
				schema.Registration{
					Email:          random.Email(4),
					HashedPassword: random.String(10),
					Username:       random.String(5),
					FirstName:      random.String(5),
					LastName:       random.String(5),
					Patronymic:     random.String(5),
				},
			},
		},
		{
			name:       "valid JSON with phone",
			method:     http.MethodPost,
			url:        "/signup",
			statusCode: http.StatusPermanentRedirect,
			args: args{
				schema.Registration{
					PhoneNumber:    random.Phone(),
					HashedPassword: random.String(8),
					Username:       random.String(9),
					FirstName:      random.String(4),
					LastName:       random.String(6),
					Patronymic:     random.String(4),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// services
			tokenSrv := token.New(random.String(5), nil)
			registrationSrv := registration.New(nil, tokenSrv)

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

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON), map[string]string{})
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}
