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

type stat struct {
	sent    int
	lost    int
	loss    float64
	lastErr error
	lastRTT time.Duration
	min     time.Duration
	max     time.Duration
	avg     time.Duration
	mean    time.Duration
	stddev  time.Duration
	rtts    []time.Duration
	size    uint16
	total   time.Duration
}

func (s *stat) addRTT(rtt time.Duration, err error) {
	s.sent++
	s.lastRTT = rtt
	s.lastErr = err

	if err != nil {
		s.lost++
	}
	s.loss = float64(s.lost) / float64(s.sent) * 100
	if err != nil {
		return
	}

	if rtt < s.min {
		s.min = rtt
	}
	if rtt > s.max {
		s.max = rtt
	}

	s.rtts = append(s.rtts, rtt)
	s.total += rtt
	s.avg = s.total / time.Duration(len(s.rtts))
}

type uiApp struct {
	header  *tview.TextView
	body    *tview.TextView
	summary *tview.TextView
	app     *tview.Application
}

func newUIApp(host, ipaddr string, payloadSize uint16) *uiApp {
	ui := &uiApp{
		header:  tview.NewTextView(),
		body:    tview.NewTextView(),
		summary: tview.NewTextView(),
	}

	ui.header.SetText(fmt.Sprintf("PING %s [%s] with %d bytes of data:", host, ipaddr, payloadSize))
	ui.body.SetTextColor(tcell.ColorLightGray)

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		AddItem(ui.header, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.body, 1, 0, 1, 1, 0, 0, true).
		AddItem(ui.summary, 2, 0, 1, 1, 0, 0, false)
	ui.app = tview.NewApplication().
		EnableMouse(true).
		SetRoot(grid, true).
		SetFocus(ui.body)

	return ui
}

func (ui *uiApp) update(s *stat) {
	bodyLine := ""
	if s.lastErr != nil {
		bodyLine = fmt.Sprintln(s.lastErr)
	} else {
		bodyLine = fmt.Sprintf("icmp_seq=%d time=%d ms\n", s.sent, s.lastRTT.Milliseconds())
	}
	ui.body.Write([]byte(bodyLine))

	sumLine := fmt.Sprintf("Packets: %d sent, %d received, %d lost (%.f%% loss). RTT: min = %d ms, max = %d ms, avg = %d ms",
		s.sent, s.sent-s.lost, s.lost, s.loss, s.min.Milliseconds(), s.max.Milliseconds(), s.avg.Milliseconds())
	ui.summary.SetText(sumLine)

	ui.app.Draw()
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: pings <host>")
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	host := flag.Arg(0)

	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		panic(err)
	}

	pinger, err := ping.New("0.0.0.0", "")
	if err != nil {
		panic(err)
	}
	defer pinger.Close()

	ui := newUIApp(host, ipAddr.String(), pinger.PayloadSize())
	s := &stat{}

	go func() {
		for {
			s.addRTT(pinger.Ping(ipAddr, 3*time.Second))
			ui.update(s)
			time.Sleep(1 * time.Second)
		}
	}()

	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}
