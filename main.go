package main 

import (
		// "log"
		"fmt"
		"goIperf/Iperf"
		)

func main() {
	fmt.Println("start main")
	test_udp()
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

	udpClient.Run()
}

func start_upd_server(ipAddress string, port int){
	udpSrv, _ := Iperf.NewIperfUdpServer(ipAddress, port)
	defer udpSrv.Close()

	udpSrv.Run()
}