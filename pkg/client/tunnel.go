package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/sosedoff/pgweb/pkg/connection"
)

const (
	PORT_START = 29168
	PORT_LIMIT = 500
)

type Tunnel struct {
	TargetHost string
	TargetPort string

	SshHost     string
	SshPort     string
	SshUser     string
	SshPassword string
	SshKey      string

	Config *ssh.ClientConfig
	Client *ssh.Client
}

func privateKeyPath() string {
	return os.Getenv("HOME") + "/.ssh/id_rsa"
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(buff)
}

func makeConfig(user, password, keyPath string) (*ssh.ClientConfig, error) {
	methods := []ssh.AuthMethod{}

	if keyPath != "" {
		key, err := parsePrivateKey(keyPath)
		if err != nil {
			return nil, err
		}

		methods = append(methods, ssh.PublicKeys(key))
	}

	methods = append(methods, ssh.Password(password))

	return &ssh.ClientConfig{User: user, Auth: methods}, nil
}

func (tunnel *Tunnel) sshEndpoint() string {
	return fmt.Sprintf("%s:%v", tunnel.SshHost, tunnel.SshPort)
}

func (tunnel *Tunnel) targetEndpoint() string {
	return fmt.Sprintf("%v:%v", tunnel.TargetHost, tunnel.TargetPort)
}

func (tunnel *Tunnel) Start() error {
	config, err := makeConfig(tunnel.SshUser, tunnel.SshPassword, tunnel.SshKey)
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", tunnel.sshEndpoint(), config)
	if err != nil {
		return err
	}
	defer client.Close()

	port, err := connection.AvailablePort(PORT_START, PORT_LIMIT)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go tunnel.handleConnection(conn, client)
	}
}

func (tunnel *Tunnel) copy(wg *sync.WaitGroup, writer, reader net.Conn) {
	defer wg.Done()
	if _, err := io.Copy(writer, reader); err != nil {
		log.Println("Tunnel copy error:", err)
	}
}

func (tunnel *Tunnel) handleConnection(local net.Conn, sshClient *ssh.Client) {
	remote, err := sshClient.Dial("tcp", tunnel.targetEndpoint())
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go tunnel.copy(&wg, local, remote)
	go tunnel.copy(&wg, remote, local)

	wg.Wait()
}
