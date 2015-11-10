package Iperf

import (
		// "time"
		"sync"
		"fmt"
		)

type Statistics struct {
	TotalBytesReceived    uint64
    IntervalBytesReceived uint64
    TotalBytesSend        uint64
    IntervalBytesSend     uint64
	TotalTime     uint64
	IntervalTime  int
	ExpireTime    int
}

type HostMap struct {
	lock sync.RWMutex
	LengthLimit int //limit the HostMap's length if it's equal to 0 there's no limit
	hostmap map[string]*Statistics
}

func NewHostMap() (hm *HostMap){
	hm = &HostMap{
		hostmap : make(map[string]*Statistics),
	}
	return
}
func (self *HostMap)AddStatistics(host string, interval,expire int) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if stat, ok := self.hostmap[host]; ok {
		stat.TotalBytesReceived = 0
		stat.IntervalBytesReceived = 0
		stat.TotalBytesSend = 0
		stat.TotalTime = 0
		stat.IntervalTime = interval
		stat.ExpireTime = expire
		
	} else {
		if self.LengthLimit > 0 && self.LengthLimit<= len(self.hostmap) {
			return
		}
		nb := &Statistics{
			ExpireTime : expire,
			IntervalTime : interval,
		}

		self.hostmap[host] = nb
	}
}

func (self * HostMap)AddTrafficStat(host string, intervalBytesSend, intervalBytesRecv uint64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if stat, ok := self.hostmap[host]; ok {
		stat.IntervalBytesReceived += intervalBytesRecv
		stat.TotalBytesReceived    += intervalBytesRecv
		stat.IntervalBytesSend     += intervalBytesSend
		stat.TotalBytesSend        += intervalBytesSend
	}else {
		fmt.Printf("%s not in stat map please add it at first\n", host)
	}
}

func (self *HostMap)GetTrafficStat(host string)(sendBytes, recvBytes uint64){
	self.lock.Lock()
	defer self.lock.Unlock()
    
    if stat, ok := self.hostmap[host]; ok {
    	sendBytes = stat.IntervalBytesSend
    	recvBytes = stat.IntervalBytesReceived

    	stat.IntervalBytesReceived = 0
    	stat.IntervalBytesSend = 0
    } else {
    	fmt.Printf("%s not in stat map please add it at first\n", host)
    }
    return 
}


