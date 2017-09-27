package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/shared"
)

const (
	PORT_START = 29168
	PORT_LIMIT = 500
)

type Tunnel struct {
	TargetHost string
	TargetPort string
	Port       int
	SSHInfo    *shared.SSHInfo
	Config     *ssh.ClientConfig
	Client     *ssh.Client
	Listener   *net.TCPListener
}

func privateKeyPath() string {
	return os.Getenv("HOME") + "/.ssh/id_rsa"
}

func expandKeyPath(path string) string {
	home := os.Getenv("HOME")
	if home == "" {
		return path
	}
	return strings.Replace(path, "~", home, 1)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(buff)
}

func makeConfig(info *shared.SSHInfo) (*ssh.ClientConfig, error) {
	methods := []ssh.AuthMethod{}

	// Try to use user-provided key, fallback to system default key
	keyPath := info.Key
	if keyPath == "" {
		keyPath = privateKeyPath()
	} else {
		keyPath = expandKeyPath(keyPath)
	}

	if fileExists(keyPath) {
		key, err := parsePrivateKey(keyPath)
		if err != nil {
			return nil, err
		}

		methods = append(methods, ssh.PublicKeys(key))
	}

	methods = append(methods, ssh.Password(info.Password))

	cfg := &ssh.ClientConfig{
		User:    info.User,
		Auth:    methods,
		Timeout: time.Second * 10,
	}

	return cfg, nil
}

func (tunnel *Tunnel) sshEndpoint() string {
	return fmt.Sprintf("%s:%v", tunnel.SSHInfo.Host, tunnel.SSHInfo.Port)
}

func (tunnel *Tunnel) targetEndpoint() string {
	return fmt.Sprintf("%v:%v", tunnel.TargetHost, tunnel.TargetPort)
}

func (tunnel *Tunnel) copy(wg *sync.WaitGroup, writer, reader net.Conn) {
	defer wg.Done()
	if _, err := io.Copy(writer, reader); err != nil {
		log.Println("Tunnel copy error:", err)
	}
}

func (tunnel *Tunnel) handleConnection(local net.Conn) {
	remote, err := tunnel.Client.Dial("tcp", tunnel.targetEndpoint())
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go tunnel.copy(&wg, local, remote)
	go tunnel.copy(&wg, remote, local)

	wg.Wait()
	local.Close()
}

func (tunnel *Tunnel) Close() {
	if tunnel.Client != nil {
		tunnel.Client.Close()
	}

	if tunnel.Listener != nil {
		tunnel.Listener.Close()
	}
}

func (tunnel *Tunnel) Configure() error {
	config, err := makeConfig(tunnel.SSHInfo)
	if err != nil {
		return err
	}
	tunnel.Config = config

	client, err := ssh.Dial("tcp", tunnel.sshEndpoint(), config)
	if err != nil {
		return err
	}
	tunnel.Client = client

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", tunnel.Port))
	if err != nil {
		return err
	}
	tunnel.Listener = listener.(*net.TCPListener)

	return nil
}

func (tunnel *Tunnel) Start() {
	defer tunnel.Close()

	for {
		conn, err := tunnel.Listener.Accept()
		if err != nil {
			return
		}

		go tunnel.handleConnection(conn)
	}
}

func NewTunnel(sshInfo *shared.SSHInfo, dbUrl string) (*Tunnel, error) {
	uri, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}

	listenPort, err := connection.AvailablePort(PORT_START, PORT_LIMIT)
	if err != nil {
		return nil, err
	}

	chunks := strings.Split(uri.Host, ":")
	host := chunks[0]
	port := "5432"

	if len(chunks) == 2 {
		port = chunks[1]
	}

	tunnel := &Tunnel{
		Port:       listenPort,
		SSHInfo:    sshInfo,
		TargetHost: host,
		TargetPort: port,
	}

	return tunnel, nil
}
