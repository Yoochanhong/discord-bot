package pkg

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
	"time"
)

func MessageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {

	if s.State.User.ID == msg.Author.ID {
		return
	}

	words := strings.Fields(strings.ToLower(msg.Content))

	if words[0] == "!ping" {
		log.Printf("%s이 %s라고 말함", msg.Author.Username, msg.Content)
		s.ChannelMessageSend(msg.ChannelID, "pong!")
	}

	if words[0] == "!pong" {
		log.Printf("%s이 %s라고 말함", msg.Author.Username, msg.Content)
		s.ChannelMessageSend(msg.ChannelID, "ping!")
	}

	r1, _ := regexp.Compile("[니닉닠늬뉘ㄴ][\\d\\s\\W\\t\\r\\n ]*([ㄱ거겨가]+|그롱)")
	r2, _ := regexp.Compile("(?i)(n[ei]+g+[ea]*r*|n word)")
	if r1.MatchString(msg.Content) || r2.MatchString(msg.Content) {
		log.Printf("%s이 %s라고 말함", msg.Author.Username, msg.Content)
		s.ChannelMessageSend(msg.ChannelID, "목화재배!")
	}

	if words[0] == "!port" {
		log.Printf("%s의 %s 포트 스캐닝 요청", msg.Author.Username, words[1])
		scanner := MonitorScanner{}
		scanner.StartScan(words[1], 1, 65535, 500*time.Millisecond, s, msg.ChannelID, msg.Reference())
	}
}
