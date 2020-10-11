package notify

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	Token = "NzY0NzY0MDczMjY2MzgwODAx.X4LACw.LE1RE6hQCKkd6h3Ye_HkSBsRbps"
)

func Notificate() {
	dg, err := discordgo.New("Bot "  + Token)
	if err != nil {
		log.Printf("Error logging in: %v", err)
		return
	}
	dg.ChannelMessageSend("764763681224654870", "買えたかも、メール要確認！")
}
