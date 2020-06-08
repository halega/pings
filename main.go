package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/sparrc/go-ping"
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
	//header.SetBackgroundColor(tcell.ColorDarkBlue)
	body := tview.NewTextView()
	body.SetTextColor(tcell.ColorLightGray)
	summary := tview.NewTextView()
	//summary.SetBackgroundColor(tcell.ColorDarkBlue)
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(body, 1, 0, 1, 1, 0, 0, true).
		AddItem(summary, 2, 0, 1, 1, 0, 0, false)
	app := tview.NewApplication().
		EnableMouse(true).
		SetRoot(grid, true).
		SetFocus(body)

	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	pinger.Size = 32
	pinger.Count = -1
	pinger.Interval = 1 * time.Second
	pinger.Timeout = 9223372036854775807
	pinger.SetPrivileged(true)
	pinger.OnRecv = func(p *ping.Packet) {
		line := fmt.Sprintf(bodyTmpl, p.Nbytes, p.Addr, p.IPAddr, p.Seq, p.Ttl, p.Rtt.Milliseconds())
		body.Write([]byte(line))
		stat := pinger.Statistics()
		sumTxt := fmt.Sprintf(summaryTmpl, stat.PacketsSent, stat.PacketsRecv, stat.PacketsSent-stat.PacketsRecv,
			int64(math.Round(stat.PacketLoss)), stat.MinRtt.Milliseconds(), stat.MaxRtt.Milliseconds(),
			stat.AvgRtt.Milliseconds())
		summary.SetText(sumTxt)
		app.Draw()
	}
	go pinger.Run()

	header.SetText(fmt.Sprintf(headerTmpl, pinger.Addr(), pinger.IPAddr(), pinger.Size))
	if err := app.Run(); err != nil {
		panic(err)
	}
}
