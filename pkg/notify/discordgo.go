package notify

import (
	"errors"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var Token string

func Notificator() error {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		fmt.Printf("failed to open env: %v\n", err)
	}
	Token = os.Getenv("DISCORD_TOKEN")

	if len(Token) == 0 {
		return errors.New("there is no token at env")
	}
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		return errors.New(fmt.Sprintf("Error logging in: %v", err))
	}
	if _, err := dg.ChannelMessageSend("764763681224654870", "買えたかも、メール要確認！"); err != nil {
		return errors.New(fmt.Sprintf("Failed to send a message: %v", err))
	}
	return nil
}
