package connection

import (
	"errors"
	"fmt"
	"net"
)

// IsPortAvailable returns true if there's no listeners on a given port
func IsPortAvailable(port int) bool {
	// Try to bind the port instead of dialing it and matching the error text:
	// the "connection refused" message is OS-specific (Windows reports "actively
	// refused it"), which made every port look busy there.
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		return false
	}

	listener.Close()
	return true
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
