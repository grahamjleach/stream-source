package stream

import (
	"net"
	"os/exec"

	"github.com/stunndard/goicy/config"
	"github.com/stunndard/goicy/network"
)

type StreamClient struct {
	sock   net.Conn
	config *config.Config
	cmd    *exec.Cmd
}

func NewClient(icyConfig *config.Config) (client *StreamClient, err error) {
	sock, err := network.ConnectServer(icyConfig.Host, icyConfig.Port, 0, 0, 0)
	if err != nil {
		return
	}

	client = &StreamClient{
		sock:   sock,
		config: icyConfig,
	}

	return
}
