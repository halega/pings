package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const headerTmpl = "PING %s [%s] with %d bytes of data:"
const bodyTmpl = "%d bytes from %s (%s): icmp_seq=%d ttl=%d time=%d ms\n"
const summaryTmpl = "Packets: %d sent, %d received, %d lost (%d%% loss). RTT: min = %d ms, max = %d ms, avg = %d ms"

var host string

func main() {
	flag.Usage = func() {
		fmt.Println("ping destination")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		return
	}
	host = flag.Arg(0)

	header := tview.NewTextView()
	body := tview.NewTextView()
	body.SetTextColor(tcell.ColorLightGray)
	summary := tview.NewTextView()
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(body, 1, 0, 1, 1, 0, 0, true).
		AddItem(summary, 2, 0, 1, 1, 0, 0, false)
	app := tview.NewApplication().
		EnableMouse(true).
		SetRoot(grid, true).
		SetFocus(body)

	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		panic(err)
	}

	pinger, err := ping.New("0.0.0.0", "")
	if err != nil {
		panic(err)
	}
	defer pinger.Close()

	go ping(pinger)

	header.SetText(fmt.Sprintf(headerTmpl, pinger.Addr(), pinger.IPAddr(), pinger.Size))
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func ping(pinger *ping.Pinger, addr *net.IPAddr) {
	for {
		pinger.Ping(addr, 5*time.Second)
	}
}

func onRecv(p *ping.Packet) {
	line := fmt.Sprintf(bodyTmpl, p.Nbytes, p.Addr, p.IPAddr, p.Seq, p.Ttl, p.Rtt.Milliseconds())
	body.Write([]byte(line))
	stat := pinger.Statistics()
	sumTxt := fmt.Sprintf(summaryTmpl, stat.PacketsSent, stat.PacketsRecv, stat.PacketsSent-stat.PacketsRecv,
		int64(math.Round(stat.PacketLoss)), stat.MinRtt.Milliseconds(), stat.MaxRtt.Milliseconds(),
		stat.AvgRtt.Milliseconds())
	summary.SetText(sumTxt)
	app.Draw()
}
