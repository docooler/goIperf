package Iperf

import (
		"net"
		"strconv"
		)

type IperfUdpClient struct {
	srvIp       string
	srvPort     int 
	host        string
	conn *net.UDPConn
}

func NewIperfUdpClient(srvAddress string, srvPortNum int) (client *IperfUdpClient, err error){
	host := srvAddress + ":" + strconv.Itoa(srvPortNum) 
	udpAddr, err := net.ResolveUDPAddr("udp4", host)
	HandleError(err, 0, "NewIperfUdpClient ResolveUDPAddr")

	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	HandleError(err, 0, "NewIperfUdpClient DialUDP")

	client = &IperfUdpClient{
		srvIp   : srvAddress,
		srvPort : srvPortNum,
		host    : srvAddress+":send",
		conn    : udpConn,
	}
	return
}

func (client *IperfUdpClient )Write(buff []byte)(n int, err error) {
	n, err = client.Write(buff)
	return 
}

func (client *IperfUdpClient)Run() {
	
	if stats.hostChan != nil {
			hi := HostInfo{
			hostName  : client.host,
			intervalTime : 3,
			expireTime : 36000,
			
		}
		stats.hostChan <-hi
	}

	for {

		content := "testtesttesttesttesttesttesttest"
		length := len(content)

		client.conn.Write([]byte(content))

		

		if stats.updateChan != nil {
			u := UpdateTraffic{
				hostName : client.host,
				SendBytes : uint64(length),
				RecvBytes : 0,
			}
			stats.updateChan <-u
		}

	}
}

func (client *IperfUdpClient)Close() {
	client.conn.Close()
}