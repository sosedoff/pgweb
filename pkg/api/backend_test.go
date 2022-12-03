package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBackendFetchCredential(t *testing.T) {
	examples := []struct {
		name         string
		backend      Backend
		resourceName string
		cred         *BackendCredential
		reqCtx       *gin.Context
		ctx          func() (context.Context, context.CancelFunc)
		err          error
	}{
		{
			name:    "Bad auth token",
			backend: Backend{Endpoint: "http://localhost:5555/unauthorized"},
			err:     errors.New("backend credential fetch received HTTP status code 401"),
		},
		{
			name:    "Backend timeout",
			backend: Backend{Endpoint: "http://localhost:5555/timeout"},
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Millisecond*100)
			},
			err: errors.New("Unable to connect to the auth backend"),
		},
		{
			name:    "Empty response",
			backend: Backend{Endpoint: "http://localhost:5555/empty-response"},
			err:     errors.New("Connection string is required"),
		},
		{
			name:    "Missing header",
			backend: Backend{Endpoint: "http://localhost:5555/pass-header"},
			err:     errors.New("backend credential fetch received HTTP status code 400"),
		},
		{
			name: "Require header",
			backend: Backend{
				Endpoint:    "http://localhost:5555/pass-header",
				PassHeaders: []string{"x-foo"},
			},
			reqCtx: &gin.Context{
				Request: &http.Request{
					Header: http.Header{
						"X-Foo": []string{"bar"},
					},
				},
			},
			cred: &BackendCredential{DatabaseURL: "postgres://hostname/bar"},
		},
		{
			name:    "Success",
			backend: Backend{Endpoint: "http://localhost:5555/success"},
			cred:    &BackendCredential{DatabaseURL: "postgres://hostname/dbname"},
		},
	}

	srvCtx, srvCancel := context.WithTimeout(context.Background(), time.Minute)
	defer srvCancel()

	startTestBackend(srvCtx, "localhost:5555")

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			if ex.ctx != nil {
				ctx, cancel = ex.ctx()
			}
			defer cancel()

			reqCtx := ex.reqCtx
			if reqCtx == nil {
				reqCtx = &gin.Context{
					Request: &http.Request{},
				}
			}

			cred, err := ex.backend.FetchCredential(ctx, ex.resourceName, reqCtx)
			assert.Equal(t, ex.err, err)
			assert.Equal(t, ex.cred, cred)
		})
	}
}

func startTestBackend(ctx context.Context, listenAddr string) {
	router := gin.New()

	router.Use(func(c *gin.Context) {
		if c.GetHeader("content-type") != "application/json" {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	})

	router.POST("/unauthorized", func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})

	router.POST("/timeout", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.JSON(http.StatusOK, gin.H{})
	})

	router.POST("/empty-response", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	router.POST("/pass-header", func(c *gin.Context) {
		req := BackendRequest{}
		if err := c.BindJSON(&req); err != nil {
			panic(err)
		}

		header := req.Headers["x-foo"]
		if header == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"database_url": "postgres://hostname/" + header,
		})
	})

	router.POST("/success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"database_url": "postgres://hostname/dbname",
		})
	})

	server := &http.Server{Addr: listenAddr, Handler: router}
	mustStartServer(server)

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func mustStartServer(server *http.Server) {
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	if err := waitForServer(server.Addr, 5); err != nil {
		panic(err)
	}
}

func waitForServer(addr string, n int) error {
	var lastErr error

	for i := 0; i < n; i++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}

		lastErr = err
		time.Sleep(time.Millisecond * 100)
	}

	return lastErr
}
