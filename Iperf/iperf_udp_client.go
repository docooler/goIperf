package Iperf

import (
		"net"
		"strconv"
		)

type IperfUdpClient struct {
	srvIp       string
	srvPort     int 
	conn *net.UDPConn
}

func NewIperfUdpClient(srvAddress string, srvPortNum int) (client *IperfUdpClient, err error){
	udpAddr, err := net.ResolveUDPAddr("udp4", srvAddress + ":" + strconv.Itoa(srvPortNum))
	HandleError(err, 0, "NewIperfUdpClient ResolveUDPAddr")

	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	HandleError(err, 0, "NewIperfUdpClient DialUDP")

	client = &IperfUdpClient{
		srvIp   : srvAddress,
		srvPort : srvPortNum,
		conn    : udpConn,
	}
	return
}

func (client *IperfUdpClient )Write(buff []byte)(n int, err error) {
	n, err = client.Write(buff)
	return 
}

func (client *IperfUdpClient)Run() {
	for {
		client.conn.Write([]byte("testtesttesttesttesttesttesttest"))
	}
}

func (client *IperfUdpClient)Close() {
	client.conn.Close()
}