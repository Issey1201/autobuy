package notify

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var Token string

func Notificator() {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")));
		err != nil {
		// .env読めなかった場合の処理
		log.Fatalf("failed to open env: %v", err)
	}
	Token = os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot "  + Token)
	if err != nil {
		log.Printf("Error logging in: %v", err)
		return
	}
	dg.ChannelMessageSend("764763681224654870", "買えたかも、メール要確認！")
}
