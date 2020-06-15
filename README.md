# Pinger

![pinger.png](https://i.imgur.com/g3AIoq0.png)

## Installation & Usage

```
go get -u github.com/halega/pings
~/go/bin/pings www.google.com
```

## Libraries

- https://github.com/gdamore/tcell
- https://github.com/rivo/tview
- https://github.com/digineo/go-ping


## Research

### Libraries & Projects

- https://github.com/gdamore/tcell
- https://github.com/rivo/tview
  - https://rocketnine.space/post/tview-and-you/
  - https://flak.tedunangst.com/post/package-of-the-moment-tview-and-tcell
- https://gitlab.com/tslocum/cview
- https://github.com/gcla/gowid
- https://github.com/digineo/go-ping
- https://github.com/sparrc/go-ping
- https://github.com/glinton/ping
- https://github.com/zyedidia/micro
- https://github.com/gcla/termshark

### sparrc/go-ping API

```go
pinger, _ := ping.NewPinger(host)
defer pinger.Stop()

pinger.OnRecv = func(pkt *ping.Packet) {
	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
}
pinger.OnFinish = func(stats *ping.Statistics) {
  fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
  fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
    stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
  fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
    stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
}
pinger.Count = *count
pinger.Interval = *interval
pinger.Timeout = *timeout
pinger.SetPrivileged(*privileged)

fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
pinger.Run()
```

### digineo/go-ping API

```go
ip, _ := net.ResolveIPAddr("ip4", host)
pinger, _ := ping.New("0.0.0.0", "")
defer pinger.Close()
rtt, err := pinger.Ping(ip, 3*time.Second)
```

### glinton/ping API

```go
ipAddr, _ := net.ResolveIPAddr("ip4", host)
data := bytes.Repeat([]byte{1}, 56)
c := &ping.Client{}
req := ping.Request{
  Dst:  net.ParseIP(host.String()),
  Src:  net.ParseIP(getAddr(*iface)),
  Data: data,
}
ctx, _ := context.WithTimeout(context.Background(), time.Duration(*timeout*float64(time.Second)))
resp, err := c.Do(ctx, &req)
if err != nil {
  fmt.Println("failed to ping:", err)
  return
}
```

### System's Ping Outputs

```
root@longt:~# ping ya.ru
PING ya.ru(ya.ru (2a02:6b8::2:242)) 56 data bytes
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=1 ttl=52 time=21.5 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=2 ttl=52 time=21.5 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=3 ttl=52 time=21.4 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=4 ttl=52 time=21.6 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=5 ttl=52 time=21.5 ms
^C
--- ya.ru ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4005ms
rtt min/avg/max/mdev = 21.485/21.570/21.695/0.199 ms
```

Summary bar:

```
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=1 ttl=52 time=21.5 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=2 ttl=52 time=21.5 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=3 ttl=52 time=21.4 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=4 ttl=52 time=21.6 ms
64 bytes from ya.ru (2a02:6b8::2:242): icmp_seq=5 ttl=52 time=21.5 ms
.............................

-------------------------------------------------------------------------------
5 packets transmitted, 5 received, 0% packet loss, time 4005ms
rtt min/avg/max/mdev = 21.485/21.570/21.695/0.199 ms
```