package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler_Registration_InvalidJSON(t *testing.T) {
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
			h := &handler{
				registration: nil,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			body := strings.NewReader(tt.body)

			resp := testRequest(t, ts, tt.method, tt.url, body)
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}

func Test_handler_Registration_JSON_Email(t *testing.T) {
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
					Email:          "test@mail.ru",
					HashedPassword: "test1",
					Username:       "test",
					FullName:       "test",
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
			h := &handler{
				registration: nil,
			}

			r, err := runTestServer(*h)
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			bodyJSON, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			resp := testRequest(t, ts, tt.method, tt.url, bytes.NewReader(bodyJSON))
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}
