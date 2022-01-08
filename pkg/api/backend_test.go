package api

import (
	"context"
	"errors"
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
			err:     errors.New("received HTTP status code 401"),
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
			err:     errors.New("received HTTP status code 400"),
		},
		{
			name: "Require header",
			backend: Backend{
				Endpoint:    "http://localhost:5555/pass-header",
				PassHeaders: "x-foo",
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
			name:         "Success",
			resourceName: "default",
			backend:      Backend{Endpoint: "http://localhost:5555/success"},
			cred:         &BackendCredential{DatabaseURL: "postgres://hostname/dbname"},
		},
	}

	testCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	go startTestBackend(testCtx, "localhost:5555")

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
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		server.Shutdown(context.Background())
	}
}
