package Iperf

import (
		"net"
		"strconv"
		"log"
		)

type IperfUdpServer struct {
	ip       string
	port     int 
	host     string 
	start    int
	listener *net.UDPConn
}


func NewIperfUdpServer(ipAddr string, portNum int) (server *IperfUdpServer, err error){
    
    host := ipAddr + ":" + strconv.Itoa(portNum)
	udpAddr, err := net.ResolveUDPAddr("udp4", host)
	HandleError(err, 0, "NewIperfUdpServer ResolveUDPAddr")
    
    udpListener, err := net.ListenUDP("udp4", udpAddr)
    HandleError(err, 0, "NewIperfUdpServer udpListener")
    log.Printf("listen on %s:%d\n", ipAddr, portNum)

    server = &IperfUdpServer{
    	ip   : ipAddr,
    	port : portNum,
    	host : host,
    	listener : udpListener,
    }
    return
}

func (m *IperfUdpServer)Run() {
	m.HandlerMessage()
}
func (m *IperfUdpServer)HandlerMessage()  {
	buff := make([]byte, 1024)
	for  {
		n, addr, err := m.listener.ReadFromUDP(buff)
		HandleError(err, 0, "HandlerMessage ReadFromUDP")
		if n > 0{
			// log.Printf("recevie from %s : %d byte \n", addr,  n)
			m.anlyzeMessage(buff, uint64(n))
			addr = addr
			
		}
	}
}

func  (m *IperfUdpServer)anlyzeMessage(buff []byte, length uint64) {
	if m.start == 0 {
		m.start += 1
		hi := HostInfo{
			hostName  : m.host,
			intervalTime : 3,
			expireTime : 36000,
			
		}
		if stats.hostChan != nil {
			stats.hostChan <-hi
		}
		return
	} else {
		u := UpdateTraffic{
			hostName : m.host,
			SendBytes : 0,
			RecvBytes : length,
		}

		if stats.updateChan != nil {
			
			stats.updateChan <-u
		}
		return 
	}
	
}

func (m *IperfUdpServer)Close() {
	m.listener.Close()
}