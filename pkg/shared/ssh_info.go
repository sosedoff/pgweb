package shared

import (
	"fmt"
)

type SSHInfo struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Key      string `json:"key,omitempty"`
}

func (info SSHInfo) String() string {
	return fmt.Sprintf("%s@%s:%s", info.User, info.Host, info.Port)
}
