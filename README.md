# Pinger

![pinger.png](https://i.imgur.com/GIXWTf6.png)

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

## Libraries

- https://github.com/gdamore/tcell
- https://github.com/rivo/tview
- https://gitlab.com/tslocum/cview
- https://github.com/gcla/gowid
- https://github.com/digineo/go-ping
- https://github.com/sparrc/go-ping
- https://github.com/glinton/ping
- https://github.com/zyedidia/micro
- https://github.com/gcla/termshark

## Use

```
ping ya.ru | pings
pings ya.ru 192.168.31.1 google.com
```