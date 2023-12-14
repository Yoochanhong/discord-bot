package pkg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"regexp"
	"strings"
)

var gif = os.Getenv("GIF")

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
		content := fmt.Sprintf("목화재배!\n %s", gif)
		s.ChannelMessageSend(msg.ChannelID, content)
	}
}
