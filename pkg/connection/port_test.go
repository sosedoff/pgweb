package connection

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_portAvailable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIXME")
	}

	assert.Equal(t, true, portAvailable(8081))

	serv, err := net.Listen("tcp", "127.0.0.1:8081")
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

	assert.Equal(t, false, portAvailable(8081))
	assert.Equal(t, true, portAvailable(8082))
}

func Test_getAvailablePort(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIXME")
	}

	port, err := AvailablePort(8081, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 8081, port)

	serv, err := net.Listen("tcp", "127.0.0.1:8081")
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

	port, err = AvailablePort(8081, 0)
	assert.EqualError(t, err, "No available port")
	assert.Equal(t, -1, port)

	port, err = AvailablePort(8081, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 8082, port)
}
