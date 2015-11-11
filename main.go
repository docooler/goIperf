package main 

import (
		// "log"
		// "fmt"
		"goIperf/Iperf"
		"flag"
		"sync"
		)



var (
	port       = flag.Int("p", 9005, " server port to listen on/connect to")
	hostAddress = flag.String("c", "", "run in client mode, connecting to <host>")
	serverMode = flag.Bool("s", false, "run in server mode") 
	udpMode    = flag.Bool("u", false, "use UDP rather than TCP")
	bind       = flag.String("B", "", "bind to <host>, an interface or multicast address")
	download   = flag.String("mode", "download", "upload package to server or download package from server only for client")
	stats      = flag.String("stats", "", "stats webUI server 127.0.0.1:8090")
	pthreadNum = flag.Int("P", 1, " number of parallel client threads to run" )
	mergeOut   = flag.Bool("m", false, "merge mutil client output")
)


var wg sync.WaitGroup

func main() {
	flag.Parse()

	
	ipAddress := ""

	if *bind != "" {
		ipAddress += *bind
	}

	
	portNum := *port

    //enable webui
	webui := false
	if *stats != "" {
		webui = true
	}
    //init statistics
	Iperf.InitStats(*stats, webui)
	
	if *udpMode {
		if *hostAddress != "" {
			start_mutil_udp_client(*pthreadNum, *hostAddress, portNum)
			
		}else{

			if *serverMode == false {
				flag.Usage()
				return
			}
			
			start_upd_server(ipAddress, portNum)
		}
	} else {
		if *hostAddress != "" {
			mode := Iperf.DOWNLOAD
			switch *download{
			case "download":
				mode = Iperf.DOWNLOAD
			case "upload":
				mode = Iperf.UPLOAD
			}
			start_mutil_tcp_client(*pthreadNum, *hostAddress, portNum, mode)
		} else{
			if *serverMode == false {
				flag.Usage()
				return
			}
			start_tcp_server(ipAddress, portNum)
		}
	}

	// test_udp()
	// test_tcp()
}



func test_udp(){
	srvAddress := "127.0.0.1"
	srvPort    := 5503
	go start_upd_server(srvAddress, srvPort)

	start_upd_client(srvAddress, srvPort)
}

func start_upd_client(ipAddress string, port int) {
	udpClient,_ := Iperf.NewIperfUdpClient(ipAddress, port)
	defer udpClient.Close()
    defer wg.Add(-1)
	udpClient.Run()
}

func start_mutil_udp_client(threadNum int, ipAddress string, port int) {
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		start_upd_client(ipAddress,port)
	}
	wg.Wait()
}
func start_upd_server(ipAddress string, port int){
	udpSrv, _ := Iperf.NewIperfUdpServer(ipAddress, port)
	defer udpSrv.Close()
	defer wg.Add(-1)

	udpSrv.Run()
}

func test_tcp(){
	srvAddress := "127.0.0.1"
	srvPort    := 5009

	go start_tcp_server(srvAddress, srvPort)

	start_tcp_client(srvAddress, srvPort, Iperf.DOWNLOAD, 0)
}

func start_mutil_tcp_client(threadNum int, ipAddress string, port int, model int) {
	clientSerails := 0 
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		if *mergeOut == false {
			clientSerails += 1
		}
		go start_tcp_client(ipAddress, port, model, clientSerails)
	}
	wg.Wait()
}

func start_tcp_client(ipAddress string, port int, model int, serials int ) {
	tcpClient, _ := Iperf.NewIperfTcpClient(ipAddress, port, model, serials)
	defer tcpClient.Close()
	defer wg.Add(-1)
	tcpClient.Run()
}

func start_tcp_server(ipAddress string, port int) {
	tcpSrv,_ := Iperf.NewIperfTcpServer(ipAddress, port)
	defer tcpSrv.Close()
	tcpSrv.Run()
}