package Iperf

import (
		 "net"
		 "strconv"
		 "log"
		)

type IperfTcpServer struct {
	ip       string
	port     int 
	host     string 
	start    int
	listener *net.TCPListener
	conn      *net.TCPConn
	Model    int
}

const DOWNLOAD = 1
const UPLOAD   = 2


func NewIperfTcpServer(ipAddr string, portNum int) (server *IperfTcpServer, err error){
	listen,err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ipAddr), portNum, ""})
	HandleError(err, 1,"NewIperfTcpServer ListenTCP")
	log.Printf("listen on TCP %s:%d\n", ipAddr, portNum)

	server = &IperfTcpServer{
		ip    : ipAddr,
		port  : portNum,
		host  : ipAddr + strconv.Itoa(portNum),
		listener : listen,
	}
	return 
}

func (this *IperfTcpServer)Run() {
	for  {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
			HandleError(err, 0, "(this *IperfTcpServer)Run AcceptTCP")
			continue
		}

		go this.HandlerMessage(conn)
		
	}
}

func (this *IperfTcpServer)HandlerMessage(conn *net.TCPConn) {
	data := make([]byte, 1024)
	for {
		 n, err := conn.Read(data)
		 HandleError(err, 0, "HandlerMessage ReadFromTCP")
         if n > 0 {
         	this.anlyzeMessage(data, n)
         }

         switch this.Model {
         	case DOWNLOAD:
         		//SEND DATA
         		this.loopSend(conn)
         	case UPLOAD:
         		//RECV DATA
         		this.loopRecv(conn)
         }
	}
}

func (this * IperfTcpServer)loopSend(conn *net.TCPConn){
	defer conn.Close()
	sendData := "send Data Test Server to Client DUMMY DUMMY"
	
	errCount := 0
	for {
		n, err := conn.Write([]byte(sendData))
		if err != nil {
			errCount += 1
			if errCount > 5 {
				HandleError(err, 0 , "loopSend error 5 times")
				break 
			}
		} else {
			errCount = 0
		}
		// log.Println("send ok")
		this.anlyzeMessage([]byte(sendData), n)
	}
}

func (this *IperfTcpServer)loopRecv(conn *net.TCPConn) {
	defer conn.Close()
	dataBuf := make([]byte, 1024)
	for  {
		n, err := this.conn.Read(dataBuf)
		if err != nil {
			HandleError(err, 0, "loopRecv ReadFromTCP")
			continue
		}
		
		this.anlyzeMessage(dataBuf, n)
	}
}

func  (this *IperfTcpServer)anlyzeMessage(data []byte, length int){
	if this.start == 0{
		this.start += 1

		msgs := splitMsg(data, length)
		model := ""

		switch msgs[0] {
		case "download":
			this.Model = DOWNLOAD
			model += "send"
		case "upload":
			this.Model = UPLOAD
			model += "recv"
		default :
			this.Model = UPLOAD
		}

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
		model := ""
		if this.Model == DOWNLOAD {
			sendBytes = uint64(length)
			model += "send"
		} else {
			recvBytes = uint64(length)
			model += "recv"
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


//消息解析，[]byte -> []string
func splitMsg(buff []byte, len int) ([]string) {
	analMsg := make([]string, 0)
	strNow := ""
	for i := 0; i < len; i++ {
		if string(buff[i:i + 1]) == ":" {
			analMsg = append(analMsg, strNow)
			strNow = ""
		} else {
			strNow += string(buff[i:i + 1])
		}
	}
	analMsg = append(analMsg, strNow)
	return analMsg
}

func (this *IperfTcpServer)Close() {
	defer this.conn.Close()
}
