package connection

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// Check if the TCP port available on localhost
func portAvailable(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))

	if err != nil {
		if strings.Index(err.Error(), "connection refused") > 0 {
			return true
		}
		return false
	}

	conn.Close()
	return false
}

// Get available TCP port on localhost by trying available ports in a range
func AvailablePort(start int, limit int) (int, error) {
	for i := start; i <= (start + limit); i++ {
		if portAvailable(i) {
			return i, nil
		}
	}
	return -1, errors.New("No available port")
}
