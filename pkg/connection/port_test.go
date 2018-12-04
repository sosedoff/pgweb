package connection

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPortAvailable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIXME")
	}

	assert.Equal(t, true, IsPortAvailable(30000))

	serv, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to start test tcp listener:", err)
		t.Fail()
		return
	}
	defer serv.Close()

	go func() {
		for {
			conn, err := serv.Accept()
			if err == nil {
				conn.Close()
			}
			serv.Close()
		}
	}()

	assert.Equal(t, false, IsPortAvailable(30000))
	assert.Equal(t, true, IsPortAvailable(30001))
}

func TestFindAvailablePort(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIXME")
	}

	port, err := FindAvailablePort(30000, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 30000, port)

	serv, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to start test tcp listener:", err)
		t.Fail()
		return
	}
	defer serv.Close()

	go func() {
		for {
			conn, err := serv.Accept()
			if err == nil {
				conn.Close()
			}
		}
	}()

	port, err = FindAvailablePort(30000, 0)
	assert.EqualError(t, err, "No available port")
	assert.Equal(t, -1, port)

	port, err = FindAvailablePort(30000, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 30001, port)
}
