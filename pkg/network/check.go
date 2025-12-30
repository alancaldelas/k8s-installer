package network

import (
	"fmt"
	"net"
	"time"
)

func CheckConnectivity(ip string, port int) error {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
