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
	StartScan(ip string, f, l int, timeout time.Duration, s *MonitorScanner)
}

type MonitorScanner struct {
	openPorts []string
}

func (s *MonitorScanner) TCPScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			s.TCPScanPort(ip, port, timeout)
		}
		return
	}

	s.openPorts = append(s.openPorts, target)

	defer conn.Close()
}

func (s *MonitorScanner) StartWithMonitor(ip string, f, l int, timeout time.Duration) {
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
			s.TCPScanPort(ip, port, timeout)

			ps.lock.Lock()
			ps.available++
			ps.cond.Signal()
			ps.lock.Unlock()
		}(port)
	}
}

func (s *MonitorScanner) SendOpenPorts(session *discordgo.Session, channelID string, msgRef *discordgo.MessageReference) {

	openPortsStr := strings.Join(s.openPorts, "\n")

	embed := &discordgo.MessageEmbed{
		Title: "Open Ports",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Ports",
				Value:  openPortsStr,
				Inline: true,
			},
		},
	}

	session.ChannelMessageSendEmbedReply(channelID, embed, msgRef)
}

func (s *MonitorScanner) StartScan(ip string, f, l int, timeout time.Duration) {
	s.StartWithMonitor(ip, f, l, timeout)
}
