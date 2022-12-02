package api

import (
	"sort"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestSessionManager(t *testing.T) {
	t.Run("return ids", func(t *testing.T) {
		manager := NewSessionManager(nil)
		assert.Equal(t, []string{}, manager.IDs())

		manager.sessions["foo"] = &client.Client{}
		manager.sessions["bar"] = &client.Client{}

		ids := manager.IDs()
		sort.Strings(ids)
		assert.Equal(t, []string{"bar", "foo"}, ids)
	})

	t.Run("get session", func(t *testing.T) {
		manager := NewSessionManager(nil)
		assert.Nil(t, manager.Get("foo"))

		manager.sessions["foo"] = &client.Client{}
		assert.NotNil(t, manager.Get("foo"))
	})

	t.Run("set session", func(t *testing.T) {
		manager := NewSessionManager(nil)
		assert.Nil(t, manager.Get("foo"))

		manager.Add("foo", &client.Client{})
		assert.NotNil(t, manager.Get("foo"))
	})

	t.Run("remove session", func(t *testing.T) {
		manager := NewSessionManager(nil)
		assert.Nil(t, manager.Get("foo"))

		manager.Add("foo", &client.Client{})
		assert.NotNil(t, manager.Get("foo"))
		assert.True(t, manager.Remove("foo"))
		assert.False(t, manager.Remove("foo"))
		assert.Nil(t, manager.Get("foo"))
	})

	t.Run("return len", func(t *testing.T) {
		manager := NewSessionManager(nil)
		manager.sessions["foo"] = &client.Client{}
		manager.sessions["bar"] = &client.Client{}

		assert.Equal(t, 2, manager.Len())
	})

	t.Run("clean up stale sessions", func(t *testing.T) {
		defer func() {
			command.Opts.ConnectionIdleTimeout = 0
		}()

		manager := NewSessionManager(logrus.New())
		conn := &client.Client{}
		manager.Add("foo", conn)

		command.Opts.ConnectionIdleTimeout = 0
		assert.Equal(t, 1, manager.Len())
		assert.Equal(t, 0, manager.Cleanup())
		assert.Equal(t, 1, manager.Len())

		command.Opts.ConnectionIdleTimeout = 1
		assert.Equal(t, 1, manager.Cleanup())
		assert.Equal(t, 0, manager.Len())
		assert.True(t, conn.IsClosed())
	})
}
