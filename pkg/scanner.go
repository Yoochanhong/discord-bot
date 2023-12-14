package pkg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	ip        string
	lock      sync.Mutex
	cond      *sync.Cond
	available int64
}

type Scanner interface {
	StartScan(ip string, f, l int, timeout time.Duration)
}

type MonitorScanner struct{}

func TCPScanPort(ip string, port int, timeout time.Duration, session *discordgo.Session, channelID string, msgRef *discordgo.MessageReference) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			TCPScanPort(ip, port, timeout, session, channelID, msgRef)
		}
		return
	}

	msg := fmt.Sprintf("%s 열린포트는 %s 입니다", ip, port)
	session.ChannelMessageSendReply(channelID, msg, msgRef)

	defer conn.Close()
}

func (s *MonitorScanner) StartWithMonitor(ip string, f, l int, timeout time.Duration, session *discordgo.Session, channelID string, ref *discordgo.MessageReference) {
	ps := &PortScanner{
		ip:        ip,
		available: 256,
	}
	ps.cond = sync.NewCond(&ps.lock)

	for port := f; port <= l; port++ {
		ps.lock.Lock()
		for ps.available <= 0 {
			ps.cond.Wait()
		}
		ps.available--
		ps.lock.Unlock()

		go func(port int) {
			TCPScanPort(ip, port, timeout, session, channelID, ref)

			ps.lock.Lock()
			ps.available++
			ps.cond.Signal()
			ps.lock.Unlock()
		}(port)
	}
}

func (s *MonitorScanner) StartScan(ip string, f, l int, timeout time.Duration, session *discordgo.Session, channelID string, msgRef *discordgo.MessageReference) {
	s.StartWithMonitor(ip, f, l, timeout, session, channelID, msgRef)
}
