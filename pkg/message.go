package pkg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"image/color"
	"image/png"
	"log"
	"os"
	"regexp"
	"strings"
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
		f, _ := os.Open("./image/dallas.gif")
		s.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
			Content: "목화재배!",
			File: &discordgo.File{
				Name:        "cotton.gif",
				ContentType: "image/gif",
				Reader:      f,
			},
		})
	}

	if words[0] == "!지지관계" {
		var username string
		if words[1] == "김현석" || words[1] == "khs0826" || words[1] == "pokabook" || strings.ToLower(words[1]) == "kimhyeonseok" {
			username = "도원준"
		} else {
			username = words[1]
		}

		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("오늘부터 지지관계에서 벗어나\n쓱규와(과) %s은(는) 한몸으로 일체가 된다\n쓱규에 대한 공격은 %s에 대한 공격으로 간주한다\n\n", username, username))

		for _, fans := range []string{"70억", "1억", "천만", "백", "한"} {
			builder.WriteString(fmt.Sprintf("세상에 %s 명의 쓱규 팬이 있다면, %s은(는) 그들 중 한 명일 것이다.\n", fans, username))
		}
		builder.WriteString(fmt.Sprintf("세상에 단 한 명의 쓱규 팬도 없다면, %s은(는) 그제서야 이 세상에 없는 것이다.\n\n", username))

		for _, attribute := range []string{"사랑", "빛", "어둠", "삶", "기쁨", "슬픔", "안식", "영혼"} {
			builder.WriteString(fmt.Sprintf("쓱규, %s의 %s.\n", username, attribute))
		}
		builder.WriteString(fmt.Sprintf("쓱규, %s.\n", username))

		embed := &discordgo.MessageEmbed{
			Title:       "지지관계",
			Description: builder.String(),
			Color:       0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("-%s-", username),
			},
		}
		s.ChannelMessageSendEmbed(msg.ChannelID, embed)
	}
}

func test() {

	catFile, err := os.Open("/home/arkaprabham/Documents/Journal_Dev/Golang/github.com/image-op/cat.png")
	if err != nil {
		log.Fatal(err)
	}
	defer catFile.Close()

	// Consider using the general image.Decode as it can sniff and decode any registered image format.
	img, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(img)

	levels := []string{" ", "░", "▒", "▓", "█"}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			level := c.Y / 51 // 51 * 5 = 255
			if level == 5 {
				level--
			}
			fmt.Print(levels[level])
		}
		fmt.Print("\n")
	}
}
