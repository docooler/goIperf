package Iperf

import (
		"log"
		"time"
		"fmt"
		)

const  queueLength = 1000
const  flushTimer  = 3
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
					fmt.Printf("%5d send %5d Bytes \n", count, send)
				}
				if recv >0 {
					fmt.Printf("%5d recv %5d Bytes %.2fM/s\n", count, recv,float64(recv>>17)/float64(hi.intervalTime))
				}
				// if send>0 || recv>0 {
				// 	fmt.Printf("%5d  %9d %6")
				// }
				if hi.expireTime>0 && count>=hi.expireTime {
					fmt.Printf("%s end test \n", hi.hostName)
					return 
				}
		}
		
	}
}

var stats *TrafficStats

func init() {
	stats = NewTrafficStats()
	go stats.StatsLoop()
}