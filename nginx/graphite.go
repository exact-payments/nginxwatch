package nginx

import (
	"fmt"
	"net"
	"time"
)

func connectToGraphite(server string) (conn *net.UDPConn, addr *net.UDPAddr, err error) {
	if addr, err = net.ResolveUDPAddr("udp", server); err != nil {
		return nil, nil, err
	}

	if conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}); err != nil {
		return nil, nil, err
	}

	return conn, addr, nil
}

func writeData(name string, value float64, conn *net.UDPConn, addr *net.UDPAddr) (err error) {

	out := fmt.Sprintf("%s %f %d", name, value, time.Now().Unix())
	go conn.WriteToUDP([]byte(out), addr)

	return nil
}
