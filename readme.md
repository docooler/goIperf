Welcome to the goIperf wiki!
this is iperf implemented by golang. good luck!

inspired from https://github.com/esnet/iperf

#Why goIperf? because iperf need both sides have a public IP address(or they can find each other directly). But when some of tester Ip is a private address. iperf will not work. goIperf is coming!

Let UE acts as a hot spot and PC_B connects to UE via Wifi. PC_B requests PC_A to send data to UE by goIperf. By this way, UE control is not mandatory and there would be no dependency on UE type. 

![net](https://github.com/docooler/goIperf/blob/master/net.png)

Monitor traffic by graph on webpage on whichever PC that could access ordered website. And we could check traffic of several UEs at the same time.
![monitor](https://github.com/docooler/goIperf/blob/master/Picture1.png)




