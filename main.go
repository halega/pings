package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/digineo/go-ping"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type stat struct {
	mu      sync.Mutex
	sent    int
	lost    int
	loss    float64
	lastErr error
	lastRTT time.Duration
	min     time.Duration
	max     time.Duration
	avg     time.Duration
	total   time.Duration
	start   time.Time
	uptime  time.Duration
	timeout time.Duration
}

func newStat(timeout time.Duration) *stat {
	return &stat{
		start:   time.Now(),
		min:     timeout + time.Hour,
		timeout: timeout,
	}
}

func (s *stat) update(rtt time.Duration, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.start.IsZero() {
		s.uptime = time.Since(s.start)
	}

	s.sent++
	s.lastRTT = rtt
	s.lastErr = err

	if err != nil {
		s.lost++
	} else {
		if s.min > rtt {
			s.min = rtt
		}
		if s.max < rtt {
			s.max = rtt
		}
		s.total += rtt
		s.avg = s.total / time.Duration(s.sent-s.lost)
	}
	s.loss = float64(s.lost) / float64(s.sent) * 100
}

func (s *stat) summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	sum := fmt.Sprintf("Packets: %d sent, %d received, %d lost (%.0f%% loss).",
		s.sent, s.sent-s.lost, s.lost, s.loss)
	if s.sent != s.lost {
		sum += fmt.Sprintf(" RTT: min = %d ms, max = %d ms, avg = %d ms.",
			s.min.Milliseconds(), s.max.Milliseconds(), s.avg.Milliseconds())
	}
	if s.uptime != 0 {
		sum += fmt.Sprintf(" Uptime: %v", s.uptime.Round(time.Second))
	}
	return sum
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
	ui.summary.SetText(s.summary())
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
		fmt.Println(err)
		return
	}

	pinger, err := ping.New("0.0.0.0", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pinger.Close()

	ui := newUIApp(host, ipAddr.String(), pinger.PayloadSize())
	s := newStat(3 * time.Second)

	go func() {
		for {
			s.update(pinger.Ping(ipAddr, s.timeout))
			ui.update(s)
			time.Sleep(1 * time.Second)
		}
	}()

	if err := ui.app.Run(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("PING %s (%s) %d bytes of data.\n", host, ipAddr, pinger.PayloadSize())
	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Println(s.summary())
}
