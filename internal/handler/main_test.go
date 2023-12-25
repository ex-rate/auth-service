package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader, headers map[string]string) *http.Response {

	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Close = true
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "PostmanRuntime/7.32.3")
	require.NoError(t, err)

	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	ts.Client()

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func runTestServer(handler Handler) *gin.Engine {
	r := gin.Default()

	// registration
	r.GET("/signup", handler.GetRegistration)
	r.POST("/signup", handler.Registration)

	// confirm registration
	r.GET("/confirm", handler.GetConfirm)
	r.POST("/confirm", handler.Confirm)

	// restore token
	r.PUT("/restore_token", handler.RestoreToken)

	// authorization
	r.GET("/login", handler.GetAuth)
	r.POST("/login", handler.Auth)

	// authorization via code
	r.GET("/code", handler.GetCode)
	r.POST("/code", handler.Code)

	// authorization via password
	r.GET("/password", handler.GetPassword)
	r.POST("/password", handler.Password)

	return r
}
