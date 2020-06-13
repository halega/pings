package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/digineo/go-ping"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const headerTmpl = "PING %s [%s] with %d bytes of data:"
const bodyTmpl = "%d bytes from %s (%s): icmp_seq=%d ttl=%d time=%d ms\n"
const summaryTmpl = "Packets: %d sent, %d received, %d lost (%d%% loss). RTT: min = %d ms, max = %d ms, avg = %d ms"

var host string

type stat struct {
	pktSent int
	pktLoss float64
	err     string
	last    time.Duration
	best    time.Duration
	worst   time.Duration
	mean    time.Duration
	stddev  time.Duration
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: pings <host>")
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

	ipaddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		panic(err)
	}

	pinger, err := ping.New("0.0.0.0", "")
	if err != nil {
		panic(err)
	}
	defer pinger.Close()

	go func() {
		for {
			rtt, err := pinger.Ping(ipaddr, 3*time.Second)
			var line string
			if err != nil {
				line = "Request timed out.\n"
			} else {
				line = fmt.Sprintf(bodyTmpl, pinger.PayloadSize(), host, ipaddr, 0, 0, rtt.Milliseconds())
			}
			body.Write([]byte(line))
			app.Draw()
			time.Sleep(1 * time.Second)
		}
	}()

	header.SetText(fmt.Sprintf(headerTmpl, host, ipaddr, pinger.PayloadSize()))
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// func onRecv(p *ping.Packet) {
// 	line := fmt.Sprintf(bodyTmpl, p.Nbytes, p.Addr, p.IPAddr, p.Seq, p.Ttl, p.Rtt.Milliseconds())
// 	body.Write([]byte(line))
// 	stat := pinger.Statistics()
// 	sumTxt := fmt.Sprintf(summaryTmpl, stat.PacketsSent, stat.PacketsRecv, stat.PacketsSent-stat.PacketsRecv,
// 		int64(math.Round(stat.PacketLoss)), stat.MinRtt.Milliseconds(), stat.MaxRtt.Milliseconds(),
// 		stat.AvgRtt.Milliseconds())
// 	summary.SetText(sumTxt)
// 	app.Draw()
// }
