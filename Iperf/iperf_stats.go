package Iperf

import (
		"log"
		"time"
		"fmt"
		"net/rpc"
		"github.com/nf/stat"
		)

const  queueLength = 1000
const  flushTimer  = 3

// const  addr = "127.0.0.1:8090"
var    webUI = 0


type HostInfo struct {
	hostName         string
	intervalTime int
    expireTime   int
}

type UpdateTraffic struct{
	hostName string
	SendBytes  uint64
	RecvBytes  uint64
}

type TrafficStats struct {
	hostmap *HostMap 
	hostChan chan HostInfo
	updateChan chan UpdateTraffic 
}

func  NewTrafficStats() (ts *TrafficStats){
     ts = &TrafficStats{
     	hostmap : NewHostMap(),
     	hostChan : make(chan HostInfo, queueLength),
     	updateChan : make(chan UpdateTraffic, queueLength),
     }
     return 
}
func (self *TrafficStats)StatsLoop() {
	
	for {
		select {
		case h := <-self.hostChan:
			log.Println("add a host")
			//send data to same host only need only statistics
			if self.hostmap.IsExist(h.hostName) {
				log.Println("we don't create other servers")
				return
			}
			log.Println("create one servers")
			self.hostmap.AddStatistics(h.hostName, h.intervalTime, h.expireTime)
			go self.statHandler(h)
		case u := <-self.updateChan:
			// log.Println("recv updateChan")
			self.hostmap.AddTrafficStat(u.hostName, u.SendBytes, u.RecvBytes)
		}
	}
}

func (self *TrafficStats)statHandler(hi HostInfo) {
	t := time.NewTicker(time.Duration(hi.intervalTime) * time.Second)
	count := 0
	fmt.Printf("times SendBytes SendRate RecvBytes RecvRate\n")
	for  {
		select {
			case <-t.C:
				count ++
				send, recv := self.hostmap.GetTrafficStat(hi.hostName)
				if send >0 {
					fmt.Printf("%5d send %5d Bytes %.2fM\n", count, send, float64(send>>17)/float64(hi.intervalTime))
					update(hi.hostName, "send", send>>7/uint64(hi.intervalTime))
				}
				if recv >0 {
					fmt.Printf("%5d recv %5d Bytes %.2fM/s\n", count, recv,float64(recv>>17)/float64(hi.intervalTime))
					update(hi.hostName, "recv", recv>>7/uint64(hi.intervalTime))
				}
				

				if hi.expireTime>0 && count>=hi.expireTime {
					fmt.Printf("%s end test \n", hi.hostName)
					return 
				}
		}
		
	}
}

var stats *TrafficStats
var client   *rpc.Client 

func InitStats(addr string, webui bool) {
	stats = NewTrafficStats()
	if webui {
		webUI = 1
	}

	if webUI == 1 {
		var err error
		client, err = rpc.DialHTTP("tcp", addr)
		if err != nil {
			fmt.Println(err)
		}
	}
	
	go stats.StatsLoop()
}


func update(host string, dataType string, value uint64) {
	if webUI != 1 {
		return
	}
	if client != nil {
		err := client.Call("Server.Update", &stat.Point{host, dataType, int64(value)}, &struct{}{})
		if err != nil {
			fmt.Println("stat update:", err)
		}
	}
	
}

