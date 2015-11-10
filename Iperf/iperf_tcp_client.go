package Iperf

import (
		"net"
		"strconv"
		)

type IperfTcpClient struct {
	ip       string
	port     int 
	host     string 
	start    int
	conn      net.Conn
	Model    int
}

func NewIperfTcpClient(ipAddress string ,portNum int, model int )(client *IperfTcpClient, err error) {
	host := ipAddress + ":" + strconv.Itoa(portNum)
	conn, err := net.Dial("tcp", host)
	HandleError(err, 1, "NewIperfTcpClient connect server failed")
	
	client = &IperfTcpClient{
		ip : ipAddress,
		port : portNum,
		host : host,
		conn : conn,
		Model: model,
	}
	return
}

func (this *IperfTcpClient)Run() {
	switch this.Model{
	case DOWNLOAD:
		//send handshare to server
		//分号后面的内容留作协议扩充
		this.conn.Write([]byte("download:"))
		this.loopRecv()
	case UPLOAD:
		this.conn.Write([]byte("upload:"))
		this.loopSend()
	}
}

func (this *IperfTcpClient)loopSend() {
	sendData := "send Data Test Client to server DUMMY DUMMY"
	dataLen  := len(sendData)
	for {
		this.conn.Write([]byte(sendData))
		this.anlyzeMessage([]byte(sendData), int64(dataLen))
	}
}

func (this *IperfTcpClient)loopRecv() {
	dataBuf := make([]byte, 1024)
	for  {
		n, err := this.conn.Read(dataBuf)
		if err != nil {
			HandleError(err, 0, "loopRecv ReadFromTCP")
			continue
		}
		
		this.anlyzeMessage(dataBuf, int64(n))
	}
}

func (this * IperfTcpClient)anlyzeMessage(data []byte, length int64) {

	model := ""
	if this.Model == DOWNLOAD {
		model += "recv"
	} else{
		model += "send"
	}

	if this.start == 0{
		this.start += 1
		hi := HostInfo{
			hostName  : this.host + model,
			intervalTime : 3,
			expireTime : 36000,
			
		}

		if stats.hostChan != nil {
			stats.hostChan <-hi
		}

		return
	} else {
		var sendBytes uint64
		var recvBytes uint64
		
		if this.Model == DOWNLOAD {
			recvBytes = uint64(length)
		} else {
			sendBytes = uint64(length)
		}
		u := UpdateTraffic{
			hostName : this.host + model,
			SendBytes : sendBytes,
			RecvBytes : recvBytes,
		}

		if stats.updateChan != nil {
			
			stats.updateChan <-u
		}
		return 
	}
}

func (this * IperfTcpClient)Close() {
	this.conn.Close()
}