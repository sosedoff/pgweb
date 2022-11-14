package connection

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// IsPortAvailable returns true if there's no listeners on a given port
func IsPortAvailable(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		return strings.Index(err.Error(), "connection refused") > 0
	}

	conn.Close()
	return false
}

// FindAvailablePort returns the first available TCP port in the range
func FindAvailablePort(start int, limit int) (int, error) {
	for i := start; i <= (start + limit); i++ {
		if IsPortAvailable(i) {
			return i, nil
		}
	}
	return -1, errors.New("No available port")
}
